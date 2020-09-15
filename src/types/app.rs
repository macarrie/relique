use actix_web::web;
use anyhow::Result;
use chan::chan_select;
use chan::Receiver;
use chan_signal::Signal;
use ticker::Ticker;

use crate::types::config;
use log::*;
use std::sync::RwLock;
use std::time::Duration;

#[derive(Debug)]
pub enum Stopping {
    /// The run loop should halt
    Yes,

    /// The run loop should continue
    No,
}

/// The application; domain-specific program logic
pub trait ReliqueApp: Sized {
    /// Create a new instance given the options and config
    fn new(_: config::Config) -> Result<Self>;

    /// Called repeatedly in the main loop of the application.
    fn loop_func(&mut self) -> Result<Stopping>;

    /// Which signal the application is interested in receiving.
    /// By default, only INT and TERM are blocked and handled.
    fn signals() -> &'static [Signal] {
        static SIGNALS: &[Signal] = &[Signal::INT, Signal::TERM];
        SIGNALS
    }

    /// Handle a received signal
    fn received_signal(&mut self, _: Signal) -> Result<Stopping> {
        Ok(Stopping::Yes)
    }

    /// Called when the application is shutting down
    fn shutdown(&mut self) -> Result<()> {
        Ok(())
    }
}

/// Run an Application of type T
///
/// `run` creates an application from `opts` and `config`. A run loop is entered
/// where `run_once` is repeatedly called on the `T`. Between calls, any
/// arriving signals are checked for and passed to the application via
/// `received_signal`.
pub fn run<T>(app: web::Data<RwLock<T>>, signal: Receiver<Signal>) -> Result<()>
where
    T: ReliqueApp,
{
    let mut ticker = Ticker::new(0.., Duration::from_secs(10)).into_iter();
    'main: loop {
        if let Stopping::Yes = app.write().unwrap().loop_func()? {
            break;
        }

        // Handle any and all pending signals.
        loop {
            chan_select! {
                default => { break; },
                signal.recv() -> sig => {
                    let stopping = sig.map(|s| app.write().unwrap().received_signal(s));
                    if let Some(s) = stopping {
                        if let Stopping::Yes = s.unwrap_or(Stopping::Yes) {
                            info!("Signal ({:?}) received. Shutting down app", sig.unwrap());
                            break 'main;
                        }
                    }
                },
            }
        }

        ticker.next();
    }

    app.write().unwrap().shutdown()?;
    Ok(())
}
