package options

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kingpin/v2"
)

var (
	opnsenseProtocol = kingpin.Flag(
		"opnsense.protocol",
		"Protocol to use to connect to OPNsense API. One of: [http, https]",
	).Envar("OPNSENSE_EXPORTER_OPS_PROTOCOL").Required().String()
	opnsenseAPI = kingpin.Flag(
		"opnsense.address",
		"Hostname or IP address of OPNsense API",
	).Envar("OPNSENSE_EXPORTER_OPS_API").Required().String()
	opnsenseAPIKey = kingpin.Flag(
		"opnsense.api-key",
		"API key to use to connect to OPNsense API. This flag/ENV or the OPS_API_KEY_FILE my be set.",
	).Default("").Envar("OPNSENSE_EXPORTER_OPS_API_KEY").String()
	opnsenseAPISecret = kingpin.Flag(
		"opnsense.api-secret",
		"API secret to use to connect to OPNsense API. This flag/ENV or the OPS_API_SECRET_FILE my be set.",
	).Default("").Envar("OPNSENSE_EXPORTER_OPS_API_SECRET").String()
	opnsenseInsecure = kingpin.Flag(
		"opnsense.insecure",
		"Disable TLS certificate verification",
	).Envar("OPNSENSE_EXPORTER_OPS_INSECURE").Default("false").Bool()
)

// ReadFirstLine opens a file and reads its first line.
// It returns the first line as a string and any error encountered.
func getLineFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err // Return an empty string and the error
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func opsAPISecret() (string, error) {
	if env, ok := os.LookupEnv("OPS_API_SECRET_FILE"); ok {
		apiKey, err := getLineFromFile(env)
		if err != nil {
			return "", errors.Join(fmt.Errorf("failed to read OPS_API_SECRET_FILE"), err)
		}
		if len(apiKey) > 0 {
			return apiKey, nil
		}
	}
	if *opnsenseAPIKey == "" {
		return "", fmt.Errorf("opnsense.api-secret or OPS_API_SECRET_FILE must be set")
	}

	return *opnsenseAPISecret, nil
}

func opsAPIKey() (string, error) {
	if env, ok := os.LookupEnv("OPS_API_KEY_FILE"); ok {
		apiSecret, err := getLineFromFile(env)
		if err != nil {
			return "", errors.Join(fmt.Errorf("failed to read OPS_API_KEY_FILE"), err)
		}
		if len(apiSecret) > 0 {
			return apiSecret, nil
		}
	}
	if *opnsenseAPISecret == "" {
		return "", fmt.Errorf("opnsense.api-key or OPS_API_KEY_FILE must be set")
	}

	return *opnsenseAPIKey, nil
}

// OPNSenseConfig holds the configuration for the OPNsense API.
type OPNSenseConfig struct {
	Protocol  string
	Host      string
	APIKey    string
	APISecret string
	Insecure  bool
}

// Validate checks if the configuration is valid.
// returns an error on any missing value
func (c *OPNSenseConfig) Validate() error {
	if c.Protocol != "http" && c.Protocol != "https" {
		return fmt.Errorf("protocol must be one of: [http, https]")
	}
	if c.Host == "" {
		return fmt.Errorf("host must be set")
	}
	if c.APIKey == "" {
		return fmt.Errorf("api-key must be set")
	}
	if c.APISecret == "" {
		return fmt.Errorf("api-secret must be set")
	}
	return nil
}

func OPNSense() (*OPNSenseConfig, error) {
	apiKey, err := opsAPIKey()
	if err != nil {
		return nil, err
	}
	apiSecret, err := opsAPISecret()
	if err != nil {
		return nil, err
	}
	conf := &OPNSenseConfig{
		Protocol:  strings.TrimSpace(*opnsenseProtocol),
		Host:      strings.TrimSpace(*opnsenseAPI),
		APIKey:    apiKey,
		APISecret: apiSecret,
		Insecure:  *opnsenseInsecure,
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return conf, nil
}
