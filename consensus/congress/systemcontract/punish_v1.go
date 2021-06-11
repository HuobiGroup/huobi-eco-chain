package systemcontract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/congress/vmcaller"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math"
	"math/big"
)

const (
	punishV1Code = "0x608060405234801561001057600080fd5b50600436106101375760003560e01c8063741579b1116100b8578063d93d2cb91161007c578063d93d2cb914610243578063e0d8ea5314610260578063ea7221a114610268578063ec0cb33614610172578063f3b1cc6714610172578063f62af26c1461028e57610137565b8063741579b1146102025780638129fc1c1461020a5780639001eed814610214578063c967f90f1461021c578063cb1ea7251461023b57610137565b806332f3c17f116100ff57806332f3c17f1461018257806344c1aa99146101a857806344f99900146101b057806363e1d451146101d457806371a1bb75146101fa57610137565b806303fab4f61461013c578063158ef93e1461015657806315de360e146101725780632897183d1461017a5780632e4f67e414610172575b600080fd5b6101446102ab565b60408051918252519081900360200190f35b61015e6102b7565b604080519115158252519081900360200190f35b6101446102c0565b6101446102c5565b6101446004803603602081101561019857600080fd5b50356001600160a01b03166102cb565b6101446102e6565b6101b86102ec565b604080516001600160a01b039092168252519081900360200190f35b61015e600480360360208110156101ea57600080fd5b50356001600160a01b03166102f2565b6101b861059a565b6101446105a0565b6102126105ac565b005b61014461061b565b610224610627565b6040805161ffff9092168252519081900360200190f35b61014461062c565b6102126004803603602081101561025957600080fd5b5035610632565b6101446108d5565b6102126004803603602081101561027e57600080fd5b50356001600160a01b03166108db565b6101b8600480360360208110156102a457600080fd5b5035610cd3565b671bc16d674ec8000081565b60005460ff1681565b606481565b60035481565b6001600160a01b031660009081526004602052604090205490565b60025481565b61f00681565b6000805460ff16610339576040805162461bcd60e51b815260206004820152600c60248201526b139bdd081a5b9a5d081e595d60a21b604482015290519081900360640190fd5b604080516365f69f9760e01b81526001600160a01b03841660048201529051339161f005916365f69f9791602480820192602092909190829003018186803b15801561038457600080fd5b505afa158015610398573d6000803e3d6000fd5b505050506040513d60208110156103ae57600080fd5b50516001600160a01b03161461040b576040805162461bcd60e51b815260206004820152601860248201527f56616c696461746f72206e6f7420726567697374657265640000000000000000604482015290519081900360640190fd5b6001600160a01b03821660009081526004602052604090205415610443576001600160a01b0382166000908152600460205260408120555b6001600160a01b03821660009081526004602052604090206002015460ff16801561046f575060055415155b15610592576005546001600160a01b0383166000908152600460205260409020600101546000199091011461053957600580546000919060001981019081106104b457fe5b60009182526020808320909101546001600160a01b03868116845260049092526040909220600101546005805492909316935083929181106104f257fe5b600091825260208083209190910180546001600160a01b0319166001600160a01b039485161790558583168252600490526040808220600190810154949093168252902001555b600580548061054457fe5b60008281526020808220830160001990810180546001600160a01b03191690559092019092556001600160a01b038416825260049052604081206001810191909155600201805460ff191690555b506001919050565b61f00581565b670de0b6b3a764000081565b60005460ff16156105fa576040805162461bcd60e51b8152602060048201526013602482015272105b1c9958591e481a5b9a5d1a585b1a5e9959606a1b604482015290519081900360640190fd5b6018600181815560306002556003919091556000805460ff19169091179055565b6729a2241af62c000081565b600581565b60015481565b334114610673576040805162461bcd60e51b815260206004820152600a6024820152694d696e6572206f6e6c7960b01b604482015290519081900360640190fd5b4360009081526007602052604090205460ff16156106cc576040805162461bcd60e51b8152602060048201526011602482015270105b1c9958591e48191958dc99585cd959607a1b604482015290519081900360640190fd5b60005460ff16610712576040805162461bcd60e51b815260206004820152600c60248201526b139bdd081a5b9a5d081e595d60a21b604482015290519081900360640190fd5b8080438161071c57fe5b0615610762576040805162461bcd60e51b815260206004820152601060248201526f426c6f636b2065706f6368206f6e6c7960801b604482015290519081900360640190fd5b436000908152600760205260409020805460ff19166001179055600554610788576108d1565b60005b6005548110156108a657600354600254816107a257fe5b0460046000600584815481106107b457fe5b60009182526020808320909101546001600160a01b03168352820192909252604001902054111561086557600354600254816107ec57fe5b0460046000600584815481106107fe57fe5b60009182526020808320909101546001600160a01b0316835282019290925260400181205460058054939091039260049291908590811061083b57fe5b60009182526020808320909101546001600160a01b0316835282019290925260400190205561089e565b6000600460006005848154811061087857fe5b60009182526020808320909101546001600160a01b031683528201929092526040019020555b60010161078b565b506040517f181d51be54e8e8eaca6eae0eab32d4162099236bd519e7238d015d0870db464190600090a15b5050565b60055490565b33411461091c576040805162461bcd60e51b815260206004820152600a6024820152694d696e6572206f6e6c7960b01b604482015290519081900360640190fd5b60005460ff16610962576040805162461bcd60e51b815260206004820152600c60248201526b139bdd081a5b9a5d081e595d60a21b604482015290519081900360640190fd5b4360009081526006602052604090205460ff16156109ba576040805162461bcd60e51b815260206004820152601060248201526f105b1c9958591e481c1d5b9a5cda195960821b604482015290519081900360640190fd5b436000908152600660209081526040808320805460ff191660011790556001600160a01b0384168352600490915290206002015460ff16610a6357600580546001600160a01b038316600081815260046020526040812060018082018590558085019095557f036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db090930180546001600160a01b0319168317905552600201805460ff191690911790555b6001600160a01b03811660009081526004602052604090208054600101908190556002549081610a8f57fe5b06610b8d57600061f0056001600160a01b03166365f69f97836040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b158015610ae557600080fd5b505afa158015610af9573d6000803e3d6000fd5b505050506040513d6020811015610b0f57600080fd5b50516040805163209b4f7b60e21b815290519192506001600160a01b0383169163826d3dec9160048082019260009290919082900301818387803b158015610b5657600080fd5b505af1158015610b6a573d6000803e3d6000fd5b505050506001600160a01b03821660009081526004602052604081205550610c91565b6001546001600160a01b03821660009081526004602052604090205481610bb057fe5b06610c9157600061f0056001600160a01b03166365f69f97836040518263ffffffff1660e01b815260040180826001600160a01b0316815260200191505060206040518083038186803b158015610c0657600080fd5b505afa158015610c1a573d6000803e3d6000fd5b505050506040513d6020811015610c3057600080fd5b50516040805163ba26d9ff60e01b815290519192506001600160a01b0383169163ba26d9ff9160048082019260009290919082900301818387803b158015610c7757600080fd5b505af1158015610c8b573d6000803e3d6000fd5b50505050505b6040805142815290516001600160a01b038316917f770e0cca42c35d00240986ce8d3ed438be04663c91dac6576b79537d7c180f1e919081900360200190a250565b60058181548110610ce057fe5b6000918252602090912001546001600160a01b031690508156fea26469706673582212206e0fbda1c640ec5461f4d4f2a854284412f5a8279bbbd98f71df5a8c51e0777e64736f6c634300060c0033"
)

type hardForkPunishV1 struct {
}

func (s *hardForkPunishV1) GetName() string {
	return PunishV1ContractName
}

func (s *hardForkPunishV1) Update(config *params.ChainConfig, height *big.Int, state *state.StateDB) (err error) {
	contractCode := common.FromHex(punishV1Code)

	//write code to sys contract
	state.SetCode(PunishV1ContractAddr, contractCode)
	log.Debug("Write code to system contract account", "addr", PunishV1ContractAddr.String(), "code", punishV1Code)

	return
}

func (s *hardForkPunishV1) Execute(state *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (err error) {
	// initialize v1 contract
	method := "initialize"
	data, err := GetInteractiveABI()[s.GetName()].Pack(method)
	if err != nil {
		log.Error("Can't pack data for initialize", "error", err)
		return err
	}

	msg := types.NewMessage(header.Coinbase, &PunishV1ContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	_, err = vmcaller.ExecuteMsg(msg, state, header, chainContext, config)

	return
}
