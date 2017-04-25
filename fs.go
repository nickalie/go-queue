package queue

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var errFound = errors.New("Found")
var errNotFound = errors.New("Not Found")

// FSBackend uses file system to manage queues.
// Suitable for multithreaded and multi-process environments.
type FSBackend struct {
	path     string
	codec    Codec
	key      string
	interval time.Duration
}

// NewFSBackend creates new FSBackend.
func NewFSBackend(path string) (*FSBackend, error) {
	return &FSBackend{path: path, codec: NewGOBCodec(), interval: time.Second}, nil
}

// Codec sets codec to encode/decode objects in queues. GOBCodec is default.
func (b *FSBackend) Codec(c Codec) *FSBackend {
	b.codec = c
	return b
}

// Interval sets interval to poll new queue element. Default value is one second.
func (b *FSBackend) Interval(interval time.Duration) *FSBackend {
	b.interval = interval
	return b
}

// Put adds value to the end of a queue.
func (b *FSBackend) Put(queueName string, value interface{}) error {
	path := filepath.Join(b.path, queueName)

	err := os.MkdirAll(path, 0777)

	if err != nil {
		return err
	}

	b.key = increaseString(b.key)
	fileName := filepath.Join(path, b.key)

	l, err := getLock(path)

	if err != nil {
		return err
	}

	data, err := b.codec.Marshal(value)

	if err != nil {
		l.unlock()
		return err
	}

	err = ioutil.WriteFile(fileName, data, 0777)

	if err != nil {
		l.unlock()
		return err
	}

	return l.unlock()
}

// Get removes the first element from a queue and put it in the value pointed to by v
func (b *FSBackend) Get(queueName string, value interface{}) error {
	dir := filepath.Join(b.path, queueName)

	var fileName string
	var l *lock

	for {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() || filepath.Ext(path) == ".lock" {
				return nil
			}

			l, err = getLock(path)

			if err != nil {
				return nil
			}

			fileName = path
			return errFound
		})

		if err == errFound {
			err = b.readFile(l, fileName, value)
		} else {
			err = errNotFound
		}

		if err != nil {
			time.Sleep(time.Second)
		} else {
			return nil
		}
	}
}

func (b *FSBackend) readFile(l *lock, fileName string, value interface{}) error {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		l.unlock()
		return err
	}

	err = l.unlock()

	if err != nil {
		return err
	}

	err = os.Remove(fileName)

	if err != nil {
		return err
	}

	return b.codec.Unmarshal(data, value)
}

var errLocked = errors.New("Locked")

type lock struct {
	path string
	id   []byte
}

func getLock(path string) (*lock, error) {
	lockFile := path + ".lock"

	_, err := os.Stat(lockFile)

	if err == nil {
		return nil, errLocked
	}

	id := []byte(strconv.FormatUint(rand.Uint64(), 10))
	err = ioutil.WriteFile(lockFile, id, 0777)

	if err != nil {
		os.Remove(lockFile)
		return nil, err
	}

	return &lock{path: lockFile, id: id}, nil
}

func (l *lock) unlock() error {
	id, err := ioutil.ReadFile(l.path)

	if err != nil {
		os.Remove(l.path)
		return err
	}

	if !bytes.Equal(l.id, id) {
		os.Remove(l.path)
		return errLocked
	}

	return os.Remove(l.path)
}
