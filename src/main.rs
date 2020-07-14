use crate::types::daemon_type::DaemonType;
use log::*;
use std::path::Path;

mod cli;
mod client;
mod config;
mod lib;
mod logger;
mod server;
mod types;

fn main() {
    let relique = cli::get_app();
    let matches = relique.get_matches();

    let debug_enabled = matches.is_present("debug");

    if let Some(server_cmd) = matches.subcommand_matches("server") {
        logger::init(debug_enabled, DaemonType::Server).unwrap();
        let config_path = matches.value_of("config").unwrap_or("server.toml");
        let cfg = config::load(Path::new(config_path), true).unwrap();

        if let Some(_start_cmd) = server_cmd.subcommand_matches("start") {
            server::start(cfg).unwrap_or_else(|err| {
                error!("Could not start relique server: {}", err);
                std::process::exit(exitcode::SOFTWARE);
            });
        }

        return;
    }

    if let Some(client_cmd) = matches.subcommand_matches("client") {
        logger::init(debug_enabled, DaemonType::Client).unwrap();
        let config_path = matches.value_of("config").unwrap_or("client.toml");
        let cfg = config::load(Path::new(config_path), false).unwrap();

        if let Some(_start_cmd) = client_cmd.subcommand_matches("start") {
            client::start(cfg).unwrap_or_else(|err| {
                error!("Could not start relique client: {}", err);
                std::process::exit(exitcode::SOFTWARE);
            });
        }

        return;
    }

    cli::get_app().print_long_help().unwrap();
}
