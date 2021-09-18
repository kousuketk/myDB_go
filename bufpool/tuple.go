package bufpool

import (
	"encoding/json"

	"github.com/kousuketk/myDB_go/pkg"
)

type TupleDataType int32

const (
	INT TupleDataType = iota
	STRING
)

type TupleData struct {
	Type TupleDataType `json:"type"`
	Num  int32         `json:"num"`
	Str  string        `json:"str"`
}

type Tuple struct {
	Data []TupleData `json:"tupleData"`
}

func NewTuple(values []interface{}) *Tuple {
	var t Tuple

	var td *TupleData
	for _, v := range values {
		switch concrete := v.(type) {

		case int:
			td = &TupleData{
				Type: INT,
				Num:  int32(concrete),
			}

		case string:
			td = &TupleData{
				Type: STRING,
				Str:  string(concrete),
			}
		}

		t.Data = append(t.Data, *td)
	}

	return &t
}

func (m *Tuple) Less(than pkg.Item) bool {
	t, ok := than.(*Tuple)
	if !ok {
		return false
	}

	// FIXME
	left := m.Data[0].Num
	right := t.Data[0].Num

	return left < right
}

func SerializeTuple(t *Tuple) ([128]byte, error) {
	out, err := json.Marshal(t)

	if err != nil {
		return [128]byte{}, err
	}

	var b [128]byte
	copy(b[:], out)

	return b, nil
}

func DeserializeTuple(b [128]byte) (*Tuple, error) {
	var t Tuple

	err := json.Unmarshal(b[:], &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (m *Tuple) Equal(order int, s string, n int) bool {
	tupleData := m.Data[order]

	if tupleData.Type == STRING {
		return tupleData.Str == s
	} else if tupleData.Type == INT {
		return tupleData.Num == int32(n)
	}

	return false
}
