package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "testchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_testchain"
)

var (
	ParamsKey = collections.NewPrefix("p_testchain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
