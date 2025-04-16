package print

import (
	"fmt"
	"strings"
)

// This is just to print with pretty format for console using ASCII UniCode

// ANSI color codes
const (
	Reset        = "\033[0m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Purple       = "\033[35m"
	Cyan         = "\033[36m"
	Gray         = "\033[37m"
	DarkGray     = "\033[90m"
	LightRed     = "\033[91m"
	LightGreen   = "\033[92m"
	LightYellow  = "\033[93m"
	LightBlue    = "\033[94m"
	LightMagenta = "\033[95m"
	LightCyan    = "\033[96m"
	White        = "\033[97m"
)

// Color represents a color code for terminal output, for readibility
type Color struct {
	Code string
}

// Predefined colors for convenience as variables (can be changed)
// These are the colors that can be used directly
// You can use these colors directly in your code
// Example: fmt.Println(ColorRed.Code + "This is red text" + Reset)
// You can also use the Color struct to create your own colors
// Example: myColor := Color{Code: "\033[38;5;82m"} // Custom color
var (
	ColorRed          = Color{Red}
	ColorGreen        = Color{Green}
	ColorYellow       = Color{Yellow}
	ColorBlue         = Color{Blue}
	ColorPurple       = Color{Purple}
	ColorCyan         = Color{Cyan}
	ColorGray         = Color{Gray}
	ColorDarkGray     = Color{DarkGray}
	ColorLightRed     = Color{LightRed}
	ColorLightGreen   = Color{LightGreen}
	ColorLightYellow  = Color{LightYellow}
	ColorLightBlue    = Color{LightBlue}
	ColorLightMagenta = Color{LightMagenta}
	ColorLightCyan    = Color{LightCyan}
	ColorWhite        = Color{White}
	ColorNothing      = Color{""}
	ColorReset        = Color{Reset}
)

// BoxChars represents the characters used to draw a box
// This is used to draw a box around the text
// You can change these characters to use different styles
// Example: BoxChars{TopLeft: "╔", TopRight: "╗", BottomLeft: "╚", BottomRight: "╝", Horizontal: "═", Vertical: "║"}
// Example: BoxChars{TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "┘", Horizontal: "─", Vertical: "│"}
// Example: BoxChars{TopLeft: "╭", TopRight: "╮", BottomLeft: "╰", BottomRight: "╯", Horizontal: "─", Vertical: "│"}
// Example: BoxChars{TopLeft: "┏", TopRight: "┓", BottomLeft: "┗", BottomRight: "┛", Horizontal: "━", Vertical: "┃"}
// Example: BoxChars{TopLeft: "╓", TopRight: "╖", BottomLeft: "╙", BottomRight: "╜", Horizontal: "─", Vertical: "│"}
type BoxChars struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
}

// UnicodeBox is the default box characters used for drawing boxes
// You can change this to use different styles
// Example: UnicodeBox = BoxChars{TopLeft: "╔", TopRight: "╗", BottomLeft: "╚", BottomRight: "╝", Horizontal: "═", Vertical: "║"}
// Example: UnicodeBox = BoxChars{TopLeft: "╭", TopRight: "╮", BottomLeft: "╰", BottomRight: "╯", Horizontal: "─", Vertical: "│"}
// Example: UnicodeBox = BoxChars{TopLeft: "┏", TopRight: "┓", BottomLeft: "┗", BottomRight: "┛", Horizontal: "━", Vertical: "┃"}
// Example: UnicodeBox = BoxChars{TopLeft: "╓", TopRight: "╖", BottomLeft: "╙", BottomRight: "╜", Horizontal: "─", Vertical: "│"}
// Example: UnicodeBox = BoxChars{TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "┘", Horizontal: "─", Vertical: "│"}
var UnicodeBox = BoxChars{
	TopLeft:     "┌",
	TopRight:    "┐",
	BottomLeft:  "└",
	BottomRight: "┘",
	Horizontal:  "─",
	Vertical:    "│",
}

// Colored is a function that takes a string and a color code and returns the string with the color code applied
// This is used to colorize the text
// Example: Colored("Hello World", ColorRed.Code) // This will return the string with red color code
func Colored(text string, color Color) string {
	return color.Code + text + Reset
}

// Used to print 2 columns of key-value pairs in a box, between left and right border
// KeyValue represents a key-value pair
// TODO: change this to KeyValueStrings and maybe put them in object?
type KeyValue struct {
	Key       string
	Value     string
	OneColumn bool // If true, this entry will use the full width (one column)
	EmptyLine bool // If true, this will render as an empty line with just borders
}

// Mapping for Print, Println  and Printf
func PrintColor(col Color, a ...any) {
	fmt.Print(col.Code)
	fmt.Print(a...)
	fmt.Print(Reset)
}

func PrintlnColor(message string, col Color) {
	fmt.Println(Colored(message, col))
}

func PrintfColor(format string, col Color, a ...any) {
	fmt.Printf(Colored(format, col), a...)
}

// Just alias to make the code shorter
func Content(singleColumn, empty bool, k interface{}, v interface{}) KeyValue {
	keyStr := fmt.Sprintf("%v", k) // Using fmt.Sprintf to convert interface{} to string
	valueStr := fmt.Sprintf("%v", v)

	return KeyValue{Key: keyStr, Value: valueStr, OneColumn: singleColumn, EmptyLine: empty}
}

func PrintBoxHeadingContent(heading []string, headingColors []Color, content []KeyValue, keyColor Color, valueColor Color) {
	// Find the longest key for consistent spacing.
	longestKey := 0
	for _, kv := range content {
		if len(kv.Key) > longestKey && !kv.EmptyLine {
			longestKey = len(kv.Key)
		}
	}

	// Calculate the box width.
	boxWidth := 80 // Fixed box width

	// Print the top border.
	fmt.Println(UnicodeBox.TopLeft + strings.Repeat(UnicodeBox.Horizontal, boxWidth-2) + UnicodeBox.TopRight)

	// Print the heading.
	for i, line := range heading {
		headingPadding := (boxWidth - len(line) - 2) / 2
		if headingPadding < 0 {
			headingPadding = 0
		}
		var currentColor Color
		if i < len(headingColors) {
			currentColor = headingColors[i]
		} else {
			currentColor = ColorWhite // Default color
		}
		fmt.Printf("%s%s%s%s%s\n",
			UnicodeBox.Vertical,
			strings.Repeat(" ", headingPadding),
			Colored(line, currentColor),
			strings.Repeat(" ", boxWidth-len(line)-2-headingPadding),
			UnicodeBox.Vertical)
	}

	// extra line between heading and content
	fmt.Println(UnicodeBox.Vertical + strings.Repeat(" ", boxWidth-2) + UnicodeBox.Vertical) // just space

	// Calculate fixed positions
	firstColStart := 2               // After border and space
	secondColStart := boxWidth/2 - 3 // Fixed position for second column

	// Process content items
	i := 0
	for i < len(content) {
		if content[i].EmptyLine {
			// Empty line - just print borders with spaces
			fmt.Println(UnicodeBox.Vertical + strings.Repeat(" ", boxWidth-2) + UnicodeBox.Vertical)
			i++
			continue
		}

		if content[i].OneColumn || i+1 >= len(content) || content[i+1].OneColumn || content[i+1].EmptyLine {
			// Single column mode (either explicitly marked or last item or next is one-column/empty)
			kv := content[i]
			line := UnicodeBox.Vertical + " " // Start with left border and space

			// Full width for key-value
			line += Colored(kv.Key, keyColor) + strings.Repeat(" ", ZeroIfNegative(longestKey-len(kv.Key))) + " : " + Colored(kv.Value, valueColor)

			// Pad to right border
			remainingSpace := max(0, boxWidth-firstColStart-longestKey-len(kv.Value)-3-1) // -3 for " : ", -1 for border
			line += strings.Repeat(" ", remainingSpace)

			line += UnicodeBox.Vertical
			fmt.Println(line)
			i++
		} else {
			// Two column mode
			kv1 := content[i]
			kv2 := content[i+1]

			line := UnicodeBox.Vertical + " " // Start with left border and space

			// First column
			line += Colored(kv1.Key, keyColor) + strings.Repeat(" ", ZeroIfNegative(longestKey-len(kv1.Key))) + " : " + Colored(kv1.Value, valueColor)

			// Calculate padding between first column and second column
			currentPos := firstColStart + longestKey + 3 + len(kv1.Value) // 3 for " : "
			middlePadding := max(0, secondColStart-currentPos)
			line += strings.Repeat(" ", middlePadding)

			// Second column
			line += Colored(kv2.Key, keyColor) + strings.Repeat(" ", ZeroIfNegative(longestKey-len(kv2.Key))) + " : " + Colored(kv2.Value, valueColor)

			// Add padding to right border
			currentPos = secondColStart + longestKey + 3 + len(kv2.Value) // 3 for " : "
			rightPadding := ZeroIfNegative(boxWidth - currentPos - 1)     // -1 for border
			line += strings.Repeat(" ", rightPadding)

			line += UnicodeBox.Vertical
			fmt.Println(line)
			i += 2
		}
	}

	// Print the bottom border.
	fmt.Println(UnicodeBox.BottomLeft + strings.Repeat(UnicodeBox.Horizontal, boxWidth-2) + UnicodeBox.BottomRight)
}

// BytesToHumanReadable converts bytes to human-readable format (KB, MB, GB)
// Output: "12 MB" or "12 GB" or "890 KB"
func BytesToHumanReadable(bytes int64, separator string) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d%sB", bytes, separator)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	value := float64(bytes) / float64(div)
	return fmt.Sprintf("%.1f%s%cB", value, separator, "KMGTPE"[exp])
}

// Helper function to get maximum value (to ensure no negative padding)
func ZeroIfNegative(a int) int {
	if a < 0 {
		return 0
	}
	return a
}

// Print box, then some heading (aligned center) then some K-V fields displayed in 2 columns.
// Using unicode
// func PrintBoxHeadingContentOLD(heading []string, headingColors []Color, content []KeyValue, keyColor Color, valueColor Color) {
// 	// Find the longest key for consistent spacing.
// 	longestKey := 0
// 	for _, kv := range content {
// 		if len(kv.Key) > longestKey {
// 			longestKey = len(kv.Key)
// 		}
// 	}

// 	// Calculate the box width.
// 	boxWidth := 80 // Fixed box width

// 	// Print the top border.
// 	fmt.Println(UnicodeBox.TopLeft + strings.Repeat(UnicodeBox.Horizontal, boxWidth-2) + UnicodeBox.TopRight)

// 	// Print the heading.
// 	for i, line := range heading {
// 		headingPadding := (boxWidth - len(line) - 2) / 2
// 		if headingPadding < 0 {
// 			headingPadding = 0
// 		}
// 		var currentColor Color
// 		if i < len(headingColors) {
// 			currentColor = headingColors[i]
// 		} else {
// 			currentColor = ColorNothing // Default color
// 		}
// 		fmt.Printf("%s%s%s%s%s\n",
// 			UnicodeBox.Vertical,
// 			strings.Repeat(" ", headingPadding),
// 			Colored(line, currentColor),
// 			strings.Repeat(" ", boxWidth-len(line)-2-headingPadding),
// 			UnicodeBox.Vertical)
// 	}

// 	// extra line between heading and content
// 	// fmt.Println(UnicodeBox.Vertical + strings.Repeat(UnicodeBox.Horizontal, boxWidth-2) + UnicodeBox.Vertical) // add line
// 	fmt.Println(UnicodeBox.Vertical + strings.Repeat(" ", boxWidth-2) + UnicodeBox.Vertical) // just space

// 	// Calculate fixed positions
// 	firstColStart := 2               // After border and space
// 	secondColStart := boxWidth/2 - 3 // Fixed position for second column

// 	// Print content in pairs
// 	for i := 0; i < len(content); i += 2 {
// 		kv1 := content[i]
// 		line := UnicodeBox.Vertical + " " // Start with left border and space

// 		// First column
// 		line += Colored(kv1.Key, keyColor) + strings.Repeat(" ", longestKey-len(kv1.Key)) + " : " + Colored(kv1.Value, valueColor)

// 		if i+1 < len(content) {
// 			// Second column exists
// 			kv2 := content[i+1]

// 			// Calculate padding between first column and second column
// 			currentPos := firstColStart + longestKey + 3 + len(kv1.Value) // 3 for " : "
// 			middlePadding := secondColStart - currentPos
// 			line += strings.Repeat(" ", middlePadding)

// 			// Second column
// 			line += Colored(kv2.Key, keyColor) + strings.Repeat(" ", longestKey-len(kv2.Key)) + " : " + Colored(kv2.Value, valueColor)

// 			// Add padding to right border
// 			currentPos = secondColStart + longestKey + 3 + len(kv2.Value) // 3 for " : "
// 			rightPadding := boxWidth - currentPos - 1                     // -1 for border
// 			line += strings.Repeat(" ", rightPadding)
// 		} else {
// 			// Single column - pad to right border
// 			remainingSpace := boxWidth - firstColStart - longestKey - len(kv1.Value) - 3 - 1 // -3 for " : ", -1 for border
// 			line += strings.Repeat(" ", remainingSpace)
// 		}

// 		line += UnicodeBox.Vertical
// 		fmt.Println(line)
// 	}

// 	// Print the bottom border.
// 	fmt.Println(UnicodeBox.BottomLeft + strings.Repeat(UnicodeBox.Horizontal, boxWidth-2) + UnicodeBox.BottomRight)
// }

// Here is the example to call this.
func TestBoxPrint() {
	appName := []string{"Box Print", "With Colored Text", "Configuration"}
	headingColors := []Color{
		ColorCyan,
		ColorGreen,
		ColorNothing,
	}

	// Content defined in order
	appSettings := []KeyValue{
		Content(false, false, "Version", "1.2.3"),
		Content(false, false, "Environment", "Production"),
		Content(false, false, "Database", "PostgreSQL"),
		Content(false, false, "API Key", "********************"),
		Content(false, false, "Server IP", "192.168.1.100"),
		Content(false, false, "Port", "8080"),
		Content(false, false, "Debug", "true"),
		Content(false, false, "Timeout", "10s"),
		Content(false, false, "Extra", "value"),
	}

	keyColor := ColorYellow
	valueColor := ColorWhite

	PrintBoxHeadingContent(appName, headingColors, appSettings, keyColor, valueColor)
}
