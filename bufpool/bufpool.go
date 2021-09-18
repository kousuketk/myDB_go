package bufpool

import (
	"crypto/md5"
	"encoding/binary"

	"github.com/kousuketk/myDB_go/pkg"
)

type bufferPool struct {
	lru   *pkg.Lru
	btree map[string]*pkg.BTree
}

type bufferTag struct {
	tableName string
	pgid      uint64
}

func newBufferTag(tableName string, pgid uint64) *bufferTag {
	return &bufferTag{
		tableName: tableName,
		pgid:      pgid,
	}
}

func (b *bufferTag) hash() [16]byte {
	from := []byte(b.tableName)
	pgidByte := make([]byte, 8)
	binary.BigEndian.PutUint64(pgidByte, b.pgid)
	from = append(from, pgidByte...)
	hash := md5.Sum(from)
	return hash
}

type pageDescriptor struct {
	tableName string
	pgid      uint64
	dirty     bool
	ref       uint64
	page      *Page
}

func newBufferPool() *bufferPool {
	return &bufferPool{
		lru:   pkg.NewLru(1000),
		btree: make(map[string]*pkg.BTree),
	}
}

func (b *bufferPool) toPgid(tid uint64) uint64 {
	return tid / TupleNumber
}

func (b *bufferPool) readPage(tableName string, pgid uint64) (*Page, error) {
	bt := newBufferTag(tableName, pgid)

	hash := bt.hash()
	p := b.lru.Get(hash)
	if p == nil {
		return nil, nil
	}

	pd := p.(*pageDescriptor)
	return pd.page, nil
}

func (b *bufferPool) appendTuple(tableName string, t *Tuple) bool {
	bt := newBufferTag(tableName, NewPgid(tableName))

	hash := bt.hash()
	p := b.lru.Get(hash)
	if p == nil {
		return false
	}

	pd := p.(*pageDescriptor)
	pd.dirty = true

	for i := range pd.page.Tuples {
		pd.page.Tuples[i] = *t
		break
	}

	return true
}

func (b *bufferPool) putPage(tableName string, pgid uint64, p *Page) (bool, *Page) {
	bt := newBufferTag(tableName, pgid)

	pd := &pageDescriptor{
		tableName: tableName,
		pgid:      pgid,
		page:      p,
		ref:       0,
		dirty:     false,
	}

	hash := bt.hash()
	victimPage := b.lru.Insert(hash, pd)
	if victimPage == nil {
		return false, nil
	}

	victim := victimPage.(*pageDescriptor)
	return victim.dirty, victim.page
}

func (b *bufferPool) readIndex(indexName string) (bool, *pkg.BTree) {
	tree, found := b.btree[indexName]
	return found, tree
}
