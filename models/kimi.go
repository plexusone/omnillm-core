package models

// Kimi (Moonshot AI) Model Documentation
const (
	// KimiModelsURL is the official Moonshot AI models documentation page.
	// Use this to check for new models, deprecations, and model updates.
	KimiModelsURL = "https://platform.moonshot.cn/docs/api/chat"

	// KimiAPIURL is the Kimi API reference page.
	KimiAPIURL = "https://platform.moonshot.cn/docs"
)

// Kimi K2.5 series - Most intelligent and versatile multimodal model
const (
	// KimiK2_5 is the most intelligent model with native multimodal architecture.
	// Supports visual/text input, thinking/non-thinking modes, 256k context.
	KimiK2_5 = "kimi-k2.5"
)

// Kimi K2 series - MoE foundation models with 1T total params, 32B activated
const (
	// KimiK2_0905 is an enhanced agentic coding model with improved frontend code
	// quality and context understanding, 256k context window.
	KimiK2_0905 = "kimi-k2-0905-preview"

	// KimiK2_0711 is the MoE base model with powerful code and agent capabilities,
	// 128k context window.
	KimiK2_0711 = "kimi-k2-0711-preview"

	// KimiK2Turbo is the high-speed version of K2-0905, 60-100 tokens/sec output speed,
	// 256k context window.
	KimiK2Turbo = "kimi-k2-turbo-preview"

	// KimiK2Thinking is a long-term thinking model with multi-step tool usage
	// and deep reasoning, 256k context window.
	KimiK2Thinking = "kimi-k2-thinking"

	// KimiK2ThinkingTurbo is the high-speed thinking model, 60-100 tokens/sec,
	// excels at deep reasoning, 256k context window.
	KimiK2ThinkingTurbo = "kimi-k2-thinking-turbo"
)

// Moonshot V1 series - General text generation models
const (
	// MoonshotV1_8K is suitable for generating short texts with 8k context window.
	MoonshotV1_8K = "moonshot-v1-8k"

	// MoonshotV1_32K is suitable for generating long texts with 32k context window.
	MoonshotV1_32K = "moonshot-v1-32k"

	// MoonshotV1_128K is suitable for generating very long texts with 128k context window.
	MoonshotV1_128K = "moonshot-v1-128k"
)

// Moonshot V1 Vision series - Multimodal models with image understanding
const (
	// MoonshotV1_8KVision understands image content and outputs text, 8k context.
	MoonshotV1_8KVision = "moonshot-v1-8k-vision-preview"

	// MoonshotV1_32KVision understands image content and outputs text, 32k context.
	MoonshotV1_32KVision = "moonshot-v1-32k-vision-preview"

	// MoonshotV1_128KVision understands image content and outputs text, 128k context.
	MoonshotV1_128KVision = "moonshot-v1-128k-vision-preview"
)
