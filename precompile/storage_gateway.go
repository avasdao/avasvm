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
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	_                           StatefulPrecompileConfig = &StorageGatewayConfig{}
	StorageGatewayPrecompile StatefulPrecompiledContract = createStorageGatewayPrecompile()

	getDataSignature         = CalculateFunctionSelector("getData(string)")
	getDataWithPathSignature = CalculateFunctionSelector("getData(string,string)")
	getDataByKeySignature    = CalculateFunctionSelector("getDataByKey(string,string)")
	setRecipientSignature    = CalculateFunctionSelector("setRecipient(string,string)")

	nameKey      = common.BytesToHash([]byte("recipient"))
	initialValue = common.BytesToHash([]byte("world"))

	/* Set web gateway. */
	WebGateway = ".ipfs.dweb.link"
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
 * Pack Storage Gateway Input
 *
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
 * Unpack Storage Gateway Input
 *
 * Unpacks the recipient string from the Storage Gateway input.
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

func getUrl(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("GET error: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Status error: %v", resp.StatusCode)
    }

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("Read body: %v", err)
    }

    return string(data), nil
}

/**
 * Get Data
 */
func getData(
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
	log.Info("\n[getData] input->", string(input), nil)

	encodedString := hex.EncodeToString(input)
	log.Info("\n[getData] encodedString->\n", string(encodedString), nil)

	/* Calculate remaining gas. */
	if remainingGas, err = deductGas(suppliedGas, ReadStorageCost); err != nil {
		return nil, 0, err
	}

	/* Create a regex to filter only want letters and numbers. */
    reg, err := regexp.Compile("[^a-zA-Z0-9]+")

	/* Handle errors. */
    if err != nil {
		return nil, remainingGas, err
    }

    cid := reg.ReplaceAllString(string(input), "")
	log.Info("\n[getData] cid->", cid, nil)

	/* Set (data) path. */
	path := "/readme"

	/* Set data target. */
	target := "https://" + cid + WebGateway + path
	log.Info("\n[getData] target->", string(target), nil)

	/* Request URL data. */
	data, err := getUrl(target)

	/* Handle error. */
    if err != nil {
		return nil, remainingGas, err
    }

	/* Handle response. */
	response := []byte(string(data))

	/* Return response. */
	return response, remainingGas, nil
}

/**
 * Get Data With Path
 */
func getDataWithPath(
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
	log.Info("\n[getDataWithPath] input->", string(input), nil)

	/* Calculate remaining gas. */
	if remainingGas, err = deductGas(suppliedGas, ReadStorageCost); err != nil {
		return nil, 0, err
	}

	/* Handle error. */
    if err != nil {
		return nil, remainingGas, err
    }

	encodedString := hex.EncodeToString(input)
	log.Info("\n[getDataWithPath] encodedString->\n", string(encodedString), nil)

	param1Pos := common.TrimLeftZeroes(input[:32])
	param1PosHex := hex.EncodeToString(param1Pos)
	param1PosDec, _ := strconv.ParseInt(param1PosHex, 16, 64)
	log.Info("\n[getDataWithPath] param1PosHex->", string(param1PosHex), nil)

	param2Pos := common.TrimLeftZeroes(input[32:64])
	param2PosHex := hex.EncodeToString(param2Pos)
	param2PosDec, _ := strconv.ParseInt(param2PosHex, 16, 64)
	log.Info("\n[getDataWithPath] param2PosHex->", string(param2PosHex), nil)

	param1Len := common.TrimLeftZeroes(input[param1PosDec:(param1PosDec + 32)])
	param1LenHex := hex.EncodeToString(param1Len)
	param1LenDec, _ := strconv.ParseInt(param1LenHex, 16, 64)
	log.Info("\n[getDataWithPath] param1LenHex->", string(param1LenHex), nil)

	param2Len := common.TrimLeftZeroes(input[param2PosDec:(param2PosDec + 32)])
	param2LenHex := hex.EncodeToString(param2Len)
	param2LenDec, _ := strconv.ParseInt(param2LenHex, 26, 64)
	log.Info("\n[getDataWithPath] param2LenHex->", string(param2LenHex), nil)

	param1 := common.TrimRightZeroes(input[(param1PosDec + 32):(param1PosDec + 32 + param1LenDec)])
	log.Info("\n[getDataWithPath] param1->", string(param1), nil)

	param2 := common.TrimRightZeroes(input[(param2PosDec + 32):(param2PosDec + 32 + param2LenDec)])
	log.Info("\n[getDataWithPath] param2->", string(param2), nil)

	cid := string(param1)
	log.Info("\n[getDataWithPath] cid->", cid, nil)

	/* Set (data) path. */
	// NOTE: Add forward slash prefix.
	path := "/" + string(param2)
	log.Info("\n[getDataWithPath] path->", path, nil)

	/* Set data target. */
	target := "https://" + cid + WebGateway + path
	log.Info("\n[getDataWithPath] target->", string(target), nil)

	/* Request URL data. */
	data, err := getUrl(target)

	/* Handle error. */
    if err != nil {
		return nil, remainingGas, err
    }

	/* Handle response. */
	response := []byte(string(data))

	/* Return response. */
	return response, remainingGas, nil
}

/**
 * Get Data By Key
 */
func getDataByKey(
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
	log.Info("\n[getDataByKey] input->", string(input), nil)

	/* Calculate remaining gas. */
	if remainingGas, err = deductGas(suppliedGas, ReadStorageCost); err != nil {
		return nil, 0, err
	}

	/* Set CID. */
	cid := "bafybeie5nqv6kd3qnfjupgvz34woh3oksc3iau6abmyajn7qvtf6d2ho34"

	/* Set (data) path. */
	path := "/readme"

	/* Set data target. */
	target := "https://" + cid + WebGateway + path

	/* Request URL data. */
	data, err := getUrl(target)

	/* Handle response. */
	response := []byte(string(data))

	/* Return response. */
	return response, remainingGas, err
}

/**
 * Set Recipient
 *
 * Is the execution function of "setRecipient(name string)"
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
	log.Info("\n[setRecipient] input->", input, nil)

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
func createStorageGatewayPrecompile() StatefulPrecompiledContract {
	/* Construct the contract without a fallback function. */
	storageGatewayFuncs := []*statefulPrecompileFunction {
		/* Get data. */
		newStatefulPrecompileFunction(
			getDataSignature,
			getData,
		),

		/* Get data w/ path. */
		newStatefulPrecompileFunction(
			getDataWithPathSignature,
			getDataWithPath,
		),

		/* Get data by key. */
		newStatefulPrecompileFunction(
			getDataByKeySignature,
			getDataByKey,
		),

		// TODO
		newStatefulPrecompileFunction(
			setRecipientSignature,
			setRecipient,
		),
	}

	/* Construct the contract without a fallback function. */
	contract := newStatefulPrecompileWithFunctionSelectors(
		nil,
		storageGatewayFuncs,
	)

	/* Return contract. */
	return contract
}
