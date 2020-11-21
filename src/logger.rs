use anyhow::Result;
use std::fs::File;

use crate::types::daemon_type::DaemonType;
use simplelog::*;

pub fn init(debug_enabled: bool, daemon_type: Option<DaemonType>) -> Result<()> {
    let log_level = if daemon_type.is_none() {
        LevelFilter::Warn
    } else {
        if debug_enabled {
            LevelFilter::Debug
        } else {
            LevelFilter::Info
        }
    };

    let log_config = ConfigBuilder::new().set_time_to_local(true).build();
    let log_file_path = match daemon_type {
        Some(DaemonType::Server) => "/var/log/relique/relique-server.log",
        Some(DaemonType::Client) => "/var/log/relique/relique-client.log",
        None => "/dev/null",
    };

    CombinedLogger::init(vec![
        TermLogger::new(log_level, log_config.clone(), TerminalMode::Mixed),
        WriteLogger::new(
            log_level,
            log_config,
            File::create(log_file_path).expect("Could not create log file"),
        ),
    ])?;

    Ok(())
}
