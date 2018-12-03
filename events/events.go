package events

import (
	"../plasmacontract"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

type EventCoinDeposited struct {
	Who         common.Address `json:who`
	Amount      *big.Int       `json:amount`
	BlockNumber uint64         `json:blockNumber`
}

var eventGroup = make([]EventCoinDeposited, 0)
var currentBlock uint64 = 0

func GetEvent() (bool,error) {
	client, err := ethclient.Dial("http://localhost:9545")
	if err != nil {
		return false, err
	}

	maxBlock, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return false, err
	}

	contractAddress := common.HexToAddress("0xa86a2c6B81C22d25D8EBAf8cE52895F46A263348")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(currentBlock)),
		ToBlock:   big.NewInt(int64(checker(currentBlock, maxBlock.Number.Uint64()))),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return false, err
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(store.StoreABI)))
	if err != nil {
		return false, err
	}


	Sig := []byte("CoinDeposited(address,uint256)")
	SigHash := crypto.Keccak256Hash(Sig)

	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case SigHash.Hex():
			var depositEvent EventCoinDeposited
			err := contractAbi.Unpack(&depositEvent, "CoinDeposited", vLog.Data)
			if err != nil {
				return false, err
			}
			depositEvent.Who = common.HexToAddress(vLog.Topics[1].Hex())
			depositEvent.BlockNumber = vLog.BlockNumber
			PutEventsToGroup(depositEvent)

		}
	}

	if currentBlock <= maxBlock.Number.Uint64() + 10 {
		SetLastBlock(currentBlock + 10)
	}

	if currentBlock >= maxBlock.Number.Uint64() {
		return false, nil
	}
	return true, nil
}

func SetLastBlock(v uint64) {
	currentBlock = v
}

func checker(current, final uint64) uint64 {
	if current + 10 <= final {
		current = current+10
	} else {
		delta := final - current
		current = current + delta
	}
	return current
}


func PutEventsToGroup(e EventCoinDeposited) {
	eventGroup = append(eventGroup, e)
}

func ShowGroup() {
	for i := range eventGroup {
		fmt.Printf("BlockNumber: %d\n", eventGroup[i].BlockNumber)
		fmt.Printf("Amount: %s\n", eventGroup[i].Amount.String())
		fmt.Printf("Who: %s\n", eventGroup[i].Who.String())
	}
}
