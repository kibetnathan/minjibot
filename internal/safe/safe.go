// Package safe provides panic-recovery helpers for Discord gateway handlers and
// background goroutines. discordgo does not recover panics raised inside handler
// callbacks, so an unhandled panic in any handler (or in a goroutine it spawns)
// would otherwise crash the entire bot process.
package safe

import (
	"log/slog"
	"runtime/debug"
)

// Recover recovers from a panic in the calling goroutine and logs it with a
// stack trace instead of letting it crash the process. Use it as the first
// deferred call in a handler or goroutine:
//
//	defer safe.Recover(logger, "onMessageCreate")
//
// A nil logger falls back to slog.Default().
func Recover(logger *slog.Logger, where string) {
	if r := recover(); r != nil {
		if logger == nil {
			logger = slog.Default()
		}
		logger.Error("recovered from panic",
			"where", where,
			"panic", r,
			"stack", string(debug.Stack()),
		)
	}
}

// Go runs fn in a new goroutine guarded by Recover, so a panic inside fn is
// logged rather than crashing the process.
func Go(logger *slog.Logger, where string, fn func()) {
	go func() {
		defer Recover(logger, where)
		fn()
	}()
}
