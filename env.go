// Environment helpers re-exported at the root package for backward
// compatibility with consumers that import the root package (e.g. the
// suresql engine, which predates the utils/ subpackage). New code should
// import github.com/medatechnology/goutil/utils directly.
package goutil

import (
	"time"

	"github.com/medatechnology/goutil/utils"
)

// LoadEnvEach loads environment variables from the given files.
func LoadEnvEach(envFiles ...string) { utils.LoadEnvEach(envFiles...) }

// LoadEnvAll loads and overrides environment variables from the given files.
func LoadEnvAll(envFiles ...string) { utils.LoadEnvAll(envFiles...) }

// ReloadEnvEach reloads environment variables from the given files.
func ReloadEnvEach(additionalFiles ...string) { utils.ReloadEnvEach(additionalFiles...) }

// GetEnv returns the value of the environment variable or the default
// (legacy alias kept for pre-utils consumers such as gosuresql).
func GetEnv(key string, defaultValue string) string {
	return utils.GetEnvString(key, defaultValue)
}

// GetEnvString returns the value of the environment variable or the default.
func GetEnvString(key string, defaultValue string) string {
	return utils.GetEnvString(key, defaultValue)
}

// GetEnvInt returns the integer value of the environment variable or the default.
func GetEnvInt(key string, defaultValue int) int {
	return utils.GetEnvInt(key, defaultValue)
}

// GetEnvBool returns the boolean value of the environment variable or the default.
func GetEnvBool(key string, defaultValue bool) bool {
	return utils.GetEnvBool(key, defaultValue)
}

// GetEnvDuration returns the duration value of the environment variable or the default.
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	return utils.GetEnvDuration(key, defaultValue)
}
