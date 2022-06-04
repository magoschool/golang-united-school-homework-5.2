package cache

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func (c Cache) String() string {
	lNow := time.Now()
	lKeys := c.Keys()

	lNeedSeparator := false
	lBuffer := strings.Builder{}
	lBuffer.WriteString("[")
	for _, lKey := range lKeys {
		lCacheItem := c.items[lKey]

		if lCacheItem.isAlive(lNow) {
			if lNeedSeparator {
				lBuffer.WriteString(", ")
			}

			lBuffer.WriteString(fmt.Sprintf("%v:%v", lKey, lCacheItem.itemValue))
			lNeedSeparator = true
		}
	}

	lBuffer.WriteString("]")
	return lBuffer.String()
}

func testCache(t *testing.T, aCache Cache, aValue string) {
	lExpected := fmt.Sprintf("%v", aCache)

	if lExpected != aValue {
		t.Errorf("Incorrect cache: %v expected but %v found", lExpected, aValue)
	}
}

func testGetResult(t *testing.T, aValue string, aResult bool, aExpectedValue string, aExpectedResult bool) {
	if aResult != aExpectedResult {
		t.Errorf("Incorrect Get result: %v expected but %v found", aExpectedResult, aResult)
	} else if aValue != aExpectedValue {
		t.Errorf("Incorrect cache value: %v expected but %v found", aExpectedValue, aValue)
	}
}

func testKeysResult(t *testing.T, aKeys []string, aExpectedKeysText string) {
	lKeysText := fmt.Sprintf("%v", aKeys)

	if aExpectedKeysText != lKeysText {
		t.Errorf("Incorrect cache: %v expected but %v found", aExpectedKeysText, lKeysText)
	}
}

func TestNewCache(t *testing.T) {
	lCache := NewCache()
	testCache(t, lCache, "[]")
}

func TestPut(t *testing.T) {
	lCache := NewCache()

	lCache.Put("1", "a")
	testCache(t, lCache, "[1:a]")

	lCache.Put("2", "b")
	testCache(t, lCache, "[1:a, 2:b]")

	lCache.Put("1", "c")
	testCache(t, lCache, "[1:c, 2:b]")
}

func TestGet(t *testing.T) {
	lCache := NewCache()

	lCache.Put("1", "a")
	lCache.Put("2", "b")
	lCache.Put("3", "c")

	lValue, lResult := lCache.Get("2")
	testGetResult(t, lValue, lResult, "b", true)

	lValue, lResult = lCache.Get("4")
	testGetResult(t, lValue, lResult, "", false)
}

func TestKeys(t *testing.T) {
	lCache := NewCache()

	testKeysResult(t, lCache.Keys(), "[]")

	lCache.Put("3", "c")
	lCache.Put("2", "b")
	lCache.Put("1", "a")

	testKeysResult(t, lCache.Keys(), "[1 2 3]")
}

func TestPutTill(t *testing.T) {
	lCache := NewCache()
	lDeadline := time.Now().Add(2 * time.Second)

	lCache.Put("1", "a")
	lCache.PutTill("2", "b", lDeadline)
	lCache.Put("3", "c")

	lValue, lResult := lCache.Get("2")
	testGetResult(t, lValue, lResult, "b", true)
	testCache(t, lCache, "[1:a, 2:b, 3:c]")
	testKeysResult(t, lCache.Keys(), "[1 2 3]")

	time.Sleep(2 * time.Second)

	lValue, lResult = lCache.Get("2")
	testGetResult(t, lValue, lResult, "", false)
	testCache(t, lCache, "[1:a, 3:c]")
	testKeysResult(t, lCache.Keys(), "[1 3]")
}
