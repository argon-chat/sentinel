package config

import "fmt"

var Instance *Config

type Config struct {
	Projects map[string]string `json:"projects"`
	Port     string            `json:"port"`
	Route    string            `json:"route"`
}

func Parse(apps any, server any) {
	projects, err := parseApps(apps)
	if err != nil {
		panic(fmt.Errorf("failed to parse apps: %w", err))
	}
	port, route, err := parseServer(server)
	if err != nil {
		panic(fmt.Errorf("failed to parse server: %w", err))
	}

	Instance = &Config{Projects: projects, Port: port, Route: route}
}

func parseServer(val any) (string, string, error) {
	server, ok := val.(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("invalid type for server: expected map[string]interface{}, got %T", val)
	}
	port, ok := server["port"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid type for server.port: expected string, got %T", server["port"])
	}
	route, ok := server["route"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid type for server.route: expected string, got %T", server["route"])
	}
	for _, err := range []error{
		validatePort(port), validateRoute(route),
	} {
		if err != nil {
			return "", "", fmt.Errorf("server configuration error: %w", err)
		}
	}

	return port, route, nil
}

func validateRoute(route string) error {
	if route == "" {
		return fmt.Errorf("server.route cannot be empty")
	}
	if route[0] != '/' {
		return fmt.Errorf("server.route must start with '/'")
	}
	if route[len(route)-1] == '/' {
		return fmt.Errorf("server.route must not end with '/'")
	}
	if len(route) > 1 && route[len(route)-1] == '/' {
		return fmt.Errorf("server.route must not end with '/'")
	}
	return nil
}

func validatePort(port string) error {
	if port == "" {
		return fmt.Errorf("server.port cannot be empty")
	}
	return nil
}

func parseApps(val any) (map[string]string, error) {
	projects, ok := val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type for projects: expected map[string]interface{}, got %T", val)
	}

	result := make(map[string]string)
	for k, v := range projects {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}

	return result, nil
}
