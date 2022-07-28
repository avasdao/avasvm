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

interface IStorageGateway {
    /* Get Data */
    function getData(string calldata) external view returns (string memory);

    /* Get Data By Key*/
    function getDataByKey(string calldata, string calldata) external view returns (string memory);

    // setRecipient
    function setRecipient(string calldata, string calldata) external returns (string memory);
}
