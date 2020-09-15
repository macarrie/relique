use actix_web::dev::Server;
use actix_web::middleware::Logger;
use actix_web::{web, App, HttpServer};
use log::*;

use crate::types;
use crate::types::app::ReliqueApp;
use crate::{config, lib};

use anyhow::Result;
use server_daemon::ServerDaemon;
use std::sync::RwLock;
use std::thread;

mod routes;
mod server_daemon;

#[actix_rt::main]
pub async fn start(cfg: types::config::Config) -> Result<()> {
    let cfg_checks = config::check(&cfg);
    let cfg_critical_errors: Vec<&types::config::Error> = cfg_checks
        .iter()
        .filter(|e| e.level == types::config::ErrorLevel::Critical)
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

    let app = web::Data::new(RwLock::new(ServerDaemon::new(cfg.clone())?));
    let signal = chan_signal::notify(ServerDaemon::signals());

    let app_state = web::Data::clone(&app);
    let app_thread = thread::spawn(move || {
        types::app::run::<ServerDaemon>(app_state, signal).unwrap();
    });

    let http_state = web::Data::clone(&app);

    let http_server = start_http_server::<ServerDaemon>(&cfg, http_state)?;

    app_thread.join().unwrap();
    http_server.stop(true).await;

    Ok(())
}

pub fn start_http_server<T: 'static>(
    cfg: &types::config::Config,
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
                    .service(routes::api::backup_jobs_register)
                    .service(routes::api::update_backup_jobs_status),
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
