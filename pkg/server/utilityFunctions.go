package server

import (
	"custodyEthereum/internal/models"
	"errors"
)

func (s *Server) preliminaryRequestChecks(storeName string, secretName string, role string, create bool) (*models.ServerStore, error) {
	if !s.AvailableStores[storeName] {
		return nil, errors.New("Store does not exist")
	}

	store := s.Stores[storeName]

	allowedRole := false
	for _, storeRole := range store.Roles {
		if storeRole == role {
			allowedRole = true
			break
		}
	}

	if !allowedRole {
		return nil, errors.New("User role not allowed to access store")
	}

	path := storeName + "/" + secretName

	if create && s.RoleSecrets[role][path].Path != "" {
		return nil, errors.New("Secret already exists")
	} else if !create && s.RoleSecrets[role][path].Path == "" {
		return nil, errors.New("Secret does not exist")
	}

	return store, nil
}

func (s *Server) removeStoreFromServer(store *models.ServerStore) {
	delete(s.AvailableStores, store.Store.ID)
	delete(s.Stores, store.Store.ID)

	for _, role := range s.Roles {
		for _, secret := range store.Store.Secrets {
			delete(s.RoleSecrets[role], secret.Path)
		}
	}
}
