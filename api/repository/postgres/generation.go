package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"gorm.io/gorm"
)

type generationRepository struct {
	db *gorm.DB
}

func NewGenerationRepository(db lib.Database) repository.GenerationRepository {
	return &generationRepository{db: db.DB}
}

func (r *generationRepository) Create(ctx context.Context, gen *model.Generation) error {
	return r.db.WithContext(ctx).Create(gen).Error
}

func (r *generationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Generation, error) {
	var gen model.Generation
	result := r.db.WithContext(ctx).First(&gen, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &gen, result.Error
}

func (r *generationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.GenerationStatus, errMsg *string) error {
	updates := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	if errMsg != nil {
		updates["error"] = *errMsg
	}
	return r.db.WithContext(ctx).Model(&model.Generation{}).Where("id = ?", id).Updates(updates).Error
}

func (r *generationRepository) ListByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*model.Generation, int, error) {
	var gens []*model.Generation
	var total int64

	q := r.db.WithContext(ctx).Model(&model.Generation{}).Where("account_id = ?", accountID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&gens).Error; err != nil {
		return nil, 0, err
	}

	return gens, int(total), nil
}
