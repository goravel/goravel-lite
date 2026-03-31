package config

import (
	"goravel/app/facades"

	"github.com/goravel/framework/ai/openai"
	"github.com/goravel/framework/contracts/ai"
)

func init() {
	config := facades.Config()
	config.Add("ai", map[string]any{
		// Default AI Provider
		//
		// This option controls the default AI provider that will be used.
		"default": "openai",

		// AI Providers
		//
		// Here you may configure each AI provider used by your application.
		// A variety of drivers are available, and each provider may also
		// configure the models available to your application.
		"providers": map[string]any{
			"openai": map[string]any{
				"key": "sk-Ycb6dZrNcaGSqNUaEoTAL7OPg3WxD0yX",
				"url": "https://codex-api.packycode.com/v1",
				"via": func() (ai.Provider, error) {
					return openai.NewOpenAI(config, "openai")
				},
			},
		},
	})
}
