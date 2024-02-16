package types

const (
	// ModuleName defines the module name
	ModuleName = "nibtdfvas"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_nibtdfvas"
)

var ParamsKey = []byte{0x00} // Prefix for params key

func KeyPrefix(p string) []byte {
	return []byte(p)
}
