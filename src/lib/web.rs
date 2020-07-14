use crate::types;

use anyhow::Result;
use openssl::ssl::{SslAcceptor, SslAcceptorBuilder, SslFiletype, SslMethod};
use reqwest::ClientBuilder;
use std::fs::File;
use std::io::Read;

pub fn get_actix_ssl_builder(cfg: &types::config::Config) -> Result<SslAcceptorBuilder> {
    let mut builder = SslAcceptor::mozilla_intermediate(SslMethod::tls())?;
    builder.set_private_key_file(cfg.ssl_key.as_ref().unwrap(), SslFiletype::PEM)?;
    builder.set_certificate_chain_file(cfg.ssl_cert.as_ref().unwrap())?;

    Ok(builder)
}

pub fn get_http_client(
    ssl_cert_path: Option<String>,
    check_ssl_cert: Option<bool>,
) -> Result<reqwest::Client> {
    let mut buf = Vec::new();
    let mut file =
        File::open(ssl_cert_path.unwrap_or_else(|| String::from("/etc/relique/cert.pem")))?;
    file.read_to_end(&mut buf)?;
    let cert = reqwest::Certificate::from_pem(&buf)?;

    let client_builder: ClientBuilder = if !check_ssl_cert.unwrap_or(false) {
        ClientBuilder::new()
            .danger_accept_invalid_certs(true)
            .add_root_certificate(cert)
    } else {
        ClientBuilder::new()
    };

    let client = client_builder.build()?;

    Ok(client)
}
