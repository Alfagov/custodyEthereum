package encryptedStore

import (
	"custodyEthereum/configs"
	"encoding/hex"
	"encoding/json"
	secretbox "github.com/GoKillers/libsodium-go/cryptosecretbox"
	"github.com/GoKillers/libsodium-go/randombytes"
	"github.com/hashicorp/vault/shamir"
	"os"
)

func EncryptKVStore(data map[string]string, key []byte) EncryptedKVStore {
	encStore := EncryptedKVStore{Data: make(map[string][]byte)}
	nonce := randombytes.RandomBytes(24)
	encStore.Nonce = hex.EncodeToString(nonce)

	for k, v := range data {
		b := []byte(v)
		ciphertext, err := secretbox.CryptoSecretBoxEasy(b, nonce, key)
		if err != 1 {
			panic(err)
		}

		encStore.Data[k] = ciphertext
	}
	return encStore
}

func DecryptKVStore(encStore EncryptedKVStore, key []byte) map[string]string {
	decryptedData := make(map[string]string)
	nonce, _ := hex.DecodeString(encStore.Nonce)

	for k, v := range encStore.Data {
		plaintext, err := secretbox.CryptoSecretBoxOpenEasy(v, nonce, key)
		if err != 1 {
			panic(err)
		}

		decryptedData[k] = string(plaintext)
	}
	return decryptedData
}

func SaveEncryptedKVStore(encStore EncryptedKVStore, name string) {
	storeBytes, err := json.Marshal(encStore)
	if err != nil {
		panic(err)
	}

	path := configs.GlobalViper.GetString("server.basepath") + name + ".json"

	err = os.WriteFile(path, storeBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func LoadEncryptedKVStore(storeName string) EncryptedKVStore {

	path := configs.GlobalViper.GetString("server.basepath") + storeName + ".json"
	storeBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var encStore EncryptedKVStore
	err = json.Unmarshal(storeBytes, &encStore)
	if err != nil {
		panic(err)
	}

	return encStore
}

func CreateNewStore(name string, threshold int, total int) [][]byte {
	key := randombytes.RandomBytes(32)
	data := map[string]string{}

	encStore := EncryptKVStore(data, key)
	SaveEncryptedKVStore(encStore, name)

	shares, err := shamir.Split(key, threshold, total)
	if err != nil {
		panic(err)
	}

	return shares
}
