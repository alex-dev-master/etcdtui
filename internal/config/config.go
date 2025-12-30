package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	// DefaultConfigDir is the default config directory name
	DefaultConfigDir = "etcdtui"

	// DefaultConfigFile is the default config file name
	DefaultConfigFile = "config.yaml"

	// DefaultFileMode is the file permission for config file (owner read/write only)
	DefaultFileMode = 0600

	// DefaultDirMode is the directory permission for config directory
	DefaultDirMode = 0700
)

// Config represents the application configuration
type Config struct {
	// Profiles is a list of connection profiles
	Profiles []*Profile `yaml:"profiles" mapstructure:"profiles"`

	// ActiveProfile is the name of the currently active profile
	ActiveProfile string `yaml:"active_profile,omitempty" mapstructure:"active_profile"`
}

// Manager handles configuration loading and saving
type Manager struct {
	v          *viper.Viper
	configPath string
	config     *Config
}

// NewManager creates a new config manager
func NewManager() *Manager {
	return &Manager{
		v:      viper.New(),
		config: &Config{},
	}
}

// GetConfigDir returns the config directory path
func GetConfigDir() (string, error) {
	// Try XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, DefaultConfigDir), nil
	}

	// Fall back to ~/.config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".config", DefaultConfigDir), nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, DefaultConfigFile), nil
}

// Load loads configuration from the default location
func (m *Manager) Load() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	return m.LoadFromPath(configPath)
}

// LoadFromPath loads configuration from a specific path
func (m *Manager) LoadFromPath(path string) error {
	m.configPath = path

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return empty config, not an error
		m.config = &Config{}
		return nil
	}

	m.v.SetConfigFile(path)
	m.v.SetConfigType("yaml")

	if err := m.v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := m.v.Unmarshal(m.config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// Save saves configuration to the config file
func (m *Manager) Save() error {
	if m.configPath == "" {
		var err error
		m.configPath, err = GetConfigPath()
		if err != nil {
			return err
		}
	}

	return m.SaveToPath(m.configPath)
}

// SaveToPath saves configuration to a specific path
func (m *Manager) SaveToPath(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, DefaultDirMode); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set all values in viper
	m.v.Set("profiles", m.config.Profiles)
	m.v.Set("active_profile", m.config.ActiveProfile)

	// Write config
	if err := m.v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Set file permissions
	if err := os.Chmod(path, DefaultFileMode); err != nil {
		return fmt.Errorf("failed to set config file permissions: %w", err)
	}

	return nil
}

// GetProfiles returns all profiles
func (m *Manager) GetProfiles() []*Profile {
	return m.config.Profiles
}

// GetProfile returns a profile by name
func (m *Manager) GetProfile(name string) (*Profile, error) {
	for _, p := range m.config.Profiles {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, ErrProfileNotFound
}

// GetDefaultProfile returns the default profile
func (m *Manager) GetDefaultProfile() (*Profile, error) {
	// First, try active profile
	if m.config.ActiveProfile != "" {
		if p, err := m.GetProfile(m.config.ActiveProfile); err == nil {
			return p, nil
		}
	}

	// Then, try profile marked as default
	for _, p := range m.config.Profiles {
		if p.Default {
			return p, nil
		}
	}

	// Finally, return first profile if exists
	if len(m.config.Profiles) > 0 {
		return m.config.Profiles[0], nil
	}

	return nil, ErrNoDefaultProfile
}

// AddProfile adds a new profile
func (m *Manager) AddProfile(profile *Profile) error {
	if err := profile.Validate(); err != nil {
		return err
	}

	// Check if profile with this name already exists
	for i, p := range m.config.Profiles {
		if p.Name == profile.Name {
			// Update existing profile
			m.config.Profiles[i] = profile
			return nil
		}
	}

	// If this is the first profile or marked as default, make it default
	if len(m.config.Profiles) == 0 || profile.Default {
		// Clear other defaults if this one is default
		if profile.Default {
			for _, p := range m.config.Profiles {
				p.Default = false
			}
		}
	}

	m.config.Profiles = append(m.config.Profiles, profile)
	return nil
}

// DeleteProfile removes a profile by name
func (m *Manager) DeleteProfile(name string) error {
	for i, p := range m.config.Profiles {
		if p.Name == name {
			m.config.Profiles = append(m.config.Profiles[:i], m.config.Profiles[i+1:]...)
			return nil
		}
	}
	return ErrProfileNotFound
}

// SetActiveProfile sets the active profile by name
func (m *Manager) SetActiveProfile(name string) error {
	if _, err := m.GetProfile(name); err != nil {
		return err
	}
	m.config.ActiveProfile = name
	return nil
}

// GetActiveProfileName returns the active profile name
func (m *Manager) GetActiveProfileName() string {
	return m.config.ActiveProfile
}

// SetDefaultProfile sets a profile as the default
func (m *Manager) SetDefaultProfile(name string) error {
	found := false
	for _, p := range m.config.Profiles {
		if p.Name == name {
			p.Default = true
			found = true
		} else {
			p.Default = false
		}
	}

	if !found {
		return ErrProfileNotFound
	}
	return nil
}

// HasProfiles returns true if any profiles are configured
func (m *Manager) HasProfiles() bool {
	return len(m.config.Profiles) > 0
}

// CreateDefaultConfig creates a default configuration with a local profile
func (m *Manager) CreateDefaultConfig() {
	m.config = &Config{
		Profiles: []*Profile{
			{
				Name:      "local",
				Endpoints: []string{"localhost:2379"},
				Default:   true,
			},
		},
	}
}
