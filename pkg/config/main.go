package config

var Instance *Config

type Config struct {
	Projects map[string]string `json:"projects"`
}

func Parse(val any) {
	projects, ok := val.(map[string]interface{})
	if !ok {
		panic("Invalid configuration format for projects")
	}

	result := make(map[string]string)
	for k, v := range projects {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}

	Instance = &Config{Projects: result}
}
