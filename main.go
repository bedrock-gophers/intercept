package main

import (
	"fmt"
	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/pelletier/go-toml"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"log"
	"log/slog"
	"os"
)

type handler struct{}

func (handler) HandleClientPacket(ctx *intercept.Context, pk packet.Packet) {
	fmt.Printf("client packet: %T\n", pk)
	if h, ok := ctx.Val().Handle(); ok {
		h.ExecWorld(func(tx *world.Tx, e world.Entity) {
			p := e.(*player.Player)
			p.Messagef("Received packet: %T", pk)
		})
	}
}
func (handler) HandleServerPacket(ctx *intercept.Context, pk packet.Packet) {
	fmt.Printf("server packet: %T\n", pk)
}

func main() {
	intercept.Hook(handler{})

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(slog.Default())
	if err != nil {
		log.Fatalln(err)
	}

	conf.Listeners = intercept.WrapListeners(conf.Listeners)

	srv := conf.New()
	intercept.Start(srv)
	srv.CloseOnProgramEnd()

	srv.Listen()
	for p := range srv.Accept() {
		_ = p
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
