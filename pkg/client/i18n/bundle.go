// bundle.go — a minimal YAML-bundle-backed Translator for pkg/client.
//
// Round-335 §11.4 CONST-046 sweep (2026-05-19). The Translator
// interface (translator.go) is deliberately tiny so any consumer can
// plug a go-i18n / gotext resolver. BundleTranslator is a batteries-
// included implementation: it reads a go-i18n-shaped YAML bundle
// (messageID -> {other: "text"}) and resolves message IDs against it.
//
// The English bundle is embedded so the library is self-sufficient
// out of the box — calling DefaultTranslator() and wiring it via
// client.SetTranslator yields a fully localised (en) corpus with no
// external file. Additional locales load via LoadBundle.
//
// CONST-051(B): no project-specific context — BundleTranslator works
// for any consumer of GandalfSolutions.
package i18n

import (
	"context"
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed bundles/active.en.yaml
var embeddedBundles embed.FS

// bundleEntry is the go-i18n-shaped value for one message ID.
type bundleEntry struct {
	Other string `yaml:"other"`
}

// BundleTranslator resolves message IDs against an in-memory map
// loaded from one or more YAML bundles. It satisfies Translator.
type BundleTranslator struct {
	messages map[string]string
}

// NewBundleTranslator returns an empty BundleTranslator. Load bundles
// via LoadBundle / LoadBundleBytes before use.
func NewBundleTranslator() *BundleTranslator {
	return &BundleTranslator{messages: make(map[string]string)}
}

// DefaultTranslator returns a BundleTranslator pre-loaded with the
// embedded English (en) bundle. This is the zero-configuration path
// for consumers that just want a localised default corpus.
func DefaultTranslator() (*BundleTranslator, error) {
	bt := NewBundleTranslator()
	data, err := embeddedBundles.ReadFile("bundles/active.en.yaml")
	if err != nil {
		return nil, fmt.Errorf("read embedded en bundle: %w", err)
	}
	if err := bt.LoadBundleBytes(data); err != nil {
		return nil, fmt.Errorf("load embedded en bundle: %w", err)
	}
	return bt, nil
}

// LoadBundleBytes parses a go-i18n-shaped YAML bundle and merges its
// entries into the translator. Later loads overwrite earlier keys.
func (bt *BundleTranslator) LoadBundleBytes(data []byte) error {
	var raw map[string]bundleEntry
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse bundle yaml: %w", err)
	}
	if bt.messages == nil {
		bt.messages = make(map[string]string)
	}
	for id, entry := range raw {
		bt.messages[id] = entry.Other
	}
	return nil
}

// T resolves messageID against the loaded bundle. Unknown IDs return
// an error so tr() degrades to the loud message-ID echo — never a
// silent empty string (§11.4 PASS-bluff guard).
func (bt *BundleTranslator) T(_ context.Context, messageID string, _ map[string]any) (string, error) {
	if bt == nil || bt.messages == nil {
		return "", fmt.Errorf("translator not initialised")
	}
	v, ok := bt.messages[messageID]
	if !ok || v == "" {
		return "", fmt.Errorf("unknown message id %q", messageID)
	}
	return v, nil
}

// Len reports how many message IDs are currently loaded.
func (bt *BundleTranslator) Len() int {
	if bt == nil {
		return 0
	}
	return len(bt.messages)
}
