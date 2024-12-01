package cache

import "sync"

// local storage for information that the user’s command is already being processed by the service
type BotCache struct {
	m       sync.Mutex
	storage map[int64]bool
}

func New() *BotCache {
	return &BotCache{
		m:       sync.Mutex{},
		storage: make(map[int64]bool),
	}
}

func (b *BotCache) Add(userID int64) {
	b.m.Lock()
	b.storage[userID] = true
	b.m.Unlock()
}

func (b *BotCache) IsExist(userID int64) bool {
	_, exist := b.storage[userID]
	return exist
}

func (b *BotCache) Delete(userID int64) {
	b.m.Lock()
	delete(b.storage, userID)
	b.m.Unlock()
}
