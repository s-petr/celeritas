package cache

import (
	"testing"
)

func TestBadgerCache_Has(t *testing.T) {
	err := testBadgerCache.Forget("test")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("test")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("test found in cache but it should not be there")
	}

	err = testBadgerCache.Set("test", "hello world")
	if err != nil {
		t.Error(err)
	}

	inCache, err = testBadgerCache.Has("test")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("\"test\" not found in cache but it should be there")
	}

	err = testBadgerCache.Forget("test")
	if err != nil {
		t.Error(err)
	}
}

func TestBadgerCache_Get(t *testing.T) {
	err := testBadgerCache.Set("test", "hello world")
	if err != nil {
		t.Error(err)
	}

	x, err := testBadgerCache.Get("test")
	if err != nil {
		t.Error(err)
	}

	if x != "hello world" {
		t.Error("did not receive correct value from cache")
	}
}

func TestBadgerCache_Forget(t *testing.T) {
	err := testBadgerCache.Set("one", "two")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Forget("one")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("one")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("\"one\" found in cache but it should not be there")
	}
}

func TestBadgerCache_Empty(t *testing.T) {
	err := testBadgerCache.Set("one", "two")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("one")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("\"one\" found in cache but it should not be there")
	}
}

func TestBadgerCache_EmptyByMatch(t *testing.T) {
	err := testBadgerCache.Set("one", "two")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Set("three", "four")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.EmptyByMatch("o")
	if err != nil {
		t.Error(err)
	}

	err = testBadgerCache.Forget("one")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testBadgerCache.Has("one")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("\"one\" found in cache but it should not be there")
	}

	inCache, err = testBadgerCache.Has("three")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("\"three\" not found in cache but it should be there")
	}
}
