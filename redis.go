package queue

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// RedisBackend provides redis-based backend to manage queues.
// Suitable for multi-host, multi-process and multithreaded environment
type RedisBackend struct {
	pool  *redis.Pool
	codec Codec
}

// NewRedisBackend creates new RedisBackend
func NewRedisBackend(redisURL string) (*RedisBackend, error) {
	return NewRedisBackendWithPool(&redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.DialURL(redisURL) },
	}), nil
}

// NewRedisBackendWithPool creates new RedisBackend with provided redis.Pool
func NewRedisBackendWithPool(pool *redis.Pool) *RedisBackend {
	b := &RedisBackend{pool: pool}
	return b.Codec(NewGOBCodec())
}

// Codec sets codec to encode/decode objects in queues. GOBCodec is default.
func (b *RedisBackend) Codec(c Codec) *RedisBackend {
	b.codec = c
	return b
}

// Put adds value to the end of a queue.
func (b *RedisBackend) Put(queueName string, value interface{}) error {
	data, err := b.codec.Marshal(value)

	if err != nil {
		return err
	}

	_, err = b.pool.Get().Do("RPUSH", queueName, data)
	return err
}

// Get removes the first element from a queue and put it in the value pointed to by v
func (b *RedisBackend) Get(queueName string, v interface{}) error {
	d, err := redis.ByteSlices(b.pool.Get().Do("BLPOP", queueName, 0))

	if err != nil {
		return err
	}

	return b.codec.Unmarshal(d[1], v)
}
