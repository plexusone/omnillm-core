package models

// Qwen (Alibaba Cloud Tongyi Qianwen) Model Documentation
const (
	// QwenModelsURL is the official Qwen model documentation page.
	// Use this to check for new models, deprecations, and model updates.
	QwenModelsURL = "https://help.aliyun.com/zh/model-studio/getting-started/models"

	// QwenAPIURL is the Qwen DashScope API reference page.
	QwenAPIURL = "https://help.aliyun.com/zh/model-studio/developer-reference/use-qwen-by-calling-api"
)

// Qwen3 Flagship series - Wide availability (International/Global/China)
const (
	// Qwen3Max is the latest flagship reasoning model with thinking capability.
	// Strong instruction following and complex task handling.
	Qwen3Max = "qwen3-max"

	// Qwen3MaxPreview is the preview version with extended thinking capabilities.
	Qwen3MaxPreview = "qwen3-max-preview"

	// QwenMax is the flagship reasoning model (International/China regions only).
	QwenMax = "qwen-max"
)

// Qwen3.5 series - Balanced models
const (
	// Qwen3_5Plus is an advanced balanced model with thinking mode support.
	Qwen3_5Plus = "qwen3.5-plus"

	// QwenPlus is the balanced performance model with thinking mode.
	// Wide availability including US region.
	QwenPlus = "qwen-plus"

	// Qwen3_5Flash is an ultra-fast lightweight model for high-throughput.
	Qwen3_5Flash = "qwen3.5-flash"

	// QwenFlash is the fast lightweight model with context caching.
	// Wide availability including US region.
	QwenFlash = "qwen-flash"

	// QwenTurbo is the fastest lightweight model (deprecated, use QwenFlash).
	// Deprecated: use QwenFlash instead.
	QwenTurbo = "qwen-turbo"
)

// QwQ series - Deep reasoning models
const (
	// QwQPlus is the deep reasoning model with extended chain-of-thought.
	QwQPlus = "qwq-plus"

	// QwQ32B is the open-source 32B reasoning model with powerful logic capabilities.
	QwQ32B = "qwq-32b"
)

// Open source Qwen3.5 series
const (
	// Qwen3_5_397B is the largest open-source model with 397B parameters.
	Qwen3_5_397B = "qwen3.5-397b-a17b"

	// Qwen3_5_122B is the large open-source model with 122B parameters.
	Qwen3_5_122B = "qwen3.5-122b-a10b"

	// Qwen3_5_27B is the medium open-source model with 27B parameters.
	Qwen3_5_27B = "qwen3.5-27b"

	// Qwen3_5_35B is the efficient open-source model with 35B parameters.
	Qwen3_5_35B = "qwen3.5-35b-a3b"
)

// Open source Qwen3 series
const (
	// Qwen3_235B is the dual-mode 235B model supporting thinking and non-thinking.
	Qwen3_235B = "qwen3-235b-a22b"

	// Qwen3_32B is the versatile 32B model with dual-mode capabilities.
	Qwen3_32B = "qwen3-32b"

	// Qwen3_30B is the efficient 30B model with MoE architecture.
	Qwen3_30B = "qwen3-30b-a3b"

	// Qwen3_14B is the medium-sized 14B model with good performance-cost balance.
	Qwen3_14B = "qwen3-14b"

	// Qwen3_8B is the compact 8B model optimized for efficiency.
	Qwen3_8B = "qwen3-8b"
)

// Open source Qwen2.5 series
const (
	// Qwen2_5_72B is the large 72B instruction-following model.
	Qwen2_5_72B = "qwen2.5-72b-instruct"

	// Qwen2_5_32B is the medium 32B instruction-following model.
	Qwen2_5_32B = "qwen2.5-32b-instruct"

	// Qwen2_5_14B is the compact 14B instruction-following model.
	Qwen2_5_14B = "qwen2.5-14b-instruct"

	// Qwen2_5_7B is the small 7B instruction-following model.
	Qwen2_5_7B = "qwen2.5-7b-instruct"

	// Qwen2_5_14B_1M is the extended context 14B model supporting up to 1M tokens.
	Qwen2_5_14B_1M = "qwen2.5-14b-instruct-1m"

	// Qwen2_5_7B_1M is the extended context 7B model supporting up to 1M tokens.
	Qwen2_5_7B_1M = "qwen2.5-7b-instruct-1m"
)
