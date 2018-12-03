package db

import (
	"log"

	"./leveldb"
)

func Tx(dbName string) *leveldb.Table {
	db, err := leveldb.Connect(dbName, 16, 16, 1024)
	txTable := leveldb.NewTable(db, "tx-")
	if err != nil {
		panic(err)
	}
	return txTable
}

func Event(dbName string) *leveldb.Table {
	db, err := leveldb.Connect(dbName, 16, 16, 1024)
	event := leveldb.NewTable(db, "event-")
	if err != nil {
		log.Println(err)
	}
	return event
}

func Block(dbName string) *leveldb.Table {
	db, err := leveldb.Connect(dbName, 16, 16, 1024)
	block := leveldb.NewTable(db, "block-")
	if err != nil {
		panic(err)
	}
	return block
}
