package idgen

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"
)

// Snowflake ID 生成器（生成 9 位数字）
// 参考 duidui 项目实现，使用标准 Snowflake + 混淆函数

const (
	// 起始时间戳 (2024-01-01 00:00:00 +0800 CST) 毫秒
	epoch int64 = 1704038400000

	// 位数分配（标准 Snowflake）
	timestampBits = 41
	workerBits    = 10
	randomBits    = 12

	// 最大值
	workerMax    = -1 ^ (-1 << workerBits)    // 1023
	randomMax    = -1 ^ (-1 << randomBits)    // 4095
	timestampMax = -1 ^ (-1 << timestampBits) // 2199023255551

	// 左移位数
	timestampShift = randomBits + workerBits // 22
	workerShift    = randomBits              // 12

	// 用于混淆的大质数
	prime1 = 2862933555777941757
	prime2 = 3037000493
)

type Snowflake struct {
	mu       sync.Mutex
	workerId int64
}

// NewSnowflake 创建雪花算法生成器
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > workerMax {
		return nil, fmt.Errorf("机器ID必须在 0-%d 之间", workerMax)
	}

	return &Snowflake{
		workerId: machineID,
	}, nil
}

// generateRandom 生成加密级随机数
func (s *Snowflake) generateRandom() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		// 如果随机数生成失败，使用时间的纳秒部分
		return time.Now().UnixNano() & randomMax
	}
	return int64(binary.BigEndian.Uint64(b[:])) & randomMax
}

// scramble 混淆函数，使用大质数和位运算打散 ID，生成9位数字
func (s *Snowflake) scramble(id int64) int64 {
	// 使用多个质数和位运算进行混淆
	id = ((id ^ prime1) * prime2) & math.MaxInt64

	// 确保生成的是9位数
	// 999999999 是最大的9位数
	// 100000000 是最小的9位数
	id = id % 900000000 // 限制在 [0, 900000000) 范围内
	id += 100000000     // 加上最小9位数，确保是9位

	return id
}

// NextID 生成下一个ID（9位数字）
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixMilli() - epoch
	if timestamp > timestampMax {
		// 如果时间戳超出范围，取模处理
		timestamp = timestamp % timestampMax
	}

	// 生成随机数
	random := s.generateRandom()

	// 组合 ID：时间戳 + 机器ID + 随机数（标准 Snowflake 格式）
	rawId := (timestamp << timestampShift) |
		(s.workerId << workerShift) |
		random

	// 混淆 ID 生成9位数字
	return s.scramble(rawId)
}
