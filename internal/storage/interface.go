package storage

import (
	"context"
)

// Storage represents a generic storage interface.
type EntityStorage interface {
	Create(ctx context.Context, tableName string, entity *Entity) (*Entity, error)
	GetByID(ctx context.Context, tableName string, id string) (*Entity, error)
	Update(ctx context.Context, tableName string, entity *Entity) (*Entity, error)
	Delete(ctx context.Context, tableName string, id string) (bool, error)
	Archive(ctx context.Context, tableName string, id string) (bool, error)
	List(ctx context.Context, tableName string) ([]*Entity, error)
}
