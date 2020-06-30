package plugins

import (
	"context"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/tfplugin5"
	"github.com/hashicorp/terraform-provider-kubernetes-alpha/provider"
	"google.golang.org/grpc"
)

// func init() {
// 	defer provider.InitDevLog()()
// 	provider.Dlog.Printf("Starting up")
// 	provider.Serve()
// }

func init() {
	defer provider.InitDevLog()()
	handshake := goplugin.HandshakeConfig{
		ProtocolVersion: 5,
		// The magic cookie values should NEVER be changed.
		MagicCookieKey:   "TF_PLUGIN_MAGIC_COOKIE",
		MagicCookieValue: "d602bf8f470bc67ca7faa0386276bbdd4330efaf76d1a219cb4d6991ca9872b2",
	}

	serveConfig := &goplugin.ServeConfig{
		HandshakeConfig: handshake,
		GRPCServer:      goplugin.DefaultGRPCServer,
		Plugins: goplugin.PluginSet{
			"provider": &grpcPlugin{
				providerServer: &provider.RawProviderServer{},
			},
		},
	}

	exec := func() {
		plugin.Serve(&plugin.ServeOpts{
			ServeConfig: serveConfig,
		})
	}
	KnownPlugins["terraform-provider-kubernetes-alpha"] = exec
}

type grpcPlugin struct {
	goplugin.Plugin
	providerServer *provider.RawProviderServer
}

func (p *grpcPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	tfplugin5.RegisterProviderServer(s, p.providerServer)
	return nil
}

func (p *grpcPlugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	panic("This is a plugin - it cannot implement GRPCClient")
}
