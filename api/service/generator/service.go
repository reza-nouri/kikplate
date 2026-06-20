package generator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
)

type GenerateInput struct {
	PlateSlug string
	Values    map[string]any
}

type GenerateResult struct {
	GenerationID uuid.UUID
	ZipBytes     []byte
}

type GeneratorService interface {
	Generate(ctx context.Context, accountID *uuid.UUID, input GenerateInput) (*GenerateResult, error)
	GetSchema(ctx context.Context, slug string) (*PlateYAML, error)
}

type generatorService struct {
	plates      repository.PlateRepository
	generations repository.GenerationRepository
	logger      lib.Logger
}

func NewGeneratorService(
	plates repository.PlateRepository,
	generations repository.GenerationRepository,
	logger lib.Logger,
) GeneratorService {
	return &generatorService{
		plates:      plates,
		generations: generations,
		logger:      logger,
	}
}

func (s *generatorService) Generate(ctx context.Context, accountID *uuid.UUID, input GenerateInput) (*GenerateResult, error) {
	plate, err := s.plates.GetBySlug(ctx, input.PlateSlug)
	if err != nil {
		return nil, err
	}
	if plate == nil {
		return nil, ErrPlateNotFound
	}
	if plate.Status != model.PlateStatusApproved {
		return nil, ErrPlateNotApproved
	}
	if plate.RepoURL == nil || *plate.RepoURL == "" {
		return nil, ErrNoRepoURL
	}

	branch := "main"
	if plate.Branch != nil && *plate.Branch != "" {
		branch = *plate.Branch
	}

	py, err := fetchPlateYAML(*plate.RepoURL, branch)
	if err != nil {
		return nil, err
	}

	if input.Values == nil {
		input.Values = map[string]any{}
	}

	applyDefaults(py, input.Values)

	valJSON, err := json.Marshal(input.Values)
	if err != nil {
		return nil, ErrInvalidInput
	}

	genID := uuid.New()
	gen := &model.Generation{
		ID:        genID,
		PlateID:   plate.ID,
		AccountID: accountID,
		Status:    model.GenerationStatusPending,
		Values:    valJSON,
	}

	if err := s.generations.Create(ctx, gen); err != nil {
		return nil, err
	}

	zipBytes, renderErr := renderProject(py, input.Values)

	now := time.Now()
	if renderErr != nil {
		errMsg := renderErr.Error()
		if dbErr := s.generations.UpdateStatus(ctx, genID, model.GenerationStatusFailed, &errMsg); dbErr != nil {
			s.logger.Warnf("failed to mark generation %s as failed: %v", genID, dbErr)
		}
		return nil, renderErr
	}

	_ = now
	if dbErr := s.generations.UpdateStatus(ctx, genID, model.GenerationStatusComplete, nil); dbErr != nil {
		s.logger.Warnf("failed to mark generation %s as complete: %v", genID, dbErr)
	}

	return &GenerateResult{
		GenerationID: genID,
		ZipBytes:     zipBytes,
	}, nil
}

func applyDefaults(py *PlateYAML, values map[string]any) {
	for key, field := range py.Schema {
		if _, ok := values[key]; !ok && field.Default != nil {
			values[key] = field.Default
		}
	}
}

func (s *generatorService) GetSchema(ctx context.Context, slug string) (*PlateYAML, error) {
	plate, err := s.plates.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if plate == nil {
		return nil, ErrPlateNotFound
	}
	if plate.Status != model.PlateStatusApproved {
		return nil, ErrPlateNotApproved
	}
	if plate.RepoURL == nil || *plate.RepoURL == "" {
		return nil, ErrNoRepoURL
	}

	branch := "main"
	if plate.Branch != nil && *plate.Branch != "" {
		branch = *plate.Branch
	}

	return fetchPlateYAML(*plate.RepoURL, branch)
}
