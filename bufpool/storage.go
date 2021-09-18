package bufpool

import (
	"log"

	"github.com/kousuketk/myDB_go/pkg"
)

type Storage struct {
	buffer *bufferPool
	disk   *diskManager
	prefix string
}

func NewStorage(home string) *Storage {
	return &Storage{
		buffer: newBufferPool(),
		disk:   newDiskManager(),
		prefix: home,
	}
}

func (s *Storage) insertPage(tableName string) {
	pg := NewPage()
	pgid := NewPgid(tableName)
	isNeedPersist, victim := s.buffer.putPage(tableName, pgid, pg)

	if isNeedPersist {
		// if a victim is dirty, its data must be persisted in the disk now.
		if err := s.disk.persist(s.prefix, tableName, pgid, victim); err != nil {
			log.Println(err)
		}
	}
}

func (s *Storage) InsertTuple(tablename string, t *Tuple) {
	for !s.buffer.appendTuple(tablename, t) {
		// if not exist in buffer, put a pageDescriptor to lru-cache
		s.insertPage(tablename)
		// res := s.buffer.appendTuple(tablename, t)
		// fmt.Println("res:", res)
	}
}

func (s *Storage) CreateIndex(indexName string) (*pkg.BTree, error) {
	btree := pkg.NewBTree()
	s.buffer.btree[indexName] = btree
	return btree, nil
}

func (s *Storage) InsertIndex(indexName string, item pkg.Item) error {
	btree, err := s.ReadIndex(indexName)
	if err != nil {
		return err
	}

	btree.Insert(item)
	return nil
}

func (s *Storage) ReadIndex(indexName string) (*pkg.BTree, error) {
	found, tree := s.buffer.readIndex(indexName)

	if found {
		return tree, nil
	}

	tree, err := s.disk.readIndex(indexName)
	if err != nil || tree == nil {
		return s.CreateIndex(indexName)
	}

	return tree, nil
}

func (s *Storage) ReadTuple(tableName string, tid uint64) (*Tuple, error) {
	pgid := s.buffer.toPgid(tid)

	pg, err := s.readPage(tableName, pgid)
	if err != nil {
		return nil, err
	}

	return &pg.Tuples[tid%TupleNumber], nil
}

func (s *Storage) readPage(tableName string, pgid uint64) (*Page, error) {
	pg, err := s.buffer.readPage(tableName, pgid)

	if err != nil {
		return nil, err
	}

	// if a page exists in the buffer, return that.
	if pg != nil {
		return pg, nil
	}

	pg, err = s.disk.fetchPage(s.prefix, tableName, pgid)

	if err != nil {
		return nil, err
	}
	s.buffer.putPage(tableName, pgid, pg)

	return pg, nil
}

func (s *Storage) Terminate() error {
	items := s.buffer.lru.GetAll()
	for _, item := range items {
		pd := item.(*pageDescriptor)
		if pd.dirty {
			err := s.disk.persist(s.prefix, pd.tableName, pd.pgid, pd.page)
			if err != nil {
				return err
			}
		}
	}

	for key, val := range s.buffer.btree {
		err := s.disk.writeIndex(s.prefix, key, val)
		if err != nil {
			return err
		}
	}

	return nil
}
