use actix_web::dev::Server;
use actix_web::middleware::Logger;
use actix_web::{web, App, HttpServer};
use log::*;

use crate::client::client_daemon::ClientDaemon;
use crate::lib;
use crate::types;
use crate::types::app::ReliqueApp;

use anyhow::Result;
use std::sync::RwLock;
use std::thread;

mod client_daemon;
mod routes;

#[actix_rt::main]
pub async fn start(cfg: types::config::Config) -> Result<()> {
    info!(
        "Starting relique client on port {}",
        cfg.port.unwrap_or_default()
    );

    let app = web::Data::new(RwLock::new(ClientDaemon::new(cfg.clone(), None).unwrap()));
    let signal = chan_signal::notify(ClientDaemon::signals());

    let app_state = web::Data::clone(&app);
    let app_thread = thread::spawn(move || {
        types::app::run::<ClientDaemon>(app_state, signal).unwrap();
    });

    let http_state = web::Data::clone(&app);

    let http_server = start_http_server::<ClientDaemon>(&cfg, http_state).unwrap();

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
            .service(
                web::scope("/api/v1")
                    .service(routes::api::get_config_version)
                    .service(routes::api::config),
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
