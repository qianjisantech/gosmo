//go:build pro

package gosmo

// PRO this value indicates if goreplay is running in PRO mode..
// it must not be modified explicitly in production
var PRO = true

// SettingsHook is intentionally left as a no-op
var SettingsHook = func(*AppSettings) {}
