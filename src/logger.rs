use anyhow::Result;
use std::fs::File;

use crate::types::daemon_type::DaemonType;
use simplelog::*;

pub fn init(debug_enabled: bool, daemon_type: DaemonType) -> Result<()> {
    let log_level = if debug_enabled {
        LevelFilter::Debug
    } else {
        LevelFilter::Info
    };

    let log_config = ConfigBuilder::new().set_time_to_local(true).build();
    let log_file_path = match daemon_type {
        DaemonType::Server => "/var/log/relique/relique-server.log",
        DaemonType::Client => "/var/log/relique/relique-client.log",
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
