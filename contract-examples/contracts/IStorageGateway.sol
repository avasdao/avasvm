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

    /* Get Data With Path */
    function getData(string calldata, string calldata) external view returns (string memory);

    /* Get Data By Key */
    function getDataByKey(string calldata, string calldata) external view returns (string memory);

    /* Save Data */
    function saveData(string calldata) external returns (bool, string memory);
}
