package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Errorf("unexpected error: %v", err)

	}
	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("unexpected error validating: %v", err)
	}
	if gotID != userID {
		t.Errorf("expected %v, got %v", userID, gotID)

	}
}
func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret"
	token, _ := MakeJWT(userID, secret, -time.Hour)
	_, err := ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expected an error for expired token but got none")
	}
}
func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "secret"

	token, _ := MakeJWT(userID, secret, time.Hour)
	_, err := ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Errorf("expected an error but got none")
	}
}
