package constants

import (
	"time"
)

// User roles
const (
	// RegularUser - regular user
	RegularUser = 0
	// Moderator - moderator
	Moderator = 1
	// AdminUser - admin user
	AdminUser = 2
)

const (
	// AccessTokenLifetime is a lifetime for access token
	AccessTokenLifetime = time.Hour * 2
	// RefreshTokenLifetime is a lifetime for refresh token
	RefreshTokenLifetime = time.Hour * 24 * 30
)
