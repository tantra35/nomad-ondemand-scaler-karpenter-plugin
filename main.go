package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func main() {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:        "karpenterprovidergrpc-0.29.2",
		DisableTime: true,
		Level:       hclog.Info,
	})
	hclog.SetDefault(appLogger)

	lclustername := os.Args[1]

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			// This isn't required when using VersionedPlugins
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: map[string]plugin.Plugin{
			"grpc": NewK8sKapenterProviderPlugin(lclustername),
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
