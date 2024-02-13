// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.12;

import { BytesLib } from "eigenlayer-contracts/src/contracts/libraries/BytesLib.sol";
import { OperatorStateRetriever } from "eigenlayer-middleware/src/OperatorStateRetriever.sol";
import { ServiceManagerBase } from "eigenlayer-middleware/src/ServiceManagerBase.sol";
import { IDelegationManager } from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import { IRegistryCoordinator } from "eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import { IStakeRegistry } from "eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import { IOmniPortal } from "../interfaces/IOmniPortal.sol";

// avs contract on L1
contract OmniDelegationAVS is ServiceManagerBase, OperatorStateRetriever {
    using BytesLib for bytes;

    IOmniPortal public omni;

    bytes public constant QUORUM_NUMBERS = hex"00";

    constructor(
        IDelegationManager delegationManager_,
        IRegistryCoordinator registryCoordinator_,
        IStakeRegistry stakeRegistry_
    ) ServiceManagerBase(delegationManager_, registryCoordinator_, stakeRegistry_) { }

    function initialize(address owner_, address omniPortal_) external virtual initializer {
        _transferOwnership(owner_);
        omni = IOmniPortal(omniPortal_);
    }

    // relayer calls this
    function feeForSync() external view returns (uint256) {
        Operator[][] memory allOperatorInfo =
            getOperatorState(_registryCoordinator, QUORUM_NUMBERS, uint32(block.number));
        return omni.feeFor(OmniChains.OMNI, abi.encodeWithSelector(OmniDelegations.sync.selector, allOperatorInfo));
    }

    function getOperatorState() public view returns (Operator[][] memory) {
        return OperatorStateRetriever.getOperatorState(_registryCoordinator, QUORUM_NUMBERS, uint32(block.number));
    }

    function getOperatorState(uint256 blockNumber) public view returns (Operator[][] memory) {
        return OperatorStateRetriever.getOperatorState(_registryCoordinator, QUORUM_NUMBERS, uint32(blockNumber));
    }

    // calls this with msg.value == feeForSync()
    function syncWithOmni() external payable {
        Operator[][] memory allOperatorInfo =
            getOperatorState(_registryCoordinator, QUORUM_NUMBERS, uint32(block.number));

        omni.xcall{ value: msg.value }(
            OmniChains.OMNI,
            OmniPredeploys.OMNI_DELEGTION_MANAGER,
            abi.encodeWithSelector(OmniDelegations.sync.selector, allOperatorInfo)
        );
    }
}

library OmniChains {
    uint8 public constant ETH = 1;
    uint8 public constant OMNI = 165; // omni evm
    uint8 public constant OMNI_C = 166; // omni cchain
}

library OmniPredeploys {
    address public constant OMNI_DELEGTION_MANAGER = address(0x1234);
}

// predeployed on our OmniEVM - at known address
contract OmniDelegations {
    // matches OperatorStateRetriever

    event DelegationSync(Operator[][] operatorState);

    struct Operator {
        address operator;
        bytes32 operatorId; // will we need this?
        uint96 stake;
    }

    IOmniPortal public omni;
    address public avs;

    // what's the portal address
    // what's the address of the AVS contract on L1

    function initialize(address omni_, address avs_) external {
        require(address(omni) == address(0) && avs == address(0), "already initialized");
        require(omni_ != address(0), "invalid omni");
        require(avs_ != address(0), "invalid avs");

        omni = IOmniPortal(omni_);
        avs = avs_;
    }

    function sync(Operator[][] calldata operatorState) external {
        require(msg.sender == address(omni), "only omni");
        require(omni.isXCall(), "only xcall");

        // checks that is from l1 and is the address of the AVS contracts
        // unclear how this will work. we need away of setting OMNI_DELEGTION_AVS
        // if this is predeployed, we can do it in genesis. we also need to set portal contract
        //
        require(omni.xmsg().sourceChainId == OmniChains.ETH, "only omni");
        require(omni.xmsg().sender == avs, "only avs");

        emit DelegationSync(operatorState);
    }
}
