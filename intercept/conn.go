package intercept

import (
    "fmt"
    "github.com/df-mc/dragonfly/server/event"
    "github.com/df-mc/dragonfly/server/player"
    "github.com/df-mc/dragonfly/server/session"
    "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type conn struct {
    session.Conn
    p *player.Player
}

func (c *conn) ReadPacket() (packet.Packet, error) {
    pkt, err := c.Conn.ReadPacket()
    if err != nil {
        return pkt, err
    }

    ctx := event.C()
    for _, h := range handlers {
        h.HandleClientPacket(ctx, c.p, pkt)
    }

    if ctx.Cancelled() {
        return NopPacket{}, nil
    }
    return pkt, err
}

func (c *conn) WritePacket(pk packet.Packet) error {
    ctx := event.C()
    for _, h := range handlers {
        h.HandleServerPacket(ctx, c.p, pk)
    }

    if ctx.Cancelled() {
        return nil
    }
    return c.Conn.WritePacket(pk)
}

func Intercept(p *player.Player) {
    s := playerSession(p)

    c := fetchPrivateField[session.Conn](s, "conn")
    cn := &conn{c, p}
    updatePrivateField[session.Conn](s, "conn", cn)
}
