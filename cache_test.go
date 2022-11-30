package cache_test

import (
	"testing"
	"time"

	"github.com/jeromelesaux/cache"
)

func TestDnsCache(t *testing.T) {
	dns := cache.New()
	dns.Set("apple.com", "17.253.144.10")
	ip, exists := dns.Get("apple.com")
	if !exists {
		t.Error("apple.com was not found")
	}
	if ip == nil {
		t.Error("dns[apple.com] is nil")
	}
	if ip != "17.253.144.10" {
		t.Error("dns[apple.com] != 17.253.144.10")
	}
}

func TestRemove(t *testing.T) {
	fruits := cache.New()
	fruits.Set("Apple", 1.39)
	applePrice, exists := fruits.Get("Apple")
	if !exists {
		t.Error("Apple price was not set")
	}
	if applePrice == nil {
		t.Error("Apple price is nil")
	}
	if applePrice != 1.39 {
		t.Error("Apple price expected to be 1.39")
	}
	fruits.Remove("Apple")
	applePrice, exists = fruits.Get("Apple")
	if exists {
		t.Error("Apple price was not removed")
	}
	if applePrice != nil {
		t.Error("Apple price is not nil after removal")
	}
}

func TestThreading(t *testing.T) {
	bloom := cache.New()
	bloom.Set("Granniwinkle", 4.8)
	go func() {
		for {
			_, _ = bloom.Get("Granniwinkle")
		}
	}()
	time.Sleep(1 * time.Second)
	go func() {
		for {
			bloom.Set("Maude", 4.6)
		}
	}()
	time.Sleep(5 * time.Second)
}

func Benchmark10(b *testing.B) {
	c := cache.New()
	for i := 0; i < b.N; i++ {
		c.Set("value", i)
		v, ok := c.Get("value")
		if ok {
			if v != i {
				b.Logf("expected value %d and gets %d\n", i, v)
			}
		}
	}
}
