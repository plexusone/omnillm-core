# Refactor: Thin/Thick Provider Architecture

## Overview

Restructure OmniLLM into a layered architecture with thin (stdlib-only) providers in core and thick (official SDK) providers as optional overlays. This supersedes the approach in `FEAT_PROVIDERSDK_PLAN.md`.

## Motivation

### Current State

```
omnillm-core/           # NEW: Interfaces only (zero deps) - TO BE DISCARDED
├── provider.go
├── types.go
└── errors.go

omnillm/                # Batteries-included with thin providers
├── provider/           # Interfaces
├── providers/          # 8 thin (native HTTP) implementations
│   ├── openai/
│   ├── anthropic/
│   ├── gemini/         # Uses google.golang.org/genai (PROBLEM)
│   └── ...
├── client.go           # Multi-provider client
├── memory.go           # Conversation memory
├── cache.go            # Response caching
└── fallback.go         # Fallback logic

omnillm-openai/         # Thick provider (openai-go SDK)
omnillm-anthropic/      # Thick provider (anthropic-sdk-go)
```

### Problems with Current Approach

1. **omnillm-core is too minimal** — Interfaces-only package has limited standalone value
2. **Gemini pollutes omnillm** — `google.golang.org/genai` is pulled in even if user only wants OpenAI
3. **Thin providers are valuable** — Native HTTP implementations have zero deps and work fine for most use cases
4. **Confusing layering** — omnillm has thin providers, separate repos have thick providers, unclear which to use

### Target State

```
omnillm-core/           # RENAMED from omnillm (thin providers + features)
├── provider/           # Interfaces
├── providers/          # 7 thin (native HTTP) implementations
│   ├── openai/         # Thin: stdlib only
│   ├── anthropic/      # Thin: stdlib only
│   ├── xai/            # Thin: stdlib only
│   ├── glm/            # Thin: stdlib only
│   ├── kimi/           # Thin: stdlib only
│   ├── qwen/           # Thin: stdlib only
│   └── ollama/         # Thin: stdlib only
├── client.go           # Multi-provider client
├── memory.go           # Conversation memory
├── cache.go            # Response caching
├── fallback.go         # Fallback logic
└── registry.go         # Provider registry with override support

omnillm-openai/         # Thick provider (openai-go SDK) - EXISTS
omnillm-anthropic/      # Thick provider (anthropic-sdk-go) - EXISTS
omnillm-gemini/         # Thick provider (google genai) - NEW, extracted

omnillm/                # NEW: Batteries-included (core + thick providers)
├── omnillm.go          # Re-exports + thick provider registration
└── go.mod              # Imports core + all thick providers
```

## Design Decisions

### D1: Rename omnillm → omnillm-core

The current `omnillm` becomes `omnillm-core` because:

- It contains the **core** functionality: thin providers + advanced features
- Users who want minimal dependencies use this directly
- Zero external dependencies except grokify/mogo (logging) and grokify/sogo (KVS)

### D2: Discard New omnillm-core (Interfaces-Only)

The recently created interfaces-only `omnillm-core` is discarded because:

- Interfaces-only package has limited standalone value
- Users want working providers, not just interfaces
- Thin providers ARE the core value proposition

### D3: Extract Gemini to omnillm-gemini

Gemini is the only provider requiring a heavy SDK (`google.golang.org/genai`). Extract it so:

- omnillm-core has zero heavy dependencies
- Users who need Gemini explicitly import `omnillm-gemini`
- Gemini is only available as a thick provider (no thin implementation)

### D4: Thick Providers Override Thin

When thick providers are imported, they register themselves to override thin implementations:

```go
// omnillm-core/registry.go
var providers = map[string]ProviderFactory{}

func RegisterProvider(name string, factory ProviderFactory, priority int) {
    if existing, ok := providers[name]; !ok || priority > existing.priority {
        providers[name] = factory
    }
}

// omnillm-core/providers/openai/init.go (thin, priority=0)
func init() {
    RegisterProvider("openai", NewThinProvider, 0)
}

// omnillm-openai/init.go (thick, priority=10)
func init() {
    core.RegisterProvider("openai", NewThickProvider, 10)
}
```

### D5: New omnillm as Aggregator

The new `omnillm` package simply imports everything:

```go
package omnillm

import (
    _ "github.com/plexusone/omnillm-core"
    _ "github.com/plexusone/omnillm-openai"
    _ "github.com/plexusone/omnillm-anthropic"
    _ "github.com/plexusone/omnillm-gemini"
)

// Re-export core types for convenience
type Provider = core.Provider
type Message = core.Message
// etc.
```

## Provider Matrix

| Provider | omnillm-core (thin) | Thick Package | SDK |
|----------|---------------------|---------------|-----|
| OpenAI | ✅ Native HTTP | omnillm-openai | openai-go |
| Anthropic | ✅ Native HTTP | omnillm-anthropic | anthropic-sdk-go |
| Gemini | ❌ None | omnillm-gemini | google genai |
| X.AI | ✅ Native HTTP | — | — |
| GLM | ✅ Native HTTP | — | — |
| Kimi | ✅ Native HTTP | — | — |
| Qwen | ✅ Native HTTP | — | — |
| Ollama | ✅ Native HTTP | — | — |

## Usage Patterns

### Lightweight (Most Users)

```go
import "github.com/plexusone/omnillm-core"

// Uses thin providers, minimal dependencies
client := omnillm.NewClient(omnillm.Config{
    Providers: []string{"openai", "anthropic"},
})
```

### Lightweight + Gemini

```go
import (
    "github.com/plexusone/omnillm-core"
    _ "github.com/plexusone/omnillm-gemini" // Registers gemini provider
)

client := omnillm.NewClient(omnillm.Config{
    Providers: []string{"openai", "gemini"},
})
```

### Full SDK Support

```go
import (
    "github.com/plexusone/omnillm-core"
    _ "github.com/plexusone/omnillm-openai"    // Thick overrides thin
    _ "github.com/plexusone/omnillm-anthropic" // Thick overrides thin
    _ "github.com/plexusone/omnillm-gemini"
)

// OpenAI and Anthropic now use official SDKs
client := omnillm.NewClient(omnillm.Config{
    Providers: []string{"openai", "anthropic", "gemini"},
})
```

### Batteries Included (Simplest)

```go
import "github.com/plexusone/omnillm"

// Everything included, all thick providers active
client := omnillm.NewClient(omnillm.Config{
    Providers: []string{"openai", "anthropic", "gemini"},
})
```

## Implementation Phases

### Phase 1: Prepare omnillm for Rename

**Goal**: Clean up omnillm before renaming to omnillm-core

1. Remove `providers/gemini/` (will become omnillm-gemini)
2. Add provider registry with priority-based override support
3. Update thin providers to register at priority 0
4. Verify all tests pass without gemini
5. Update go.mod to remove `google.golang.org/genai`

**Verification**:

- `go test ./...` passes (skip gemini tests)
- No gemini imports remain
- Registry supports override mechanism

### Phase 2: Create omnillm-gemini

**Goal**: Extract gemini as standalone thick provider

1. Create `github.com/plexusone/omnillm-gemini` repository
2. Move gemini provider code from omnillm
3. Update to import omnillm-core (use replace directive temporarily)
4. Register at priority 10 (thick)
5. Tag v0.1.0

**Verification**:

- Gemini integration tests pass
- Imports omnillm-core correctly

### Phase 3: Rename omnillm → omnillm-core

**Goal**: Rename repository and update module path

1. Rename GitHub repository: omnillm → omnillm-core
2. Update go.mod: `module github.com/plexusone/omnillm-core`
3. Update all internal imports
4. Update README and documentation
5. Tag v0.1.0

**Verification**:

- `go build ./...` succeeds
- `go test ./...` passes
- All 7 thin providers work

### Phase 4: Update Thick Providers

**Goal**: Update omnillm-openai and omnillm-anthropic to import omnillm-core

1. Update omnillm-openai:
   - Change import from omnillm-core (interfaces) to omnillm-core (renamed omnillm)
   - Register at priority 10
   - Remove replace directives
   - Tag new version

2. Update omnillm-anthropic:
   - Same changes as omnillm-openai
   - Tag new version

3. Update omnillm-gemini:
   - Remove replace directive
   - Tag new version

**Verification**:

- All thick providers import omnillm-core correctly
- Override mechanism works (thick replaces thin)

### Phase 5: Create New omnillm

**Goal**: Create batteries-included aggregator package

1. Create new `github.com/plexusone/omnillm` repository
2. Import omnillm-core + all thick providers
3. Re-export types for convenience
4. Add documentation
5. Tag v0.1.0

**Verification**:

- Single import gives all providers
- Thick providers are active by default
- All integration tests pass

### Phase 6: Discard Old omnillm-core

**Goal**: Clean up the interfaces-only package

1. Archive or delete `omnillm-core` (interfaces-only version)
2. Update any documentation references

## Migration Guide

### For Users of Current omnillm

**Before**:
```go
import "github.com/plexusone/omnillm"
```

**After (lightweight)**:
```go
import "github.com/plexusone/omnillm-core"
```

**After (batteries-included, same as before)**:
```go
import "github.com/plexusone/omnillm"
```

### For Users Who Want Gemini

**Before**:
```go
import "github.com/plexusone/omnillm"
// Gemini included automatically
```

**After**:
```go
import (
    "github.com/plexusone/omnillm-core"
    _ "github.com/plexusone/omnillm-gemini"
)
// Or just use omnillm for batteries-included
```

## Dependency Comparison

| Package | Dependencies |
|---------|-------------|
| omnillm-core | grokify/mogo, grokify/sogo |
| omnillm-openai | omnillm-core, openai-go |
| omnillm-anthropic | omnillm-core, anthropic-sdk-go |
| omnillm-gemini | omnillm-core, google genai |
| omnillm | all of the above |

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking existing imports | High | Clear migration guide, type aliases |
| Registry complexity | Medium | Keep registry simple, well-tested |
| Circular dependencies | Medium | Careful interface design |
| Gemini users surprised | Low | Clear documentation |

## Success Criteria

1. **omnillm-core works standalone** — All 7 thin providers functional
2. **Thick override works** — Importing thick provider replaces thin
3. **Gemini extracted** — No google genai in omnillm-core
4. **Backwards compatible** — omnillm import still works
5. **Clear documentation** — Users understand thin vs thick

## File Changes Summary

### Deleted

- `omnillm-core/` (interfaces-only version) — entire repo discarded

### Renamed

- `omnillm/` → `omnillm-core/`

### Created

- `omnillm-gemini/` — extracted from omnillm/providers/gemini
- `omnillm/` (new) — aggregator package

### Modified

- `omnillm-openai/go.mod` — import omnillm-core
- `omnillm-anthropic/go.mod` — import omnillm-core
- `omnillm-core/providers/gemini/` — deleted (moved to omnillm-gemini)
