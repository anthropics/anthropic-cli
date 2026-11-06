package main

import (
	"log/slog"
	"strings"
)

// warnHTTP2Debug logs one WARN line through logger when godebug turns on Go's
// HTTP/2 debug logging. godebug is the value of the GODEBUG environment
// variable.
//
// With the logging on, net/http and golang.org/x/net/http2 log each header of
// every HTTP/2 request the dispatcher sends, values included. They also log up
// to 256 bytes of each body chunk it receives. The dispatcher's environment
// key, its Kubernetes token and work item secrets can then reach standard
// error. At http2debug=2 they also log each header it receives and the start
// of each body chunk it sends. Pod tokens and nonces then reach standard
// error as well.
func warnHTTP2Debug(logger *slog.Logger, godebug string) {
	if !http2DebugOn(godebug) {
		return
	}
	logger.Warn("GODEBUG http2debug is on: the HTTP/2 debug lines on standard error can include credentials from headers and bodies; remove http2debug from GODEBUG wherever real credentials are in use")
}

// http2DebugOn reports whether godebug turns on Go's HTTP/2 debug logging.
// It applies the substring test that net/http and golang.org/x/net/http2
// both apply at startup, for "http2debug=1" or "http2debug=2".
//
// An "http2debug=0" elsewhere in godebug does not undo "http2debug=1".
// "http2debug=10" also matches, as it does in Go.
func http2DebugOn(godebug string) bool {
	return strings.Contains(godebug, "http2debug=1") || strings.Contains(godebug, "http2debug=2")
}
