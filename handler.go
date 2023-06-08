package packethandler

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Handler interface {
	HandleServerPacket(ctx *event.Context, pk packet.Packet)
	HandleClientPacket(ctx *event.Context, pk packet.Packet)
}

type NopHandler struct{}

// Comp time check to ensure that NopHandler implements Handler.
var _ Handler = (*NopHandler)(nil)

func (h NopHandler) HandleServerPacket(*event.Context, packet.Packet) {}
func (h NopHandler) HandleClientPacket(*event.Context, packet.Packet) {}

type Conn struct {
	*minecraft.Conn
	h Handler
}

func NewConn(c *minecraft.Conn) *Conn {
	cn := &Conn{Conn: c, h: NopHandler{}}
	return cn
}

func (c *Conn) Handle(h Handler) {
	if h == nil {
		h = NopHandler{}
	}
	c.h = h
}

func (c *Conn) WritePacket(pk packet.Packet) error {
	ctx := event.C()
	c.h.HandleServerPacket(ctx, pk)
	if ctx.Cancelled() {
		return nil
	}
	return c.Conn.WritePacket(pk)
}

func (c *Conn) ReadPacket() (packet.Packet, error) {
	pk, err := c.Conn.ReadPacket()
	if err != nil {
		return nil, err
	}
	ctx := event.C()
	c.h.HandleClientPacket(ctx, pk)
	if ctx.Cancelled() {
		return nil, nil
	}
	return pk, nil
}
