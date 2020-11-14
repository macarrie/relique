use crate::types::rsync::Delta;
use crate::types::rsync::Signature;
use anyhow::Result;
use std::fs::{rename, File};
use std::io::{BufReader, BufWriter, Cursor};
use std::path::Path;

pub fn get_signature(path: &Path) -> Result<Signature> {
    let file = File::open(path)?;
    let mut buf_reader = BufReader::new(file);

    let mut sig = Vec::new();
    librsync::whole::signature(&mut buf_reader, &mut sig)?;

    Ok(sig)
}

pub fn signatures_match(sig1: &Signature, sig2: &Signature) -> Result<bool> {
    let matching = sig1
        .iter()
        .zip(sig2.iter())
        .filter(|&(a, b)| a == b)
        .count();

    Ok(matching == sig1.len() && matching == sig2.len())
}

pub fn get_delta(path_str: String, sig: Signature) -> Result<Delta> {
    let path = Path::new(&path_str);
    let file = File::open(path)?;
    let mut buf_reader = BufReader::new(file);

    let mut dlt = Vec::new();
    librsync::whole::delta(&mut buf_reader, &mut Cursor::new(sig), &mut dlt)?;

    Ok(dlt)
}

pub fn apply_delta(path_str: String, dlt: Delta, out_path: &Path) -> Result<()> {
    let path = Path::new(&path_str);
    let file = File::open(path)?;
    let mut buf_reader = BufReader::new(file);

    let tmp_path: &Path = Path::new("/tmp/bkp_delta.tmp");
    let tmp_file = File::create(tmp_path)?;
    let mut buf_writer = BufWriter::new(tmp_file);
    librsync::whole::patch(&mut buf_reader, &mut Cursor::new(dlt), &mut buf_writer)?;

    rename(tmp_path, out_path).unwrap();

    Ok(())
}
