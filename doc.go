// Package astryx provides composable native Go UI components for MyGo.
//
// Create a UI builder with New during each view build, or use Scope to apply
// a theme to an entire subtree. Builders, elements and component parts belong
// to that build pass; keep model state and ui.Handle values in your app.
// Stateful components should receive their identity through Context.Key before
// construction. The module path is local until a publishing repository is chosen.
package astryx
