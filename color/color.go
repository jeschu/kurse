package color

import "fmt"

//goland:noinspection GoUnusedGlobalVariable
var (
	Reset = "\033[0m"

	/////////////
	// Special //
	/////////////

	Bold      = "\033[1m"
	Underline = "\033[4m"

	/////////////////
	// Text colors //
	/////////////////

	Black  = "\033[30m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"

	///////////////////////
	// Background colors //
	///////////////////////

	BlackBackground  = "\033[40m"
	RedBackground    = "\033[41m"
	GreenBackground  = "\033[42m"
	YellowBackground = "\033[43m"
	BlueBackground   = "\033[44m"
	PurpleBackground = "\033[45m"
	CyanBackground   = "\033[46m"
	GrayBackground   = "\033[47m"
	WhiteBackground  = "\033[107m"
)

func colorize(color string, s any) string {
	switch s := s.(type) {
	case string:
		return color + s + Reset
	default:
		return color + fmt.Sprint(s) + Reset
	}
}

func InRed(s any) string    { return colorize(Red, s) }
func InGreen(s any) string  { return colorize(Green, s) }
func InYellow(s any) string { return colorize(Yellow, s) }
func InGray(s any) string   { return colorize(Gray, s) }

func ByAmount(amount float64, str string) string {
	if amount < 0 {
		return InRed(str)
	} else {
		return InGreen(str)
	}
}
