// translator.go — CONST-046 message-ID resolver seam for pkg/client.
//
// Round-335 §11.4 anti-bluff sweep (2026-05-19). Mirrors the
// "consumer defines its own Translator + tr() helper" pattern used
// by every other CONST-046-migrated package in the HelixCode
// codebase. GandalfSolutions stays fully decoupled (CONST-051(B)):
// this seam carries no project-specific context — a consumer wires
// any Translator implementation it likes.
package client

import (
	"context"

	"digital.vasic.gandalfsolutions/pkg/client/i18n"
)

// translator resolves CONST-046 message IDs for every user-facing
// string emitted by this package (seed-corpus descriptions and
// solution fragments).
//
// The default is the embedded English BundleTranslator so the
// library ships a fully-localised (en) corpus out of the box with no
// consumer wiring. If the embedded bundle ever fails to load, the
// default degrades to i18n.NoopTranslator{} (loud message-ID echo) —
// never a silent empty corpus (§11.4 PASS-bluff guard). A consuming
// binary overrides the locale via SetTranslator.
var translator i18n.Translator = defaultTranslator()

// defaultTranslator builds the embedded-English BundleTranslator,
// degrading loudly to NoopTranslator on any load failure.
func defaultTranslator() i18n.Translator {
	bt, err := i18n.DefaultTranslator()
	if err != nil {
		return i18n.NoopTranslator{}
	}
	return bt
}

// SetTranslator wires a CONST-046-compliant Translator. Passing nil
// resets to the embedded-English default (never silently disables
// translation lookup, which would be a §11.4 PASS-bluff at the i18n
// injection layer). Call before New / NewFromConfig so the seeded
// corpus is resolved against the wired locale.
func SetTranslator(tr i18n.Translator) {
	if tr == nil {
		translator = defaultTranslator()
		return
	}
	translator = tr
}

// tr is the internal CONST-046 resolver used by every user-facing
// string emission in this package. It NEVER returns an error to the
// caller — translation failures degrade to the message ID itself
// (matching NoopTranslator behaviour) so output remains loud +
// obvious instead of silently empty.
func tr(ctx context.Context, msgID string, data map[string]any) string {
	if translator == nil {
		translator = i18n.NoopTranslator{}
	}
	out, err := translator.T(ctx, msgID, data)
	if err != nil || out == "" {
		return msgID
	}
	return out
}
