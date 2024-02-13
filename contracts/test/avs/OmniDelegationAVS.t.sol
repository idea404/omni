// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.12;

import { OmniDelegationAVSBase } from "./OmniDelegationAVSBase.sol";
import { OperatorStateRetriever } from "eigenlayer-middleware/src/OperatorStateRetriever.sol";
import { console } from "forge-std/console.sol";

contract OmniDelegationAVS_Test is OmniDelegationAVSBase {
    function test_getOperatorState() public {
        OperatorStateRetriever.Operator[][] memory allOperatorInfo = omniDelegationAVS.getOperatorState();
        assertEq(allOperatorInfo.length, 1);
        for (uint i = 0; i < allOperatorInfo[0].length; i++) {
            OperatorStateRetriever.Operator memory operator = allOperatorInfo[0][i];
            console.log("operator: ", operator.operator);
            console.log("stakes: ", operator.stake);
        }
    }
}
