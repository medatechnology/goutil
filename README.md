# GoUtil - Golang Utility Library

[![Go Report Card](https://goreportcard.com/badge/github.com/medatechnology/goutil)](https://goreportcard.com/report/github.com/medatechnology/goutil)
[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/medatechnology/goutil)](go.mod)

GoUtil is a comprehensive collection of utility functions for Go projects. This library aims to reduce boilerplate code and provide common functionality that can be easily imported into any Golang project.

## Table of Contents

- [Installation](#installation)
- [Features](#features)
- [Usage](#usage)
  - [Time and Date Utilities](#time-and-date-utilities)
  - [TTL Map](#ttl-map)
  - [Simple Logging](#simple-logging)
  - [Text Formatting](#text-formatting)
  - [String and Number Handling](#string-and-number-handling)
  - [Object Conversion](#object-conversion)
  - [Error Handling](#error-handling)
  - [HTTP Client](#http-client)
  - [Filesystem Operations](#filesystem-operations)
  - [Environment Variables](#environment-variables)
  - [JWT and Encryption](#jwt-and-encryption)
  - [Random Generators](#random-generators)
  - [Performance Metrics](#performance-metrics)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/medatechnology/goutil
```

## Features

GoUtil includes the following key features:

- **Time/Date Handling**: Extensive formatting, conversion, and manipulation functions
- **TTL Map**: In-memory key-value store with automatic expiration (like a simple Redis)
- **Logging**: Simple, configurable logging utilities
- **Text Formatting**: Console text coloring and box formatting
- **Type Handling**: Generic number and string manipulation functions
- **Object Conversion**: Struct to map and map to struct conversions
- **Error Handling**: Customizable error types
- **HTTP Client**: Simplified HTTP operations
- **Filesystem**: Common file operations
- **Environment**: Environment variable handling
- **JWT & Encryption**: Auth and encryption utilities
- **Random Generation**: Secure random number and token generators

## Usage

### Time and Date Utilities

The `timedate` package provides extensive time formatting and manipulation functions:

```go
import "github.com/medatechnology/goutil/timedate"

// Converting time to string
timeStr := timedate.ConvertGoLangTimeToPBString(time.Now())

// Converting string to time
timeObj := timedate.ConvertDateTimeStringToGoLangTime("2024-04-03 15:30:45.000Z")

// Format duration as human-readable
duration := time.Hour*24 + time.Minute*45
humanReadable := timedate.DurationHumanReadableLong(duration)
// Output: "1 day, 45 minutes"

// Get days between dates
days := timedate.GetDaysFromDateTimeBToA(startTime, endTime)
```

### TTL Map

The `medattlmap` package provides an in-memory key-value store with automatic expiration:

```go
import (
    "github.com/medatechnology/goutil/medattlmap"
    "time"
)

// Create TTL Map with 5 second default expiration
ttlMap := medattlmap.NewTTLMap(5*time.Second, 1*time.Second)

// Store value with default expiration
ttlMap.Put("key1", 0, "value1")

// Store value with custom expiration
ttlMap.Put("key2", 10*time.Second, "value2")

// Get value
if val, ok := ttlMap.Get("key1"); ok {
    fmt.Println("Value:", val.(string)) // Output: Value: value1
}

// Wait for expiration
time.Sleep(6 * time.Second)

// Key will be automatically removed
if _, ok := ttlMap.Get("key1"); !ok {
    fmt.Println("Key expired") // Output: Key expired
}

// Don't forget to stop the cleanup goroutine
defer ttlMap.Stop()
```

### Simple Logging

The `simplelog` package provides easy-to-use logging functions:

```go
import "github.com/medatechnology/goutil/simplelog"

// Set debug level (0 = no debug, higher number = more verbosity)
simplelog.DEBUG_LEVEL = 5

// Basic logging with auto-detected caller function name
simplelog.LogThis("This is a log message")
// Output: # FunctionName: This is a log message

// Formatted logging
simplelog.LogFormat("Value: %d, String: %s", 42, "test")
// Output: # FunctionName: Value: 42, String: test

// Log with explicit function name and level
simplelog.LogInfoStr("MyFunction", 3, "Processing item", "ID: 123")
// Output: ### MyFunction: Processing item, ID: 123

// Error logging (only prints if err != nil)
err := someFunction()
simplelog.LogErr(err, "Failed during processing")
// Output (if error): !!! FunctionName: Failed during processing ERROR: error details
```

### Text Formatting

The `print` package provides text formatting utilities including colored console output:

```go
import "github.com/medatechnology/goutil/print"

// Print colored text
print.PrintlnColor("Success!", print.ColorGreen)
print.PrintfColor("Number: %d\n", print.ColorCyan, 42)

// Create a formatted box with heading and key-value pairs
heading := []string{"User Information", "Details"}
headingColors := []print.Color{print.ColorBlue, print.ColorCyan}

content := []print.KeyValue{
    print.Content(false, false, "Username", "john_doe"),
    print.Content(false, false, "Email", "john@example.com"),
    print.Content(false, false, "Status", "Active"),
    print.Content(false, false, "Role", "Admin"),
    print.Content(true, false, "Notes", "VIP customer with premium subscription"),
}

print.PrintBoxHeadingContent(heading, headingColors, content, print.ColorYellow, print.ColorWhite)
```

Output:
```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              User Information                                 │
│                                  Details                                      │
│                                                                              │
│ Username : john_doe               Email : john@example.com                   │
│ Status : Active                   Role : Admin                               │
│ Notes : VIP customer with premium subscription                               │
└──────────────────────────────────────────────────────────────────────────────┘
```

### String and Number Handling

Generic type functions for string and number operations:

```go
import (
    "github.com/medatechnology/goutil/object"
)

// Absolute value for any number type
absInt := object.Abs(-42)    // 42
absFloat := object.Abs(-3.14) // 3.14

// Fast string to integer conversion
num := object.Int("1,234.56", false) // 1234

// String operations
lastPart := object.LastStringAfterDelimiter("path/to/file.txt", "/") // "file.txt"
combined := object.CombineToStringWithSpace("Hello", "World", 42) // "Hello World 42"

// Check if array contains string
found := object.ArrayAContainsBString([]string{"apple", "orange", "banana"}, "orange") // true
```

### Object Conversion

The `object` package provides utility functions for object conversions:

```go
import "github.com/medatechnology/goutil/object"

type User struct {
    ID        int       `json:"id" db:"user_id"`
    Name      string    `json:"name" db:"user_name"`
    Email     string    `json:"email" db:"email"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    Active    bool      `json:"active" db:"is_active"`
}

user := User{
    ID:        123,
    Name:      "John Doe",
    Email:     "john@example.com",
    CreatedAt: time.Now(),
    Active:    true,
}

// Convert struct to map using JSON tags
userMap := object.StructToMapJSON(user)
// Output: map[string]interface{}{
//   "id": 123,
//   "name": "John Doe",
//   "email": "john@example.com",
//   "created_at": "2024-04-14T12:00:00Z",
//   "active": true,
// }

// Convert struct to map using DB tags
dbMap := object.StructToMapDB(user)
// Output: map[string]interface{}{
//   "user_id": 123,
//   "user_name": "John Doe",
//   "email": "john@example.com",
//   "created_at": "2024-04-14T12:00:00Z",
//   "is_active": true,
// }

// Custom options for conversion
opts := object.DefaultNoSkipMapOptions() // Don't skip zero values
opts.SkipFalseBools = true  // Skip false boolean values
customMap := object.StructToMapWithOptions(user, opts)

// Convert map back to struct
newUser := object.MapToStruct[User](userMap)
```

### Error Handling

The `medaerror` package provides a customizable error type:

```go
import "github.com/medatechnology/goutil/medaerror"

// Create simple error
err := medaerror.NewString("File not found")
fmt.Println(err.Error()) // "File not found"

// Create formatted error
err = medaerror.Errorf("Invalid ID: %d", 123)

// Create error with code and data
err = medaerror.NewMedaErr(404, "Resource not found", "The requested resource does not exist", map[string]string{"id": "123"})

// Get different representations
fmt.Println(err.Error())    // "Resource not found"
fmt.Println(err.Response()) // "The requested resource does not exist"
fmt.Println(err.String())   // "code:404;message:Resource not found;response:The requested resource does not exist;data:map[id:123]"
```

### HTTP Client

The `httpclient` package provides a simple HTTP client:

```go
import "github.com/medatechnology/goutil/httpclient"

client := httpclient.NewHttp()

// Configure headers
client.SetHeader(map[string][]string{
    "Content-Type": {"application/json"},
    "Authorization": {"Bearer token123"},
})

// Set query parameters
client.SetQueryParams(map[string]string{
    "page": "1",
    "limit": "10",
})

// Set basic auth
client.SetBasicAuth("username", "password")

// GET request
var response struct {
    Data []string `json:"data"`
}
statusCode, err := client.Get("https://api.example.com/items", &response, nil)
if err == nil && statusCode == 200 {
    fmt.Println("Data:", response.Data)
}

// POST request
payload := map[string]interface{}{
    "name": "New Item",
    "quantity": 5,
}
statusCode, err = client.Post("https://api.example.com/items", payload, &response, nil)
```

### Filesystem Operations

The `filesystem` package provides common file operations:

```go
import "github.com/medatechnology/goutil/filesystem"

// Get file name from path
name := filesystem.FileName("/path/to/file.txt") // "file.txt"

// Get directory path
dir := filesystem.DirPath("/path/to/file.txt") // "/path/to"

// List files in a directory
files := filesystem.Dir("/path/to", ".go") // Lists all .go files

// Check if file exists
exists := filesystem.FileExists("/path/to/file.txt")

// Check if directory exists
isDir := filesystem.DirFileExist("/path/to", true)

// Read text file content
lines := filesystem.More("/path/to/file.txt") // Returns []string with each line
```

### Environment Variables

The `utils` package provides environment variable handling:

```go
import "github.com/medatechnology/goutil/utils"

// Load .env files
utils.LoadEnvEach(".env", ".env.local")

// Get environment variables with defaults
port := utils.GetEnvInt("PORT", 8080)
debug := utils.GetEnvBool("DEBUG", false)
apiKey := utils.GetEnvString("API_KEY", "default-key")
timeout := utils.GetEnvDuration("TIMEOUT", 30*time.Second)

// Value or default with comparison function
maxConnections := utils.ValueOrDefault(maxConn, 100, utils.IntBiggerThanZero)
```

### JWT and Encryption

The `encryption` package provides JWT and encryption utilities:

```go
import "github.com/medatechnology/goutil/encryption"

// Parse JWT from auth header
authType, token := encryption.GetAuthorizationFromHeader("Bearer eyJhbGciOiJ...")

// Get claims from JWT token
claims, err := encryption.GetJWTClaimMapFromTokenString(token, "your-jwt-secret")

// Hash a password
hashedPin, err := encryption.HashPin("1234", "salt1", "salt2")

// Encrypt/decrypt data
encrypted, err := encryption.EncryptWithKey("sensitive data", "16-character-key!")
decrypted, err := encryption.DecryptWithKey(encrypted, "16-character-key!")

// JWE encryption
jweToken, err := encryption.CreateJWE([]byte(`{"user_id": 123}`), []byte("32-character-encryption-key-here!"))
plaintext, err := encryption.ParseJWE(jweToken, []byte("32-character-encryption-key-here!"))
```

### Random Generators

The `encryption` package also provides secure random generators:

```go
import "github.com/medatechnology/goutil/encryption"

// Generate secure random number
otp, err := encryption.GenerateSecureRandomNumber(6) // "847291"

// Generate formatted OTP
formattedOtp := encryption.GenerateOTP(3, 2, "-") // "847-291"

// Generate random token (short UUID)
token := encryption.NewRandomToken() // "KwSysDpxcBU9FNhGkn2dCf"

// Generate longer random token
longToken := encryption.NewRandomTokenIterate(3) // Concatenates 3 random tokens
```

### Performance Metrics

The `metrics` package provides a stopwatch utility for measuring execution time of operations:

```go
import "github.com/medatechnology/goutil/metrics"

// Start a timer with message and visual ticker (dots every 250ms)
watch := metrics.StartTimeIt("Processing data", 250)
// Output starts: Processing data....

// Do some work
time.Sleep(1 * time.Second)

// Stop timer and print message with elapsed time
elapsed := metrics.StopTimeItPrint(watch, "Complete")
// Output: Processing data.... Complete (1.002s)

// For silent timing (no output during execution)
watch = metrics.StartTimeIt("", 0)

// Do work
complexOperation()

// Get elapsed time without printing
elapsed = metrics.StopTimeIt(watch)
fmt.Printf("Operation took %v\n", elapsed)

// Or print custom message with elapsed time
watch = metrics.StartTimeIt("", 0)
operation()
elapsed = metrics.StopTimeItPrint(watch, "Operation finished")
// Output: Operation finished (243ms)
```

The stopwatch utility is thread-safe and can be used to track multiple concurrent operations.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.