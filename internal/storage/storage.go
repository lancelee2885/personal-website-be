package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PostgresStore struct {
	log zerolog.Logger
	db  *goqu.Database
}

var _ EntityStorage = (*PostgresStore)(nil)

func NewPostgresStore(log zerolog.Logger, db *sql.DB) (*PostgresStore, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	return &PostgresStore{
		log: log,
		db:  goqu.New("postgres", db),
	}, nil
}

func (p *PostgresStore) Create(ctx context.Context, tableName string, entity *Entity) (*Entity, error) {

	uuid := uuid.New().String()
	entity.ID = uuid
	now := time.Now()
	entity.CreatedAt = now
	entity.UpdatedAt = now

	entities := make([]*Entity, 0)

	err := p.db.Insert(tableName).Returning("*").Rows(entity).Executor().ScanStructsContext(ctx, &entities)

	if err != nil {
		p.log.Error().Err(err).Msg("Failed to create entity")
		return nil, err
	}

	if len(entities) != 1 {
		p.log.Error().Err(err).Msg("Failed to create entity")
		return nil, fmt.Errorf("failed to create entity, more than one entity created")
	}

	return entities[0], err
}

func (p *PostgresStore) GetByID(ctx context.Context, tableName string, id string) (*Entity, error) {

	entities := make([]*Entity, 0)

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	err = p.db.From(tableName).Where(goqu.Ex{"id": uuid, "archived": false}).ScanStructsContext(ctx, &entities)
	if err != nil {
		return nil, err
	}
	if len(entities) != 1 {
		return nil, fmt.Errorf("entity with ID %s not found", id)
	}
	return entities[0], nil
}

func (p *PostgresStore) Update(ctx context.Context, tableName string, entity *Entity) (*Entity, error) {

	entities := make([]*Entity, 0)

	entity.UpdatedAt = time.Now()

	err := p.db.Update(tableName).Set(entity).Where(goqu.Ex{"id": entity.ID}).Returning("*").Executor().ScanStructsContext(ctx, &entities)

	if err != nil {
		p.log.Error().Err(err).Msg("Failed to update entity")
		return nil, err
	}

	if len(entities) != 1 {
		p.log.Error().Err(err).Msg("Failed to update entity")
		return nil, fmt.Errorf("failed to update entity, more than one entity updated")
	}

	return entities[0], err
}

func (p *PostgresStore) Delete(ctx context.Context, tableName string, id string) (bool, error) {
	_, err := p.db.Delete(tableName).Where(goqu.Ex{"id": id}).Executor().ExecContext(ctx)
	if err != nil {
		log.Printf("Failed to delete entity: %v", err)
		return false, err
	}
	return true, nil
}

func (p *PostgresStore) Archive(ctx context.Context, tableName string, id string) (bool, error) {

	entity, err := p.GetByID(ctx, tableName, id)
	if err != nil {
		return false, err
	}

	now := time.Now()

	entity.UpdatedAt = now

	_, err = p.db.Update(tableName).Set(entity).Where(goqu.Ex{"id": id}).Executor().ExecContext(ctx)
	if err != nil {
		log.Printf("Failed to archive entity: %v", err)
		return false, err
	}

	return true, nil
}

func (p *PostgresStore) List(ctx context.Context, tableName string) ([]*Entity, error) {

	entities := make([]*Entity, 0)

	err := p.db.From(tableName).Where(goqu.Ex{"archived": false}).ScanStructsContext(ctx, &entities)
	if err != nil {
		log.Printf("Failed to list entities: %v", err)
		return nil, err
	}

	return entities, nil
}
