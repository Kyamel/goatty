

default: build

build:
	./scripts/build.sh

test:
	go test ./...

# Escape sequence conformance against esctest2, ratcheted by a checked-in
# baseline of known failures.
esctest:
	./scripts/esctest.sh

esctest-update:
	./scripts/esctest.sh --update

# The parser is an untrusted-input boundary; this looks for panics and hangs.
fuzz:
	go test ./internal/app/darktile/termutil -run '^$$' -fuzz FuzzParser -fuzztime 300s

# Write results to a file and compare two runs with benchstat before and after
# a change.
bench:
	go test ./internal/app/darktile/termutil -run '^$$' -bench Parser -benchmem -count=8

.PHONY: default build test esctest esctest-update fuzz bench
