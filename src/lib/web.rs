use crate::types;

use actix_web::client::{ClientBuilder, Connector};
use anyhow::Result;
use openssl::ssl::{
    SslAcceptor, SslAcceptorBuilder, SslConnector, SslFiletype, SslMethod, SslVerifyMode,
};

pub fn get_actix_ssl_builder(cfg: &types::config::Config) -> Result<SslAcceptorBuilder> {
    let mut builder = SslAcceptor::mozilla_intermediate(SslMethod::tls())?;
    builder.set_private_key_file(cfg.ssl_key.as_ref().unwrap(), SslFiletype::PEM)?;
    builder.set_certificate_chain_file(cfg.ssl_cert.as_ref().unwrap())?;

    Ok(builder)
}

pub fn get_http_client(check_ssl_cert: Option<bool>) -> Result<actix_web::client::Client> {
    let mut ssl_connector_builder = SslConnector::builder(SslMethod::tls_client())?;
    if !check_ssl_cert.unwrap_or(true) {
        ssl_connector_builder.set_verify(SslVerifyMode::NONE);
    }
    let ssl_connector = ssl_connector_builder.build();

    let connector = Connector::new().ssl(ssl_connector).finish();
    let client = ClientBuilder::new().connector(connector).finish();

    Ok(client)
}
