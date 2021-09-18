package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kousuketk/myDB_go/bufpool"
	"github.com/kousuketk/myDB_go/query"
)

func main() {
	fmt.Println("hello myDB")

	home := ".mydb/"
	if _, err := os.Stat(home); os.IsNotExist(err) {
		err := os.Mkdir(home, 0777)
		if err != nil {
			panic(err)
		}
	}
	storage := bufpool.NewStorage(home)
	executor := query.NewExecutor(storage)
	msg, err := executor.CreateTable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msg)
	msg, err = executor.InsertTable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msg)
	msgs, err := executor.SelectTable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msgs)
}
