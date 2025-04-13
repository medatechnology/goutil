package simplelog

import (
	"fmt"
	"log"
	"runtime"

	"strings"
)

const (
	LOG_ERROR   = 1
	LOG_NORMAL  = 5
	LOG_VERBOSE = 10

	DB_ARRAY_FIELD_SEPARATOR = "|"

	ERROR_PREFIX   = "!!!"
	WARNING_PREFIX = "---"

	COUNTRY_PHONE_SEPARATOR = "."
	COUNTRY_PHONE_PREFIX    = "+"
)

var (
	DEBUG_LEVEL int = 10 // 0 means no DEBUGGING info, higher number higher info
	// DEBUG_LEVEL     = 5
	DEBUG_PREFIX        string = "#"
	DEBUG_SUFFIX        string = ""
	DEBUG_VALUE         string = ": "
	DEBUG_SEPARATOR     string = ", "
	DEBUG_TRIM_FUNCTION bool   = true // for GetCallerFunction, do we remove/trim the leading filename info etc?
	// JoinAlias        func([]string, string) string = strings.Join
	// Aliases
	LogInfo func(string, int, ...string) = LogInfoStr
)

// No need to use this anymore because it was used only to seed math random generator
// Now since go 1.20 there is no need to do this anymore. If we need to initialize util
// in the future then use this again
func InitUtils() {
	// Cannot call encryption package from here because that package also import utils, which create circular independence
	// encryption.Init()

	// Was for loading env variables and some settings, but already moved that to configs package
	LogInfoAny("Utils", 13, "Utils initializing ... empty")
	LoadDefaults()
}

// NOTE: NOT USED
func LoadDefaults() {
	LogInfoAny("Utils", 13, "Loading defaults ... empty")
}

// Default strings join function with DEBUG_SEPARATOR
func JoinDebug(elems []string) string {
	return strings.Join(elems[:], DEBUG_SEPARATOR)
}

// Simplest log model, just like the default go lang log package or println
// Basically called with function-stack default (-1) the previous caller of this function
// and using the current DEBUG_LEVEL
func LogThis(msgs ...string) {
	caller := GetCallerFunctionName(0)
	LogInfoStr(caller, DEBUG_LEVEL, msgs...)
}

func LogFormat(format string, v ...any) {
	caller := GetCallerFunctionName(0)
	str := fmt.Sprintf(format, v...)
	LogInfoStr(caller, DEBUG_LEVEL, str)
}

func LogAny(msgs ...any) {
	caller := GetCallerFunctionName(0)
	LogInfoAny(caller, DEBUG_LEVEL, msgs...)
}

// Use this to log/print with debugging level
// NOTE: NOT USED cannot remove this for backward compatibility only
// func LogInfo(fn string, level int, msgs ...string) {
// 	if DEBUG_LEVEL >= level {
// 		log.Println(strings.Repeat(DEBUG_PREFIX, level) + fn + ": " + JoinDebug(msgs[:]))
// 	}
// }

// Only log, but not returning anything, message = string
func LogInfoStr(fn string, level int, msgs ...string) {
	if DEBUG_LEVEL >= level {
		log.Println(strings.Repeat(DEBUG_PREFIX, level) + fn + ": " + JoinDebug(msgs[:]))
	}
}

// Only log, but not returning anything, message = any
func LogInfoAny(fn string, level int, msgs ...any) {
	if DEBUG_LEVEL >= level {
		log.Println(strings.Repeat(DEBUG_PREFIX, level)+fn+": ", msgs)
	}
}

// Log MULTIPLE message = string and return error object
// NOTE: NOT USED
func LogInfoStrReturn(fn string, level int, err error, msgs ...string) error {
	LogInfo(fn, level, JoinDebug(msgs[:]))
	return err
}

// Log MULTIPLE messages = any, and return error
func LogInfoAnyReturn(fn string, level int, err error, msgs ...any) error {
	LogInfoAny(fn, level, msgs)
	return err
}

// NOTE: Pocketbase apierrors are not returning the data for security reason
//       To set manual do below
// apiErr := apis.NewNotFoundError("didn't find", err)
// apiErr.Data["something"] = err
// return apiErr

// msgs [0] = message
// msgs [1] = internal log data 'any'
// msgs [2] = api error response data 'any
// NOT USED YET!
// func LogAndReturnForbidden(fn string, level int, msgs ...any) error {
// 	LogInfoAny(fn, level, msgs[1:])
// 	return apis.NewForbiddenError(fn+":", msgs[2:])
// }

// Make console log (more readble code). Remember this doesn't use DEBUG_LEVEL, it always print as long there is error
// NOTE: not printing anything if there is no error
func LogError(fn string, err error, msgs ...string) {
	if err != nil {
		log.Println(ERROR_PREFIX+fn+": "+JoinDebug(msgs[:])+DEBUG_SEPARATOR+"ERROR: ", err)
	}
}

// Make console log (more readble code). Remember this doesn't use DEBUG_LEVEL, it always print as long there is error
// NOTE: not printing anything if there is no error
func LogErrorStr(fn string, err error, msg string) {
	if err != nil {
		log.Println(ERROR_PREFIX+fn+": "+msg+" ERROR: ", err)
	}
}

// Make console log (more readble code). Remember this doesn't use DEBUG_LEVEL, it always print as long there is error
// NOTE: not printing anything if there is no error
func LogErrorAny(fn string, err error, msgs ...any) {
	if err != nil {
		log.Println(ERROR_PREFIX+fn+DEBUG_VALUE, msgs, DEBUG_SEPARATOR+"ERROR: ", err)
	}
}

// Default log error, which we get the caller function-stack default(1)
// NOTE: this disregard DEBUG_LEVEL and only print if there is error (!= nil)
func LogErr(err error, msgs ...any) {
	caller := GetCallerFunctionName(0)
	LogErrorAny(caller, err, msgs...)
	// log.Println(ERROR_PREFIX + fn + ": " + message)
}

// NOT USED
func LogErrorStrReturn(fn string, err error, msg string) error {
	LogErrorStr(fn, err, msg)
	return err
}

// NOTE: Phone has to be informat +[country_code].[phone number without leading 0]
// This convert phone +62.812039124 to username = "62.81231241"
func ConvertPhoneToUsername(phone string) string {
	// if idx := strings.Index(phone, COUNTRY_PHONE_PREFIX); idx != -1 {
	// 	return phone[idx+1:]
	// }
	// return phone
	return strings.ReplaceAll(phone, "+", "")
}

// NOT USED
// NOTE: maybe just removing the + sign, below is wrong!
// Convert username "62.8123512" to phone "+628123027521"
func ConvertUsernameToPhones(uname string) (string, string) {
	if idx := strings.Index(uname, COUNTRY_PHONE_SEPARATOR); idx != -1 {
		return COUNTRY_PHONE_PREFIX + uname[:idx], uname[idx+1:]
	}
	return "", uname
}

// This is the format of E.164 phone: +[country_code][complete phone]
// Our format now is +[country_code].[complete phone]
// complete phone must be all number, cannot contain any symbols
func ConverPhoneToE164(phone string) string {
	return strings.ReplaceAll(phone, ".", "")
}

// To get previous caller need to pass 1 for stack, but since this is used inside the function
// that need to know the caller, we need to have 2.
// To get default, value use 0 as argument
func GetCallerFunctionName(iStack int) string {
	// Get the caller frame from the stack
	if iStack < 1 {
		iStack = 2 // default is, this is called inside a function A that need to get the caller, ie: B, so it's 2 level up in stack
	}
	caller_pc, _, _, ok := runtime.Caller(iStack)
	if !ok {
		return ""
	}

	// Get the function name from the caller frame
	caller_func := runtime.FuncForPC(caller_pc).Name()

	// Extract the function name without the package path
	if DEBUG_TRIM_FUNCTION {
		// From internet, need to trim from last slash
		// last_slash := strings.LastIndexByte(caller_func, '/')
		// if last_slash >= 0 {
		// 	caller_func = caller_func[last_slash+1:]
		// }
		// I think it's better form last dot/period
		last_dot := strings.LastIndexByte(caller_func, '.')
		if last_dot >= 0 {
			caller_func = caller_func[last_dot+1:]
		}
	}
	return caller_func
}
