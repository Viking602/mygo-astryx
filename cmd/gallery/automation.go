//go:build e2e

package main

// The local automation server is included only in explicit e2e builds and
// listens only when MYGO_AUTOMATION specifies a socket path.
import _ "github.com/egoist/mygo/plugins/automation"
