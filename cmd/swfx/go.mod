module github.com/b7c/swfx/cmd/swfx

go 1.21.1

require (
	github.com/gabriel-vasile/mimetype v1.4.3
	github.com/spf13/cobra v1.7.0
	github.com/b7c/swfx v0.0.0-00010101000000-000000000000
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.17.0 // indirect
)

replace github.com/b7c/swfx/cmd => ./cmd

replace github.com/b7c/swfx => ../../
