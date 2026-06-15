package models

// AWS Bedrock Model Documentation
const (
	// BedrockModelsURL is the official AWS Bedrock models documentation page.
	// Use this to check for new models, deprecations, and model updates.
	BedrockModelsURL = "https://docs.aws.amazon.com/bedrock/latest/userguide/models-supported.html"

	// BedrockAPIURL is the AWS Bedrock API reference page.
	BedrockAPIURL = "https://docs.aws.amazon.com/bedrock/latest/APIReference/welcome.html"
)

// Bedrock Claude Fable 5 Models
const (
	// BedrockClaudeFable5 is Claude Fable 5 on AWS Bedrock.
	BedrockClaudeFable5 = "anthropic.claude-fable-5"
)

// Bedrock Claude 4.8/4.7/4.6 Models (Dateless format)
const (
	// BedrockClaudeOpus4_8 is Claude Opus 4.8 on AWS Bedrock.
	BedrockClaudeOpus4_8 = "anthropic.claude-opus-4-8"

	// BedrockClaudeOpus4_7 is Claude Opus 4.7 on AWS Bedrock.
	BedrockClaudeOpus4_7 = "anthropic.claude-opus-4-7"

	// BedrockClaudeOpus4_6 is Claude Opus 4.6 on AWS Bedrock.
	BedrockClaudeOpus4_6 = "anthropic.claude-opus-4-6-v1"

	// BedrockClaudeSonnet4_6 is Claude Sonnet 4.6 on AWS Bedrock.
	BedrockClaudeSonnet4_6 = "anthropic.claude-sonnet-4-6"
)

// Bedrock Claude 4.5 Models
const (
	// BedrockClaudeOpus4_5 is Claude Opus 4.5 on AWS Bedrock.
	BedrockClaudeOpus4_5 = "anthropic.claude-opus-4-5-20251101-v1:0"

	// BedrockClaudeSonnet4_5 is Claude Sonnet 4.5 on AWS Bedrock.
	BedrockClaudeSonnet4_5 = "anthropic.claude-sonnet-4-5-20250929-v1:0"

	// BedrockClaudeHaiku4_5 is Claude Haiku 4.5 on AWS Bedrock.
	BedrockClaudeHaiku4_5 = "anthropic.claude-haiku-4-5-20251001-v1:0"
)

// Bedrock Claude 4.x Models (Deprecated)
const (
	// BedrockClaudeOpus4_1 is Claude Opus 4.1 on AWS Bedrock (deprecated).
	BedrockClaudeOpus4_1 = "anthropic.claude-opus-4-1-20250805-v1:0"

	// BedrockClaudeOpus4 is Claude Opus 4 on AWS Bedrock (deprecated).
	BedrockClaudeOpus4 = "anthropic.claude-opus-4-20250514-v1:0"

	// BedrockClaudeSonnet4 is Claude Sonnet 4 on AWS Bedrock (deprecated).
	BedrockClaudeSonnet4 = "anthropic.claude-sonnet-4-20250514-v1:0"
)

// Bedrock Claude 3 Models
const (
	// BedrockClaude3Opus is Claude 3 Opus on AWS Bedrock.
	BedrockClaude3Opus = "anthropic.claude-3-opus-20240229-v1:0"

	// BedrockClaude3Sonnet is Claude 3 Sonnet on AWS Bedrock.
	BedrockClaude3Sonnet = "anthropic.claude-3-sonnet-20240229-v1:0"
)

// Bedrock Amazon Titan Models
const (
	// BedrockTitan is Amazon Titan Text Express on AWS Bedrock.
	BedrockTitan = "amazon.titan-text-express-v1"
)
