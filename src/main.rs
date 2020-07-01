use std::path::Path;

use log::*;

mod app;
mod client;
mod config;
mod logger;
mod server;
mod types;

fn main() {
    let relique = app::get_app();
    let matches = relique.get_matches();

    let debug_enabled = matches.is_present("debug");
    logger::init(debug_enabled).unwrap();

    if let Some(server_cmd) = matches.subcommand_matches("server") {
        let config_path = matches.value_of("config").unwrap_or("server.toml");
        let cfg = config::load(Path::new(config_path), true).unwrap();

        if let Some(start_cmd) = server_cmd.subcommand_matches("start") {
            server::start(cfg).unwrap();
        }

        return;
    }

    if let Some(client_cmd) = matches.subcommand_matches("client") {
        let config_path = matches.value_of("config").unwrap_or("client.toml");
        let cfg = config::load(Path::new(config_path), false).unwrap();

        if let Some(start_cmd) = client_cmd.subcommand_matches("start") {
            client::start(cfg).unwrap();
        }

        return;
    }

    app::get_app().print_long_help().unwrap();
}
