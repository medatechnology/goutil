package timedate

import (
	"fmt"
	"strings"
	"time"
)

type TimeConvertSetting string

// TODO: This was created originally for PocketBase, we need to change
// this to more standard format.
const (
	// not used
	DATE_ONLY     TimeConvertSetting = "DateOnly"
	TIME_ONLY     TimeConvertSetting = "TimeOnly"
	DATE_AND_TIME TimeConvertSetting = "DateAndTime"

	// FORMAT_PB_DATETIME string = "2006-01-02 15:04:05.000Z"
	FORMAT_PB_DATETIME string = "2006-01-02 15:04:05.999Z"
	FORMAT_PB_TIME     string = "15:04:05.000Z"
	FORMAT_H_M         string = "15:04"
	EMPTY_DATE_STRING  string = "0001-01-01 00:00:0.000Z" // not sure we need this, but just in case we need to define "empty" date while variable must have value, then this is the value to denote empty date
	BIG_DATE_STRING    string = "3000-01-01 00:00:0.000Z" // sometime we need to check if date is less than but the end-date is empty

	nanosPerMicro = 1e3
	nanosPerMilli = 1e6
	nanosPerSec   = 1e9
	nanosPerMin   = 60 * nanosPerSec
	nanosPerHour  = 60 * nanosPerMin
	nanosPerDay   = 24 * nanosPerHour
	nanosPerWeek  = 7 * nanosPerDay
	nanosPerMonth = 30 * nanosPerDay
	nanosPerYear  = 365 * nanosPerDay
)

// ShortCut Only, Not really needed
var (
	// NOTE: WARNING: IMPORTANT: PLEASE USE TimeToString before saving to PB if not it will not work!
	TimeToString       func(time.Time) string          = ConvertGoLangTimeToPBString
	StringToTime       func(string) time.Time          = ConvertDateTimeStringToGoLangTime
	DateStrToTime      func(string) time.Time          = ConvertDateOnlyStringToGoLangTime // convert "yyyy-mm-dd" to time.Time
	TimeStrToTime      func(string) time.Time          = ConvertTimeOnlyStringToGoLangTime
	EMPTY_DATE         time.Time                       = TimeStrToTime(EMPTY_DATE_STRING)
	ParseStringnToTime func(string) (time.Time, error) = ConvertManyTimeFormatToGoLangTime
)

type HumanReadableOptions struct {
	NanoSecondString   string
	MicroSecondString  string
	MilliSecondString  string
	SecondString       string
	MinuteString       string
	HourString         string
	DayString          string
	WeekString         string
	MonthString        string
	YearString         string
	ValueUnitSeparator string
	WithSpaceUnit      bool // space between units like: 20 days[this-space]20 hour
	WithComma          bool // comma between units like: 20 days, 20 hour
	WithAndEnd         bool // with "and" at the end: 20 days and 20 hour
	AutoSingular       bool // if true: if the value==1 then remove the last "s" at the end
}

func DefaultLongtHumanReadableOptions() HumanReadableOptions {
	return HumanReadableOptions{
		NanoSecondString:   "nanoseconds",
		MicroSecondString:  "microseconds",
		MilliSecondString:  "milliseconds",
		SecondString:       "seconds",
		MinuteString:       "minute",
		HourString:         "hours",
		DayString:          "days",
		WeekString:         "weeks",
		MonthString:        "months",
		YearString:         "years",
		ValueUnitSeparator: " ",
		WithSpaceUnit:      true,
		WithComma:          true,
		WithAndEnd:         true,
		AutoSingular:       true,
	}
}

func DefaultShortHumanReadableOptions() HumanReadableOptions {
	return HumanReadableOptions{
		NanoSecondString:   "ns",
		MicroSecondString:  "µs",
		MilliSecondString:  "ms",
		SecondString:       "s",
		MinuteString:       "m",
		HourString:         "h",
		DayString:          "d",
		WeekString:         "w",
		MonthString:        "m",
		YearString:         "y",
		ValueUnitSeparator: "",
		WithSpaceUnit:      false,
		WithComma:          true,
		WithAndEnd:         false,
		AutoSingular:       false,
	}
}

// Convert time.Time to string using default FORMAT_PB_DATETIME format
func ConvertGoLangTimeToPBString(tt time.Time) string {
	return ConvertGoLangTimeToStringWithFormat(tt, "")
}

// Convert time.Time to string with format (if empty return using FORMAT_PB_DATETIME)
func ConvertGoLangTimeToStringWithFormat(timeObj time.Time, format string) string {
	// Use provided format string if specified
	if format != "" {
		return timeObj.Format(format)
	}
	// Use PocketBase datetime format by default
	return timeObj.Format(FORMAT_PB_DATETIME)
}

// This is to convert string like 'hh:mm' | 'hh:mm:ss' | 'hh:mm:ss:nnnZ' format to time.Time struct
func ConvertTimeOnlyStringToGoLangTime(s string) time.Time {
	tt, err := time.Parse(time.TimeOnly, s)
	if err != nil {
		tt, err = time.Parse(FORMAT_H_M, s)
		if err != nil {
			tt, _ = time.Parse(FORMAT_PB_TIME, s)
			return tt
		}
	}
	return tt
}

// Convert string to time.Time (both date and time) Only accept yyyy-mm-dd format, else return .IsZero()
func ConvertDateOnlyStringToGoLangTime(s string) time.Time {
	tt, _ := time.Parse(time.DateOnly, s)
	return tt
}

// Convert string to time.Time (both date and time)
func ConvertDateTimeStringToGoLangTime(s string) time.Time {
	tt, err := time.Parse(FORMAT_PB_DATETIME, s)
	if err != nil {
		// Try again with RFC3339
		tt, err = time.Parse(time.RFC3339, s)
		if err != nil {
			tt, _ = time.Parse(time.DateTime, s)
			return tt
		}
	}
	return tt
}

// Convert string (time and date) to GoLangTime
// parseTime tries to parse a string into a time.Time using common formats.
func ConvertManyTimeFormatToGoLangTime(str string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",                     // Format for "2025-03-16 09:02:21"
		"2006-01-02 15:04",                        // Format for "2025-03-16 09:02"
		"2006-01-02 15:04:05.999999999 -0700 MST", // Format for "2025-03-16 09:02:21.438859256 +0000 UTC"
		"2006-01-02",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.Kitchen,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", str)
}

// // Convert string to time.Time with settings DATE_ONLY | TIME_ONLY | DATE_AND_TIME
// // using ConvertPBDateTimeStringToGoLangTime() which basically try PB format '...0000Z' OR RFC3339
// func ConvertDateTimeStringToDateTimeWithSetting(s string, c TimeConvertSetting) time.Time {
// 	// tt, err := ConvertPBDateTimeStringToGoLangTime(s)
// 	tt, err := time.Parse(FORMAT_PB_DATETIME, s)
// 	if err != nil {
// 		// Try again with RFC3339
// 		tt, err = time.Parse(time.RFC3339, s)
// 		if err != nil {
// 			return time.Time{}
// 		}
// 	}
// 	switch c {
// 	case DATE_ONLY:
// 		return time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, tt.Location())
// 	case TIME_ONLY:
// 		return time.Date(1, time.January, 1, tt.Hour(), tt.Minute(), tt.Second(), 0, tt.Location())
// 	default:
// 		return time.Date(tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), tt.Second(), 0, tt.Location())
// 	}
// 	// return time.Date(date.Year(), date.Month(), date.Day(), res.Hour(), res.Minute(), res.Second(), 0, date.Location())
// }

// Add days into string date, can do negative days as well to substract dateStr by days
// dateStr = mm/dd/yyyy  and return format is the same
func AddDaysToStringDate(dateStr string, days int) string {
	t := DateStrToTime(dateStr)
	t = t.AddDate(0, 0, days)
	return t.Format(time.DateOnly)
}

// Get n-days previous and n-days next from calcDate. Accept the dateOnly string
// Also return dateOnly format as string only, not time.Time
func GetRangeFromDate(calcDate string, n int) (string, string) {
	// d := DateStrToTime(calcDate)
	// prev := d.AddDate(0, 0, -n)
	// next := d.AddDate(0, 0, n)
	return AddDaysToStringDate(calcDate, -n), AddDaysToStringDate(calcDate, n)
}

// Add minutes to datetimeStr and return as time
func AddMinutesToDateTimeStringAsGoTime(datetimeStr string, minutes int) time.Time {
	t := ConvertDateTimeStringToGoLangTime(datetimeStr)
	t = t.Add(time.Duration(minutes) * time.Minute)
	return t
}

// Convert int duration to time string : now + duration
// Output is RFC3339
func TimeStringFromNow(duration int) string {
	return time.Now().Add(time.Second * time.Duration(int64(duration))).Format(time.RFC3339)
}

// Add minutes to date-time string (PB format) and return PB format as well.
// Can use negative value for minutes to substract
func AddMinutesToDateTimeString(datetimeStr string, minutes int) string {
	// t := ConvertDateTimeStringToGoLangTime(datetimeStr)
	// t = t.Add(time.Duration(minutes) * time.Minute)
	// return t.Format(FORMAT_PB_DATETIME)
	return AddMinutesToDateTimeStringAsGoTime(datetimeStr, minutes).Format(FORMAT_PB_DATETIME)
}

// Compare date_only format for strings ie "2024-04-03" and "2024-04-20" then true
func DateStringABeforeB(a, b string) bool {
	ta := DateStrToTime(a)
	tb := DateStrToTime(b)
	return ta.Before(tb)
}

// Is timestring expired
// Check if expiresAt(string) is already expired (before) compared to time.now
func IsTimeStringExpired(expiresAt string) bool {
	nowStr := TimeToString(time.Now())
	return DateStringABeforeB(expiresAt, nowStr)
	// Pocketbase use RFC3339
	// expTime, _ := time.Parse(time.RFC3339, expiresAt)
	// now := time.Now()
	// return expTime.Before(now)
}

// compare if string x is in between string a and string b and only compare the date.
// This takes string in dateOnly format 'yyyy-mm-dd'
func IsDateStrInRange(x, a, b string) bool {
	xStr := DateStrToTime(x)
	aStr := DateStrToTime(a)
	bStr := DateStrToTime(b)
	return (xStr.After(aStr) || xStr.Equal(aStr)) && (xStr.Before(bStr) || xStr.Equal(bStr))
}

func PBDateTimeStringABeforeB(a, b string) bool {
	ta := ConvertDateTimeStringToGoLangTime(a)
	tb := ConvertDateTimeStringToGoLangTime(b)
	return ta.Before(tb)
}

// days.TimeStringToTimeInstanceWithDate cannot be used because it parse
// timeonly as date.TimeOnly ("hh:mm:ss") format. While timeonly is usually
func SetDateToTimeOnlyGoLangTime(timeonly, dateToSet time.Time) time.Time {
	return time.Date(dateToSet.Year(), dateToSet.Month(), dateToSet.Day(), timeonly.Hour(), timeonly.Minute(), timeonly.Second(), 0, dateToSet.Location())
}

// Format that from PB usually is 2006-01-02 15:04:05.999Z we only need the first 10 chars
func ExtractDateFromString(thedate string) string {
	if len(thedate) < 10 {
		return "" // String is too short to contain a valid date
	}
	return thedate[:10]
}

// format from PB to get time only
func ExtractTimeFromString(thedate string) string {
	if len(thedate) < 17 {
		return ""
	}
	return thedate[11:16]
}

// Format for printing only mostly
func ExtractDateTimeFromString(thedate string) string {
	if len(thedate) < 17 {
		return ""
	}
	return thedate[:16]
}

// Format that from PB usually is 2006-01-02 15:04:05.999Z we only need the first 4 chars
func ExtractYearFromString(thedate string) string {
	if len(thedate) < 4 {
		return "" // String is too short to contain a valid date
	}
	return thedate[:4]
}

// Get start and end date of the year in time.Time (UTC)
func GetFirstAndLastDayOfYear(year int) (time.Time, time.Time) {
	return GetFirstDayOfYear(year), GetLastDayOfYear(year)
}

// Get first day of the year
func GetFirstDayOfYear(year int) time.Time {
	return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func GetLastDayOfYear(year int) time.Time {
	return time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)
}

// Meaning from time B --- to time A ==> meaning = B - A (if posistive means B is later, negative means B is before A)
func GetMinutesFromDateTimeBToA(a, b time.Time) float64 {
	// Calculate the duration between t1 and t2
	if a.IsZero() || b.IsZero() {
		return 0
	}
	duration := b.Sub(a)
	// Convert duration to minutes NOTE: maybe if next time need to be rounded, then maybe return float64
	// minutes := int(duration.Minutes())
	// return minutes
	return duration.Minutes()
}

// get number of days between time B to time A (B - A). Except time.Time format and only dateStr
func GetDaysFromDateTimeBToA(a, b time.Time) int {
	// Calculate the duration between t1 and t2
	duration := b.Sub(a)
	// Convert duration to days
	days := int(duration.Hours() / 24)
	// fmt.Println("Get difference between time  ", " time A=", a, " time B=", b, " duration=", duration, " days=", days)
	return days
}

// Ignore error and return 0 instead. Except "10:32" and "13:20" format
func GetMinutesFromTimeStringBToA(timeStr1, timeStr2 string) float64 {
	t1 := TimeStrToTime(timeStr1)
	t2 := TimeStrToTime(timeStr2)
	if t1.IsZero() || t2.IsZero() {
		return 0
	}
	// Calculate the duration between t1 and t2
	duration := t2.Sub(t1)
	// Convert duration to minutes
	// minutes := int(duration.Minutes())
	return duration.Minutes()
}

// Ignore error, accept "2024-12-03 10:23:00.000Z" format. B minus A
func GetMinutesFromDateTimeStringBtoA(a, b string) float64 {
	t1 := StringToTime(a)
	t2 := StringToTime(b)
	if t1.IsZero() || t2.IsZero() {
		return 0
	}
	// Calculate the duration between t1 and t2
	duration := t2.Sub(t1)
	return duration.Minutes()
}

// ============= THIS IS UTILS functions if not yet exist
// getNumberOfWeeksInMonth calculates the number of weeks in a month
func GetNumberOfWeeksInMonth(year int, month time.Month) int {
	// Get the first and last day of the month
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	// Get the ISO week numbers for the first and last day
	_, firstWeek := firstOfMonth.ISOWeek()
	_, lastWeek := lastOfMonth.ISOWeek()

	// If the month starts and ends in the same ISO week, there's only 1 week
	if firstWeek == lastWeek {
		return 1
	}

	// Otherwise, consider if the last day belongs to the next month's first week
	if lastWeek == 52 && month == time.December {
		return 5
	} else if lastWeek == 53 && (month == time.November || month == time.December) {
		return 5
	}

	// Account for leftover days by adding 1
	return lastWeek - firstWeek + 1
}

// 2024, April =

// getWeekNumberInMonth determines the week number for a specific date
func GetWeekNumberInMonth(date time.Time) int {
	// Get the first day of the month
	firstOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Get the ISO week numbers for the date and first day
	_, thisWeek := date.ISOWeek()
	_, beginningWeek := firstOfMonth.ISOWeek()

	// Account for leftover days by adding 1
	return 1 + thisWeek - beginningWeek
}

func DurationHUmanReadableWithOptions(uptime time.Duration, opt HumanReadableOptions) string {
	ns := uptime.Nanoseconds()

	var parts []string
	addUnit := func(value int64, unit string) {
		if unit == "" || value == 0 {
			return
		}
		if opt.AutoSingular && value == 1 {
			unit = strings.TrimSuffix(unit, "s")
		}
		parts = append(parts, fmt.Sprintf("%d%s%s", value, opt.ValueUnitSeparator, unit))
	}

	years := ns / nanosPerYear
	addUnit(years, opt.YearString)
	ns %= nanosPerYear

	months := ns / nanosPerMonth
	addUnit(months, opt.MonthString)
	ns %= nanosPerMonth

	weeks := ns / nanosPerWeek
	addUnit(weeks, opt.WeekString)
	ns %= nanosPerWeek

	days := ns / nanosPerDay
	addUnit(days, opt.DayString)
	ns %= nanosPerDay

	hours := ns / nanosPerHour
	addUnit(hours, opt.HourString)
	ns %= nanosPerHour

	minutes := ns / nanosPerMin
	addUnit(minutes, opt.MinuteString)
	ns %= nanosPerMin

	seconds := ns / nanosPerSec
	addUnit(seconds, opt.SecondString)
	ns %= nanosPerSec

	millis := ns / nanosPerMilli
	addUnit(millis, opt.MilliSecondString)
	ns %= nanosPerMilli

	micros := ns / nanosPerMicro
	addUnit(micros, opt.MicroSecondString)
	ns %= nanosPerMicro

	addUnit(ns, opt.NanoSecondString)

	if len(parts) == 0 {
		return ""
	}

	// Determine separators
	var mainSeparator string
	if opt.WithComma {
		mainSeparator = ","
	}
	if opt.WithSpaceUnit {
		mainSeparator += " "
	}

	switch {
	case len(parts) == 1:
		return parts[0]

	case len(parts) == 2 && opt.WithAndEnd:
		return strings.Join(parts, " and ")

	case len(parts) >= 2 && opt.WithAndEnd:
		return strings.Join(parts[:len(parts)-1], mainSeparator) + " and " + parts[len(parts)-1]

	default:
		return strings.Join(parts, mainSeparator)
	}
}

// Default long format for duration string, prints every unit from year-to-nanoseconds
// Output: 2 days, 3 hours and 15 minutes
func DurationHumanReadableLong(uptime time.Duration) string {
	return DurationHUmanReadableWithOptions(uptime, DefaultLongtHumanReadableOptions())
}

// Default long format for Uptime duration string
// Output: remove the seconds and smaller
func DurationUptimeLong(uptime time.Duration) string {
	options := DefaultLongtHumanReadableOptions()
	options.SecondString = ""
	options.MilliSecondString = ""
	options.MicroSecondString = ""
	options.NanoSecondString = ""
	return DurationHUmanReadableWithOptions(uptime, options)
}

// Default short format for Uptime duration string
// Output: remove the seconds and smaller
// 2d,20h,12m  or just remove the comma as well: 2d20h12m
func DurationUptimeShort(uptime time.Duration) string {
	options := DefaultShortHumanReadableOptions()
	options.SecondString = ""
	options.MilliSecondString = ""
	options.MicroSecondString = ""
	options.NanoSecondString = ""
	return DurationHUmanReadableWithOptions(uptime, options)
}

// This is usually to print the elapsed time for some function or db access
// usually is small units, so definitely no years and all.
// Output: 23 milliseconds,12 microseconds,39 nanoseconds
func DurationElapsedLong(uptime time.Duration) string {
	options := DefaultLongtHumanReadableOptions()
	options.YearString = ""
	options.MonthString = ""
	options.WeekString = ""
	options.DayString = ""
	options.HourString = ""
	options.WithAndEnd = false
	options.WithSpaceUnit = false
	return DurationHUmanReadableWithOptions(uptime, options)
}

// This is usually to print the elapsed time for some function or db access
// usually is small units, so definitely no years -- until hours
// Output: 23ms,12µs,39ns
func DurationElapsedShort(uptime time.Duration) string {
	options := DefaultShortHumanReadableOptions()
	options.YearString = ""
	options.MonthString = ""
	options.WeekString = ""
	options.DayString = ""
	options.HourString = ""
	return DurationHUmanReadableWithOptions(uptime, options)
}
