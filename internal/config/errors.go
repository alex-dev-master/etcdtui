package config

import "errors"

var (
	// ErrProfileNameRequired is returned when profile name is empty
	ErrProfileNameRequired = errors.New("profile name is required")

	// ErrEndpointsRequired is returned when no endpoints are specified
	ErrEndpointsRequired = errors.New("at least one endpoint is required")

	// ErrProfileNotFound is returned when profile doesn't exist
	ErrProfileNotFound = errors.New("profile not found")

	// ErrConfigNotFound is returned when config file doesn't exist
	ErrConfigNotFound = errors.New("config file not found")

	// ErrNoProfiles is returned when no profiles are configured
	ErrNoProfiles = errors.New("no profiles configured")

	// ErrNoDefaultProfile is returned when no default profile is set
	ErrNoDefaultProfile = errors.New("no default profile set")
)
