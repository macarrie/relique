use actix_web::dev::Server;
use actix_web::middleware::Logger;
use actix_web::{web, App, HttpServer};
use log::*;

use crate::config as cfgutils;
use crate::db;
use crate::lib;
use crate::types::app;
use crate::types::app::ReliqueApp;
use crate::types::config;

use anyhow::Result;
use server_daemon::ServerDaemon;
use std::sync::RwLock;
use std::thread;

mod routes;
mod server_daemon;

#[actix_rt::main]
pub async fn start(cfg: config::Config) -> Result<()> {
    let cfg_checks = cfgutils::check(&cfg);
    let cfg_critical_errors: Vec<&config::error::Error> = cfg_checks
        .iter()
        .filter(|e| e.level == config::error::ErrorLevel::Critical)
        .collect();

    if cfg_critical_errors.is_empty() {
        info!("Configuration checks passed");
    } else {
        error!("Fatal configuration errors found. Exiting relique");
        std::process::exit(exitcode::CONFIG);
    }

    info!(
        "Starting relique server on port {}",
        cfg.port.unwrap_or_default()
    );

    // TODO: Add logs and check db conn
    let db_pool = db::server::get_pool().await.unwrap_or_else(|e| {
        error!("Could not get database connection pool: '{}'", e);
        std::process::exit(exitcode::DATAERR);
    });
    let app_data = web::Data::new(RwLock::new(ServerDaemon::new(cfg.clone(), Some(db_pool))?));
    let signal = chan_signal::notify(ServerDaemon::signals());

    let app_state = web::Data::clone(&app_data);
    let app_thread = thread::spawn(move || {
        app::run::<ServerDaemon>(app_state, signal).unwrap();
    });

    let http_state = web::Data::clone(&app_data);

    let http_server = start_http_server::<ServerDaemon>(&cfg, http_state)?;

    app_thread.join().unwrap();
    http_server.stop(true).await;

    Ok(())
}

pub fn start_http_server<T: 'static>(
    cfg: &config::Config,
    state: web::Data<RwLock<T>>,
) -> Result<Server>
where
    T: ReliqueApp + Send + Sync,
{
    let builder = lib::web::get_actix_ssl_builder(cfg)?;
    let http_server = HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(web::Data::clone(&state))
            .service(web::scope("/ui").service(routes::ui::index))
            .service(
                web::scope("/api/v1")
                    .service(routes::api::post_backup_jobs_register)
                    .service(routes::api::put_backup_jobs_id_status)
                    .service(routes::api::get_backup_jobs_id_signature)
                    .service(routes::api::post_backup_jobs_id_delta),
            )
    })
    .bind_openssl(
        format!(
            "{bind_addr}:{port}",
            bind_addr = cfg.bind_addr.as_ref().unwrap(),
            port = cfg.port.unwrap_or_default()
        ),
        builder,
    )?
    .run();

    Ok(http_server)
}
