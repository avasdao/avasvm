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

/**
 * Storage Gateway
 *
 * A precompiled contract address at:
 * 0x0000 (0)
 *
 * This precompile is a part of the Leet Suite of Subnet contracts.
 */
address constant STORAGE_GATEWAY_ADDRESS = 0x0359000000000000000000000000000000000000;

contract StorageGateway {
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
        bool,
        string memory
    ) {
        // result = storageGateway.getData(_cid);

        (bool success, bytes memory data) = address(storageGateway).staticcall(
            abi.encodeWithSignature("getData(string)", _cid)
        );

        /* Convert response data to string. */
        string memory response = string(data);

        return (success, response);
    }

    /**
     * Get Data By Key
     */
    function getDataByKey(
        string calldata _cid,
        string calldata _cid2
    ) external view returns (
        bool,
        string memory
    ) {
        // result = storageGateway.getData(_cid);

        (bool success, bytes memory data) = address(storageGateway).staticcall(
            abi.encodeWithSignature("getDataByKey(string,string)", _cid, _cid2)
        );

        /* Convert response data to string. */
        string memory response = string(data);

        return (success, response);
    }

    // setRecipient
    function setRecipient(
        string calldata _firstName,
        string calldata _lastName
    ) external returns (
        string memory result
    ) {
        result = storageGateway
            .setRecipient(_firstName, _lastName);

        emit LogString(result);
    }
}
