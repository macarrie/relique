use std::error::Error;
use std::fs::File;

use simplelog::*;

pub fn init(debug_enabled: bool) -> Result<(), Box<dyn Error>> {
    let log_level = if debug_enabled {
        LevelFilter::Debug
    } else {
        LevelFilter::Info
    };

    let log_config = ConfigBuilder::new().set_time_to_local(true).build();

    CombinedLogger::init(vec![
        TermLogger::new(log_level, log_config.clone(), TerminalMode::Mixed),
        WriteLogger::new(
            log_level,
            log_config.clone(),
            File::create("/var/log/relique/relique.log").unwrap(),
        ),
    ])?;

    Ok(())
}
