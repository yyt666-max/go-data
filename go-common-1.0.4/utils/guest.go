package utils

import (
	"context"
	"fmt"

	"github.com/eolinker/eosc/env"
	"github.com/google/uuid"
)

const (
	EnvGuestMode         = "GUEST_MODE"
	EnvGuestID           = "GUEST_ID"
	EnvGuestUser         = "GUEST_USER"
	EnvGuestPassword     = "GUEST_PASSWORD"
	DefaultGuestUser     = "guest"
	DefaultGuestPassword = "12345678"
)

var (
	defaultUserId = uuid.New().String()
	userId        string
	guestAllow    bool
	guestUser     string
	guestPassword string
)

func init() {
	guestMode, ok := env.GetEnv(EnvGuestMode)
	if ok {
		guestAllow = guestMode == "allow"
	}
	guestUser, ok = env.GetEnv(EnvGuestUser)
	if !ok {
		guestUser = DefaultGuestUser
	}
	guestPassword, ok = env.GetEnv(EnvGuestPassword)
	if !ok {
		guestPassword = DefaultGuestPassword
	}
	userId, ok = env.GetEnv(EnvGuestID)
	if !ok {
		userId = defaultUserId
	}
}

func GuestAllow() bool {
	return guestAllow
}

func IsGuest(ctx context.Context) bool {
	uid := UserId(ctx)
	return uid == userId
}

func GuestUser() string {
	return guestUser
}

func GuestLogin(ctx context.Context, user, password string) (string, error) {
	if user == guestUser && password == guestPassword {
		return userId, nil
	}
	return "", fmt.Errorf("invalid user or password")
}
