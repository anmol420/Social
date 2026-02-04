package cache

import (
	"context"

	"github.com/anmol420/Social/internal/store"
)

func NewMockCacheStore() Storage {
	return Storage{
		Users: &MockUserCacheStore{},
	}
}

type MockUserCacheStore struct{}

func (m *MockUserCacheStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	return &store.User{}, nil
}

func (m *MockUserCacheStore) Set(ctx context.Context, user *store.User) error {
	return nil
}

func (m *MockUserCacheStore) Delete(ctx context.Context, userID int64) {}
