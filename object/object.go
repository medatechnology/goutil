package object

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/medatechnology/goutil/timedate"
)

var (
	DEBUG_KEY_SEPARATOR string = ": "
)

// MapOptions defines options for structure to map conversion
type MapOptions struct {
	SkipZeroNumbers  bool   // Skip integers/floats that are 0
	SkipEmptyStrings bool   // Skip strings that are empty
	SkipFalseBools   bool   // Skip booleans that are false
	SkipZeroTimes    bool   // Skip time.Time that are zero (Jan 1, 0001 UTC)
	SkipEmptySlices  bool   // Skip empty slices and arrays
	SkipEmptyMaps    bool   // Skip empty maps
	SkipNilPointers  bool   // Skip nil pointers (always true, included for completeness)
	TimeFormat       string // Format to use for time.Time values (empty string uses time.RFC3339)
}

// DefaultSkipMapOptions returns the default options for map conversion
func DefaultSkipMapOptions() MapOptions {
	return MapOptions{
		SkipZeroNumbers:  true,
		SkipEmptyStrings: true,
		SkipFalseBools:   false, // Default to including false booleans
		SkipZeroTimes:    true,
		SkipEmptySlices:  true,
		SkipEmptyMaps:    true,
		SkipNilPointers:  true, // Always true
		TimeFormat:       time.RFC3339,
	}
}

// Can also use this for non-skipping the zero and all, usually we want to set 0 to
// the integer like saving in DB, so usage would be:
//
// resultMap := StructToMapDBWithOptions[StructName](struct_object, DefaultNoSkipMapOptions())
//
// DefaultSkipMapOptions returns the default options for map conversion
func DefaultNoSkipMapOptions() MapOptions {
	return MapOptions{
		SkipZeroNumbers:  false,
		SkipEmptyStrings: false,
		SkipFalseBools:   false, // Default to including false booleans
		SkipZeroTimes:    true,
		SkipEmptySlices:  true,
		SkipEmptyMaps:    true,
		SkipNilPointers:  true, // Always true
		TimeFormat:       time.RFC3339,
	}
}

// Get json only tag
// First try to get JSON tag, if not exist get the DB tag.
// If not exist, return empty string.
func GetJSONTag(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	tagName := strings.Split(tag, ",")[0] // Extract the tag name before any comma
	return tagName
}

func GetDBTag(field reflect.StructField) string {
	tag := field.Tag.Get("db")
	tagName := strings.Split(tag, ",")[0] // Extract the tag name before any comma
	return tagName
}

// Get json or db tag from a struct.
// First try to get JSON tag, if not exist get the DB tag.
// If not exist, return empty string.
func GetJSONOrDBTag(field reflect.StructField) string {
	tagName := GetJSONTag(field)
	if tagName == "" {
		tagName = GetDBTag(field)
	}
	return tagName
}

// structToMapWithOptions converts a struct to a map using provided options
func structToMapWithOptions[T any](input T, tagfunc func(reflect.StructField) string, opts MapOptions) map[string]interface{} {
	result := make(map[string]interface{})

	inputValue := reflect.ValueOf(input)
	// If it's a pointer, dereference it
	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}
	inputType := inputValue.Type()

	if inputType.Kind() == reflect.Struct {
		for i := 0; i < inputType.NumField(); i++ {
			field := inputType.Field(i)
			tagName := tagfunc(field)
			if tagName != "" {
				value := inputValue.Field(i)

				// Handle different types with special cases
				switch value.Kind() {
				case reflect.Ptr:
					// Always skip nil pointers
					if value.IsNil() {
						continue
					}
					result[tagName] = reflect.Indirect(value).Interface()

				case reflect.String:
					if opts.SkipEmptyStrings && value.String() == "" {
						continue
					}
					result[tagName] = value.String()

				case reflect.Struct:
					// Special handling for time.Time
					if value.Type() == reflect.TypeOf(time.Time{}) {
						timeValue := value.Interface().(time.Time)
						if opts.SkipZeroTimes && timeValue.IsZero() {
							continue
						}

						// Format time according to options
						format := opts.TimeFormat
						if format == "" {
							format = time.RFC3339
						}
						result[tagName] = timeValue.UTC().Format(format)
					} else {
						// For other structs, check if they implement Stringer
						if stringer, ok := value.Interface().(fmt.Stringer); ok {
							result[tagName] = stringer.String()
						} else {
							// Convert other structs to maps recursively
							nestedMap := structToMapWithOptions(value.Interface(), tagfunc, opts)
							if len(nestedMap) > 0 || !opts.SkipEmptyMaps {
								result[tagName] = nestedMap
							}
						}
					}

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if opts.SkipZeroNumbers && value.Int() == 0 {
						continue
					}
					result[tagName] = value.Int()

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if opts.SkipZeroNumbers && value.Uint() == 0 {
						continue
					}
					result[tagName] = value.Uint()

				case reflect.Float32, reflect.Float64:
					if opts.SkipZeroNumbers && value.Float() == 0 {
						continue
					}
					result[tagName] = value.Float()

				case reflect.Bool:
					if opts.SkipFalseBools && !value.Bool() {
						continue
					}
					result[tagName] = value.Bool()

				case reflect.Slice, reflect.Array:
					if opts.SkipEmptySlices && value.Len() == 0 {
						continue
					}
					result[tagName] = value.Interface()

				case reflect.Map:
					if opts.SkipEmptyMaps && value.Len() == 0 {
						continue
					}
					result[tagName] = value.Interface()

				default:
					// For any other types, include regardless of value
					result[tagName] = value.Interface()
				}
			}
		}
	}
	return result
}

// Updated functions with options parameter for more flexibility
//
// Using default options (skips zeros, empty strings, empty slices, etc.)
// userMap := StructToMap(user)
//
// Using custom options
// opts := DefaultMapOptions()
// opts.SkipFalseBools = true       // Skip false booleans
// opts.SkipZeroNumbers = false     // Include zero numbers
// opts.TimeFormat = time.RFC1123   // Use a different time format
//
// customMap := StructToMapWithOptions(user, opts)

// StructToMap converts a struct to a map using generics with JSON or DB tags and default options
func StructToMap[T any](input T) map[string]interface{} {
	return structToMapWithOptions(input, GetJSONOrDBTag, DefaultSkipMapOptions())
}

// StructToMapWithOptions converts a struct to a map using generics with JSON or DB tags and custom options
func StructToMapWithOptions[T any](input T, opts MapOptions) map[string]interface{} {
	return structToMapWithOptions(input, GetJSONOrDBTag, opts)
}

// StructToMapDB converts a struct to a map with DB tags and default options
func StructToMapDB[T any](input T) map[string]interface{} {
	return structToMapWithOptions(input, GetDBTag, DefaultSkipMapOptions())
}

// StructToMapDBWithOptions converts a struct to a map with DB tags and custom options
func StructToMapDBWithOptions[T any](input T, opts MapOptions) map[string]interface{} {
	return structToMapWithOptions(input, GetDBTag, opts)
}

// StructToMapJSON converts a struct to a map with JSON tags and default options
func StructToMapJSON[T any](input T) map[string]interface{} {
	return structToMapWithOptions(input, GetJSONTag, DefaultSkipMapOptions())
}

// StructToMapJSONWithOptions converts a struct to a map with JSON tags and custom options
func StructToMapJSONWithOptions[T any](input T, opts MapOptions) map[string]interface{} {
	return structToMapWithOptions(input, GetJSONTag, opts)
}

// func structToMapUniversal[T any](input T, tagfunc func(reflect.StructField) string) map[string]interface{} {
// 	result := make(map[string]interface{})

// 	inputValue := reflect.ValueOf(input)
// 	// If it's a pointer, dereference it
// 	if inputValue.Kind() == reflect.Ptr {
// 		inputValue = inputValue.Elem()
// 	}
// 	inputType := inputValue.Type()

// 	if inputType.Kind() == reflect.Struct {
// 		for i := 0; i < inputType.NumField(); i++ {
// 			field := inputType.Field(i)
// 			// tagName := GetJSONOrDBTag(field)
// 			tagName := tagfunc(field)
// 			if tagName != "" {
// 				value := inputValue.Field(i)
// 				// NOTE: modified this to get the kind, because sometimes we have struct that
// 				//       has pointer member and if nil, then we skip creating the map! Same with
// 				//       empty string.  This then we can use like *bool and *int to denote
// 				//       empty-value then it won't be created as map[key]=default_value
// 				//       which is 0 for int and false for bool.
// 				switch value.Kind() {
// 				case reflect.Ptr:
// 					if !value.IsNil() {
// 						result[tagName] = reflect.Indirect(value).Interface()
// 					}
// 				case reflect.String:
// 					if value.String() != "" {
// 						result[tagName] = value.String()
// 					}
// 				case reflect.Struct:
// 					// try to call the fmt.Stringer first, the String() function
// 					if a, ok := fmt.Print(field); ok != nil {
// 						result[tagName] = a
// 					} else if a, ok = fmt.Print(value); ok != nil {
// 						result[tagName] = a
// 					} else {
// 						result[tagName] = fmt.Sprintf("%v", value)
// 					}
// 				default:
// 					// value := inputValue.Field(i).Interface()
// 					result[tagName] = value
// 				}
// 			}
// 		}
// 	}
// 	return result
// }

// // StructToMap converts a struct to a map using generics using `json` tags or DB tag
// func StructToMap[T any](input T) map[string]interface{} {
// 	return structToMapUniversal(input, GetJSONOrDBTag)
// }

// // StructToMap that convert struct to a map with key defined in `db` tag, if not exist then skip the field
// func StructToMapDB[T any](input T) map[string]interface{} {
// 	return structToMapUniversal(input, GetDBTag)
// }

// // StructToMap that convert struct to a map with key defined in `json` tag, if not exist then skip the field
// func StructToMapJSON[T any](input T) map[string]interface{} {
// 	return structToMapUniversal(input, GetJSONTag)
// }

// This can be written using reflect and loop through the fields and json tag.
// But marshall works as well
func MapToStruct[T any](dict map[string]interface{}) T {
	var fStruct T
	fType := reflect.TypeOf(fStruct).String()
	jsonbody, err := json.Marshal(dict)
	if err != nil {
		// do error check
		log.Println("!!!"+fType+".FromMap: marshal error - ERROR:", err)
		return fStruct
	}
	if err := json.Unmarshal(jsonbody, &fStruct); err != nil {
		// do error check
		log.Println("!!!"+fType+".FromMap: unmarshall error - ERROR:", err)
		return fStruct
	}
	return fStruct
}

// MapToStruct that convert map to struct with key defined in `json` tag
func MapToStructSlow[T any](dict map[string]interface{}) T {
	return mapToStructUniversal[T](dict, GetJSONOrDBTag)
}

// MapToStruct that convert map to struct with key defined in `json` tag
func MapToStructSlowDB[T any](dict map[string]interface{}) T {
	return mapToStructUniversal[T](dict, GetDBTag)
}

// MapToStruct that convert map to struct with key defined in `json` tag
func MapToStructSlowJSON[T any](dict map[string]interface{}) T {
	return mapToStructUniversal[T](dict, GetJSONTag)
}

// MapToStruct converts a map to a struct, handling time.Time and nested structs.
func mapToStructUniversal[T any](dict map[string]interface{}, tagfunc func(reflect.StructField) string) T {
	var result T
	rv := reflect.ValueOf(&result).Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)
		tagName := tagfunc(field)

		if tagName == "" {
			tagName = field.Name
		}

		if value, ok := dict[tagName]; ok {
			if fieldValue.CanSet() {
				setFieldValue(fieldValue, value, tagfunc)
			}
		}
	}

	return result
}

// Utility function for comprehensive and flexible MapToStruct
// setFieldValue sets the value of a struct field.
func setFieldValue(fieldValue reflect.Value, value interface{}, tagfunc func(reflect.StructField) string) {
	switch fieldValue.Kind() {
	case reflect.Ptr:
		if value == nil {
			fieldValue.Set(reflect.Zero(fieldValue.Type()))
		} else {
			newValue := reflect.New(fieldValue.Type().Elem())
			setFieldValue(newValue.Elem(), value, tagfunc)
			fieldValue.Set(newValue)
		}
	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			// Handle time.Time fields
			if str, ok := value.(string); ok {
				parsedTime, err := timedate.ParseStringnToTime(str)
				if err != nil {
					log.Printf("Error parsing time: %v", err)
				} else {
					fieldValue.Set(reflect.ValueOf(parsedTime))
				}
			}
		} else {
			// Handle nested structs
			if nestedMap, ok := value.(map[string]interface{}); ok {
				nestedStruct := reflect.New(fieldValue.Type()).Interface()
				populateStruct(nestedMap, nestedStruct, tagfunc)
				fieldValue.Set(reflect.ValueOf(nestedStruct).Elem())
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Handle int types
		switch v := value.(type) {
		case float64:
			fieldValue.SetInt(int64(v)) // Convert float64 to int
		case int:
			fieldValue.SetInt(int64(v))
		case int64:
			fieldValue.SetInt(v)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Handle uint types
		switch v := value.(type) {
		case float64:
			fieldValue.SetUint(uint64(v)) // Convert float64 to uint
		case int:
			fieldValue.SetUint(uint64(v))
		case uint64:
			fieldValue.SetUint(v)
		}
	case reflect.Float32, reflect.Float64:
		// Handle float types
		if v, ok := value.(float64); ok {
			fieldValue.SetFloat(v)
		}
	case reflect.Bool:
		// Handle bool
		if v, ok := value.(bool); ok {
			fieldValue.SetBool(v)
		}
	case reflect.String:
		// Handle string
		if v, ok := value.(string); ok {
			fieldValue.SetString(v)
		}
	}
}

// populateStruct populates a struct using a map.
func populateStruct(dict map[string]interface{}, result interface{}, tagfunc func(reflect.StructField) string) {
	rv := reflect.ValueOf(result).Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)
		tagName := tagfunc(field)

		if tagName == "" {
			tagName = field.Name
		}

		if value, ok := dict[tagName]; ok {
			if fieldValue.CanSet() {
				setFieldValue(fieldValue, value, tagfunc)
			}
		}
	}
}

// Combines map A,B into A without checking if A already have the key or not, just overwrite
func CombineTwoMaps(a, b map[string]interface{}) {
	// for key, value := range b {
	//      a[key] = value
	// }
	CombineTwoMapsDistinct(true, a, b)
}

// If bwins is set means b will overwrite a for the same key
func CombineTwoMapsDistinct(bwins bool, a, b map[string]interface{}) {
	for key, value := range b {
		_, ok := a[key]
		if !ok || bwins {
			// If the not key exists
			a[key] = value
		}
	}
}

// NOTE: For debugging usually
// Use this to log/print map into key: value array of string
func MapToStringSlice(d map[string]interface{}) []string {
	flat := []string{}
	for key, value := range d {
		flat = append(flat, key+string(DEBUG_KEY_SEPARATOR)+fmt.Sprintf("%v", value))
	}
	return flat
}

// NOTE: Not USED!
func MapToString(d map[string]interface{}) string {
	flat := ""
	for key, value := range d {
		flat += key + string(DEBUG_KEY_SEPARATOR) + fmt.Sprintf("%v", value)
	}
	return flat
}

// Get Struct Type (or any types, but mostly using this to get what type of struct) in generics functions
// removeDots will remove the filename.[structName] and return only structName
func GetType(a interface{}, removeDots bool) string {
	if removeDots {
		return LastStringAfterDelimiter(fmt.Sprintf("%T", a), ".")
	}
	return fmt.Sprintf("%T", a)
}
