package bindings

import (
	_ "embed"
)

const (
	StakingDeployedBytecode = "0x6080604052600436106100c25760003560e01c806384768b7a1161007f578063c6a2aac811610059578063c6a2aac8146101f8578063cf8e629a1461020d578063d146fd1b14610222578063f2fde38b1461023c57600080fd5b806384768b7a1461017d5780638da5cb5b146101bd578063a5a470ad146101e557600080fd5b8063117407e3146100c757806311bcd830146100e95780633f0b1edf1461011957806359bcddde146101395780635c19a95c14610155578063715018a614610168575b600080fd5b3480156100d357600080fd5b506100e76100e23660046106f2565b61025c565b005b3480156100f557600080fd5b5061010668056bc75e2d6310000081565b6040519081526020015b60405180910390f35b34801561012557600080fd5b506100e76101343660046106f2565b6102d1565b34801561014557600080fd5b50610106670de0b6b3a764000081565b6100e7610163366004610767565b610341565b34801561017457600080fd5b506100e7610438565b34801561018957600080fd5b506101ad610198366004610767565b60666020526000908152604090205460ff1681565b6040519015158152602001610110565b3480156101c957600080fd5b506033546040516001600160a01b039091168152602001610110565b6100e76101f3366004610797565b61044c565b34801561020457600080fd5b506100e76105a2565b34801561021957600080fd5b506100e76105b9565b34801561022e57600080fd5b506065546101ad9060ff1681565b34801561024857600080fd5b506100e7610257366004610767565b6105cd565b610264610646565b60005b818110156102cc57600160666000858585818110610287576102876107f7565b905060200201602081019061029c9190610767565b6001600160a01b031681526020810191909152604001600020805460ff1916911515919091179055600101610267565b505050565b6102d9610646565b60005b818110156102cc576000606660008585858181106102fc576102fc6107f7565b90506020020160208101906103119190610767565b6001600160a01b031681526020810191909152604001600020805460ff19169115159190911790556001016102dc565b670de0b6b3a7640000341161039d5760405162461bcd60e51b815260206004820152601d60248201527f5374616b696e673a20696e73756666696369656e74206465706f73697400000060448201526064015b60405180910390fd5b336001600160a01b038216146103f55760405162461bcd60e51b815260206004820152601d60248201527f5374616b696e673a206f6e6c792073656c662064656c65676174696f6e0000006044820152606401610394565b6040513481526001600160a01b0382169033907f510b11bb3f3c799b11307c01ab7db0d335683ef5b2da98f7697de744f465eacc9060200160405180910390a350565b610440610646565b61044a60006106a0565b565b60655460ff16158061046d57503360009081526066602052604090205460ff165b6104b05760405162461bcd60e51b815260206004820152601460248201527314dd185ada5b99ce881b9bdd08185b1b1bddd95960621b6044820152606401610394565b602181146105005760405162461bcd60e51b815260206004820152601e60248201527f5374616b696e673a20696e76616c6964207075626b6579206c656e67746800006044820152606401610394565b68056bc75e2d631000003410156105595760405162461bcd60e51b815260206004820152601d60248201527f5374616b696e673a20696e73756666696369656e74206465706f7369740000006044820152606401610394565b336001600160a01b03167fc7abef7b73f049da6a9bc2349ba5066a39e316eabc9f671b6f9406aa9490a4538383346040516105969392919061080d565b60405180910390a25050565b6105aa610646565b6065805460ff19166001179055565b6105c1610646565b6065805460ff19169055565b6105d5610646565b6001600160a01b03811661063a5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610394565b610643816106a0565b50565b6033546001600160a01b0316331461044a5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610394565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6000806020838503121561070557600080fd5b823567ffffffffffffffff8082111561071d57600080fd5b818501915085601f83011261073157600080fd5b81358181111561074057600080fd5b8660208260051b850101111561075557600080fd5b60209290920196919550909350505050565b60006020828403121561077957600080fd5b81356001600160a01b038116811461079057600080fd5b9392505050565b600080602083850312156107aa57600080fd5b823567ffffffffffffffff808211156107c257600080fd5b818501915085601f8301126107d657600080fd5b8135818111156107e557600080fd5b86602082850101111561075557600080fd5b634e487b7160e01b600052603260045260246000fd5b604081528260408201528284606083013760006060848301015260006060601f19601f860116830101905082602083015294935050505056fea2646970667358221220fb413f3fd014520bfeda4e23585dd938f5e4e5518b71218f8932096ccbbe1ad964736f6c63430008180033"
)

//go:embed staking_storage_layout.json
var stakingStorageLayoutJSON []byte

var StakingStorageLayout = mustGetStorageLayout(stakingStorageLayoutJSON)
