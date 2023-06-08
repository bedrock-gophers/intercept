package packethandler

import (
	"errors"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft"
	"net"
)

type PacketListener struct {
	incoming chan *Conn
}

func NewPacketListener() *PacketListener {
	return &PacketListener{
		incoming: make(chan *Conn),
	}
}

type Listener struct {
	*minecraft.Listener
	pk *PacketListener
}

func (pkt *PacketListener) Listen(conf *server.Config) {
	conf.Listeners = nil
	conf.Listeners = append(conf.Listeners, func(_ server.Config) (server.Listener, error) {
		l, err := minecraft.ListenConfig{
			StatusProvider:       minecraft.NewStatusProvider(conf.Name),
			ResourcePacks:        conf.Resources,
			TexturePacksRequired: conf.ResourcesRequired,
			AcceptedProtocols:    []minecraft.Protocol{},
			FlushRate:            -1,
		}.Listen("raknet", ":19132")
		if err != nil {
			return nil, err
		}

		return &Listener{
			Listener: l,
			pk:       pkt,
		}, nil
	})
}

var errListenerClosed = errors.New("use of closed listener")

func (pkt *PacketListener) Accept() (*Conn, error) {
	c, ok := <-pkt.incoming
	if !ok {
		return nil, errors.New("listener closed")
	}
	return c, nil
}

func (l *Listener) Accept() (session.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, &net.OpError{Op: "accept", Net: "minecraft", Addr: l.Listener.Addr(), Err: errListenerClosed}
	}
	conn := NewConn(c.(*minecraft.Conn))
	l.pk.incoming <- conn
	return conn, nil
}

// Disconnect disconnects a connection from the Listener with a reason.
func (l *Listener) Disconnect(conn session.Conn, reason string) error {
	return l.Listener.Disconnect(conn.(*Conn).Conn, reason)
}

// Close closes the Listener.
func (l *Listener) Close() error {
	_ = l.Listener.Close()
	close(l.pk.incoming)
	return nil
}
