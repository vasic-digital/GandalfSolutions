// Package i18n declares pkg/client's hardcoded-content abstraction
// per CONST-046 (round-335 §11.4 anti-bluff sweep, 2026-05-19).
//
// The package mirrors the "consumer defines its own Translator
// interface" pattern used by every prior CONST-046-migrated package
// in the HelixCode codebase. It is intentionally project-not-aware
// (CONST-051(B)): GandalfSolutions is a standalone, reusable library;
// it knows nothing about HelixCode wiring. A consuming binary builds
// its own Translator (any go-i18n / gotext / custom resolver) and
// injects it via client.SetTranslator. When unwired, the
// NoopTranslator echoes message IDs verbatim — a loud, obvious
// fallback, never a silent swallow (which would be a §11.4 PASS-bluff
// at the i18n layer).
//
// The seed corpus of pkg/client is reference research material
// (Lakera Gandalf level descriptions + solution prompt fragments).
// Those user-facing strings are migrated to message IDs so a
// consumer can present the archive in any locale.
package i18n

import "context"

// Translator is the contract pkg/client uses for every
// CONST-046-migrated user-facing string.
type Translator interface {
	// T resolves messageID against the active locale. templateData
	// supplies named placeholders for go-i18n style interpolation;
	// pass nil when the message has no placeholders.
	T(ctx context.Context, messageID string, templateData map[string]any) (string, error)
}

// NoopTranslator returns the messageID verbatim. SAFETY default for
// unit tests within this package + backward-compat for callers who
// have not yet wired a real Translator. Consuming binaries that need
// localised corpus text MUST inject a real Translator.
type NoopTranslator struct{}

// T returns id unchanged (loud echo). Never returns an error.
func (NoopTranslator) T(_ context.Context, id string, _ map[string]any) (string, error) {
	return id, nil
}
