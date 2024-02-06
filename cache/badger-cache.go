package cache

import (
	"time"

	"github.com/dgraph-io/badger/v4"
)

type BadgerCache struct {
	Conn   *badger.DB
	Prefix string
}

func (b *BadgerCache) Has(str string) (bool, error) {
	if _, err := b.Get(str); err != nil {
		return false, nil
	}
	return true, nil
}

func (b *BadgerCache) Get(str string) (any, error) {
	var fromCache []byte

	if err := b.Conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(str))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			fromCache = append([]byte{}, val...)
			return nil
		})

		return err
	}); err != nil {
		return nil, err
	}

	decoded, err := decode(string(fromCache))
	if err != nil {
		return nil, err
	}

	item := decoded[str]

	return item, nil
}

func (b *BadgerCache) Set(str string, value any, ttl ...int) error {
	entry := Entry{}
	entry[str] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(ttl) > 0 {
		if err := b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded).WithTTL(time.Second * time.Duration(ttl[0]))
			err = txn.SetEntry(e)
			return err
		}); err != nil {
			return err
		}
	} else {
		if err := b.Conn.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded)
			err = txn.SetEntry(e)
			return err
		}); err != nil {
			return err
		}
	}
	return nil
}

func (b *BadgerCache) Forget(str string) error {
	return b.Conn.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(str))
	})
}

func (b *BadgerCache) EmptyByMatch(str string) error {
	return b.emptyByMatch(str)
}

func (b *BadgerCache) Empty() error {
	return b.emptyByMatch("")
}

func (b *BadgerCache) emptyByMatch(str string) error {

	deleteKeys := func(keysForDelete [][]byte) error {
		return b.Conn.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		})
	}

	collectSize := 100000

	return b.Conn.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		for it.Seek([]byte(str)); it.ValidForPrefix([]byte(str)); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					return err
				}
			}

		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}

		return nil
	})
}
