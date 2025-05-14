package main

import (
	"log/slog"
	"os"

	"evanmoelter.com/tailscale-authentik-webfinger/cmd"
	"gabe565.com/utils/cobrax"
)

var version = "beta"

func main() {
	root := cmd.New(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
