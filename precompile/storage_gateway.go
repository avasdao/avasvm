// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	_                            StatefulPrecompileConfig    = &StorageGatewayConfig{}
	ContractStorageGatewayPrecompile StatefulPrecompiledContract = createStorageGatewayPrecompile()

	storageGatewaySignature             = CalculateFunctionSelector("sayHello()")
	setStorageGatewayRecipientSignature = CalculateFunctionSelector("setRecipient(string)")

	nameKey      = common.BytesToHash([]byte("recipient"))
	initialValue = common.BytesToHash([]byte("world!"))
)

type StorageGatewayConfig struct {
	BlockTimestamp *big.Int `json:"blockTimestamp"`
}

// Address returns the address of the precompile
func (h *StorageGatewayConfig) Address() common.Address { return StorageGatewayAddress }

// Return the timestamp at which the precompile is enabled or nil, if it is never enabled
func (h *StorageGatewayConfig) Timestamp() *big.Int { return h.BlockTimestamp }

func (h *StorageGatewayConfig) Configure(_ ChainConfig, stateDB StateDB, _ BlockContext) {
	// h.AllowListConfig.Configure(state, StorageGatewayAddress)
	stateDB.SetState(StorageGatewayAddress, nameKey, initialValue)
}
// func (h *StorageGatewayConfig) Configure(stateDB StateDB) {
// 	// This will be called in the first block where StorageGateway stateful precompile is enabled.
// 	// 1) If BlockTimestamp is nil, this will not be called
// 	// 2) If BlockTimestamp is 0, this will be called while setting up the genesis block
// 	// 3) If BlockTimestamp is 1000, this will be called while processing the first block whose timestamp is >= 1000
// 	//
// 	// Set the initial value under [nameKey] to "world!"
// 	stateDB.SetState(StorageGatewayAddress, nameKey, initialValue)
// }

// Return the precompile contract
func (h *StorageGatewayConfig) Contract() StatefulPrecompiledContract {
	return ContractStorageGatewayPrecompile
}

// Arguments are passed in to functions according to the ABI specification: https://docs.soliditylang.org/en/latest/abi-spec.html.
// Therefore, we maintain compatibility with Solidity by following the same specification while encoding/decoding arguments.
func PackStorageGatewayInput(name string) ([]byte, error) {
	byteStr := []byte(name)
	if len(byteStr) > common.HashLength {
		return nil, fmt.Errorf("cannot pack hello world input with string: %s", name)
	}

	input := make([]byte, common.HashLength+len(byteStr))
	strLength := new(big.Int).SetUint64(uint64(len(byteStr)))
	strLengthBytes := strLength.Bytes()
	copy(input[:common.HashLength], strLengthBytes)
	copy(input[common.HashLength:], byteStr)

	return input, nil
}

// UnpackStorageGatewayInput unpacks the recipient string from the hello world input
func UnpackStorageGatewayInput(input []byte) (string, error) {
	if len(input) < common.HashLength {
		return "", fmt.Errorf("cannot unpack hello world input with length: %d", len(input))
	}

	strLengthBig := new(big.Int).SetBytes(input[:common.HashLength])
	if !strLengthBig.IsUint64() {
		return "", fmt.Errorf("cannot unpack hello world input with stated length that is non-uint64")
	}

	strLength := strLengthBig.Uint64()
	if strLength > common.HashLength {
		return "", fmt.Errorf("cannot unpack hello world string with length: %d", strLength)
	}

	if len(input) != common.HashLength+int(strLength) {
		return "", fmt.Errorf("input had unexpected length %d with string length defined as %d", len(input), strLength)
	}

	str := string(input[common.HashLength:])
	return str, nil
}

func GetReceipient(state StateDB) string {
	value := state.GetState(StorageGatewayAddress, nameKey)
	b := value.Bytes()
	trimmedbytes := common.TrimLeftZeroes(b)
	return string(trimmedbytes)
}

// SetRecipient sets the recipient for the hello world precompile
func SetRecipient(state StateDB, recipient string) {
	state.SetState(StorageGatewayAddress, nameKey, common.BytesToHash([]byte(recipient)))
}

// sayHello is the execution function of "sayHello()"
func sayHello(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if len(input) != 0 {
		return nil, 0, fmt.Errorf("fuck")
	}
	remainingGas, err = deductGas(suppliedGas, ReadStorageCost)
	if err != nil {
		return nil, 0, err
	}

	recipient := GetReceipient(accessibleState.GetStateDB())
	return []byte(fmt.Sprintf("Hello %s!", recipient)), suppliedGas - WriteStorageCost, nil
}

// setRecipient is the execution function of "setRecipient(name string)" and sets the recipient in the string returned by hello world
func setRecipient(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	recipient, err := UnpackStorageGatewayInput(input)
	if err != nil {
		return nil, 0, err
	}
	remainingGas, err = deductGas(suppliedGas, WriteStorageCost)
	if err != nil {
		return nil, 0, err
	}

	SetRecipient(accessibleState.GetStateDB(), recipient)
	return []byte{}, remainingGas, nil
}

// createStorageGatewayPrecompile returns the StatefulPrecompile contract that implements the StorageGateway interface from solidity
func createStorageGatewayPrecompile() StatefulPrecompiledContract {
	return newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{
		newStatefulPrecompileFunction(storageGatewaySignature, sayHello),
		newStatefulPrecompileFunction(setStorageGatewayRecipientSignature, setRecipient),
	})
}
