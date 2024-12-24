package intercept

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type conn struct {
	session.Conn
	h *world.EntityHandle
}

func (c *conn) ReadPacket() (packet.Packet, error) {
	pkt, err := c.Conn.ReadPacket()
	if err != nil {
		return pkt, err
	}

	var cancelled bool
	c.h.ExecWorld(func(tx *world.Tx, e world.Entity) {
		p := e.(*player.Player)
		ctx := event.C(p)
		for _, h := range handlers {
			h.HandleClientPacket(ctx, pkt)
		}
		if ctx.Cancelled() {
			cancelled = true
		}
	})

	if cancelled {
		return NopPacket{}, nil
	}
	return pkt, nil
}

func (c *conn) WritePacket(pk packet.Packet) error {
	var cancelled bool
	c.h.ExecWorld(func(tx *world.Tx, e world.Entity) {
		p := e.(*player.Player)
		ctx := event.C(p)
		for _, h := range handlers {
			h.HandleClientPacket(ctx, pk)
		}
		if ctx.Cancelled() {
			cancelled = true
		}
	})

	if cancelled {
		return nil
	}
	return c.Conn.WritePacket(pk)
}

func Intercept(p *player.Player) {
	s := playerSession(p)

	c := fetchPrivateField[session.Conn](s, "conn")
	cn := &conn{c, p.H()}
	updatePrivateField[session.Conn](s, "conn", cn)
}
