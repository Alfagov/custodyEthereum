package encryptedStore

import (
	"crypto/rand"
	"custodyEthereum/configs"
	"custodyEthereum/internal/models"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/awnumar/memguard"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/shamir"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"os"
)

func generateRandomBytes(size int) ([]byte, error) {
	rBytes := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, rBytes)
	if err != nil {
		return nil, err
	}
	return rBytes, nil
}

func encrypt(plaintext []byte, key *[models.KeySize]byte) ([]byte, error) {
	var nonce [models.NonceSize]byte
	nonceBytes, err := generateRandomBytes(models.NonceSize)
	if err != nil {
		return nil, err
	}

	copy(nonce[:], nonceBytes)

	ciphertext := secretbox.Seal(nonce[:], plaintext, &nonce, key)
	return ciphertext, nil
}

func decrypt(ciphertext []byte, key *[models.KeySize]byte) ([]byte, error) {

	var nonce [models.NonceSize]byte
	copy(nonce[:], ciphertext[:models.NonceSize])

	plaintext, ok := secretbox.Open(nil, ciphertext[models.NonceSize:], &nonce, key)
	if !ok {
		return nil, errors.New("decryption failed")
	}
	return plaintext, nil
}

func CreateNewBoxStore(name string, threshold int, total int, description string) [][]byte {

	keyBytes, _ := generateRandomBytes(models.KeySize)
	var key [models.KeySize]byte

	copy(key[:], keyBytes)

	store := models.StoreEntry{
		ID:          uuid.New().String(),
		Secrets:     []models.Secret{},
		Description: description,
	}

	storeString, err := json.Marshal(store)
	if err != nil {
		panic(err)
	}

	encStore, err := encrypt(storeString, &key)

	shares, err := shamir.Split(keyBytes, total, threshold)
	if err != nil {
		panic(err)
	}

	path := configs.GlobalViper.GetString("server.basepath") + "/" + name + ".json"

	err = os.WriteFile(path, encStore, 0644)
	if err != nil {
		panic(err)
	}

	return shares
}

func OpenBoxStore(name string, key [models.KeySize]byte) (models.SafeStoreEntry, error) {

	path := configs.GlobalViper.GetString("server.basepath") + name + ".json"
	storeBytes, err := os.ReadFile(path)
	if err != nil {
		return models.SafeStoreEntry{}, err
	}

	decStoreString, err := decrypt(storeBytes, &key)
	if err != nil {
		return models.SafeStoreEntry{}, err
	}

	var unsafeEncStore models.StoreEntry
	err = json.Unmarshal(decStoreString, &unsafeEncStore)
	if err != nil || unsafeEncStore.ID == "" {
		return models.SafeStoreEntry{}, err
	}
	safeEncStore := models.SafeStoreEntry{
		ID:          unsafeEncStore.ID,
		Secrets:     []models.SafeSecret{},
		Description: unsafeEncStore.Description,
	}
	for _, secret := range unsafeEncStore.Secrets {
		k, _ := base64.StdEncoding.DecodeString(secret.Data)
		safeEncStore.Secrets = append(
			safeEncStore.Secrets,
			models.SafeSecret{
				ID:          secret.ID,
				Description: secret.Description,
				Data:        memguard.NewBufferFromBytes(k),
				Roles:       secret.Roles,
			})
	}

	return safeEncStore, nil
}
