use std::path::Path;

mod app;
mod config;
mod types;

fn main() {
    let relique = app::get_app();
    let matches = relique.get_matches();

    let config_path = matches.value_of("config").unwrap_or("relique.toml");
    println!("Value for config path: {}", config_path);
    let cfg = config::load(Path::new(config_path)).unwrap();
    println!("Config struct: {:?}", cfg);

    let cfg_checks = config::check(cfg);
    let cfg_critical_errors :Vec<&types::config::Error> = cfg_checks.iter()
        .filter(|e| e.level == types::config::ErrorLevel::Critical)
        .collect();

    if cfg_critical_errors.is_empty() {
        // TODO: Log OK
        println!("Configuration checks passed");
    } else {
        // TODO: Log KO
        println!("Fatal configuration errors found. Exiting relique");
        std::process::exit(exitcode::CONFIG);
    }

    if let Some(_m) = matches.subcommand_matches("server") {
        println!("Server subcommand");
        return;
    }

    if let Some(_m) = matches.subcommand_matches("client") {
        println!("Client subcommand");
        return;
    }

    app::get_app().print_long_help().unwrap();
}
