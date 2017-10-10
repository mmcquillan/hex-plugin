package hexplugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Handshake Config
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "HexBotPlugin",
	MagicCookieValue: "HexBotPlugin",
}

// call args
type Args struct {
	Debug   bool
	Command string
	Config  map[string]string
}

// Action - interface for the plugin
type Action interface {
	Perform(args Args) string
}

type ActionRPC struct{ client *rpc.Client }

func (g *ActionRPC) Perform(args Args) string {
	var resp string
	err := g.client.Call("Plugin.Action", args, &resp)
	if err != nil {
		// do something with return
		panic(err)
	}

	return resp
}

type ActionRPCServer struct {
	Impl Action
}

func (s *ActionRPCServer) Action(args Args, resp *string) error {
	*resp = s.Impl.Perform(args)
	return nil
}

type HexPlugin struct {
	Impl Action
}

func (p *HexPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ActionRPCServer{Impl: p.Impl}, nil
}

func (HexPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ActionRPC{client: c}, nil
}
