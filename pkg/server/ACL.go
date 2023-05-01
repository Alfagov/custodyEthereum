package server

import (
	"custodyEthereum/internal/models"
	"custodyEthereum/pkg/encryptedStore"
	"github.com/awnumar/memguard"
)

func (s *Server) importStoreWithACL(store *models.SafeStoreEntry, key [models.KeySize]byte) error {

	s.AvailableStores[store.ID] = true
	secureKey := memguard.NewEnclave(key[:])
	s.Stores[store.ID] = &models.ServerStore{
		Store:   store,
		SafeKey: secureKey,
	}

	for _, role := range s.Roles {
		for _, secret := range store.Secrets {
			for _, allowedRole := range secret.Roles {
				if allowedRole == role {
					s.RoleSecrets[role][secret.Path] = secret
				}
			}
		}
	}

	go encryptedStore.RunStoreRoutine(s.Stores[store.ID])

	return nil
}
