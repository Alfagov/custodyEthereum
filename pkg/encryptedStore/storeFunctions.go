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

// generateRandomBytes generates a byte slice of random bytes of size `size`.
// It returns an error if it fails to read from the random source.
func generateRandomBytes(size int) ([]byte, error) {
	// create a byte slice of given size
	rBytes := make([]byte, size)
	// fill the slice with random bytes from the random source
	_, err := io.ReadFull(rand.Reader, rBytes)
	if err != nil {
		return nil, err
	}
	return rBytes, nil
}

// Encrypt encrypts a plaintext byte slice with a given key byte array pointer
// using the NaCl secretbox.Seal function.
// It returns an error if it fails to generate a nonce or if encryption fails.
func Encrypt(plaintext []byte, key *[models.KeySize]byte) ([]byte, error) {
	// create a nonce byte array of the NaCl library's recommended size
	var nonce [models.NonceSize]byte
	// generate a random nonce byte slice using the generateRandomBytes function
	nonceBytes, err := generateRandomBytes(models.NonceSize)
	if err != nil {
		return nil, err
	}
	// copy the first N bytes of the nonce byte slice to the nonce byte array,
	// where N is the NaCl library's recommended size for nonces
	copy(nonce[:], nonceBytes)
	// Encrypt the plaintext using the nonce and key with the NaCl secretbox.Seal function
	ciphertext := secretbox.Seal(nonce[:], plaintext, &nonce, key)
	return ciphertext, nil
}

// decrypt decrypts a ciphertext byte slice with a given key byte array pointer
// using the NaCl secretbox.Open function.
// It returns an error if it fails to extract the nonce or if decryption fails.
func decrypt(ciphertext []byte, key *[models.KeySize]byte) ([]byte, error) {
	// create a nonce byte array of the NaCl library's recommended size
	var nonce [models.NonceSize]byte
	// copy the first N bytes of the ciphertext to the nonce byte array,
	// where N is the NaCl library's recommended size for nonces
	copy(nonce[:], ciphertext[:models.NonceSize])
	// decrypt the remaining ciphertext using the nonce and key with the NaCl secretbox open function
	plaintext, ok := secretbox.Open(nil, ciphertext[models.NonceSize:], &nonce, key)
	if !ok {
		return nil, errors.New("decryption failed")
	}
	return plaintext, nil
}

// CreateNewBoxStore creates a new encrypted and split store with a given name, threshold, total,
// and description, and returns an array of memory safe shares.
func CreateNewBoxStore(name string, threshold int, total int, description string) []*memguard.Enclave {
	// Generate a new random key
	keyBytes, _ := generateRandomBytes(models.KeySize)

	// Copy the key to a fixed size array
	var key [models.KeySize]byte
	copy(key[:], keyBytes)

	// Wipe the keyBytes from memory and defer the wipe of the key
	memguard.WipeBytes(keyBytes)
	defer memguard.WipeBytes(key[:])

	// Create a new unsafe store entry with a unique Path and the given description
	store := models.StoreEntry{
		ID:             uuid.New().String(),
		Secrets:        []*models.Secret{},
		Description:    description,
		SecretPathList: []string{},
	}

	// Marshal the store entry to a string
	storeString, err := json.Marshal(store)
	if err != nil {
		panic(err)
	}

	// Encrypt the store string with the key
	encStore, err := Encrypt(storeString, &key)

	// Wipe the store string from memory
	memguard.WipeBytes(storeString)

	// Split the key into shares using Shamir's Secret Sharing
	shares, err := shamir.Split(key[:], total, threshold)
	if err != nil {
		panic(err)
	}

	// Create memory safe shares from the unsafe shares
	var safeShares []*memguard.Enclave
	for i, share := range shares {
		safeShares = append(safeShares, memguard.NewEnclave(share))
		// Wipe the unsafe shares from memory
		memguard.WipeBytes(share)
		memguard.WipeBytes(shares[i])
	}

	// Save the encrypted store to disk with a filename based on the store name
	path := configs.GlobalViper.GetString("server.basepath") + "/" + name + ".json"
	err = os.WriteFile(path, encStore, 0644)
	if err != nil {
		panic(err)
	}

	return safeShares
}

// OpenBoxStore opens an existing encrypted store with a given name and key,
// and returns a memory safe version of the store.
func OpenBoxStore(name string, key [models.KeySize]byte) (*models.SafeStoreEntry, error) {
	// Construct the path to the file containing the encrypted store
	path := configs.GlobalViper.GetString("server.basepath") + name + ".json"
	// Read the encrypted store from disk
	storeBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Decrypt the store using the given key
	decStoreString, err := decrypt(storeBytes, &key)
	if err != nil {
		return nil, err
	}

	// Unmarshal the decrypted store into an unsafe store entry object
	var unsafeEncStore models.StoreEntry
	err = json.Unmarshal(decStoreString, &unsafeEncStore)
	if err != nil || unsafeEncStore.ID == "" {
		return nil, err
	}

	// Wipe the decrypted store from memory
	memguard.WipeBytes(decStoreString)

	// Create a memory safe store entry from the unsafe store entry
	safeEncStore := &models.SafeStoreEntry{
		ID:             unsafeEncStore.ID,
		Secrets:        []*models.SafeSecret{},
		Description:    unsafeEncStore.Description,
		SecretPathList: unsafeEncStore.SecretPathList,
	}

	// Create memory safe secrets from the unsafe secrets
	for _, secret := range unsafeEncStore.Secrets {
		k, _ := base64.StdEncoding.DecodeString(secret.Data)
		safeEncStore.Secrets = append(
			safeEncStore.Secrets,
			&models.SafeSecret{
				Path:        secret.ID,
				Description: secret.Description,
				Data:        memguard.NewEnclave(k),
				Roles:       secret.Roles,
			})
	}

	// Wipe the unsafe store from memory
	unsafeEncStore = models.StoreEntry{
		ID:             "",
		Secrets:        []*models.Secret{},
		Description:    "",
		SecretPathList: []string{},
	}

	return safeEncStore, nil
}
