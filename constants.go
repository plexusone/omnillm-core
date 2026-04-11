package omnillm

import "github.com/plexusone/omnillm-core/models"

const (
	EnvVarAnthropicAPIKey = "ANTHROPIC_API_KEY" // #nosec G101
	EnvVarOpenAIAPIKey    = "OPENAI_API_KEY"    // #nosec G101
	EnvVarGeminiAPIKey    = "GEMINI_API_KEY"    // #nosec G101
	EnvVarXAIAPIKey       = "XAI_API_KEY"       // #nosec G101
	EnvVarKimiAPIKey      = "KIMI_API_KEY"      // #nosec G101
	EnvVarGLMAPIKey       = "GLM_API_KEY"       // #nosec G101
	EnvVarQwenAPIKey      = "QWEN_API_KEY"      // #nosec G101
)

// ProviderName represents the different LLM provider names
type ProviderName string

const (
	ProviderNameOpenAI    ProviderName = "openai"
	ProviderNameAnthropic ProviderName = "anthropic"
	ProviderNameBedrock   ProviderName = "bedrock"
	ProviderNameOllama    ProviderName = "ollama"
	ProviderNameGemini    ProviderName = "gemini"
	ProviderNameXAI       ProviderName = "xai"
	ProviderNameKimi      ProviderName = "kimi"
	ProviderNameGLM       ProviderName = "glm"
	ProviderNameQwen      ProviderName = "qwen"
)

// Common model constants for each provider.
//
// NOTE: For new code, prefer importing "github.com/plexusone/omnillm-core/models" directly
// for better organization and documentation. These constants are maintained for
// backwards compatibility with existing code.
const (
	// Bedrock Models - Re-exported from models package
	ModelBedrockClaude3Opus   = models.BedrockClaude3Opus
	ModelBedrockClaude3Sonnet = models.BedrockClaude3Sonnet
	ModelBedrockClaudeOpus4   = models.BedrockClaudeOpus4
	ModelBedrockTitan         = models.BedrockTitan

	// Claude Models - Re-exported from models package
	ModelClaudeOpus4_1   = models.ClaudeOpus4_1
	ModelClaudeOpus4     = models.ClaudeOpus4
	ModelClaudeSonnet4   = models.ClaudeSonnet4
	ModelClaude3_7Sonnet = models.Claude3_7Sonnet
	ModelClaude3_5Haiku  = models.Claude3_5Haiku
	ModelClaude3Opus     = models.Claude3Opus
	ModelClaude3Sonnet   = models.Claude3Sonnet
	ModelClaude3Haiku    = models.Claude3Haiku

	// Gemini Models - Re-exported from models package
	ModelGemini2_5Pro       = models.Gemini2_5Pro
	ModelGemini2_5Flash     = models.Gemini2_5Flash
	ModelGeminiLive2_5Flash = models.GeminiLive2_5Flash
	ModelGemini1_5Pro       = models.Gemini1_5Pro
	ModelGemini1_5Flash     = models.Gemini1_5Flash
	ModelGeminiPro          = models.GeminiPro

	// Ollama Models - Re-exported from models package
	ModelOllamaLlama3_8B   = models.OllamaLlama3_8B
	ModelOllamaLlama3_70B  = models.OllamaLlama3_70B
	ModelOllamaMistral7B   = models.OllamaMistral7B
	ModelOllamaMixtral8x7B = models.OllamaMixtral8x7B
	ModelOllamaCodeLlama   = models.OllamaCodeLlama
	ModelOllamaGemma2B     = models.OllamaGemma2B
	ModelOllamaGemma7B     = models.OllamaGemma7B
	ModelOllamaQwen2_5     = models.OllamaQwen2_5
	ModelOllamaDeepSeek    = models.OllamaDeepSeek

	// OpenAI Models - Re-exported from models package
	ModelGPT5           = models.GPT5
	ModelGPT5Mini       = models.GPT5Mini
	ModelGPT5Nano       = models.GPT5Nano
	ModelGPT5ChatLatest = models.GPT5ChatLatest
	ModelGPT4_1         = models.GPT4_1
	ModelGPT4_1Mini     = models.GPT4_1Mini
	ModelGPT4_1Nano     = models.GPT4_1Nano
	ModelGPT4o          = models.GPT4o
	ModelGPT4oMini      = models.GPT4oMini
	ModelGPT4Turbo      = models.GPT4Turbo
	ModelGPT35Turbo     = models.GPT35Turbo

	// Vertex AI Models - Re-exported from models package
	ModelVertexClaudeOpus4 = models.VertexClaudeOpus4

	// X.AI Grok Models - Re-exported from models package
	// Grok 4.1 (Latest - November 2025)
	ModelGrok4_1FastReasoning    = models.Grok4_1FastReasoning
	ModelGrok4_1FastNonReasoning = models.Grok4_1FastNonReasoning

	// Grok 4 (July 2025)
	ModelGrok4_0709            = models.Grok4_0709
	ModelGrok4FastReasoning    = models.Grok4FastReasoning
	ModelGrok4FastNonReasoning = models.Grok4FastNonReasoning
	ModelGrokCodeFast1         = models.GrokCodeFast1

	// Grok 3
	ModelGrok3     = models.Grok3
	ModelGrok3Mini = models.Grok3Mini

	// Grok 2
	ModelGrok2_1212   = models.Grok2_1212
	ModelGrok2_Vision = models.Grok2_Vision

	// Deprecated models
	ModelGrokBeta   = models.GrokBeta
	ModelGrokVision = models.GrokVision

	// Kimi / Moonshot AI Models - Re-exported from models package
	ModelKimiK2_5        = models.KimiK2_5
	ModelKimiK2Turbo     = models.KimiK2Turbo
	ModelKimiK2Thinking  = models.KimiK2Thinking
	ModelMoonshotV1_8K   = models.MoonshotV1_8K
	ModelMoonshotV1_32K  = models.MoonshotV1_32K
	ModelMoonshotV1_128K = models.MoonshotV1_128K

	// Zhipu AI GLM Models - Re-exported from models package
	ModelGLM5         = models.GLM5
	ModelGLM4_7       = models.GLM4_7
	ModelGLM4_7FlashX = models.GLM4_7FlashX
	ModelGLM4_7Flash  = models.GLM4_7Flash
	ModelGLM4_5       = models.GLM4_5
	ModelGLM4_5Flash  = models.GLM4_5Flash

	// Alibaba Cloud Qwen Models - Re-exported from models package
	ModelQwen3Max    = models.Qwen3Max
	ModelQwenMax     = models.QwenMax
	ModelQwenPlus    = models.QwenPlus
	ModelQwenFlash   = models.QwenFlash
	ModelQwQ32B      = models.QwQ32B
	ModelQwen3_235B  = models.Qwen3_235B
	ModelQwen3_32B   = models.Qwen3_32B
	ModelQwen2_5_72B = models.Qwen2_5_72B
)
