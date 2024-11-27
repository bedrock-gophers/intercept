package intercept

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var (
	handlers []Handler
)

type Handler interface {
	HandleClientPacket(ctx *event.Context[*player.Player], pk packet.Packet)
	HandleServerPacket(ctx *event.Context[*player.Player], pk packet.Packet)
}

func Hook(handler Handler) {
	handlers = append(handlers, handler)
}
