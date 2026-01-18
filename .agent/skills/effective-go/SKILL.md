---
name: Effective Go
description: "Apply Go best practices, idioms, and conventions from golang.org/doc/effective_go. Use when writing, reviewing, or refactoring Go code to ensure idiomatic, clean, and efficient implementations."
---

# Effective Go

Apply best practices and conventions from the official [Effective Go guide](https://go.dev/doc/effective_go) to write clean, idiomatic Go code.

## When to Apply

Use this skill automatically when:
- Writing new Go code
- Reviewing Go code
- Refactoring existing Go implementations
- Debugging Go-specific issues

## Key Principles

Follow the conventions and patterns documented at https://go.dev/doc/effective_go, with particular attention to:

### Code Organization
- **Package naming**: Short, lowercase, single-word names (no underscores)
- **File structure**: Group related functionality; use `internal/` for private packages
- **Imports**: Group standard library, third-party, and local imports separately

### Formatting & Naming
- **Formatting**: Always use `gofmt` - this is non-negotiable
- **Naming**: No underscores; use MixedCaps for exported, mixedCaps for unexported
- **Getters/Setters**: Prefer `Owner()` over `GetOwner()`; use `SetOwner()` for setters

### Error Handling
- **Always check errors**: Never ignore error returns
- **Error context**: Wrap errors with context using `fmt.Errorf("context: %w", err)`
- **Don't panic**: Reserve `panic` for truly unrecoverable errors
- **Named returns**: Use sparingly, primarily for documentation or deferred cleanup

### Concurrency
- **Channels over locks**: Share memory by communicating (use channels)
- **Goroutine lifecycle**: Always ensure goroutines can exit (avoid leaks)
- **Context usage**: Pass `context.Context` as first parameter for cancellation

### Interfaces & Types
- **Small interfaces**: Keep to 1-3 methods (ideal); single-method interfaces are common
- **Accept interfaces, return structs**: Maximize flexibility in inputs, clarity in outputs
- **Embed for composition**: Prefer embedding over inheritance-style patterns

### Documentation
- **Document exports**: All exported symbols must have doc comments
- **Start with name**: Begin doc comments with the symbol name
- **Package docs**: Include package-level documentation in `doc.go`

### Performance & Idioms
- **Zero values**: Design types so their zero value is useful
- **Defer for cleanup**: Use `defer` for resource cleanup (files, locks, etc.)
- **Blank identifier**: Use `_` to explicitly ignore values or enforce interface implementation
- **Init functions**: Use sparingly; prefer explicit initialization

## Common Pitfalls to Avoid

- Copying mutexes (use pointers to structs with mutexes)
- Range loop variable capture in goroutines (create new variable or pass as parameter)
- Ignoring errors from `Close()`, `Write()`, etc.
- Using `init()` for complex initialization (prefer constructor functions)
- Over-using pointers (Go passes small structs efficiently)

## References

- Official Guide: https://go.dev/doc/effective_go
- Code Review Comments: https://github.com/golang/go/wiki/CodeReviewComments
- Standard Library: Use as reference for idiomatic patterns
- Go Proverbs: https://go-proverbs.github.io/

## Example Pattern

When suggesting code improvements, provide before/after examples:

**Before:**
```go
func getUser() (*User, error) {
    // non-idiomatic code
}
```

**After:**
```go
func User() (*User, error) {
    // idiomatic code with explanation
}
```