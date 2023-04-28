package encryptedStore

type EncryptedKVStore struct {
	Data  map[string][]byte `json:"data"`
	Nonce string            `json:"nonce"`
}

type EthereumAccount struct {
	Address string   `json:"address"`
	Key     string   `json:"key"`
	Roles   []string `json:"roles"`
}
