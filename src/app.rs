use clap::{Arg, App, SubCommand};

pub fn get_app() -> App<'static, 'static> {
    App::new("relique")
        .version(env!("CARGO_PKG_VERSION"))
        .author("macarrie")
        .about("Backup utility based on librsync")
        .arg(Arg::with_name("config")
            .short("c")
            .long("config")
            .value_name("FILE")
            .help("Sets a custom config file")
            .takes_value(true))

        .subcommand(SubCommand::with_name("server")
            .about("Controls relique server features")
            .subcommand(SubCommand::with_name("start")
                .about("Start relique server")
            )
            .subcommand(SubCommand::with_name("stop")
                .about("Stop relique server")
            )
        )

        .subcommand(SubCommand::with_name("client")
            .about("Controls relique client features")
            .subcommand(SubCommand::with_name("start")
                .about("Start relique server")
            )
            .subcommand(SubCommand::with_name("stop")
                .about("Stop relique server")
            )
        )
}