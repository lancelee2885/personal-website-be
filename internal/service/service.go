package service

import (
	"net/http"

	"github.com/lancelee2885/personal-website-be/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Service is a struct that contains the necessary fields for a service.
type Service struct {
	log     zerolog.Logger
	storage storage.EntityStorage
}

type ServiceConfig struct {
	Logger  zerolog.Logger
	Storage storage.EntityStorage
}

func NewService(config *ServiceConfig) *Service {
	return &Service{
		log:     config.Logger,
		storage: config.Storage,
	}
}

// CreateEntity creates a new entity.
func (s *Service) CreateEntity(c *gin.Context) {
	var requestBody struct {
		TableName string         `json:"tableName"`
		Entity    storage.Entity `json:"entity"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		s.log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity, err := s.storage.Create(c, requestBody.TableName, &requestBody.Entity)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to create entity")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity"})
		return
	}

	s.log.Info().Msgf("Entity created successfully: %+v", entity)
	c.JSON(http.StatusCreated, entity)
}

// GetEntity retrieves an entity by ID.
func (s *Service) GetEntity(c *gin.Context) {
	id := c.Param("id")
	tableName := c.Query("type")

	if id == "" {
		s.log.Error().Msg("Entity ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entity ID is required"})
		return
	}

	if tableName == "" {
		s.log.Error().Msg("Table name is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table name is required"})
		return
	}

	entity, err := s.storage.GetByID(c, tableName, id)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to get entity with ID: %s", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get entity"})
		return
	}

	s.log.Info().Msgf("Entity retrieved successfully: %+v", entity)
	c.JSON(http.StatusOK, entity)
}

// UpdateEntity updates an existing entity.
func (s *Service) UpdateEntity(c *gin.Context) {
	var requestBody struct {
		TableName string         `json:"tableName"`
		Entity    storage.Entity `json:"entity"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		s.log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity, err := s.storage.Update(c, requestBody.TableName, &requestBody.Entity)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to update entity")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entity"})
		return
	}

	s.log.Info().Msgf("Entity updated successfully: %+v", entity)
	c.JSON(http.StatusOK, entity)
}

// DeleteEntity deletes an entity by ID.
func (s *Service) DeleteEntity(c *gin.Context) {
	var requestBody struct {
		TableName string `json:"tableName"`
		ID        string `json:"id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		s.log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deleted, err := s.storage.Delete(c, requestBody.TableName, requestBody.ID)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to delete entity with ID: %s", requestBody.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entity"})
		return
	}

	if !deleted {
		s.log.Error().Msgf("Failed to delete entity with ID: %s", requestBody.ID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
		return
	}

	s.log.Info().Msgf("Entity deleted successfully: %s", requestBody.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Entity deleted successfully"})
}

// ArchiveEntity archives an entity by ID.
func (s *Service) ArchiveEntity(c *gin.Context) {
	var requestBody struct {
		TableName string `json:"tableName"`
		ID        string `json:"id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		s.log.Error().Err(err).Msg("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	archived, err := s.storage.Archive(c, requestBody.TableName, requestBody.ID)
	if err != nil {
		s.log.Error().Err(err).Msgf("Failed to archive entity with ID: %s", requestBody.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive entity"})
		return
	}

	if !archived {
		s.log.Error().Msgf("Failed to archive entity with ID: %s", requestBody.ID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
		return
	}

	s.log.Info().Msgf("Entity archived successfully: %s", requestBody.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Entity archived successfully"})
}

// ListEntities lists all entities.
func (s *Service) ListEntities(c *gin.Context) {
	tableName := c.Query("type")

	if tableName == "" {
		s.log.Error().Msg("Table name is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table name is required"})
		return
	}

	entities, err := s.storage.List(c, tableName)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to list entities")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list entities"})
		return
	}

	s.log.Info().Msgf("Entities listed successfully: %+v", entities)
	c.JSON(http.StatusOK, entities)
}
