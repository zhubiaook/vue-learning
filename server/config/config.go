package config

type Config struct {
	Port                 string                        `mapstructure:"port"`
	LogLevel             string                        `mapstructure:"logLevel"`
	LogMode              string                        `mapstructure:"logMode"`
	SystemPrompt         string                        `mapstructure:"systemPrompt"`
	TitleGeneratorPrompt string                        `mapstructure:"titleGeneratorPrompt"`
	LLM                  llmConfig                     `mapstructure:"llm"`
	MCPSSEServers        map[string]mcpSSEServerConfig `mapstructure:"mcpSSEServers"`
}

type llmConfig struct {
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`
	ApiKey   string `mapstructure:"apiKey"`
	Endpoint string `mapstructure:"endpoint"`
}

type mcpSSEServerConfig struct {
	URL            string `yaml:"url"`
	MaxPayloadSize int    `yaml:"maxPayloadSize"`
}
