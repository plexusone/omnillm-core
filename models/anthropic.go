package models

// Anthropic Claude Model Documentation
const (
	// AnthropicModelsURL is the official Anthropic models documentation page.
	// Use this to check for new models, deprecations, and model updates.
	AnthropicModelsURL = "https://platform.claude.com/docs/en/about-claude/models/overview"

	// AnthropicAPIURL is the Anthropic API reference page.
	AnthropicAPIURL = "https://docs.anthropic.com/en/api"
)

// Claude Fable 5 and Mythos 5 (Latest frontier)
const (
	ClaudeFable5  = "claude-fable-5"  // Claude Fable 5 - Most capable widely released model
	ClaudeMythos5 = "claude-mythos-5" // Claude Mythos 5 - Project Glasswing (limited availability)
)

// Claude 4.8/4.7/4.6 Family (Dateless format - still pinned snapshots)
const (
	ClaudeOpus4_8   = "claude-opus-4-8"   // Claude Opus 4.8 - Latest Opus-tier model
	ClaudeOpus4_7   = "claude-opus-4-7"   // Claude Opus 4.7
	ClaudeOpus4_6   = "claude-opus-4-6"   // Claude Opus 4.6
	ClaudeSonnet4_6 = "claude-sonnet-4-6" // Claude Sonnet 4.6 - Latest Sonnet-tier model
)

// Claude 4.5 Family
const (
	ClaudeOpus4_5   = "claude-opus-4-5-20251101"   // Claude Opus 4.5 (November 2025)
	ClaudeSonnet4_5 = "claude-sonnet-4-5-20250929" // Claude Sonnet 4.5 (September 2025)
	ClaudeHaiku4_5  = "claude-haiku-4-5-20251001"  // Claude Haiku 4.5 (October 2025)
)

// Claude 4.x Family (Deprecated)
const (
	ClaudeOpus4_1 = "claude-opus-4-1-20250805" // Claude Opus 4.1 (August 2025) - Deprecated, retiring Aug 5, 2026
	ClaudeOpus4   = "claude-opus-4-20250514"   // Claude Opus 4 (May 2025) - Deprecated, retiring Jun 15, 2026
	ClaudeSonnet4 = "claude-sonnet-4-20250514" // Claude Sonnet 4 (May 2025) - Deprecated, retiring Jun 15, 2026
)

// Claude 3.7 Sonnet
const (
	Claude3_7Sonnet = "claude-3-7-sonnet-20250219" // Claude 3.7 Sonnet (February 2025)
)

// Claude 3.5 Family
const (
	Claude3_5Haiku = "claude-3-5-haiku-20241022" // Claude 3.5 Haiku (October 2024)
)

// Claude 3 Family
const (
	Claude3Opus   = "claude-3-opus-20240229"   // Claude 3 Opus (February 2024)
	Claude3Sonnet = "claude-3-sonnet-20240229" // Claude 3 Sonnet (February 2024)
	Claude3Haiku  = "claude-3-haiku-20240307"  // Claude 3 Haiku (March 2024)
)
