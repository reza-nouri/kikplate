package plate

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/kickplate/api/events"
	"github.com/kickplate/api/model"
)

func (s *plateService) emitPlateSubmittedEvent(ctx context.Context, plate *model.Plate) {
	email, name, ok := s.resolveAccountEmailAndName(ctx, plate.OwnerID)
	if !ok {
		return
	}

	s.emitter.Emit(events.PlateSubmitted, events.PlateSubmittedPayload{
		Email:     email,
		Name:      name,
		PlateName: plate.Name,
	})
}

func (s *plateService) emitPlateVerifiedEvent(ctx context.Context, plate *model.Plate) {
	email, name, ok := s.resolveAccountEmailAndName(ctx, plate.OwnerID)
	if !ok {
		return
	}

	s.emitter.Emit(events.PlateVerified, events.PlateVerifiedPayload{
		Email:     email,
		Name:      name,
		PlateName: plate.Name,
	})
}

func (s *plateService) emitPlateRatedEvent(ctx context.Context, plate *model.Plate, raterAccountID uuid.UUID, rating int16) {
	email, _, ok := s.resolveAccountEmailAndName(ctx, plate.OwnerID)
	if !ok {
		return
	}

	ratedBy := s.resolveDisplayName(ctx, raterAccountID)
	s.emitter.Emit(events.PlateRated, events.PlateRatedPayload{
		Email:     email,
		PlateName: plate.Name,
		RatedBy:   ratedBy,
		Rating:    rating,
	})
}

func (s *plateService) resolveAccountEmailAndName(ctx context.Context, accountID uuid.UUID) (string, string, bool) {
	account, err := s.accounts.GetByID(ctx, accountID)
	if err != nil || account == nil || account.UserID == nil {
		return "", "", false
	}

	user, err := s.users.GetByID(ctx, *account.UserID)
	if err != nil || user == nil {
		return "", "", false
	}

	email := strings.TrimSpace(user.Email)
	if email == "" {
		return "", "", false
	}

	name := strings.TrimSpace(user.Username)
	if name == "" && account.DisplayName != nil {
		name = strings.TrimSpace(*account.DisplayName)
	}
	if name == "" {
		name = "there"
	}

	return email, name, true
}

func (s *plateService) resolveDisplayName(ctx context.Context, accountID uuid.UUID) string {
	account, err := s.accounts.GetByID(ctx, accountID)
	if err == nil && account != nil {
		if account.DisplayName != nil {
			if v := strings.TrimSpace(*account.DisplayName); v != "" {
				return v
			}
		}
		if account.UserID != nil {
			user, userErr := s.users.GetByID(ctx, *account.UserID)
			if userErr == nil && user != nil {
				if v := strings.TrimSpace(user.Username); v != "" {
					return v
				}
			}
		}
	}

	return "Someone"
}
