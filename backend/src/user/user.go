package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID             uuid.UUID
	Email            string
	PasswordHash     string
	FirstName        string
	LastName         string
	EmailVerified    bool
	IdentityVerified bool
}

func (u *User) SetPassword(p string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)

	return nil
}

func (u *User) Login(p string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(p)); err != nil {
		return false
	}

	return true
}

type BetUser struct {
	User
	IdentityVerified bool
}

func (bu BetUser) CreateVerificationRequest(id, portrait string) (IdentityVerification, error) {
	if bu.IdentityVerified {
		return IdentityVerification{}, errors.New("user already verified")
	}

	return IdentityVerification{
		UUID:                uuid.New(),
		UserUUID:            bu.UUID,
		Status:              IdentityVerificationStatusPending,
		IDPhotoBase64:       id,
		PortraitPhotoBase64: portrait,
		CreatedAt:           time.Now(),
	}, nil
}

func (bu BetUser) CanBet() bool {
	return bu.IdentityVerified
}

type IdentityVerificationStatus string

const (
	IdentityVerificationStatusPending  IdentityVerificationStatus = "pending"
	IdentityVerificationStatusRejected IdentityVerificationStatus = "rejected"
	IdentityVerificationStatusAccepted IdentityVerificationStatus = "accepted"
)

func (s IdentityVerificationStatus) Finalized() bool {
	return s != IdentityVerificationStatusPending
}

type IdentityVerification struct {
	UUID                uuid.UUID
	UserUUID            uuid.UUID
	Status              IdentityVerificationStatus
	IDPhotoBase64       string
	PortraitPhotoBase64 string
	RespondedAt         time.Time
	CreatedAt           time.Time
}

func (v *IdentityVerification) Reject() error {
	if v.Status.Finalized() {
		return errors.New("verification already finalized")
	}

	v.Status = IdentityVerificationStatusRejected
	v.RespondedAt = time.Now()

	return nil
}

func VerifyBetUserIdentity(u *BetUser, ver *IdentityVerification) error {
	if u.IdentityVerified {
		return errors.New("user already verified")
	}

	if ver.Status.Finalized() {
		return errors.New("verification already finalized")
	}

	u.IdentityVerified = true
	ver.Status = IdentityVerificationStatusAccepted
	ver.RespondedAt = time.Now()

	return nil
}

func VerifyUserEmail(u *User, ev *EmailVerification) error {
	if u.EmailVerified {
		return errors.New("user already verified")
	}

	if ev.Activated {
		return errors.New("token already activated")
	}

	u.EmailVerified = true
	ev.Activated = true

	return nil
}

type EmailVerification struct {
	UserUUID  uuid.UUID
	Token     string
	Activated bool
}
