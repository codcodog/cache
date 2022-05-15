package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	c := NewCache()
	c.Set("hello", "world", 5*time.Second)
	c.Set("int", 5, 5*time.Second)
	c.Set("float", float64(32.03), 5*time.Second)

	h, found := c.Get("hello")
	if !found {
		t.Error("test set failed: hello not found")
	}
	if !reflect.DeepEqual(h, "world") {
		t.Error("test set failed: hello not equal")
	}

	i, found := c.Get("int")
	if !found {
		t.Error("test set failed: int not found")
	}
	if !reflect.DeepEqual(i, 5) {
		t.Error("test set failed: int not equal")
	}

	f, found := c.Get("float")
	if !found {
		t.Error("test set failed: float not found")
	}
	if !reflect.DeepEqual(f, float64(32.03)) {
		t.Error("test set failed: float not equal")
	}
}

func TestExpire(t *testing.T) {
	c := NewCache()
	c.Set("hello", "world", 2*time.Second)

	h, found := c.Get("hello")
	if !found {
		t.Error("test expire failed: hello not found")
	}
	if !reflect.DeepEqual(h, "world") {
		t.Error("test expire failed: hello not equal")
	}

	time.Sleep(3 * time.Second)
	_, found = c.Get("hello")
	if found {
		t.Error("test expire failed")
	}
}

func TestDel(t *testing.T) {
	c := NewCache()
	c.Set("hello", "world", 5*time.Second)

	h, found := c.Get("hello")
	if !found {
		t.Error("test del failed: hello not found")
	}
	if !reflect.DeepEqual(h, "world") {
		t.Error("test del failed: hello not equal")
	}

	c.Del("hello")
	_, found = c.Get("hello")
	if found {
		t.Error("test del failed.")
	}
}

func TestExists(t *testing.T) {
	c := NewCache()
	c.Set("hello", "world", 5*time.Second)
	c.Set("world", "nihao", 5*time.Second)

	if !c.Exists("hello") {
		t.Error("test exists failed, key: hello")
	}
	if !c.Exists("world") {
		t.Error("test exists failed, key: world")
	}
	if c.Exists("nihao") {
		t.Error("test exists failed, key: nihao, not exists")
	}
}

func TestFlush(t *testing.T) {
	c := NewCache()
	c.Set("hello", "world", 5*time.Second)
	c.Set("int", 5, 5*time.Second)
	c.Set("float", float64(32.03), 5*time.Second)

	if !c.Exists("hello") && !c.Exists("int") && !c.Exists("float") {
		t.Error("test flush failed.")
	}

	c.Flush()

	if c.Exists("hello") || c.Exists("int") || c.Exists("float") {
		t.Error("test flush failed.")
	}
}

func TestKeys(t *testing.T) {
	c := NewCache()
	c.Set("hello", "world", 5*time.Second)
	c.Set("int", 5, 5*time.Second)
	c.Set("float", float64(32.03), 5*time.Second)

	if c.Keys() != int64(3) {
		t.Error("test keys failed.")
	}
}
