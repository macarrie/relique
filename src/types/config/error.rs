use std::fmt;

#[derive(Debug, Clone, PartialEq)]
pub enum ErrorLevel {
    Warning,
    Critical,
}

#[derive(Debug, Clone)]
pub struct Error {
    pub key: String,
    pub level: ErrorLevel,
    pub desc: String,
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.write_fmt(format_args!(
            "[{:?}] {} (key: '{}')",
            self.level, self.desc, self.key
        ))
    }
}
