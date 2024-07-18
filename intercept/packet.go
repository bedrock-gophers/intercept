package intercept

import (
    "github.com/sandertv/gophertunnel/minecraft/protocol"
    "math"
)

type NopPacket struct{}

func (NopPacket) ID() uint32 {
    return math.MaxUint32
}

func (NopPacket) Marshal(_ protocol.IO) {}
