package router

import (
	"aspen/config"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Aspen always has a global config file that is used as the single source of truth
var globalConfigFile string

// To ensure concurrent writes don't get lost, we use a RWMutex to lock the config file
// Does NOT lock changes to the globalConfigFile variable, only the contents of the file
var globalConfigLock sync.RWMutex

// SetGlobalConfigFile sets the global configuration file for the application.
// It checks if the file exists and is readable before setting it.
func SetGlobalConfigFile(file string) error {
	// Check that the file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", file)
	}

	// Check that the file is readable
	if _, err := os.ReadFile(file); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	globalConfigFile = file
	return nil
}

// ReadGlobalConfig reads the global configuration file and returns the corresponding Config struct.
func ReadGlobalConfig() (*config.Config, error) {
	globalConfigLock.RLock()
	defer globalConfigLock.RUnlock()

	return readGlobalConfigNoLock()
}

func readGlobalConfigNoLock() (*config.Config, error) {
	data, err := os.ReadFile(globalConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config config.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &config, nil
}

// ParseGlobalConfig reads and parses the global configuration file into a router.RouterInstance.
func ParseGlobalConfig() (*RouterInstance, error) {
	config, err := ReadGlobalConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	instance, err := NewRouterInstance(*config)
	if err != nil {
		return nil, fmt.Errorf("error creating instance: %w", err)
	}

	return instance, nil
}

// UpdateGlobalConfig updates the global configuration file using the provided updater function.
func UpdateGlobalConfig(updater func(config *config.Config) error) error {
	globalConfigLock.Lock()
	defer globalConfigLock.Unlock()

	config, err := readGlobalConfigNoLock()
	if err != nil {
		return fmt.Errorf("unable to read config file: %w", err)
	}

	// Call the updater function to modify the config
	if err := updater(config); err != nil {
		return fmt.Errorf("error updating config: %w", err)
	}

	// Bump the version number
	config.Version = config.Version.NextMinorVersion()

	// Verify that the new config is valid
	_, err = NewRouterInstance(*config)
	if err != nil {
		return fmt.Errorf("new config is not valid: %w", err)
	}

	// Write the updated config back to the file
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config to JSON: %w", err)
	}

	if err := os.WriteFile(globalConfigFile, data, 0644); err != nil {
		return fmt.Errorf("error writing updated config to file: %w", err)
	}

	return nil
}
