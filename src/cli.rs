use clap::{App, Arg, SubCommand};

pub fn get_app() -> App<'static, 'static> {
    App::new("relique")
        .version(env!("CARGO_PKG_VERSION"))
        .author("macarrie")
        .about("Backup utility based on librsync")
        .arg(
            Arg::with_name("debug")
                .short("d")
                .long("debug")
                .help("Sets log level to debug"),
        )
        .arg(
            Arg::with_name("config")
                .short("c")
                .long("config")
                .value_name("FILE")
                .help("Sets a custom config file")
                .takes_value(true),
        )
        .subcommand(
            SubCommand::with_name("server")
                .about("Controls relique server features")
                .subcommand(SubCommand::with_name("start").about("Start relique server")),
        )
        .subcommand(
            SubCommand::with_name("client")
                .about("Controls relique client features")
                .subcommand(SubCommand::with_name("start").about("Start relique server"))
                .subcommand(SubCommand::with_name("stop").about("Stop relique server")),
        )
        .subcommand(
            SubCommand::with_name("jobs")
                .about("View backup/restore jobs details")
                .subcommand(
                    SubCommand::with_name("list")
                        .about("List backup jobs")
                        .arg(
                            Arg::with_name("client")
                                .long("client")
                                .value_name("CLIENT_NAME")
                                .help("Client name")
                                .takes_value(true),
                        )
                        .arg(
                            Arg::with_name("module")
                                .short("m")
                                .long("module")
                                .value_name("MODULE_NAME")
                                .help("Backup module name")
                                .takes_value(true),
                        )
                        .arg(
                            Arg::with_name("backup_type")
                                .short("t")
                                .long("type")
                                .value_name("BACKUP_TYPE")
                                .help("Backup type (diff, full)")
                                .takes_value(true),
                        ),
                )
                .subcommand(
                    SubCommand::with_name("show")
                        .about("Show details about a specific jobs")
                        .arg(
                            Arg::with_name("id")
                                .long("id")
                                .value_name("JOB_ID")
                                .help("Job ID")
                                .takes_value(true)
                                .required(true),
                        ),
                ),
        )
        .subcommand(
            SubCommand::with_name("backup")
                .about("Manual backup related commands")
                .subcommand(
                    SubCommand::with_name("start")
                        .about("Start manual backup")
                        .arg(
                            Arg::with_name("client")
                                .long("client")
                                .value_name("CLIENT_NAME")
                                .help("Client to back up")
                                .takes_value(true)
                                .required(true),
                        )
                        .arg(
                            Arg::with_name("module")
                                .short("m")
                                .long("module")
                                .value_name("MODULE_NAME")
                                .help("Backup module to use")
                                .takes_value(true)
                                .required(true),
                        )
                        .arg(
                            Arg::with_name("backup_type")
                                .short("t")
                                .long("type")
                                .value_name("BACKUP_TYPE")
                                .help("Backup type (diff, full)")
                                .takes_value(true)
                                .required(true),
                        ),
                )
                .subcommand(SubCommand::with_name("list").about("List backups")),
        )
        .subcommand(
            SubCommand::with_name("restore")
                .about("Manual restore related commands")
                .subcommand(SubCommand::with_name("start").about("Start manual restore")),
        )
}
