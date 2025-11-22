package aletis

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AletisSearch/aletis/web"
)

type Config struct {
	Dev              bool
	Port             string
	OpenAIKey        string
	OpenAIURL        string
	SearxngHost      string
	Public           bool
	AIEnabled        bool
	PostgresHost     string
	PostgresPort     string
	PostgresDatabase string
	PostgresUsername string
	PostgresPassword string
}

type Option func(*Config) error

func WithDevString(dev string) Option {
	return func(c *Config) error {
		boolValue, err := strconv.ParseBool(dev)
		if err != nil {
			return fmt.Errorf("unable to parse dev environment variable: %w", err)
		}
		c.Dev = boolValue
		return nil
	}
}

func WithPort(port string) Option {
	return func(c *Config) error {
		c.Port = port
		return nil
	}
}

func WithOpenAIKey(key string) Option {
	return func(c *Config) error {
		c.OpenAIKey = key
		return nil
	}
}

func WithOpenAIURL(url string) Option {
	return func(c *Config) error {
		c.OpenAIURL = url
		return nil
	}
}

func WithSearxngHost(host string) Option {
	return func(c *Config) error {
		c.SearxngHost = host
		return nil
	}
}

func WithPublicString(public string) Option {
	return func(c *Config) error {
		boolValue, err := strconv.ParseBool(public)
		if err != nil {
			return fmt.Errorf("unable to parse public environment variable: %w", err)
		}
		c.Public = boolValue
		return nil
	}
}

func WithAIEnabledString(enabled string) Option {
	return func(c *Config) error {
		boolValue, err := strconv.ParseBool(enabled)
		if err != nil {
			return fmt.Errorf("unable to parse AI_ENABLED environment variable: %w", err)
		}
		c.AIEnabled = boolValue
		return nil
	}
}

func WithPostgresHost(host string) Option {
	return func(c *Config) error {
		c.PostgresHost = host
		return nil
	}
}

func WithPostgresPort(port string) Option {
	return func(c *Config) error {
		c.PostgresPort = port
		return nil
	}
}

func WithPostgresDatabase(database string) Option {
	return func(c *Config) error {
		c.PostgresDatabase = database
		return nil
	}
}

func WithPostgresUsername(username string) Option {
	return func(c *Config) error {
		c.PostgresUsername = username
		return nil
	}
}

func WithPostgresPassword(password string) Option {
	return func(c *Config) error {
		c.PostgresPassword = password
		return nil
	}
}

func EnvConfigOptions() []Option {
	confOptions := []Option{
		WithSearxngHost(trimGetEnv("SEARXNG_HOST")),
	}
	if dev, exist := trimLookupEnv("DEV"); exist {
		confOptions = append(confOptions, WithDevString(dev))
	}
	if addr, exist := trimLookupEnv("PORT"); exist {
		confOptions = append(confOptions, WithPort(addr))
	}
	if public, ok := trimLookupEnv("PUBLIC"); ok {
		confOptions = append(confOptions, WithPublicString(public))
	}
	// AI
	if aiEnabled, ok := trimLookupEnv("AI_ENABLED"); ok {
		confOptions = append(confOptions, WithAIEnabledString(aiEnabled))
	}
	if openaiKey, ok := trimLookupEnv("OPENAI_API_KEY"); ok {
		confOptions = append(confOptions, WithOpenAIKey(openaiKey))
	}
	if openaiURL, ok := trimLookupEnv("OPENAI_URL"); ok {
		confOptions = append(confOptions, WithOpenAIURL(openaiURL))
	}
	// PostgreSQL
	if postgresHost, ok := trimLookupEnv("POSTGRES_HOST"); ok {
		confOptions = append(confOptions, WithPostgresHost(postgresHost))
	}
	if postgresPort, ok := trimLookupEnv("POSTGRES_PORT"); ok {
		confOptions = append(confOptions, WithPostgresPort(postgresPort))
	}
	if postgresDatabase, ok := trimLookupEnv("POSTGRES_DATABASE"); ok {
		confOptions = append(confOptions, WithPostgresDatabase(postgresDatabase))
	}
	if postgresUsername, ok := trimLookupEnv("POSTGRES_USERNAME"); ok {
		confOptions = append(confOptions, WithPostgresUsername(postgresUsername))
	}
	if postgresPassword, ok := trimLookupEnv("POSTGRES_PASSWORD"); ok {
		confOptions = append(confOptions, WithPostgresPassword(postgresPassword))
	}
	return confOptions
}
func trimGetEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
func trimLookupEnv(key string) (string, bool) {
	env, ok := os.LookupEnv(key)
	env = strings.TrimSpace(env)
	if env == "" {
		return "", false
	}
	return env, ok
}

func NewConfig(options ...Option) (*Config, error) {
	conf := &Config{
		Dev:              false,
		Port:             "8080",
		Public:           true,
		AIEnabled:        false,
		PostgresPort:     "5432",
		PostgresDatabase: "aletis",
		PostgresUsername: "aletis",
	}
	for _, o := range options {
		if err := o(conf); err != nil {
			return nil, err
		}
	}

	if conf.SearxngHost == "" {
		return nil, errors.New("SEARXNG_HOST is not set")
	}

	// OpenAI configuration is only required if AI is enabled
	if conf.AIEnabled {
		if conf.OpenAIKey == "" {
			return nil, errors.New("OPENAI_API_KEY is not set")
		}
		if conf.OpenAIURL == "" {
			// Set default OpenAI URL if not provided
			conf.OpenAIURL = "https://openrouter.ai/api/v1"
		}
	}

	// PostgreSQL configuration validation
	if conf.PostgresHost == "" {
		return nil, errors.New("POSTGRES_HOST is not set")
	}
	if conf.PostgresDatabase == "" {
		return nil, errors.New("POSTGRES_DATABASE is not set")
	}
	if conf.PostgresUsername == "" {
		return nil, errors.New("POSTGRES_USERNAME is not set")
	}
	if conf.PostgresPassword == "" {
		return nil, errors.New("POSTGRES_PASSWORD is not set")
	}

	web.IsDev = conf.Dev
	return conf, nil
}
func (c *Config) DBconnStr() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUsername,
		c.PostgresPassword,
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresDatabase,
	)
}
