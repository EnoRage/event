package events

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"log"

	"../db"
	"../plasmacontract"
	"../utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventCoinDeposited struct {
	Who         common.Address `json:who`
	Amount      *big.Int       `json:amount`
	BlockNumber uint64         `json:blockNumber`
}

var eventGroup = make([]EventCoinDeposited, 0)

var currentBlock uint64 = 0

func GetEvent() bool {

	client, err := ethclient.Dial("http://localhost:9545")
	if err != nil {
		fmt.Println(err)
	}

	maxBlock, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(store.StoreABI)))
	if err != nil {
		fmt.Println(err)
	}

	Sig := []byte("CoinDeposited(address,uint256)")
	SigHash := crypto.Keccak256Hash(Sig)

	for _, vLog := range logs {
		if vLog.Topics[0].Hex() == SigHash.Hex() {
			var depositEvent EventCoinDeposited
			err := contractAbi.Unpack(&depositEvent, "CoinDeposited", vLog.Data)
			if err != nil {
				fmt.Println(err)
			}
			depositEvent.Who = common.HexToAddress(vLog.Topics[1].Hex())
			depositEvent.BlockNumber = vLog.BlockNumber
			err = PutEventsToGroup(depositEvent)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if currentBlock <= maxBlock.Number.Uint64()+10 {
		SetLastBlock(currentBlock + 10)
	}

	if currentBlock >= maxBlock.Number.Uint64() {
		return false
	}

	return true
}

func SetLastBlock(v uint64) {
	currentBlock = v
	err := db.Event("database").Put([]byte("CurrentBlock"), utils.UintToBytesArray(v))
	if err != nil {
		fmt.Println(err)
	}
}

func checker(current, final uint64) uint64 {
	if current+10 <= final {
		current = current + 10
	} else {
		delta := final - current
		current = current + delta
	}
	return current
}

func PutEventsToGroup(e EventCoinDeposited) error {

	eventGroup = append(eventGroup, e)

	j, err := json.Marshal(eventGroup)
	if err != nil {
		return err
	}

	err = db.Event("database").Put([]byte("EventGroup"), j)
	if err != nil {
		return err
	}

	return nil
}

func ShowGroup() {

	//for i := range eventGroup {
	//	fmt.Printf("BlockNumber: %d\n", eventGroup[i].BlockNumber)
	//	fmt.Printf("Amount: %s\n", eventGroup[i].Amount.String())
	//	fmt.Printf("Who: %s\n", eventGroup[i].Who.String())
	//}

	b, err := db.Event("database").Get([]byte("EventGroup"))
	if err != nil {
		fmt.Println(err)
	}

	eG := []EventCoinDeposited{}

	err = json.Unmarshal(b, &eG)

	if err != nil {
		log.Printf("error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", b)
	}

	for i := range eG {
		fmt.Printf("BlockNumber: %d\n", eG[i].BlockNumber)
		fmt.Printf("Amount: %s\n", eG[i].Amount.String())
		fmt.Printf("Who: %s\n", eG[i].Who.String())
	}

}
