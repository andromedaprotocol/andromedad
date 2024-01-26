package types

const (
	// ModuleName defines the module name
	ModuleName = "nibtdfvas"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_nibtdfvas"

	//  
	BondDenom = "ANDR"
)

var (
	ParamsKey = []byte("p_nibtdfvas")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
