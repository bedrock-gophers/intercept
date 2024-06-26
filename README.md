# Intercept
The Intercept library offers a flexible and intuitive way to handle player packets within a Minecraft server. Whether you're managing incoming players or handling events triggered by player actions, PacketHandler provides a powerful toolset for managing network communication.

# Importing Intercept into your project:

You may import Intercept by running the following command:
```bash
go get github.com/bedrock-gophers/intercept
```

## Handling Incoming Players
To handle incoming players using Intercept, you can utilize the following example code:

```go
pk := intercept.NewPacketListener()
// Replace 'c' with your Dragonfly config
pk.Listen(&c, ":19132", []minecraft.Protocol{})
go func() {
    for {
        p, err := pk.Accept()
        if err != nil {
            return
        }
        p.Handle(NewPacketHandler(p))
    }
}()
```
This code sets up a new packet listener and listens for incoming connections on port 19132. When a new connection is accepted, you may give a player your packet handler and start handling their packets.

## Handling Server and Client packets
To handle Server and Client packets using Intercept, you can use the following example code:

```go
// PacketHandler represents our custom packet handler.
type PacketHandler struct {
    c *intercept.Conn
}

// NewPacketHandler returns a new packet handler.
func NewPacketHandler(c *intercept.Conn) *PacketHandler {
    return &PacketHandler{
        c: c,
    }
}

// HandleClientPacket...
func (h *PacketHandler) HandleClientPacket(_ *event.Context, pk packet.Packet) {
    fmt.Printf("new packet sent by client: %#v", pk)
}

// HandleServerPacket ...
func (h *PacketHandler) HandleServerPacket(_ *event.Context, pk packet.Packet) {
    fmt.Printf("new packet sent from server: %#v", pk)
}
```
This code defines a packet handler structure with methods to handle client and server packets. You can use these methods to process packets sent by players or received from the server.

