package query

import (
	"encoding/json"

	"github.com/kousuketk/myDB_go/bufpool"
)

type Executor struct {
	bufpool *bufpool.Storage
}

func NewExecutor(storage *bufpool.Storage) *Executor {
	return &Executor{
		bufpool: storage,
	}
}

func FullScan(store *bufpool.Storage) []*bufpool.Tuple {
	var result []*bufpool.Tuple

	for i := uint64(0); ; i++ {
		t, err := store.ReadTuple("sampleTable", i)
		if err != nil {
			break
		}

		result = append(result, t)
	}
	return result
}

func (e *Executor) SelectTable() ([]string, error) {
	// 直接scannerを呼ぶ(一時的)
	tuples := FullScan(e.bufpool)

	var values []string
	for _, t := range tuples {
		obj, _ := json.Marshal(t)
		values = append(values, string(obj))
	}

	return values, nil
}

func (e *Executor) InsertTable(data []interface{}) error {
	t := bufpool.NewTuple(data)
	e.bufpool.InsertTuple("sampleTable", t)
	err := e.bufpool.InsertIndex("sampleTable_pkey", t)
	return err
}

func (e *Executor) CreateTable(table string, pkey string) error {
	_, err := e.bufpool.CreateIndex(table + "_" + pkey)
	return err
}
