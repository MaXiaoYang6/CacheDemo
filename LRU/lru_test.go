package LRU

import (
	"reflect"
	"testing"
	"unsafe"
)

type String string

func (s String) Size() int {
	return int(unsafe.Sizeof(s))
}

func TestGet(t *testing.T) {
	lru := New(0, nil)
	lru.Add("key1", String("1234"))
	if _, ok := lru.Get("key1"); ok {
		t.Fatalf("cache hit key1 == 1234 failed!")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed!")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	size := unsafe.Sizeof(k1) + unsafe.Sizeof(k2) + unsafe.Sizeof(v1) + unsafe.Sizeof(v2)
	lru := New(int64(size), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("k1"); ok || lru.nBytes != 64 {
		t.Fatalf("RemoveOldest key1 failed!")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callBack := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(64), callBack)
	lru.Add("k1", String("v1"))
	lru.Add("k2", String("v2"))
	lru.Add("k3", String("v3"))
	lru.Add("k4", String("v4"))

	expect := []string{"k1", "k2"}

	if !reflect.DeepEqual(keys, expect) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s, but keys are %s", expect, keys)
	}
}
