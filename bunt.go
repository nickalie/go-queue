package queue

import (
	"fmt"
	"github.com/tidwall/buntdb"
	"time"
	"strings"
)

// BuntBackend provides BuntDB-based backend to manage queues.
// https://github.com/tidwall/buntdb
// Suitable for multithreaded single process environment
type BuntBackend struct {
	db       *buntdb.DB
	codec    Codec
	interval time.Duration
	key      string
}

// NewBuntBackend creates new BuntBackend.
func NewBuntBackend(path string) (*BuntBackend, error) {
	db, err := buntdb.Open(path)

	if err != nil {
		return nil, err
	}

	return NewBuntBackendFromDB(db), nil
}

// NewBuntBackendFromDB creates new BuntBackend from bunt.DB object.
func NewBuntBackendFromDB(db *buntdb.DB) *BuntBackend {
	var key string
	var err error

	db.View(func(tx *buntdb.Tx) error {
		key, err = tx.Get("baseKey")
		return err
	})
	
	b := &BuntBackend{db: db, key: key}
	return b.Codec(NewGOBCodec()).Interval(time.Second)
}

// Codec sets codec to encode/decode objects in queues. GOBCodec is default.
func (b *BuntBackend) Codec(c Codec) *BuntBackend {
	b.codec = c
	return b
}

// Interval sets interval to poll new queue element. Default value is one second.
func (b *BuntBackend) Interval(interval time.Duration) *BuntBackend {
	b.interval = interval
	return b
}

// Put adds value to the end of a queue.
func (b *BuntBackend) Put(queueName string, value interface{}) error {
	data, err := b.codec.Marshal(value)

	if err != nil {
		return err
	}

	return b.db.Update(func(tx *buntdb.Tx) error {
		b.key = increaseString(b.key)
		key := fmt.Sprintf("%s:%s", queueName, b.key)
		tx.Set("baseKey", b.key, nil)
		_, _, err = tx.Set(key, string(data), nil)
		return err
	})
}

// Get removes the first element from a queue and put it in the value pointed to by v
func (b *BuntBackend) Get(queueName string, v interface{}) error {
	queueName += ":"
	var value string
	var k string
	for {
		found := false
		err := b.db.Update(func(tx *buntdb.Tx) error {
			err := tx.Ascend("", func(key, value string) bool {
				if strings.HasPrefix(key, queueName) {
					k = key
					found = true
					return false
				}

				return true
			})
			/*err := tx.AscendKeys(queueName+":*", func(key, value string) bool {
				k = key
				found = true
				return false
			})*/

			if err != nil {
				return err
			}

			value, err = tx.Delete(k)

			return err
		})

		if err != nil && err != buntdb.ErrNotFound {
			return err
		}

		if found {
			break
		} else {
			time.Sleep(b.interval)
		}
	}

	return b.codec.Unmarshal([]byte(value), v)
}

func (b *BuntBackend) RemoveQueue(queueName string) error {
	return nil
}

// Close closes buntdb
func (b *BuntBackend) Close() error {
	return b.db.Close()
}
