package encryptedStore

import (
	"custodyEthereum/internal/models"
	"encoding/json"
	"github.com/awnumar/memguard"
	"log"
	"os"
)

func RunStoreRoutine(store *models.ServerStore) {
	for {
		select {
		case action := <-store.ActionChannel:
			switch action.Action {
			case "add":
				store.Mux.Lock()

				// Create safesecret
				secret := &models.SafeSecret{
					Path:        action.SecretPath,
					Description: action.Description,
					Data:        action.SecretData.Seal(),
					Roles:       []string{"root", action.Role},
				}

				store.Store.Secrets = append(store.Store.Secrets, secret)

				store.Store.SecretPathList = append(store.Store.SecretPathList, action.SecretPath)

				updateStoredFile(store)
				store.Mux.Unlock()

				action.Response <- secret

			case "remove":
				store.Mux.Lock()
				for i, secret := range store.Store.Secrets {
					if secret.Path == action.SecretPath {
						store.Store.Secrets = append(store.Store.Secrets[:i], store.Store.Secrets[i+1:]...)
						break
					}
				}

				for i, path := range store.Store.SecretPathList {
					if path == action.SecretPath {
						store.Store.SecretPathList = append(store.Store.SecretPathList[:i], store.Store.SecretPathList[i+1:]...)
						break
					}
				}

				updateStoredFile(store)
				store.Mux.Unlock()

			case "update":
				store.Mux.Lock()
				for i, secret := range store.Store.Secrets {
					if secret.Path == action.SecretPath {
						store.Store.Secrets[i].Data = action.SecretData.Seal()
						if action.Description != "" {
							store.Store.Secrets[i].Description = action.Description
						}
						break
					}
				}
				updateStoredFile(store)
				store.Mux.Unlock()

			}

		case sign := <-store.SignChannel:
			store.Mux.Lock()

			secretPath := sign.TransactionStoreName + "/" + sign.TransactionSecret

			var sec *models.SafeSecret
			for _, secret := range store.Store.Secrets {
				if secret.Path == secretPath {
					sec = secret
					break
				}
			}

			if sec == nil {
				store.Mux.Unlock()
				continue
			}

		}
	}
}

func updateStoredFile(store *models.ServerStore) {
	// Marshal the store entry to a string
	storeString, err := json.Marshal(store)
	if err != nil {
		panic(err)
	}

	keyBuf, err := store.SafeKey.Open()
	if err != nil {
		log.Println(err)
		return
	}

	// Get unsafe buffer from mem-guard
	unsafeKey := keyBuf.ByteArray32()

	encStore, err := Encrypt(storeString, unsafeKey)

	memguard.WipeBytes(unsafeKey[:])
	keyBuf.Destroy()

	err = os.WriteFile(store.Path, encStore, 0644)

	memguard.WipeBytes(encStore)
	memguard.WipeBytes(storeString)
}
