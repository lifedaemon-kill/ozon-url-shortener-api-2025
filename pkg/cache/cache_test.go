package cache

import (
	"testing"
	"time"
)

func TestLFUCache(t *testing.T) {
	t.Parallel()
	cache := NewLFUCache[string, string](1*time.Second, 2)
	defer cache.Stop()

	cache.Set("a", "val1")
	cache.Set("b", "val2")

	cache.Get("a")
	cache.Get("a")
	cache.Get("b")

	cache.Set("c", "val3")

	if _, ok := cache.Get("b"); ok {
		t.Error("expected 'b' to be evicted")
	}

	if _, ok := cache.Get("a"); !ok {
		t.Error("expected 'a' to remain")
	}

	if _, ok := cache.Get("c"); !ok {
		t.Error("expected 'c' to be present")
	}

	cache.Set("x", "ttl")
	time.Sleep(1100 * time.Millisecond)
	if _, ok := cache.Get("x"); ok {
		t.Error("expected 'x' to expire")
	}

	cache.Set("a", "new")
	cache.Set("d", "val4")

	if _, ok := cache.Get("a"); !ok {
		t.Error("expected 'a' to NOT be evicted after overwrite")
	}
}
