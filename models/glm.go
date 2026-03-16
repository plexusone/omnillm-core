package models

// GLM (Zhipu AI BigModel) Model Documentation
const (
	// GLMModelsURL is the official Zhipu AI models documentation page.
	// Use this to check for new models, deprecations, and model updates.
	GLMModelsURL = "https://bigmodel.cn/dev/api/normal-model/glm-4"

	// GLMAPIURL is the GLM API reference page.
	GLMAPIURL = "https://open.bigmodel.cn/dev/howuse/introduction"
)

// GLM-5 Series - Flagship models with MoE architecture
const (
	// GLM5 is the flagship model with MoE architecture (744B/40B active).
	// Agentic engineering, 200K context, forced thinking mode.
	GLM5 = "glm-5"

	// GLM5Code is the code-specialized variant of GLM-5.
	// Optimized for programming, 200K context, forced thinking.
	GLM5Code = "glm-5-code"
)

// GLM-4.7 Series - Premium with Interleaved Thinking
const (
	// GLM4_7 is the premium model with Interleaved Thinking, 200K context.
	GLM4_7 = "glm-4.7"

	// GLM4_7FlashX is the high-speed paid version with priority GPU access.
	// 200K context, hybrid thinking, best price/performance for batch tasks.
	GLM4_7FlashX = "glm-4.7-flashx"

	// GLM4_7Flash is a free SOTA model with 200K context and hybrid thinking.
	// 1 concurrent request limit, ideal for prototyping.
	GLM4_7Flash = "glm-4.7-flash"
)

// GLM-4.6 Series - Balanced with Auto Thinking
const (
	// GLM4_6 is a balanced model with 200K context and auto-thinking.
	GLM4_6 = "glm-4.6"
)

// GLM-4.5 Series - Unified Reasoning, Coding, and Agents
const (
	// GLM4_5 is the first unified model with reasoning/coding/agent capabilities.
	// MoE 355B/32B active, 128K context, auto-thinking.
	GLM4_5 = "glm-4.5"

	// GLM4_5X is the ultra-fast premium version with lowest latency.
	// 128K context, auto-thinking.
	GLM4_5X = "glm-4.5-x"

	// GLM4_5Air is the cost-effective lightweight model.
	// MoE 106B/12B active, 128K context, auto-thinking.
	GLM4_5Air = "glm-4.5-air"

	// GLM4_5AirX is the accelerated Air version with priority GPU access.
	GLM4_5AirX = "glm-4.5-airx"

	// GLM4_5Flash is a free model with reasoning/coding/agents support.
	// 128K context, auto-thinking, function calling enabled.
	GLM4_5Flash = "glm-4.5-flash"
)

// GLM-4 Legacy - Dense Architecture
const (
	// GLM4_32B is an ultra-budget dense 32B model, 128K context, no thinking mode.
	GLM4_32B = "glm-4-32b-0414-128k"
)
