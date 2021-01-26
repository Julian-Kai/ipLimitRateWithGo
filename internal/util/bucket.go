package util

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Bucket struct {
	startTime       time.Time     // 開始時間
	fillInterval    time.Duration // 多久填充一次
	mu              sync.Mutex    //
	quantum         int64         // 填充的量
	capacity        int64         // 桶的容量
	availableTokens int64         // 可使用的令牌數量
	latestTick      int64         // 桶中最新的刻度
}

func NewBucket(fillInterval time.Duration, quantum, capacity int64) *Bucket {
	if fillInterval <= 0 {
		logrus.Fatal("token bucket fill interval is not > 0")
	}
	if quantum <= 0 {
		logrus.Fatal("token bucket quantum is not > 0")
	}
	if capacity <= 0 {
		logrus.Fatal("token bucket capacity is not > 0")
	}
	return &Bucket{
		startTime:       time.Now(),
		fillInterval:    fillInterval,
		quantum:         quantum,
		capacity:        capacity,
		availableTokens: capacity,
		latestTick:      0,
	}
}

func (bkt *Bucket) TakeAvailable(count int64) int64 {
	bkt.mu.Lock()
	defer bkt.mu.Unlock()
	return bkt.takeAvailable(time.Now(), count)
}

func (bkt *Bucket) takeAvailable(now time.Time, count int64) int64 {
	if count <= 0 {
		return 0
	}

	bkt.adjustAvailableTokens(bkt.currentTick(now))
	if bkt.availableTokens <= 0 {
		return 0
	}
	if count > bkt.availableTokens {
		count = bkt.availableTokens
	}

	bkt.availableTokens -= count
	return count
}

func (bkt *Bucket) currentTick(now time.Time) int64 {
	return int64(now.Sub(bkt.startTime) / bkt.fillInterval)
}

// 當前令牌數 = 上一次剩餘的令牌數 + (本次取令牌的時刻-上一次取令牌的時刻)/放置令牌的時間間隔 * 每次放置的令牌數
func (bkt *Bucket) adjustAvailableTokens(tick int64) {
	lastTick := bkt.latestTick
	bkt.latestTick = tick
	if bkt.availableTokens >= bkt.capacity {
		return
	}
	bkt.availableTokens += (tick - lastTick) * bkt.quantum
	if bkt.availableTokens > bkt.capacity {
		bkt.availableTokens = bkt.capacity
	}
	return
}

func (bkt *Bucket) Capacity() int64 {
	return bkt.capacity
}

func (bkt *Bucket) AvailableTokens() int64 {
	return bkt.availableTokens
}
