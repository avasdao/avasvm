// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface StorageGatewayInterface {
    /* Say hello. */
    function sayHello() external returns (bytes memory result);

    /* Set recipient. */
    function setRecipient(string calldata recipient) external returns (bytes memory result);
}
