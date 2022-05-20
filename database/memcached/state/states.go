package state

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type State struct {
	client *memcache.Client
}

func New(client *memcache.Client) State {
	return State{client: client}
}

const (
	expire = 60 * 10 // 10 minutes
)

// Store stores state, error can be only internal
func (m State) Store(state string) error {

	return m.client.Set(&memcache.Item{
		Key:        state,
		Value:      nil,
		Flags:      0,
		Expiration: expire, // 1 hour
	})
}

// Lookup returns true if state existed, error can be only internal
func (m State) Lookup(state string) (bool, error) {
	err := m.client.Delete(state)

	if err != nil {
		if err == memcache.ErrCacheMiss { // Not Found
			err = nil
		}

		return false, err
	}

	return true, nil
}
