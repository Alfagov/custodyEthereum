package server

import (
	"custodyEthereum/internal/models"
	"github.com/awnumar/memguard"
)

func (s *Server) importStoreWithACL(store models.StoreEntry, key [models.KeySize]byte) error {

	s.AvailableStores[store.ID] = true
	secureKey := memguard.NewBufferFromBytes(key[:])
	s.Stores[store.ID] = &models.ServerStore{
		Store:   store,
		SafeKey: secureKey,
	}

	for _, role := range s.Roles {
		for _, secret := range store.Secrets {
			for _, allowedRole := range secret.Roles {
				if allowedRole == role {
					s.RoleSecrets[role][secret.ID] = secret
				}
			}
		}
	}

	return nil
}
