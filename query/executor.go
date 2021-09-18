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

// SeqScan
func FullScan(store *bufpool.Storage) []*bufpool.Tuple {
	var result []*bufpool.Tuple

	for i := uint64(0); ; i++ {
		// t, err := store.ReadTuple(s.tblName, i)
		t, err := store.ReadTuple("sampleTable", i)
		if err != nil {
			// if no more pages, finish reading tuples.
			break
		}

		result = append(result, t)
	}
	return result
}

func (e *Executor) SelectTable() ([]string, error) {
	// 直接scannerを呼ぶ(一時的)
	tuples := FullScan(e.bufpool)

	// cols := []string{"id", "name"}

	var values []string
	for _, t := range tuples {
		obj, _ := json.Marshal(t.Data[0])
		values = append(values, string(obj))

		// for i, c := range cols {
		// 	obj, _ := json.Marshal(t.Data[i])
		// 	s := fmt.Sprintf(c, obj)
		// 	values = append(values, s)
		// }
	}

	return values, nil
}

func (e *Executor) InsertTable() (string, error) {
	str := []interface{}{
		"testInsert1",
		"testInsert2",
		"testInsert3",
	}
	t := bufpool.NewTuple(str)
	e.bufpool.InsertTuple("sampleTable", t)
	e.bufpool.InsertIndex("sampleTable_pkey", t)
	msg := "A row was inserted"
	return msg, nil
}

func (e *Executor) CreateTable() (string, error) {
	e.bufpool.CreateIndex("sampleTable" + "_" + "pkey")
	msg := "sampleTable" + " was created as Table"
	return msg, nil
}
