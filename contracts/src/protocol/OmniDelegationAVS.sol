// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

import { BytesLib } from "eigenlayer-contracts/src/contracts/libraries/BytesLib.sol";
import { OperatorStateRetriever } from "eigenlayer-middleware/src/OperatorStateRetriever.sol";
import { ServiceManagerBase } from "eigenlayer-middleware/src/ServiceManagerBase.sol";
import { IDelegationManager } from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import { IRegistryCoordinator } from "eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import { IStakeRegistry } from "eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import { IOmniPortal } from "../interfaces/IOmniPortal.sol";

library OmniChains {
    uint8 public constant ETH = 1;
    uint8 public constant OMNI = 165; // omni evm
    uint8 public constant OMNI_C = 166; // omni cchain
}

library OmniPredeploys {
    address public constant OMNI_DELEGTION_MANAGER = address(0x1234);
}

// predeployed - can it be predeployed with portal contract in geneis? probaly not
// how do we initialize a predeploy
// do we rely on create3? - we can, can we fail on initialization if create3 fails
// can we predeploy portal using create3 addr, would that be advised?
// can we have cchain attests to portals addrs? probably not
// okay so we deploy the RegistryCoordinator - good. we create qurom there
// we do same for all the registries
contract OmniDelegationManager is OperatorStateRetriever {
    constructor(IOmniPortal omni_) {
        omni = omni_;
    }

    function sync(Operator[][] calldata operatorState) external {
        require(msg.sender == address(omni), "only omni");

        // checks that is from l1 and is the address of the AVS contracts
        require(omni.isAVSXCall(), "only avs xcall");

        emit Synced(operatorState);
    }
}


contract OmniDelegationAVS is ServiceManagerBase, OperatorStateRetriever {
    using BytesLib for bytes;

    IOmniPortal public omni;

    constructor(
        IDelegationManager delegationManager_,
        IRegistryCoordinator registryCoordinator_,
        IStakeRegistry stakeRegistry_,
        IOmniPortal omni_
    ) ServiceManagerBase(delegationManager_, registryCoordinator_, stakeRegistry_) {
        omni = omni_;
    }

    function syncWithOmni() external {
        // how do we get quorumNumbers?
        // quorum number likely given from offchain
        // so their can only be one quorum number?
        // that quorum number either nee

        // this is is in avs incredible squaring, so we could do this - I don't think we want more quorums
        // we only use a single quorum (quorum 0) for incredible squaring
        // var QUORUM_NUMBERS = []byte{0}


        // we can also check if quorum exists?
             // * @notice Returns true iff all of the bits in `quorumBitmap` belong to initialized quorums
         // */
         // function _quorumsAllExist(uint192 quorumBitmap) internal view returns (bool) {
            // uint192 initializedQuorumBitmap = uint192((1 << quorumCount) - 1);
            // return quorumBitmap.isSubsetOf(initializedQuorumBitmap);
        // }


        bytes memory quorumNumbers = bytes("");

        Operator[][] memory allOperatorInfo = getOperatorState(registryCoordinator, quorumNumbers, block.number);

        omni.xcall{ value: msg.value }(
            OmniChains.OMNI,
            OmniPredeploys.OMNI_DELEGTION_MANAGER,
            abi.encodeWithSelector(OmniDelegationManager.sync.selector, allOperatorInfo)
        );
    }
}
