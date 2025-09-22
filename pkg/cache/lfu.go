package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(k K) (v V, ok bool)
	Set(k K, v V)
	Invalidate(k K)
	Stop()
}
type lfuEntry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
	freq      int
}

type LFUCache[K comparable, V any] struct {
	capacity int
	ttl      time.Duration
	items    map[K]*lfuEntry[K, V]
	mu       sync.RWMutex
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewLFUCache[K comparable, V any](ttl time.Duration, capacity int) *LFUCache[K, V] {
	c := &LFUCache[K, V]{
		capacity: capacity,
		ttl:      ttl,
		items:    make(map[K]*lfuEntry[K, V]),
		stopChan: make(chan struct{}),
	}
	c.wg.Add(1)
	go c.cleanupWorker()
	return c
}

func (c *LFUCache[K, V]) Set(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.items[k]; ok {
		e.value = v
		e.freq++
		e.expiresAt = time.Now().Add(c.ttl)
		return
	}

	if len(c.items) >= c.capacity {
		var minKey K
		minFreq := int(^uint(0) >> 1)
		for key, entry := range c.items {
			if entry.freq < minFreq {
				minFreq = entry.freq
				minKey = key
			}
		}
		delete(c.items, minKey)
	}
	c.items[k] = &lfuEntry[K, V]{k, v, time.Now().Add(c.ttl), 1}
}

func (c *LFUCache[K, V]) Get(k K) (v V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exist := c.items[k]
	if !exist || time.Now().After(entry.expiresAt) {
		return
	}
	entry.freq++
	return entry.value, true
}

func (c *LFUCache[K, V]) Stop() {
	close(c.stopChan)
	c.wg.Wait()
}

func (c *LFUCache[K, V]) Invalidate(k K) {
	c.mu.Lock()
	delete(c.items, k)
	c.mu.Unlock()
}

func (c *LFUCache[K, V]) cleanupWorker() {
	defer c.wg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			c.mu.Lock()
			for k, entry := range c.items {
				if now.After(entry.expiresAt) {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		case <-c.stopChan:
			return
		}
	}
}
