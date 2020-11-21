use crate::types::daemon_type::DaemonType;
use futures::executor::block_on;
use log::*;
use std::path::Path;

mod cli;
mod client;
mod config;
mod db;
mod lib;
mod logger;
mod server;
mod types;

#[actix_rt::main]
async fn main() {
    let relique = cli::get_app();
    let matches = relique.get_matches();

    let debug_enabled = matches.is_present("debug");

    match matches.subcommand() {
        ("server", Some(server_cmd)) => {
            logger::init(debug_enabled, Some(DaemonType::Server)).unwrap();
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
        ("client", Some(client_cmd)) => {
            logger::init(debug_enabled, Some(DaemonType::Client)).unwrap();
            let config_path = matches.value_of("config").unwrap_or("client.toml");
            let cfg = config::load(Path::new(config_path), false).unwrap();

            if let Some(_start_cmd) = client_cmd.subcommand_matches("start") {
                client::start(cfg).await.unwrap_or_else(|err| {
                    error!("Could not start relique client: {}", err);
                    std::process::exit(exitcode::SOFTWARE);
                });
            }

            return;
        }
        ("jobs", Some(jobs_cmd)) => {
            logger::init(debug_enabled, None).unwrap();
            let config_path = matches.value_of("config").unwrap_or("server.toml");
            let cfg = config::load(Path::new(config_path), false).unwrap();
            server::ping(cfg.clone()).unwrap_or_else(|e| {
                error!("Relique server must be running. Could not contact local relique server: '{err}'", err = e);
                std::process::exit(exitcode::SOFTWARE);
            });

            match jobs_cmd.subcommand() {
                ("list", Some(list_cmd)) => {
                    // TODO: Check if server is running
                    // TODO: Create and start backup job
                    // TODO: Check if client, module and backup type are valid
                    // TODO: Change backup_type to enum to check values automatically
                    let jobs = block_on(lib::jobs::list(
                        cfg.clone(),
                        list_cmd
                            .value_of("client_name")
                            .and_then(|val| Some(val.to_string())),
                        list_cmd
                            .value_of("module_name")
                            .and_then(|val| Some(val.to_string())),
                        list_cmd
                            .value_of("backup_type")
                            .and_then(|val| Some(val.to_string())),
                    ));
                    println!("Jobs: {:?}", jobs);
                }
                ("show", Some(show_cmd)) => println!(
                    "Job '{id}' details",
                    id = show_cmd.value_of("id").unwrap_or("none"),
                ),
                _ => cli::get_app().print_long_help().unwrap(),
            };

            return;
        }
        ("backup", Some(backup_cmd)) => {
            let config_path = matches.value_of("config").unwrap_or("server.toml");
            let cfg = config::load(Path::new(config_path), false).unwrap();

            if let Some(start_cmd) = backup_cmd.subcommand_matches("start") {
                // TODO: Check if server is running
                // TODO: Create and start backup job
                // TODO: Check if client, module and backup type are valid
                // TODO: Change backup_type to enum to check values automatically
                println!("Start backup on client '{client}', module'{module}', backup type '{type}'",
                    client = start_cmd.value_of("client").unwrap(),
                     module = start_cmd.value_of("module").unwrap(),
                     type = start_cmd.value_of("backup_type").unwrap(),
                )
            }

            return;
        }
        _ => cli::get_app().print_long_help().unwrap(),
    };
}
