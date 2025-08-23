package config

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
)

var Instance *Config

type Project struct {
	SentryProjectId string `json:"sentryProjectId"`
	SentryKey       string `json:"sentryKey"`
}

type Server struct {
	Port  int    `json:"port"`
	Route string `json:"route"`
}

type Config struct {
	Projects  map[string]Project `json:"projects"`
	Server    Server             `json:"server"`
	SentryUrl string             `json:"sentryUrl"`
}

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	if err := validateConfig(&cfg); err != nil {
		return err
	}
	Instance = &cfg
	return nil
}

func validateConfig(cfg *Config) error {
	if len(cfg.Projects) == 0 {
		return fmt.Errorf("no projects defined")
	}
	if err := validateSentryUrl(cfg.SentryUrl); err != nil {
		return err
	}
	if err := validatePort(cfg.Server.Port); err != nil {
		return err
	}
	if err := validateRoute(cfg.Server.Route); err != nil {
		return err
	}
	return nil
}

func validateSentryUrl(sentryUrl string) error {
	if sentryUrl == "" {
		return fmt.Errorf("sentryUrl cannot be empty")
	}
	parsed, err := url.Parse(sentryUrl)
	if err != nil {
		return fmt.Errorf("invalid sentryUrl: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("sentryUrl must start with http:// or https://, got %s", parsed.Scheme)
	}
	host := parsed.Hostname()
	if net.ParseIP(host) != nil {
		return fmt.Errorf("sentryUrl must not be an IP address, got %s", host)
	}
	if len(host) < 1 || len(host) > 255 {
		return fmt.Errorf("sentryUrl hostname must be between 1 and 255 characters, got %d", len(host))
	}
	return nil
}

func validateRoute(route string) error {
	if route == "" {
		return fmt.Errorf("server.route cannot be empty")
	}
	if route[0] != '/' {
		return fmt.Errorf("server.route must start with '/'")
	}
	if len(route) > 1 && route[len(route)-1] == '/' {
		return fmt.Errorf("server.route must not end with '/'")
	}
	return nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535, got %d", port)
	}
	return nil
}
