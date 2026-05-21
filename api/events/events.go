package events

const (
	UserRegistered            = "user.registered"
	UserVerificationRequested = "user.verification_requested"
	PasswordChanged           = "user.password_changed"
	UserLiked                 = "user.liked"
	PlateSubmitted            = "plate.submitted"
	PlateVerified             = "plate.verified"
	PlateRated                = "plate.rated"
	PlateSyncIssue            = "plate.sync_issue"
)

type UserRegisteredPayload struct {
	Email string
	Name  string
}

type UserVerificationRequestedPayload struct {
	Email     string
	Name      string
	VerifyURL string
}

type PasswordChangedPayload struct {
	Email string
}

type UserLikedPayload struct {
	Email   string
	LikedBy string
}

type PlateSubmittedPayload struct {
	Email     string
	Name      string
	PlateName string
}

type PlateVerifiedPayload struct {
	Email     string
	Name      string
	PlateName string
}

type PlateRatedPayload struct {
	Email     string
	PlateName string
	RatedBy   string
	Rating    int16
}

type PlateSyncIssuePayload struct {
	Email     string
	PlateName string
	Issue     string
}
