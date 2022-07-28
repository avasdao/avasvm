/*******************************************************************************
* Copyright (c) 2022 Ava's DAO
* All rights reserved.
*
* SPDX-License-Identifier: MIT
*
* https://avasdao.org
* support@avasdao.org
*/

package precompile

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	_                           StatefulPrecompileConfig = &StorageGatewayConfig{}
	StorageGatewayPrecompile StatefulPrecompiledContract = createStorageGatewayPrecompile(StorageGatewayAddress)

	sayHelloSignature     = CalculateFunctionSelector("sayHello()")

	getData1Signature      = CalculateFunctionSelector("getData1(string)")
	getData2Signature      = CalculateFunctionSelector("getData2(string,string)")
	getData3Signature      = CalculateFunctionSelector("getData3(string)")
	getData4Signature      = CalculateFunctionSelector("getData4(string,string)")

	setRecipientSignature = CalculateFunctionSelector("setRecipient(string,string)")

	nameKey      = common.BytesToHash([]byte("recipient"))
	initialValue = common.BytesToHash([]byte("world"))
)

type StorageGatewayConfig struct {
	BlockTimestamp *big.Int `json:"blockTimestamp"`
}

/* Address returns the address of the precompile. */
func (h *StorageGatewayConfig) Address() common.Address {
	return StorageGatewayAddress
}

/* Return the timestamp at which the precompile is enabled or nil, if it is never enabled. */
func (h *StorageGatewayConfig) Timestamp() *big.Int {
	return h.BlockTimestamp
}

func (h *StorageGatewayConfig) Configure(_ ChainConfig, stateDB StateDB, _ BlockContext) {
	stateDB.SetState(StorageGatewayAddress, nameKey, initialValue)
}

/* Return the precompile contract. */
func (h *StorageGatewayConfig) Contract() StatefulPrecompiledContract {
	return StorageGatewayPrecompile
}

/**
 * Arguments are passed in to functions according to the ABI specification: https://docs.soliditylang.org/en/latest/abi-spec.html.
 * Therefore, we maintain compatibility with Solidity by following the same specification while encoding/decoding arguments.
 */
func PackStorageGatewayInput(name string) ([]byte, error) {
	byteStr := []byte(name)

	if len(byteStr) > common.HashLength {
		return nil, fmt.Errorf("cannot pack Storage Gateway input with string: %s", name)
	}

	input := make([]byte, common.HashLength+len(byteStr))

	strLength := new(big.Int).SetUint64(uint64(len(byteStr)))

	strLengthBytes := strLength.Bytes()

	copy(input[:common.HashLength], strLengthBytes)

	copy(input[common.HashLength:], byteStr)

	return input, nil
}

/**
 * UnpackStorageGatewayInput unpacks the recipient string from the Storage Gateway input.
 */
func UnpackStorageGatewayInput(input []byte) (string, error) {
	log.Info("Entering UnpackStorageGatewayInput ->", string(input), nil)

	if len(input) < common.HashLength {
		return "", fmt.Errorf("cannot unpack Storage Gateway input with length: %d", len(input))
	}

	strLengthBig := new(big.Int).SetBytes(input[:common.HashLength])

	if !strLengthBig.IsUint64() {
		return "", fmt.Errorf("cannot unpack Storage Gateway input with stated length that is non-uint64")
	}

	strLength := strLengthBig.Uint64()

	if strLength > common.HashLength {
		return "", fmt.Errorf("cannot unpack Storage Gateway string with length: %d", strLength)
	}

	if len(input) != common.HashLength+int(strLength) {
		return "", fmt.Errorf("input had unexpected length %d with string length defined as %d", len(input), strLength)
	}

	str := string(input[common.HashLength:])

	return str, nil
}

func GetRecipient(state StateDB) string {
	value := state.GetState(StorageGatewayAddress, nameKey)

	b := value.Bytes()

	trimmedbytes := common.TrimLeftZeroes(b)

	return string(trimmedbytes)
}

/* SetRecipient sets the recipient for the Storage Gateway precompile. */
func SetRecipient(state StateDB, recipient string) {
	state.SetState(StorageGatewayAddress, nameKey, common.BytesToHash([]byte(recipient)))
}

/**
 * sayHello is the execution function of "sayHello()"
 */
func sayHello(
	accessibleState PrecompileAccessibleState,
	caller common.Address,
	addr common.Address,
	input []byte,
	suppliedGas uint64,
	readOnly bool,
) (
	ret []byte,
	remainingGas uint64,
	err error,
) {
	log.Info("\nEntering sayHello()")

	if len(input) != 0 {
		return nil, 0, fmt.Errorf("Oops! You cannot provide INPUT here.")
	}

	remainingGas, err = deductGas(suppliedGas, ReadStorageCost)

	if err != nil {
		return nil, 0, err
	}

	recipient := GetRecipient(accessibleState.GetStateDB())
	log.Info("\nrecipient ->", recipient, nil)

	testVal := "hi-there"

	response := []byte(string(testVal))

	// return []byte(fmt.Sprintf("Hello %s!", recipient)), remainingGas, nil
	return response, remainingGas, nil
}

func getData1(precompileAddr common.Address) RunStatefulPrecompileFunc {
	return func(
		evm PrecompileAccessibleState,
		callerAddr common.Address,
		addr common.Address,
		input []byte,
		suppliedGas uint64,
		readOnly bool,
	) (
		ret []byte,
		remainingGas uint64,
		err error,
	) {
		log.Info("input-1")
		log.Info(string(input))

		if remainingGas, err = deductGas(suppliedGas, ReadStorageCost); err != nil {
			return nil, 0, err
		}

		const storageDataLength = 128

		input = common.RightPadBytes(input, storageDataLength)
		log.Info("input-2")
		log.Info(string(input))

		// testVal := "48d3c6d9-7324-4ae5-b7fd-32e19fa3996e"
		// testVal := "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

		testVal := `{
	      "chainId": 99999,
	      "homesteadBlock": 0,
	      "eip150Block": 0,
	      "eip150Hash": "0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0",
	      "eip155Block": 0,
	      "eip158Block": 0,
	      "byzantiumBlock": 0,
	      "constantinopleBlock": 0,
	      "petersburgBlock": 0,
	      "istanbulBlock": 0,
	      "muirGlacierBlock": 0,
	      "subnetEVMTimestamp": 0,
	      "feeConfig": {
	        "gasLimit": 20000000,
	        "minBaseFee": 1000000000,
	        "targetGas": 100000000,
	        "baseFeeChangeDenominator": 48,
	        "minBlockGasCost": 0,
	        "maxBlockGasCost": 10000000,
	        "targetBlockRate": 2,
	        "blockGasCostStep": 500000
	      },
	      "storageGatewayConfig": {
	          "blockTimestamp": 0
	      }
	    }`

		response := []byte(string(testVal))

		return response, remainingGas, nil
		// return []byte(fmt.Sprintf(testVal)), remainingGas, nil
	}
}

func getData2(precompileAddr common.Address) RunStatefulPrecompileFunc {
	return func(
		evm PrecompileAccessibleState,
		callerAddr common.Address,
		addr common.Address,
		input []byte,
		suppliedGas uint64,
		readOnly bool,
	) (
		ret []byte,
		remainingGas uint64,
		err error,
	) {
		log.Info("input-1")
		log.Info(string(input))

		if remainingGas, err = deductGas(suppliedGas, ReadStorageCost); err != nil {
			return nil, 0, err
		}

		const storageDataLength = 128

		input = common.RightPadBytes(input, storageDataLength)
		log.Info("input-2")
		log.Info(string(input))

		// testVal := "48d3c6d9-7324-4ae5-b7fd-32e19fa3996e"
		// testVal := "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

		testVal := `{
	      "chainId": 99999,
	      "homesteadBlock": 0,
	      "eip150Block": 0,
	      "eip150Hash": "0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0",
	      "eip155Block": 0,
	      "eip158Block": 0,
	      "byzantiumBlock": 0,
	      "constantinopleBlock": 0,
	      "petersburgBlock": 0,
	      "istanbulBlock": 0,
	      "muirGlacierBlock": 0,
	      "subnetEVMTimestamp": 0,
	      "feeConfig": {
	        "gasLimit": 20000000,
	        "minBaseFee": 1000000000,
	        "targetGas": 100000000,
	        "baseFeeChangeDenominator": 48,
	        "minBlockGasCost": 0,
	        "maxBlockGasCost": 10000000,
	        "targetBlockRate": 2,
	        "blockGasCostStep": 500000
	      },
	      "storageGatewayConfig": {
	          "blockTimestamp": 0
	      }
	    }`

		response := []byte(string(testVal))

		return response, remainingGas, nil
		// return []byte(fmt.Sprintf(testVal)), remainingGas, nil
	}
}

func getData3(
	evm PrecompileAccessibleState,
	callerAddr common.Address,
	addr common.Address,
	input []byte,
	suppliedGas uint64,
	readOnly bool,
) (
	ret []byte,
	remainingGas uint64,
	err error,
) {
	log.Info("input-1")
	log.Info(string(input))

	if remainingGas, err = deductGas(suppliedGas, ReadStorageCost); err != nil {
		return nil, 0, err
	}

	const storageDataLength = 128

	input = common.RightPadBytes(input, storageDataLength)
	log.Info("input-2")
	log.Info(string(input))

	// testVal := "48d3c6d9-7324-4ae5-b7fd-32e19fa3996e"
	// testVal := "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	testVal := `{
      "chainId": 99999,
      "homesteadBlock": 0,
      "eip150Block": 0,
      "eip150Hash": "0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0",
      "eip155Block": 0,
      "eip158Block": 0,
      "byzantiumBlock": 0,
      "constantinopleBlock": 0,
      "petersburgBlock": 0,
      "istanbulBlock": 0,
      "muirGlacierBlock": 0,
      "subnetEVMTimestamp": 0,
      "feeConfig": {
        "gasLimit": 20000000,
        "minBaseFee": 1000000000,
        "targetGas": 100000000,
        "baseFeeChangeDenominator": 48,
        "minBlockGasCost": 0,
        "maxBlockGasCost": 10000000,
        "targetBlockRate": 2,
        "blockGasCostStep": 500000
      },
      "storageGatewayConfig": {
          "blockTimestamp": 0
      }
    }`

	response := []byte(string(testVal))

	return response, remainingGas, nil
	// return []byte(fmt.Sprintf(testVal)), remainingGas, nil
}

func getData4(
	evm PrecompileAccessibleState,
	callerAddr common.Address,
	addr common.Address,
	input []byte,
	suppliedGas uint64,
	readOnly bool,
) (
	ret []byte,
	remainingGas uint64,
	err error,
) {
	log.Info("input-1")
	log.Info(string(input))

	if remainingGas, err = deductGas(suppliedGas, ReadStorageCost); err != nil {
		return nil, 0, err
	}

	const storageDataLength = 128

	input = common.RightPadBytes(input, storageDataLength)
	log.Info("input-2")
	log.Info(string(input))

	// testVal := "48d3c6d9-7324-4ae5-b7fd-32e19fa3996e"
	// testVal := "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	testVal := `{
      "chainId": 99999,
      "homesteadBlock": 0,
      "eip150Block": 0,
      "eip150Hash": "0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0",
      "eip155Block": 0,
      "eip158Block": 0,
      "byzantiumBlock": 0,
      "constantinopleBlock": 0,
      "petersburgBlock": 0,
      "istanbulBlock": 0,
      "muirGlacierBlock": 0,
      "subnetEVMTimestamp": 0,
      "feeConfig": {
        "gasLimit": 20000000,
        "minBaseFee": 1000000000,
        "targetGas": 100000000,
        "baseFeeChangeDenominator": 48,
        "minBlockGasCost": 0,
        "maxBlockGasCost": 10000000,
        "targetBlockRate": 2,
        "blockGasCostStep": 500000
      },
      "storageGatewayConfig": {
          "blockTimestamp": 0
      }
    }`

	response := []byte(string(testVal))

	return response, remainingGas, nil
	// return []byte(fmt.Sprintf(testVal)), remainingGas, nil
}

/**
 * `setRecipient` is the execution function of "setRecipient(name string)"
 * and sets the recipient in the string returned by Storage Gateway.
*/
func setRecipient(
	accessibleState PrecompileAccessibleState,
	caller common.Address,
	addr common.Address,
	input []byte,
	suppliedGas uint64,
	readOnly bool,
) (
	ret []byte,
	remainingGas uint64,
	err error,
) {
	log.Info("Entering setRecipient()")
	recipient, err := UnpackStorageGatewayInput(input)

	if err != nil {
		return nil, 0, err
	}

	remainingGas, err = deductGas(suppliedGas, WriteStorageCost)

	if err != nil {
		return nil, 0, err
	}

	log.Info("SetRecipient-1", recipient, nil)
	log.Info("SetRecipient-2", fmt.Sprintf("recipient -> %s", recipient), nil)
	SetRecipient(accessibleState.GetStateDB(), recipient)

	return []byte{}, remainingGas, nil
}

/**
 * Create Storage Gateway Precompile
 *
 * Returns the StatefulPrecompile contract that implements
 * the StorageGateway interface from solidity.
 */
func createStorageGatewayPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	/* Construct the contract without a fallback function. */
	storageGatewayFuncs := []*statefulPrecompileFunction {
		newStatefulPrecompileFunction(
			sayHelloSignature, sayHello),

		newStatefulPrecompileFunction(
			getData1Signature, getData1(precompileAddr)),
		newStatefulPrecompileFunction(
			getData2Signature, getData2(precompileAddr)),
		newStatefulPrecompileFunction(
			getData3Signature, getData3),
		newStatefulPrecompileFunction(
			getData4Signature, getData4),

		newStatefulPrecompileFunction(
			setRecipientSignature, setRecipient),
	}

	/* Construct the contract without a fallback function. */
	contract := newStatefulPrecompileWithFunctionSelectors(
		nil, storageGatewayFuncs)

	/* Return contract. */
	return contract
}
