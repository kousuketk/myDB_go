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
	err := executor.CreateTable("sampleTable", "pkey")
	if err != nil {
		log.Fatal(err)
	}
	data := []interface{}{
		"testInsert1",
		"testInsert2",
		"testInsert3",
	}
	err = executor.InsertTable(data)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := executor.SelectTable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msgs)
}
