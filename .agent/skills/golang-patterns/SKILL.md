---
name: golang-patterns
description: Idiomatic Go patterns, best practices, and conventions for building robust, efficient, secure, and maintainable Go applications.
---

# Go Development Patterns

Idiomatic Go patterns and best practices for building robust, efficient, secure, and maintainable applications.

## When to Activate

- Writing new Go code
- Reviewing Go code
- Refactoring existing Go code
- Designing Go packages/modules
- Implementing security-critical features

## Core Principles

### 1. Simplicity and Clarity

Go favors simplicity over cleverness. Code should be obvious and easy to read.

```go
// Good: Clear and direct
func GetUser(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("get user %s: %w", id, err)
    }
    return user, nil
}

// Bad: Overly clever
func GetUser(id string) (*User, error) {
    return func() (*User, error) {
        if u, e := db.FindUser(id); e == nil {
            return u, nil
        } else {
            return nil, e
        }
    }()
}
```

### 2. Make the Zero Value Useful

Design types so their zero value is immediately usable without initialization.

```go
// Good: Zero value is useful
type Counter struct {
    mu    sync.Mutex
    count int // zero value is 0, ready to use
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}

// Good: bytes.Buffer works with zero value
var buf bytes.Buffer
buf.WriteString("hello")

// Bad: Requires initialization
type BadCounter struct {
    counts map[string]int // nil map will panic
}
```

### 3. Accept Interfaces, Return Structs

Functions should accept interface parameters and return concrete types.

```go
// Good: Accepts interface, returns concrete type
func ProcessData(r io.Reader) (*Result, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, err
    }
    return &Result{Data: data}, nil
}

// Bad: Returns interface (hides implementation details unnecessarily)
func ProcessData(r io.Reader) (io.Reader, error) {
    // ...
}
```

## Security Best Practices

### Input Validation

```go
// Always validate and sanitize user input
func ValidateUserID(id string) error {
    // Check format
    if len(id) == 0 || len(id) > 100 {
        return fmt.Errorf("invalid user ID length")
    }
    
    // Check for allowed characters (alphanumeric and hyphens)
    matched, err := regexp.MatchString(`^[a-zA-Z0-9-]+$`, id)
    if err != nil {
        return fmt.Errorf("regex validation failed: %w", err)
    }
    if !matched {
        return fmt.Errorf("user ID contains invalid characters")
    }
    
    return nil
}

func GetUser(id string) (*User, error) {
    if err := ValidateUserID(id); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Now safe to use id
    return db.FindUser(id)
}
```

### SQL Injection Prevention

```go
// Good: Use parameterized queries
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = ?"
    
    var user User
    err := db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, fmt.Errorf("query user: %w", err)
    }
    
    return &user, nil
}

// Bad: String concatenation - NEVER DO THIS
func GetUserByEmailBad(db *sql.DB, email string) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE email = '" + email + "'"
    // This is vulnerable to SQL injection!
    // ...
}
```

### Secure Password Handling

```go
import "golang.org/x/crypto/bcrypt"

// Hash passwords before storing
func HashPassword(password string) (string, error) {
    // Use bcrypt with appropriate cost
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("hash password: %w", err)
    }
    return string(hash), nil
}

// Verify passwords securely
func VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword(
        []byte(hashedPassword),
        []byte(password),
    )
}

// Never log passwords or sensitive data
func AuthenticateUser(email, password string) error {
    user, err := GetUserByEmail(email)
    if err != nil {
        // Use generic error message to prevent user enumeration
        return errors.New("invalid credentials")
    }
    
    if err := VerifyPassword(user.PasswordHash, password); err != nil {
        return errors.New("invalid credentials")
    }
    
    return nil
}
```

### Prevent Path Traversal

```go
import "path/filepath"

// Sanitize file paths
func SafeFilePath(basePath, userPath string) (string, error) {
    // Clean the user-provided path
    cleanPath := filepath.Clean(userPath)
    
    // Join with base path
    fullPath := filepath.Join(basePath, cleanPath)
    
    // Ensure result is still within base path
    if !strings.HasPrefix(fullPath, filepath.Clean(basePath)) {
        return "", fmt.Errorf("path traversal attempt detected")
    }
    
    return fullPath, nil
}
```

### Rate Limiting

```go
import (
    "golang.org/x/time/rate"
    "sync"
)

type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rate:     r,
        burst:    b,
    }
}

func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    limiter, exists := rl.limiters[key]
    if !exists {
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.limiters[key] = limiter
    }
    
    return limiter
}

// Usage in HTTP middleware
func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Use IP address as key (consider using user ID for authenticated requests)
            key := r.RemoteAddr
            limiter := rl.GetLimiter(key)
            
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## Error Handling Patterns

### Error Wrapping with Context

```go
// Good: Wrap errors with context
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("load config %s: %w", path, err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config %s: %w", path, err)
    }

    return &cfg, nil
}
```

### Custom Error Types

```go
// Define domain-specific errors
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// Sentinel errors for common cases
var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)
```

### Error Checking with errors.Is and errors.As

```go
func HandleError(err error) {
    // Check for specific error
    if errors.Is(err, sql.ErrNoRows) {
        log.Println("No records found")
        return
    }

    // Check for error type
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        log.Printf("Validation error on field %s: %s",
            validationErr.Field, validationErr.Message)
        return
    }

    // Unknown error
    log.Printf("Unexpected error: %v", err)
}
```

### Never Ignore Errors

```go
// Bad: Ignoring error with blank identifier
result, _ := doSomething()

// Good: Handle or explicitly document why it's safe to ignore
result, err := doSomething()
if err != nil {
    return fmt.Errorf("do something: %w", err)
}

// Acceptable: When error truly doesn't matter (rare)
// Document why it's safe
_ = writer.Close() // Best-effort cleanup, already logged error elsewhere
```

## Concurrency Patterns

### Worker Pool

```go
type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID int
    Value string
    Error error
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        // Process job
        result := Result{
            JobID: job.ID,
            Value: strings.ToUpper(job.Data),
        }
        results <- result
    }
}

func WorkerPool(jobs []Job, numWorkers int) []Result {
    jobChan := make(chan Job, len(jobs))
    resultChan := make(chan Result, len(jobs))
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            worker(workerID, jobChan, resultChan)
        }(i)
    }
    
    // Send jobs
    for _, job := range jobs {
        jobChan <- job
    }
    close(jobChan)
    
    // Wait and close results
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    // Collect results
    results := make([]Result, 0, len(jobs))
    for result := range resultChan {
        results = append(results, result)
    }
    
    return results
}
```

### Context for Cancellation and Timeouts

```go
func FetchWithTimeout(ctx context.Context, url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fetch %s: %w", url, err)
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### Graceful Shutdown

```go
func GracefulShutdown(server *http.Server, logger *log.Logger) error {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    <-quit
    logger.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        // Return error instead of log.Fatalf for better error handling
        return fmt.Errorf("server forced to shutdown: %w", err)
    }

    logger.Println("Server exited gracefully")
    return nil
}
```

### errgroup for Coordinated Goroutines

```go
import "golang.org/x/sync/errgroup"

func FetchAll(ctx context.Context, urls []string) ([][]byte, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([][]byte, len(urls))

    for i, url := range urls {
        i, url := i, url // Capture loop variables
        g.Go(func() error {
            data, err := FetchWithTimeout(ctx, url)
            if err != nil {
                return err
            }
            results[i] = data
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

### Avoiding Goroutine Leaks

```go
// Bad: Goroutine leak if context is cancelled
func leakyFetch(ctx context.Context, url string) <-chan []byte {
    ch := make(chan []byte)
    go func() {
        data := fetchData(url) // Placeholder function
        ch <- data // Blocks forever if no receiver
    }()
    return ch
}

// Good: Properly handles cancellation
func safeFetch(ctx context.Context, url string) <-chan []byte {
    ch := make(chan []byte, 1) // Buffered channel
    go func() {
        data := fetchData(url) // Placeholder function
        select {
        case ch <- data:
        case <-ctx.Done():
            // Context cancelled, don't send
        }
    }()
    return ch
}

// Placeholder for example
func fetchData(url string) []byte {
    return []byte("data")
}
```

## Interface Design

### Small, Focused Interfaces

```go
// Good: Single-method interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Compose interfaces as needed
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

### Define Interfaces Where They're Used

```go
// In the consumer package, not the provider
package service

// UserStore defines what this service needs
type UserStore interface {
    GetUser(id string) (*User, error)
    SaveUser(user *User) error
}

type Service struct {
    store UserStore
}

// Concrete implementation can be in another package
// It doesn't need to know about this interface
```

### Optional Behavior with Type Assertions

```go
type Flusher interface {
    Flush() error
}

func WriteAndFlush(w io.Writer, data []byte) error {
    if _, err := w.Write(data); err != nil {
        return err
    }

    // Flush if supported
    if f, ok := w.(Flusher); ok {
        return f.Flush()
    }
    return nil
}
```

## Package Organization

### Standard Project Layout

```text
myproject/
├── cmd/
│   └── myapp/
│       └── main.go           # Entry point
├── internal/
│   ├── handler/              # HTTP handlers
│   ├── service/              # Business logic
│   ├── repository/           # Data access
│   └── config/               # Configuration
├── pkg/
│   └── client/               # Public API client
├── api/
│   └── v1/                   # API definitions (proto, OpenAPI)
├── testdata/                 # Test fixtures
├── go.mod
├── go.sum
└── Makefile
```

### Package Naming

```go
// Good: Short, lowercase, no underscores
package http
package json
package user

// Bad: Verbose, mixed case, or redundant
package httpHandler
package json_parser
package userService // Redundant 'Service' suffix
```

### Avoid Package-Level State

```go
// Bad: Global mutable state
var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        panic(err) // Avoid panic in init
    }
}

// Good: Dependency injection
type Server struct {
    db *sql.DB
}

func NewServer(db *sql.DB) *Server {
    return &Server{db: db}
}
```

## Struct Design

### Functional Options Pattern

```go
type Server struct {
    addr    string
    timeout time.Duration
    logger  *log.Logger
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) {
        s.timeout = d
    }
}

func WithLogger(l *log.Logger) Option {
    return func(s *Server) {
        s.logger = l
    }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        addr:    addr,
        timeout: 30 * time.Second, // default
        logger:  log.Default(),    // default
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(":8080",
    WithTimeout(60*time.Second),
    WithLogger(customLogger),
)
```

### Embedding for Composition

```go
type Logger struct {
    prefix string
}

func (l *Logger) Log(msg string) {
    fmt.Printf("[%s] %s\n", l.prefix, msg)
}

type Server struct {
    *Logger // Embedding - Server gets Log method
    addr    string
}

func NewServer(addr string) *Server {
    return &Server{
        Logger: &Logger{prefix: "SERVER"},
        addr:   addr,
    }
}

// Usage
s := NewServer(":8080")
s.Log("Starting...") // Calls embedded Logger.Log
```

## Memory and Performance

### Preallocate Slices When Size is Known

```go
type Item struct {
    ID   int
    Name string
}

type Result struct {
    ItemID int
    Value  string
}

func processItem(item Item) Result {
    return Result{
        ItemID: item.ID,
        Value:  strings.ToUpper(item.Name),
    }
}

// Bad: Grows slice multiple times
func processItems(items []Item) []Result {
    var results []Result
    for _, item := range items {
        results = append(results, processItem(item))
    }
    return results
}

// Good: Single allocation
func processItemsOptimized(items []Item) []Result {
    results := make([]Result, 0, len(items))
    for _, item := range items {
        results = append(results, processItem(item))
    }
    return results
}
```

### Use sync.Pool for Frequent Allocations

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessRequest(data []byte) []byte {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()

    buf.Write(data)
    // Process...
    return buf.Bytes()
}
```

### Avoid String Concatenation in Loops

```go
// Bad: Creates many string allocations
func joinBad(parts []string) string {
    var result string
    for _, p := range parts {
        result += p + ","
    }
    return result
}

// Good: Single allocation with strings.Builder
func joinGood(parts []string) string {
    var sb strings.Builder
    for i, p := range parts {
        if i > 0 {
            sb.WriteString(",")
        }
        sb.WriteString(p)
    }
    return sb.String()
}

// Best: Use standard library
func joinBest(parts []string) string {
    return strings.Join(parts, ",")
}
```

## Testing Best Practices

### Table-Driven Tests

```go
func TestValidateUserID(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {
            name:    "valid ID",
            input:   "user-123",
            wantErr: false,
        },
        {
            name:    "empty ID",
            input:   "",
            wantErr: true,
        },
        {
            name:    "too long",
            input:   strings.Repeat("a", 101),
            wantErr: true,
        },
        {
            name:    "invalid characters",
            input:   "user@123",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateUserID(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateUserID() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Use Interfaces for Testability

```go
// Define interface for testing
type UserRepository interface {
    GetUser(id string) (*User, error)
}

// Production implementation
type SQLUserRepository struct {
    db *sql.DB
}

func (r *SQLUserRepository) GetUser(id string) (*User, error) {
    // Real database query
}

// Test implementation
type MockUserRepository struct {
    users map[string]*User
}

func (r *MockUserRepository) GetUser(id string) (*User, error) {
    user, ok := r.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}

// Service uses interface
type UserService struct {
    repo UserRepository
}

// Easy to test with mock
func TestUserService(t *testing.T) {
    mock := &MockUserRepository{
        users: map[string]*User{
            "1": {ID: "1", Name: "Alice"},
        },
    }
    
    service := &UserService{repo: mock}
    
    user, err := service.GetUserByID("1")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Alice" {
        t.Errorf("got name %s, want Alice", user.Name)
    }
}
```

## Go Tooling Integration

### Essential Commands

```bash
# Build and run
go build ./...
go run ./cmd/myapp

# Testing
go test ./...
go test -race ./...
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Static analysis
go vet ./...
staticcheck ./...
golangci-lint run

# Module management
go mod tidy
go mod verify
go mod download

# Formatting
gofmt -w .
goimports -w .

# Security scanning
gosec ./...
```

### Recommended Linter Configuration (.golangci.yml)

```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert
    - unparam
    - gosec        # Security focused
    - gocritic
    - revive

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  govet:
    check-shadowing: true
  gosec:
    excludes:
      - G104 # Audit errors not checked (covered by errcheck)

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

## Quick Reference: Go Idioms

| Idiom                                               | Description                                              |
| --------------------------------------------------- | -------------------------------------------------------- |
| Accept interfaces, return structs                   | Functions accept interface params, return concrete types |
| Errors are values                                   | Treat errors as first-class values, not exceptions       |
| Don't communicate by sharing memory                 | Use channels for coordination between goroutines         |
| Make the zero value useful                          | Types should work without explicit initialization        |
| A little copying is better than a little dependency | Avoid unnecessary external dependencies                  |
| Clear is better than clever                         | Prioritize readability over cleverness                   |
| gofmt is no one's favorite but everyone's friend    | Always format with gofmt/goimports                       |
| Return early                                        | Handle errors first, keep happy path unindented          |
| Validate all inputs                                 | Never trust external data                                |
| Fail securely                                       | Default to secure behavior on errors                     |

## Anti-Patterns to Avoid

```go
// Bad: Naked returns in long functions
func process() (result int, err error) {
    // ... 50 lines ...
    return // What is being returned?
}

// Bad: Using panic for control flow
func GetUser(id string) *User {
    user, err := db.Find(id)
    if err != nil {
        panic(err) // Don't do this - return errors
    }
    return user
}

// Bad: Passing context in struct
type Request struct {
    ctx context.Context // Context should be first param
    ID  string
}

// Good: Context as first parameter
func ProcessRequest(ctx context.Context, id string) error {
    // ...
    return nil
}

// Bad: Mixing value and pointer receivers
type Counter struct{ n int }
func (c Counter) Value() int { return c.n }    // Value receiver
func (c *Counter) Increment() { c.n++ }        // Pointer receiver
// Pick one style and be consistent

// Bad: Ignoring context cancellation
func ProcessData(ctx context.Context, data []Item) error {
    for _, item := range data {
        // Should check ctx.Done() in long-running loops
        process(item)
    }
    return nil
}

// Good: Respect context cancellation
func ProcessDataWithContext(ctx context.Context, data []Item) error {
    for _, item := range data {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := processItem(item); err != nil {
                return err
            }
        }
    }
    return nil
}

func processItem(item Item) error {
    // Process item
    return nil
}
```

## Security Checklist

- [ ] Validate and sanitize all user inputs
- [ ] Use parameterized queries for database operations
- [ ] Hash passwords with bcrypt or argon2
- [ ] Implement rate limiting on public endpoints
- [ ] Use HTTPS/TLS for all network communication
- [ ] Sanitize file paths to prevent directory traversal
- [ ] Set appropriate timeouts on HTTP clients and servers
- [ ] Use context for cancellation and deadlines
- [ ] Never log sensitive information (passwords, tokens, PII)
- [ ] Implement proper authentication and authorization
- [ ] Keep dependencies updated (use `go mod tidy` and security scanners)
- [ ] Use constant-time comparison for sensitive values

**Remember**: Go code should be boring in the best way - predictable, consistent, secure, and easy to understand. When in doubt, keep it simple and prioritize security.