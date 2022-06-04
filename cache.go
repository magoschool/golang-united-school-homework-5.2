package cache

import (
	"sort"
	"time"
)

type cacheItem struct {
	itemValue    string
	itemDeadline *time.Time
}
type Cache struct {
	items map[string]cacheItem
}

func (c cacheItem) isAlive(aTestTime time.Time) bool {
	return c.itemDeadline == nil || c.itemDeadline.After(aTestTime)
}

func NewCache() Cache {
	return Cache{items: map[string]cacheItem{}}
}

func (c Cache) Get(key string) (string, bool) {
	lItem, lOk := c.items[key]

	if !lOk || !lItem.isAlive(time.Now()) {
		return "", false
	}

	return lItem.itemValue, true
}

func (c *Cache) Put(key, value string) {
	c.items[key] = cacheItem{itemValue: value, itemDeadline: nil}
}

func (c Cache) Keys() []string {
	lNow := time.Now()
	lKeys := make([]string, 0, len(c.items))
	for lKey := range c.items {
		if c.items[lKey].isAlive(lNow) {
			lKeys = append(lKeys, lKey)
		}
	}

	// return sorted keys
	sort.Strings(lKeys)
	return lKeys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.items[key] = cacheItem{itemValue: value, itemDeadline: &deadline}
}
