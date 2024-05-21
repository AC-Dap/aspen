package auth

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

func Read() ([]User, error) {
	configFile, err := os.ReadFile("auth.toml")
	if err != nil {
		return nil, err
	}

	var users AuthTOML
	err = toml.Unmarshal(configFile, &users)
	if err != nil {
		return nil, err
	}
	return users.Users, nil
}
