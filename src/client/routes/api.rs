use crate::client::client_daemon::ClientDaemon;
use crate::types;
use actix_web::{get, post, web, HttpResponse, Responder};
use log::*;
use std::sync::Mutex;

#[post("/config")]
pub async fn config(
    state: web::Data<Mutex<ClientDaemon>>,
    client_cfg: web::Json<types::config::Client>,
) -> impl Responder {
    let mut state = state.lock().unwrap();
    let client_config_opt = (*state).client_config.clone();
    let current_config_version = match client_config_opt {
        None => None,
        Some(client_config) => client_config.config_version,
    };

    if let Some(version) = current_config_version {
        if Some(version) == client_cfg.config_version {
            return HttpResponse::Ok();
        }
    }

    info!("Replacing current client version with version received from relique server");
    (*state).client_config = Some(client_cfg.into_inner());

    HttpResponse::Ok()
}

#[get("/config/version")]
pub async fn get_config_version(state: web::Data<Mutex<ClientDaemon>>) -> impl Responder {
    let state = state.lock().unwrap();
    let client_config_opt = (*state).client_config.clone();
    let config_version = match client_config_opt {
        None => types::config::ConfigVersion { version: None },
        Some(client_config) => types::config::ConfigVersion {
            version: client_config.config_version,
        },
    };

    HttpResponse::Ok().json(config_version)
}
