package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mx            *sync.Mutex
	capacity      int
	queue         List
	items         map[Key]*ListItem
	reversedItems map[*ListItem]Key
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mx:            &sync.Mutex{},
		capacity:      capacity,
		queue:         NewList(),
		items:         make(map[Key]*ListItem, capacity),
		reversedItems: make(map[*ListItem]Key, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mx.Lock()
	defer c.mx.Unlock()
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = value
		return true
	}

	if c.capacity == len(c.items) {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)

		key := c.reversedItems[lastItem]
		delete(c.items, key)

		// for k, v := range c.items {
		// 	if v == lastItem {

		// 		break
		// 	}
		// }
	}

	item := c.queue.PushFront(value)
	c.items[key] = item
	c.reversedItems[item] = key
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	item, ok := c.items[key]
	if !ok {
		return nil, ok
	}
	c.queue.MoveToFront(item)
	return item.Value, ok
}

func (c *lruCache) Clear() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
