use actix_web::{get, HttpResponse, Responder};

#[get("/index")]
pub async fn index() -> impl Responder {
    HttpResponse::Ok().body("Hey there!")
}
