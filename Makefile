build:
	cargo build --verbose --all --release

check:
	cargo clippy --workspace -- -Wclippy::all -Wclippy::pedantic -D warnings

test:
	cargo test --verbose --all