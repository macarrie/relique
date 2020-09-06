build:
	cargo build --verbose --all --release

check:
	cargo clippy --all-targets --all-features -- -D warnings

test:
	cargo test --verbose --all