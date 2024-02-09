// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import {StrategyManagerMock} from "eigenlayer-contracts/src/test/mocks/StrategyManagerMock.sol";
import {DelegationManagerMock} from "eigenlayer-contracts/src/test/mocks/DelegationManagerMock.sol";
import {SlasherMock} from "eigenlayer-contracts/src/test/mocks/SlasherMock.sol";
import {EigenPodManagerMock} "eigenlayer-contracts/src/test/mocks/EigenPodManagerMock.sol";
// import "eigenlayer-contracts/src/test/utils/EigenLayerUnitTestBase.sol";

import { Test } from "forge-std/Test.sol";

contract EigenCoreBase is Test {
    StrategyManagerMock strategyManagerMock;
    DelegationManagerMock public delegationManagerMock;
    SlasherMock public slasherMock;
    EigenPodManagerMock public eigenPodManagerMock;

    function setUp() public virtual override {
        EigenLayerUnitTestBase.setUp();
        strategyManagerMock = new StrategyManagerMock();
        delegationManagerMock = new DelegationManagerMock();
        slasherMock = new SlasherMock();
        eigenPodManagerMock = new EigenPodManagerMock();
    }

}
