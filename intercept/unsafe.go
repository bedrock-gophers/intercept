package intercept

import (
    "github.com/df-mc/dragonfly/server/player"
    "github.com/df-mc/dragonfly/server/session"
    "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
    "reflect"
    "unsafe"
)

// updatePrivateField sets a private field of a session to the value passed.
func updatePrivateField[T any](s *session.Session, name string, value T) {
    reflectedValue := reflect.ValueOf(s).Elem()
    privateFieldValue := reflectedValue.FieldByName(name)

    privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

    privateFieldValue.Set(reflect.ValueOf(value))
}

// fetchPrivateField fetches a private field of a session.
func fetchPrivateField[T any](s *session.Session, name string) T {
    reflectedValue := reflect.ValueOf(s).Elem()
    privateFieldValue := reflectedValue.FieldByName(name)
    privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

    return privateFieldValue.Interface().(T)
}

// noinspection ALL
//
//go:linkname playerSession github.com/df-mc/dragonfly/server/player.(*Player).session
func playerSession(*player.Player) *session.Session

// noinspection ALL
//
//go:linkname sessionHandlePacket github.com/df-mc/dragonfly/server/session.(*Session).handlePacket
func sessionHandlePacket(s *session.Session, pk packet.Packet) error
