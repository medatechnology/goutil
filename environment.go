package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// This is used for comparing to something, some use only 1 variable (a)
// sometimes we need 2 variables.
// Ex: func(a,b) { return a> 0} -- this means if a > 0 then return true
// .   func(a,b) { return a > b} -- this meansi f a > b then return true
// This is used for comparison ValueOrDefault() function
// Which if the ComparisonFunction returns true, return a, if not return b
type ComparisonFunction[T comparable] func(a, b T) bool

// Also for Setting Options/Configs that are usually taken from environment
// Use this model below as Example
// type ClientConfig struct {
// 	ServerURL        string
// 	APIKey           string
// 	HTTPClientConfig *HTTPClientConfig // Optional HTTP client configuration
// }
// // Set the configuration option, taking the config struct pointer
// type ClientConfigOption func(*ClientConfig)
// // Declare the New Config function to include ...ClientConfigOption
// func NewClientConfig(options ...ClientConfigOption) ClientConfig {
// 	config := ClientConfig{
// 		ServerURL:   utils.GetEnv("SERVER_URL", "http://localhost:8080"),
// 		APIKey:      utils.GetEnv("API_KEY", "development_api_key"),
// 		HTTPTimeout: ValueOrDefault(time.Duration(tmpTimeout)*time.Second, DEFAULT_TIMEOUT, DurationBiggerThanZero),
// 	}
// 	// Run all the options
// 	for _, option := range options {
// 		option(&config)
// 	}
// 	return config
// }
// // Declare the function to set the options, this is great for readibility
// // Set the server URL for client
// func WithServerURL(val string) ClientConfigOption {
// 	return func(config *ClientConfig) {
// 		config.ServerURL = val
// 	}
// }

// // Set the server URL for client
// func WithApiKey(val string) ClientConfigOption {
// 	return func(config *ClientConfig) {
// 		config.APIKey = val
// 	}
// }
// // Set the HTTP client configuration
// func WithHTTPClientConfig(val *HTTPClientConfig) ClientConfigOption {
// 	return func(config *ClientConfig) {
// 		config.HTTPClientConfig = val
// 	}
// }
// -- call the New Config with options
// config := NewClientConfig(
//   WithApiKey("MY_KEY"),
//   WithServerURL("http://localhost:80"),
// )

var (
	// This is for backward compatibility, once all app using this are upgraded
	// remove below aliases.  These functions are depcrecated and renamed
	InitEnv         func(...string)             = LoadEnvEach
	LoadEnv         func(...string)             = LoadEnvAll
	LoadEnvironment func(...string)             = ReloadEnvEach
	GetEnv          func(string, string) string = GetEnvString
)

// Read one by one and error accordingly
func LoadEnvEach(envFiles ...string) {
	for _, f := range envFiles {
		err := godotenv.Load(f)
		if err != nil {
			fmt.Printf("error loading %s environment\n", f)
			return
		}
	}
}

// LoadEnvAll loads environment variables from .env file and optional additional files.
// This will load all the files in the order they are provided, skipping any existing
// variables with the same name. If a file does not exist, it will be ignored without error.
// This is the default behavior of godotenv.Load().
func LoadEnvAll(envFiles ...string) {
	err := godotenv.Load(envFiles...)
	if err != nil {
		fmt.Println("error loading environment:", err)
		return
	}
}

// LoadEnvironment loads environment variables from .env file and optional additional files.
// This use the Overload() means it will overwrite the variables if already exist, Load() won't
func ReloadEnvEach(additionalFiles ...string) {
	// Load default .env file even if it doesn't exist
	godotenv.Overload()
	// if err := godotenv.Overload(); err != nil {
	// 	fmt.Print("No .env file found, proceeding without it")
	// }

	// NOTE: we can just do godotenv.Overload(additionalFiles) and it should be fine,
	//       but if we want to get the individual file that is not exist, we need to do below loop
	// Load additional environment files if provided
	for _, file := range additionalFiles {
		if err := godotenv.Overload(file); err != nil {
			fmt.Printf("Error loading %s: %v", file, err)
		}
	}
}

// GetEnv retrieves an environment variable or returns a default value if not set.
func GetEnvString(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetEnv for integer with option for default value if not set
func GetEnvInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	val, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return int(val)
}

// GetBool retrieves a boolean environment variable or returns a default value
func GetEnvBool(key string, defaultValue bool) bool {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDuration retrieves a duration environment variable or returns a default value
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// If a is bigger than 0
func IntBiggerThanZero(a, b int) bool {
	return a > 0
}

func IntABiggerThanB(a, b int) bool {
	return a > b
}

func IntASmallerThanB(a, b int) bool {
	return a < b
}

func DurationBiggerThanZero(a, b time.Duration) bool {
	return a > 0
}

// Used to get value usually from environment if value is true via ComparisonFunction
// use value, if not use preset (or default value)
// Ex:
// ValueOrDefault(time.Duration(timeout)*time.Second, DEFAULT_TIMEOUT, DurationBiggerThanZero),
//
// Set with maxIdle if bigger than 0 or use DEFAULT_MAX_IDLE
// MaxIdleConns = ValueOrDefault(maxIdle, DEFAULT_MAX_IDLE_CONNECTIONS, IntBiggerThanZero),
//
// Set with maxIdle if it's bigger than DEFAULT_MAX_IDLE
// or use DEFAULT_MAX_IDLE
// MaxIdleConns = ValueOrDefault(maxIdle, DEFAULT_MAX_IDLE, IntABiggerThanB),
func ValueOrDefault[T comparable](value, preset T, compF ComparisonFunction[T]) T {
	if compF(value, preset) {
		return value
	}
	return preset
}
