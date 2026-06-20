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

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db lib.Database) repository.PasswordResetRepository {
	return &passwordResetRepository{db: db.DB}
}

func (r *passwordResetRepository) Create(ctx context.Context, pr *model.PasswordReset) error {
	return r.db.WithContext(ctx).Create(pr).Error
}

func (r *passwordResetRepository) GetByToken(ctx context.Context, token string) (*model.PasswordReset, error) {
	pr := &model.PasswordReset{}
	result := r.db.WithContext(ctx).
		Where("token = ? AND is_used = false AND expires_at > ?", token, time.Now()).
		First(pr)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return pr, result.Error
}

func (r *passwordResetRepository) CountByUserSince(ctx context.Context, userID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.PasswordReset{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Count(&count).Error
	return count, err
}

func (r *passwordResetRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.PasswordReset{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

func (r *passwordResetRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ? OR is_used = true", time.Now()).
		Delete(&model.PasswordReset{}).Error
}
