package idgen

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	epoch int64 = 1704038400000

	timestampBits = 41
	workerBits    = 10
	randomBits    = 12

	workerMax    = -1 ^ (-1 << workerBits)
	randomMax    = -1 ^ (-1 << randomBits)
	timestampMax = -1 ^ (-1 << timestampBits)

	timestampShift = randomBits + workerBits
	workerShift    = randomBits

	prime1 = 2862933555777941757
	prime2 = 3037000493
)

type Snowflake struct {
	mu       sync.Mutex
	workerId int64
}

func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > workerMax {
		return nil, fmt.Errorf("机器ID必须在 0-%d 之间", workerMax)
	}

	return &Snowflake{
		workerId: machineID,
	}, nil
}

func (s *Snowflake) generateRandom() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return time.Now().UnixNano() & randomMax
	}
	return int64(binary.BigEndian.Uint64(b[:])) & randomMax
}

func (s *Snowflake) scramble(id int64) int64 {
	id = ((id ^ prime1) * prime2) & math.MaxInt64
	id = id % 9000000000
	id += 1000000000
	return id
}

func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixMilli() - epoch
	if timestamp > timestampMax {
		timestamp = timestamp % timestampMax
	}

	random := s.generateRandom()

	rawId := (timestamp << timestampShift) |
		(s.workerId << workerShift) |
		random

	return s.scramble(rawId)
}
