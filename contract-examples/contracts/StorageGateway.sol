/*******************************************************************************
 * Copyright (c) 2022 Ava's DAO
 * All rights reserved.
 *
 * SPDX-License-Identifier: MIT
 *
 * https://avasdao.org
 * support@avasdao.org
 */

pragma solidity >=0.8.0;

import "./IStorageGateway.sol";

// *****************************************************************************
// Storage Gateway Address
//
// This is a unique, precompiled contract address stored on each node
// of any Validator supporting this service.
//
// This precompile is a part of the Leet Suite of Subnet contracts.
// Registered address is:
//   - 0x0539000000000000000000000000000000000001
//   - 0x01 (1)
//
address constant STORAGE_GATEWAY_ADDRESS
    = 0x0539000000000000000000000000000000000001;

contract StorageGateway {
    // FOR DEBUGGING PURPOSES ONLY
    event LogBool(bool myBool);
    event LogBytes(bytes myBytes);
    event LogString(string myString);
    event LogUint(uint256 myUint);

    /* Initialize the Storage Gateway handler. */
    IStorageGateway storageGateway = IStorageGateway(
        address(STORAGE_GATEWAY_ADDRESS));

    /**
     * Get Data
     */
    function getData(
        string calldata _cid
    ) external view returns (
        bool success,
        string memory response
    ) {
        (bool _success, bytes memory data) = address(storageGateway)
            .staticcall(
                abi.encodeWithSignature(
                    "getData(string)",
                    _cid
                )
            );

        /* Convert response data to string. */
        string memory _response = string(data);

        /* Return. */
        return (_success, _response);
    }

    /**
     * Get Data (With Path)
     */
    function getData(
        string calldata _cid,
        string calldata _path
    ) external view returns (
        bool success,
        string memory response
    ) {
        (bool _success, bytes memory data) = address(storageGateway)
            .staticcall(
                abi.encodeWithSignature(
                    "getData(string,string)",
                    _cid,
                    _path
                )
            );

        /* Convert response data to string. */
        string memory _response = string(data);

        /* Return. */
        return (_success, _response);
    }

    /**
     * Get Data By Key
     */
    function getDataByKey(
        string calldata _cid,
        string calldata _key
    ) external view returns (
        bool success,
        string memory response
    ) {
        (bool _success, bytes memory data) = address(storageGateway)
            .staticcall(
                abi.encodeWithSignature(
                    "getDataByKey(string,string)",
                    _cid,
                    _key
                )
            );

        /* Convert response data to string. */
        string memory _response = string(data);

        /* Return. */
        return (_success, _response);
    }

    /**
     * Save Data
     *
     * Provide data to be saved to immutable storage.
     *
     * Recieve back a Content Identifier (CID) for the data.
     */
    function saveData(
        string calldata _data
    ) external returns (
        bool success,
        string memory cid
    ) {
        (bool _success, bytes memory data) = address(storageGateway)
            .call(
                abi.encodeWithSignature(
                    "saveData(string)",
                    _data
                )
            );

        /* Convert response data to string. */
        string memory _cid = string(data);

        /* Emit (save) log entry. */
        emit LogString(_cid);

        /* Return. */
        return (_success, _cid);
    }
}
