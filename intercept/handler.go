package intercept

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var (
	handlers []Handler
)

type Context = event.Context[*world.EntityHandle]

type Handler interface {
	HandleClientPacket(ctx *Context, pk packet.Packet)
	HandleServerPacket(ctx *Context, pk packet.Packet)
}

func Hook(handler Handler) {
	handlers = append(handlers, handler)
}
