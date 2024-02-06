package cache

import "testing"

func TestRedisCache_Has(t *testing.T) {
	err := testRedisCache.Forget("test")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("test")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("test found in cache but it should not be there")
	}

	err = testRedisCache.Set("test", "hello world")
	if err != nil {
		t.Error(err)
	}

	inCache, err = testRedisCache.Has("test")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("\"test\" not found in cache but it should be there")
	}
}

func TestRedisCache_Get(t *testing.T) {
	err := testRedisCache.Set("test", "hello world")
	if err != nil {
		t.Error(err)
	}

	x, err := testRedisCache.Get("test")
	if err != nil {
		t.Error(err)
	}

	if x != "hello world" {
		t.Error("did not receive correct value from cache")
	}
}

func TestRedisCache_Forget(t *testing.T) {
	err := testRedisCache.Set("one", "two")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Forget("one")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("one")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("\"one\" found in cache but it should not be there")
	}
}

func TestRedisCache_Empty(t *testing.T) {
	err := testRedisCache.Set("one", "two")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Empty()
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Forget("one")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("one")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("\"one\" found in cache but it should not be there")
	}
}

func TestRedisCache_EmptyByMatch(t *testing.T) {
	err := testRedisCache.Set("one", "two")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("three", "four")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.EmptyByMatch("o")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Forget("one")
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("one")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("\"one\" found in cache but it should not be there")
	}

	inCache, err = testRedisCache.Has("three")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("\"one\" not found in cache but it should be there")
	}
}

func TestRedisCache_EncodeDecode(t *testing.T) {
	entry := Entry{}
	entry["one"] = "two"

	bytes, err := encode(entry)
	if err != nil {
		t.Error(err)
	}

	decoded, err := decode(string(bytes))
	if err != nil {
		t.Error(err)
	}

	if decoded["one"] != "two" {
		t.Errorf("Expected decoded value to be \"two\", received %+v", decoded["one"])
	}
}
