package encryptedStore

import "custodyEthereum/internal/models"

func RunStoreRoutine(store *models.ServerStore) {
	for {
		select {
		case action := <-store.ActionChannel:
			switch action.Action {
			case "add":

			case "remove":

			case "update":

			}
		}
	}
}
