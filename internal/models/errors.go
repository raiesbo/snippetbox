package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// Add a new ErrInvalidCredentials error. We'll use this later if a user
	// tried to login with an incorrect email address or password.
	ErrInvalidCrredentials = errors.New("models: invalid credentials")

	// Add a new ErrDuplicateEmail error. We'll use this later if a user
	// tried to signup with an email address that's already in use.
	ErrDuplicateEamil = errors.New("models: duplicate email")
)
