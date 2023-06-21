package config

import "os"

type AuthConfig struct {
	Iss    string
	Aud    string
	JwkUri string
}

const (
	AUTH_ISS    = "AUTH_ISS"
	AUTH_AUD    = "AUTH_AUD"
	AUTH_JWKURI = "AUTH_JWKURI"
)

func GetAuthConfig() (*AuthConfig, error) {
	data, err := getJsonReader().readConfigJson()
	if err != nil {
		return nil, err
	}

	config := AuthConfig{}

	configJson := data["auth"].(map[string]any)
	config.Iss = configJson["iss"].(string)
	config.Aud = configJson["aud"].(string)
	config.JwkUri = configJson["jwkUri"].(string)

	if os.Getenv(AUTH_ISS) != "" {
		config.Iss = os.Getenv(AUTH_ISS)
	}
	if os.Getenv(AUTH_AUD) != "" {
		config.Aud = os.Getenv(AUTH_AUD)
	}
	if os.Getenv(AUTH_JWKURI) != "" {
		config.JwkUri = os.Getenv(AUTH_JWKURI)
	}

	return &config, nil
}
