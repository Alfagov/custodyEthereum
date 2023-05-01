package models

import (
	"github.com/awnumar/memguard"
	"sync"
)

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
	Path        string            `json:"id"`
	Description string            `json:"description"`
	Data        *memguard.Enclave `json:"data"`
	Roles       []string          `json:"roles"`
}

type StoreEntry struct {
	ID             string    `json:"id"`
	Secrets        []*Secret `json:"secrets"`
	Description    string    `json:"description"`
	SecretPathList []string  `json:"secretPathList"`
}

type SafeStoreEntry struct {
	ID             string        `json:"id"`
	Secrets        []*SafeSecret `json:"secrets"`
	Description    string        `json:"description"`
	SecretPathList []string      `json:"secretPathList"`
}

type ServerStore struct {
	Store         *SafeStoreEntry `json:"store"`
	Path          string
	SafeKey       *memguard.Enclave
	ActionChannel chan *StoreAction
	Roles         []string
	Mux           sync.Mutex
}

type StoreAction struct {
	Action      string
	SecretPath  string
	SecretData  *memguard.LockedBuffer
	Description string
	Role        string
	Response    chan *SafeSecret
}

type ReqStoreAction struct {
	Action       string   `json:"action"`
	StoreName    string   `json:"storeName"`
	SecretName   string   `json:"secretName"`
	SecretData   string   `json:"secretData"`
	Description  string   `json:"description"`
	AllowedRoles []string `json:"allowedRoles"`
}

type UpdateStoreAction struct {
	StoreName string `json:"storeName"`
	Override  bool   `json:"override"`
}
