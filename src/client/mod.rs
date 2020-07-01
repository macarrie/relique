use actix_web::middleware::Logger;
use actix_web::{web, App, HttpResponse, HttpServer, Responder};
use log::*;
use std::error::Error;

use crate::config;
use crate::types;
use openssl::ssl::{SslAcceptor, SslFiletype, SslMethod};

mod routes;

#[actix_rt::main]
pub async fn start(cfg: types::config::Config) -> Result<(), Box<dyn Error>> {
    info!(
        "Starting relique client on port {}",
        cfg.port.unwrap_or_default()
    );

    let mut builder = SslAcceptor::mozilla_intermediate(SslMethod::tls()).unwrap();
    builder
        .set_private_key_file(cfg.ssl_key.unwrap_or_default(), SslFiletype::PEM)
        .unwrap();
    builder
        .set_certificate_chain_file(cfg.ssl_cert.unwrap_or_default())
        .unwrap();

    HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())
            .service(web::scope("/api/v1").service(routes::api::index))
    })
    .bind_openssl(
        format!(
            "{bind_addr}:{port}",
            bind_addr = cfg.bind_addr.unwrap_or_default(),
            port = cfg.port.unwrap_or_default()
        ),
        builder,
    )?
    .run()
    .await;

    //TODO: Stop on ctrlC

    Ok(())
}

fn stop() -> Result<(), Box<dyn Error>> {
    warn!("Stopping relique server");

    Ok(())
}
