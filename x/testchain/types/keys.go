package types

const (
	// ModuleName defines the module name
	ModuleName = "testchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_testchain"
)

var (
	ParamsKey = []byte("p_testchain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
