package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

const (
	// DefaultandromedaInstanceCost is initially set the same as in wasmd
	DefaultandromedaInstanceCost uint64 = 60_000
	// DefaultandromedaCompileCost set to a large number for testing
	DefaultandromedaCompileCost uint64 = 100
)

// andromedaGasRegisterConfig is defaults plus a custom compile amount
func andromedaGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultandromedaInstanceCost
	gasConfig.CompileCost = DefaultandromedaCompileCost

	return gasConfig
}

func NewandromedaWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(andromedaGasRegisterConfig())
}
