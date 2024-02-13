// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.12;

import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import { PauserRegistry } from "eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";
import { IStrategy } from "eigenlayer-contracts/src/contracts/interfaces/IStrategy.sol";
import { StrategyManagerMock } from "eigenlayer-contracts/src/test/mocks/StrategyManagerMock.sol";
import { DelegationManagerMock } from "eigenlayer-contracts/src/test/mocks/DelegationManagerMock.sol";
import { EigenPodManagerMock } from "eigenlayer-contracts/src/test/mocks/EigenPodManagerMock.sol";
import { SlasherMock } from "eigenlayer-contracts/src/test/mocks/SlasherMock.sol";
import { EmptyContract } from "eigenlayer-contracts/src/test/mocks/EmptyContract.sol";

import { BN254 } from "eigenlayer-middleware/src/libraries/BN254.sol";
import { RegistryCoordinatorHarness } from "eigenlayer-middleware/test/harnesses/RegistryCoordinatorHarness.t.sol";
import { StakeRegistryHarness } from "eigenlayer-middleware/test/harnesses/StakeRegistryHarness.sol";
import { BLSApkRegistryHarness } from "eigenlayer-middleware/test/harnesses/BLSApkRegistryHarness.sol";
import { DelegationMock } from "eigenlayer-middleware/test/mocks/DelegationMock.sol";
import { RegistryCoordinator } from "eigenlayer-middleware/src/RegistryCoordinator.sol";
import { BLSApkRegistry } from "eigenlayer-middleware/src/BLSApkRegistry.sol";
import { StakeRegistry } from "eigenlayer-middleware/src/StakeRegistry.sol";
import { IRegistryCoordinator } from "eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import { IBLSApkRegistry } from "eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import { IStakeRegistry } from "eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import { IndexRegistry } from "eigenlayer-middleware/src/IndexRegistry.sol";
import { IIndexRegistry } from "eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";

import { OmniDelegationAVS } from "src/protocol/OmniDelegationAVS.sol";
import { Validators } from "src/libraries/Validators.sol";
import { TestPortal } from "test/common/TestPortal.sol";
import { FeeOracleV1 } from "src/protocol/FeeOracleV1.sol";
import { Test } from "forge-std/Test.sol";

contract OmniDelegationAVSBase is Test {
    /**
     * EignBase
     */
    PauserRegistry public pauserRegistry;
    ProxyAdmin public proxyAdmin;

    address public constant proxyAdminOwner = address(111);
    address public constant registryCoordinatorOwner = address(222);
    address public constant omniAVSOwner = address(333);
    address public constant pauser = address(555);
    address public constant unpauser = address(556);

    /**
     * EigenCoreBase
     */
    StrategyManagerMock public strategyManagerMock;
    DelegationManagerMock public delegationManagerMock;
    DelegationMock public delegationMock;
    SlasherMock public slasherMock;
    EigenPodManagerMock public eigenPodManagerMock;

    address defaultOperator = address(uint160(uint256(keccak256("defaultOperator"))));
    bytes32 defaultOperatorId;
    BN254.G1Point internal defaultPubKey = BN254.G1Point(
        18_260_007_818_883_133_054_078_754_218_619_977_578_772_505_796_600_400_998_181_738_095_793_040_006_897,
        3_432_351_341_799_135_763_167_709_827_653_955_074_218_841_517_684_851_694_584_291_831_827_675_065_899
    );

    /**
     * AVS Base
     */
    EmptyContract public emptyContract;

    RegistryCoordinatorHarness public registryCoordinatorImplementation;
    StakeRegistryHarness public stakeRegistryImplementation;
    IBLSApkRegistry public blsApkRegistryImplementation;
    IIndexRegistry public indexRegistryImplementation;
    OmniDelegationAVS public omniDelegationAVSImplementation;

    RegistryCoordinatorHarness public registryCoordinator;
    StakeRegistryHarness public stakeRegistry;
    BLSApkRegistryHarness public blsApkRegistry;
    IIndexRegistry public indexRegistry;
    OmniDelegationAVS public omniDelegationAVS;

    /// @notice StakeRegistry, Constant used as a divisor in calculating weights.
    uint256 public constant WEIGHTING_DIVISOR = 1e18;

    IRegistryCoordinator.OperatorSetParam[] operatorSetParams;

    /**
     * Omni base.
     */
    TestPortal public portal;
    FeeOracleV1 public feeOracle;

    uint256 churnApproverPrivateKey = uint256(keccak256("churnApproverPrivateKey"));
    address churnApprover = vm.addr(churnApproverPrivateKey);
    bytes32 defaultSalt = bytes32(uint256(keccak256("defaultSalt")));
    address ejector = address(uint160(uint256(keccak256("ejector"))));

    uint32 defaultMaxOperatorCount = 10;
    uint16 defaultKickBIPsOfOperatorStake = 15_000;
    uint16 defaultKickBIPsOfTotalStake = 150;

    function setUp() public virtual {
        _setupBase(); // proxy admin & pausers
        _setupEigenCore(); // mock core contracts
        _setupPortal();
        _deployProxies();
        _deployImplAndUpgrade();

        vm.startPrank(proxyAdminOwner);

        // TODO: do we need this
        // set the public key for an operator, using harnessed function to bypass checks
        blsApkRegistry.setBLSPublicKey(defaultOperator, defaultPubKey);

        // setup quorum
        // uint96 minimumStake = 1000; // TODO: what should this be?
        address mockStrategy = address(444);

        // TODO: two strategy quorum - mock native eth, mock erc20 (lst)

        IStakeRegistry.StrategyParams[] memory quorumStrategiesConsideredAndMultipliers =
            new IStakeRegistry.StrategyParams[](1);

        quorumStrategiesConsideredAndMultipliers[0] = IStakeRegistry.StrategyParams({
            strategy: IStrategy(mockStrategy),
            multiplier: uint96(WEIGHTING_DIVISOR) // TODO: what should this be?
         });

        uint96[] memory minimumStakeForQuorum = new uint96[](1);
        minimumStakeForQuorum[0] = uint96(1000); // TODO: what should this be?

        registryCoordinatorImplementation =
            new RegistryCoordinatorHarness(omniDelegationAVS, stakeRegistry, blsApkRegistry, indexRegistry);
        {
            delete operatorSetParams;
            operatorSetParams.push(
                IRegistryCoordinator.OperatorSetParam({
                    maxOperatorCount: defaultMaxOperatorCount,
                    kickBIPsOfOperatorStake: defaultKickBIPsOfOperatorStake,
                    kickBIPsOfTotalStake: defaultKickBIPsOfTotalStake
                })
            );

            proxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(registryCoordinator))),
                address(registryCoordinatorImplementation),
                abi.encodeWithSelector(
                    RegistryCoordinator.initialize.selector,
                    registryCoordinatorOwner,
                    churnApprover,
                    ejector,
                    pauserRegistry,
                    0, /*initialPausedStatus*/
                    operatorSetParams,
                    minimumStakeForQuorum,
                    quorumStrategiesConsideredAndMultipliers
                )
            );
        }

        // operatorStateRetriever = new OperatorStateRetriever();

        vm.stopPrank();
    }

    // initializize proxy admin & pauser, used in many places
    function _setupBase() internal {
        address[] memory pausers = new address[](1);
        pausers[0] = pauser;
        vm.startPrank(proxyAdminOwner);
        pauserRegistry = new PauserRegistry(pausers, unpauser);
        proxyAdmin = new ProxyAdmin();
        vm.stopPrank();
    }

    // deploys mock core contracts
    function _setupEigenCore() internal {
        strategyManagerMock = new StrategyManagerMock();
        delegationManagerMock = new DelegationManagerMock();
        slasherMock = new SlasherMock();
        eigenPodManagerMock = new EigenPodManagerMock();

        strategyManagerMock.setAddresses(delegationManagerMock, eigenPodManagerMock, slasherMock);
    }

    // deploy empty proxies
    function _deployProxies() internal {
        emptyContract = new EmptyContract();

        vm.startPrank(registryCoordinatorOwner);
        registryCoordinator = RegistryCoordinatorHarness(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), ""))
        );

        stakeRegistry = StakeRegistryHarness(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), ""))
        );

        indexRegistry =
            IndexRegistry(address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), "")));

        blsApkRegistry = BLSApkRegistryHarness(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), ""))
        );

        omniDelegationAVS =
            OmniDelegationAVS(address(new TransparentUpgradeableProxy(address(emptyContract), address(proxyAdmin), "")));

        vm.stopPrank();
    }

    // deploy implementations, upgrade proxies
    function _deployImplAndUpgrade() internal {
        vm.startPrank(proxyAdminOwner);

        stakeRegistryImplementation =
            new StakeRegistryHarness(IRegistryCoordinator(registryCoordinator), delegationMock);

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(stakeRegistry))), address(stakeRegistryImplementation)
        );

        blsApkRegistryImplementation = new BLSApkRegistryHarness(registryCoordinator);

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(blsApkRegistry))), address(blsApkRegistryImplementation)
        );

        indexRegistryImplementation = new IndexRegistry(registryCoordinator);

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))), address(indexRegistryImplementation)
        );

        omniDelegationAVSImplementation = new OmniDelegationAVS(delegationMock, registryCoordinator, stakeRegistry);

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(omniDelegationAVS))), address(omniDelegationAVSImplementation)
        );

        omniDelegationAVS.initialize(omniAVSOwner, address(portal));

        vm.stopPrank();
    }

    /**
     * OmniPortal stuff
     */
    string constant valMnemonic = "test test test test test test test test test test test junk";

    uint64 constant baseValPower = 100;
    uint64 constant genesisValSetId = 1;
    uint256 constant baseFee = 1 gwei;

    address val1;
    address val2;
    address val3;
    address val4;

    mapping(uint64 => Validators.Validator[]) validatorSet;
    mapping(address => uint256) valPrivKeys;

    uint256 val1PrivKey;
    uint256 val2PrivKey;
    uint256 val3PrivKey;
    uint256 val4PrivKey;

    address portalOwner = address(10_000);

    // deploy omni portals
    function _setupPortal() internal {
        uint64 thisChainId = 1;

        (val1, val1PrivKey) = deriveRememberKey(valMnemonic, 0);
        (val2, val2PrivKey) = deriveRememberKey(valMnemonic, 1);
        (val3, val3PrivKey) = deriveRememberKey(valMnemonic, 2);
        (val4, val4PrivKey) = deriveRememberKey(valMnemonic, 3);

        valPrivKeys[val1] = val1PrivKey;
        valPrivKeys[val2] = val2PrivKey;
        valPrivKeys[val3] = val3PrivKey;
        valPrivKeys[val4] = val4PrivKey;

        Validators.Validator[] storage genVals = validatorSet[genesisValSetId];

        genVals.push(Validators.Validator(val1, baseValPower));
        genVals.push(Validators.Validator(val2, baseValPower));
        genVals.push(Validators.Validator(val3, baseValPower));
        genVals.push(Validators.Validator(val4, baseValPower));

        vm.chainId(thisChainId); // portal constructor uses block.chainid
        feeOracle = new FeeOracleV1(portalOwner, baseFee);
        portal = new TestPortal(portalOwner, address(feeOracle), genesisValSetId, validatorSet[genesisValSetId]);
    }
}
