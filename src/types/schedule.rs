use anyhow::Result;
use chrono::{Datelike, Duration, Local, NaiveTime, Weekday};
use lazy_static::*;
use log::*;
use regex::Regex;
use serde::de::{self, Visitor};
use serde::{Deserialize, Deserializer, Serialize, Serializer};
use std::fmt;
use std::vec::Vec;

#[derive(Clone, Debug)]
pub struct Bounds(Vec<(chrono::NaiveTime, chrono::NaiveTime)>);

struct BoundVisitor;

impl<'de> Visitor<'de> for BoundVisitor {
    type Value = Bounds;

    fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
        formatter.write_str("formatted 24h time ranges separated by commas (hh:mm-hh:mm)")
    }

    fn visit_str<E>(self, value: &str) -> Result<Self::Value, E>
    where
        E: de::Error,
    {
        lazy_static! {
            static ref RE_SCHEDULE_BOUNDS: Regex =
                Regex::new(r"(?P<start>[0-9]{2}:[0-9]{2})-(?P<stop>[0-9]{2}:[0-9]{2})").unwrap();
        }

        let mut bounds: Vec<(chrono::NaiveTime, chrono::NaiveTime)> = Vec::new();
        for capture in RE_SCHEDULE_BOUNDS.captures_iter(&value) {
            let start_str = capture["start"].to_string();
            let stop_str = capture["stop"].to_string();
            let start = NaiveTime::parse_from_str(&start_str, "%H:%M");
            let stop = NaiveTime::parse_from_str(&stop_str, "%H:%M");

            if start.is_err() || stop.is_err() {
                let error_msg = format!("Could not parse 24h time from schedule ranges declaration in configuration file: '{}-{}'", start_str, stop_str);
                return Err(E::custom(error_msg));
            }

            bounds.push((start.unwrap(), stop.unwrap()));
        }

        Ok(Bounds(bounds))
    }
}

impl<'de> Deserialize<'de> for Bounds {
    fn deserialize<D>(deserializer: D) -> Result<Bounds, D::Error>
    where
        D: Deserializer<'de>,
    {
        deserializer.deserialize_str(BoundVisitor)
    }
}

impl Serialize for Bounds {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let serialized = self
            .0
            .iter()
            .map(|range| format!("{}-{}", range.0.format("%H:%M"), range.1.format("%H:%M")))
            .collect::<Vec<String>>()
            .join(", ");
        serializer.serialize_str(&serialized)
    }
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Schedule {
    pub name: String,
    pub monday: Option<Bounds>,
    pub tuesday: Option<Bounds>,
    pub wednesday: Option<Bounds>,
    pub thursday: Option<Bounds>,
    pub friday: Option<Bounds>,
    pub saturday: Option<Bounds>,
    pub sunday: Option<Bounds>,
}

impl Schedule {
    pub fn is_active(&self) -> bool {
        let weekday = Local::now().date().weekday();
        let schedule_entry: &Option<Bounds> = match weekday {
            Weekday::Mon => &self.monday,
            Weekday::Tue => &self.tuesday,
            Weekday::Wed => &self.wednesday,
            Weekday::Thu => &self.thursday,
            Weekday::Fri => &self.friday,
            Weekday::Sat => &self.saturday,
            Weekday::Sun => &self.sunday,
        };

        let current_time = Local::now().time();
        let previous_loop_time = current_time - Duration::seconds(10);

        if schedule_entry.is_none() {
            return false;
        }

        for range in &schedule_entry.as_ref().unwrap().0 {
            if current_time > range.0 && current_time < range.1 {
                if !(previous_loop_time > range.0 && previous_loop_time < range.1) {
                    info!(
                        "Entering schedule '{}': '{}-{}'",
                        self.name, range.0, range.1
                    );
                }

                return true;
            } else if previous_loop_time > range.0 && previous_loop_time < range.1 {
                info!(
                    "Exiting schedule '{}': '{}-{}'",
                    self.name, range.0, range.1
                );
            }
        }

        false
    }
}
