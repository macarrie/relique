use crate::server::server_daemon::ServerDaemon;
use actix_web::{get, web, HttpResponse, Responder};
use std::sync::Mutex;

#[get("/index")]
pub async fn index(state: web::Data<Mutex<ServerDaemon>>) -> impl Responder {
    let mut state = state.lock().unwrap();
    (*state).config.port = Some(1234);

    HttpResponse::Ok().body("State changed")
}
