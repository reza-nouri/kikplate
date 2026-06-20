package postgres

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewAccountRepository),
	fx.Provide(NewEmailVerificationRepository),
	fx.Provide(NewPasswordResetRepository),
	fx.Provide(NewPlateRepository),
	fx.Provide(NewPlateTagRepository),
	fx.Provide(NewPlateMemberRepository),
	fx.Provide(NewPlateReviewRepository),
	fx.Provide(NewPlateBadgeRepository),
	fx.Provide(NewBadgeRepository),
	fx.Provide(NewOrganizationRepository),
	fx.Provide(NewGenerationRepository),
)
