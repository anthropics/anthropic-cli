package main

import (
	"context"
	"fmt"

	"dagger/ant-cli/internal/dagger"
)

// Factor represents a composable unit of work with lifecycle hooks
// Inspired by Spin SIP 021 - Spin Factors
type Factor interface {
	Name() string
	Dependencies() []string
	Execute(ctx context.Context, state *FactorState) (*dagger.Directory, error)
}

// FactorConfig holds configuration for a factor
type FactorConfig struct {
	Source  *dagger.Directory
	OS      string
	Arch    string
	Options map[string]interface{}
}

// FactorState holds shared state between factors
type FactorState struct {
	Artifacts map[string]*dagger.Directory
	Files     map[string]*dagger.File
	Binaries  map[string]*dagger.File
}

// NewFactorState creates a new factor state
func NewFactorState() *FactorState {
	return &FactorState{
		Artifacts: make(map[string]*dagger.Directory),
		Files:     make(map[string]*dagger.File),
		Binaries:  make(map[string]*dagger.File),
	}
}

// FactorRegistry manages factor composition and execution
type FactorRegistry struct {
	factors map[string]Factor
	config  *FactorConfig
}

// NewFactorRegistry creates a new factor registry
func NewFactorRegistry(config *FactorConfig) *FactorRegistry {
	return &FactorRegistry{
		factors: make(map[string]Factor),
		config:  config,
	}
}

// Register adds a factor to the registry
func (r *FactorRegistry) Register(factor Factor) {
	r.factors[factor.Name()] = factor
}

// ExecuteAll runs all registered factors in dependency order
func (r *FactorRegistry) ExecuteAll(ctx context.Context) (*dagger.Directory, error) {
	state := NewFactorState()
	output := dag.Directory()

	executed := make(map[string]bool)

	for len(executed) < len(r.factors) {
		progress := false

		for name, factor := range r.factors {
			if executed[name] {
				continue
			}

			deps := factor.Dependencies()
			depsSatisfied := true
			for _, dep := range deps {
				if !executed[dep] {
					depsSatisfied = false
					break
				}
			}

			if !depsSatisfied {
				continue
			}

			result, err := factor.Execute(ctx, state)
			if err != nil {
				return nil, fmt.Errorf("factor %s failed: %w", name, err)
			}

			// Keep result lazy — do not Sync here; WithDirectory composes the lazy refs
			output = output.WithDirectory(name, result)
			state.Artifacts[name] = result
			executed[name] = true
			progress = true
		}

		if !progress {
			return nil, fmt.Errorf("circular dependency detected in factors")
		}
	}

	// Single sync at the end materialises the full composed graph
	return output.Sync(ctx)
}
