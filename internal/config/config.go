package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	LLM     LLMConfig     `yaml:"llm"`
	Engine  EngineConfig  `yaml:"engine"`
	Tools   ToolsConfig   `yaml:"tools"`
	Session SessionConfig `yaml:"session"`
	Log     LogConfig     `yaml:"log"`
}

type ServerConfig struct {
	Host        string   `yaml:"host"`
	Port        int      `yaml:"port"`
	CORSOrigins []string `yaml:"cors_origins"`
}

type LLMConfig struct {
	Provider      string        `yaml:"provider"`
	APIURL        string        `yaml:"api_url"`
	APIKey        string        `yaml:"api_key"` // 通过环境变量 ARK_API_KEY 注入
	Model         string        `yaml:"model"`
	MaxTokens     int           `yaml:"max_tokens"`
	Temperature   float64       `yaml:"temperature"`
	Timeout       time.Duration `yaml:"timeout"`
	MaxRetries    int           `yaml:"max_retries"`
	RetryInterval time.Duration `yaml:"retry_interval"`
}

type EngineConfig struct {
	MaxReactIterations int    `yaml:"max_react_iterations"`
	PlannerModel       string `yaml:"planner_model"`
}

type ToolsConfig struct {
	Flight  ToolConfig `yaml:"flight"`
	Hotel   ToolConfig `yaml:"hotel"`
	Weather ToolConfig `yaml:"weather"`
}

type ToolConfig struct {
	ProviderType  string        `yaml:"provider_type"` // mock | real
	Timeout       time.Duration `yaml:"timeout"`
	RetryCount    int           `yaml:"retry_count"`
	RetryInterval time.Duration `yaml:"retry_interval"`
}

type SessionConfig struct {
	StoreType   string        `yaml:"store_type"` // memory | sqlite
	MaxSessions int           `yaml:"max_sessions"`
	TTL         time.Duration `yaml:"ttl"`
}

type LogConfig struct {
	Level  string `yaml:"level"`  // debug | info | warn | error
	Format string `yaml:"format"` // text | json
}

// Load 从文件加载配置，并用环境变量覆盖敏感字段
func Load(path string) (*Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 配置文件不存在时使用默认配置
			applyEnvOverrides(cfg)
			return cfg, nil
		}
		return nil, fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	applyEnvOverrides(cfg)
	return cfg, nil
}

// defaultConfig 返回默认配置
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:        "0.0.0.0",
			Port:        8080,
			CORSOrigins: []string{"http://localhost:5173"},
		},
		LLM: LLMConfig{
			Provider:      "doubao",
			APIURL:        "https://ark.cn-beijing.volces.com/api/v3/chat/completions",
			Model:         "doubao-seed-2-0-lite-260428",
			MaxTokens:     4096,
			Temperature:   0.7,
			Timeout:       30 * time.Second,
			MaxRetries:    2,
			RetryInterval: 3 * time.Second,
		},
		Engine: EngineConfig{
			MaxReactIterations: 10,
			PlannerModel:       "doubao-seed-2-0-lite-260428",
		},
		Tools: ToolsConfig{
			Flight: ToolConfig{
				ProviderType:  "mock",
				Timeout:       5 * time.Second,
				RetryCount:    1,
				RetryInterval: 2 * time.Second,
			},
			Hotel: ToolConfig{
				ProviderType:  "mock",
				Timeout:       5 * time.Second,
				RetryCount:    1,
				RetryInterval: 2 * time.Second,
			},
			Weather: ToolConfig{
				ProviderType:  "mock",
				Timeout:       5 * time.Second,
				RetryCount:    1,
				RetryInterval: 2 * time.Second,
			},
		},
		Session: SessionConfig{
			StoreType:   "memory",
			MaxSessions: 100,
			TTL:         24 * time.Hour,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

// applyEnvOverrides 用环境变量覆盖配置
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("ARK_API_KEY"); v != "" {
		cfg.LLM.APIKey = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
}
