package models

import "github.com/awnumar/memguard"

const (
	KeySize   = 32
	NonceSize = 24
)

type User struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type Secret struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Data        string   `json:"data"`
	Roles       []string `json:"roles"`
}

type SafeSecret struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Data        *memguard.LockedBuffer `json:"data"`
	Roles       []string               `json:"roles"`
}

type StoreEntry struct {
	ID          string   `json:"id"`
	Secrets     []Secret `json:"secrets"`
	Description string   `json:"description"`
}

type SafeStoreEntry struct {
	ID          string       `json:"id"`
	Secrets     []SafeSecret `json:"secrets"`
	Description string       `json:"description"`
}

type ServerStore struct {
	Store         StoreEntry `json:"store"`
	Key           [KeySize]byte
	SafeKey       *memguard.LockedBuffer
	ActionChannel chan *StoreAction
}

type StoreAction struct {
	Action string
}
