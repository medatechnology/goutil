package object

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Return the last string s after delimiter d, example: s='this;is;delimited;by;semicolon' and d=';' then this function return "semicolon"
func LastStringAfterDelimiter(s, d string) string {
	ss := strings.Split(s, d)
	return ss[len(ss)-1]
}

// Trick to combine multiple any/interface variables into 1 interface
// NOTE: NOT USED
func Combine(msgs ...any) interface{} {
	return msgs
}

// combined key-value (json usually) into string using marshal
// NOT USED
func CombineToStringWithField(msgs ...any) string {
	dat, err := json.Marshal(&msgs)
	if err != nil {
		fmt.Println("error=", err)
	}
	return string(dat)
}

// Combined with space
func CombineToStringWithSpace(msgs ...any) string {
	// aa := ""
	// for _, v := range msgs {
	//      aa = aa + fmt.Sprintf("%v ", v)
	// }
	// return strings.Trim(aa, " ")
	return CombineToStringWithSeparator(" ", msgs)
}

// Combine array of strings with seperator (sep)
func CombineToStringWithSeparator(sep string, msgs ...any) string {
	aa := ""
	for _, v := range msgs {
		aa = aa + fmt.Sprintf("%v%s", v, sep)
	}
	return strings.Trim(aa, sep)
}

// String Array A contains item (string) B
// NOTE: maybe change to generics?
func ArrayAContainsBString(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

// NOTE: Not USED!
// Use this to log/print array of string combined with DEBUG_SEPARATOR
func StringSliceDebug(msgs []string) string {
	return strings.Join(msgs[:], DEBUG_KEY_SEPARATOR)
}

// str is multiple value separated by space, ie: str="token id_token"
// cannot use contain because rtype="token" might hit id_token as well
func StringContainsSplitBySpace(str, rtype, separator string) bool {
	strArr := strings.Split(str, separator)
	return ArrayAContainsBString(strArr, rtype)
}
