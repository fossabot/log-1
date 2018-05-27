package log

import "errors"

const (
	Hour           = "H"
	Minute         = "M"
	Day            = "D"
	Midnight       = "MIDNIGHT"
	AccessFormat   = "[%D %T] [Access] %M"
	OprationFormat = "[%D %T] [%L] (%S) %M"
	MessageOnly    = "%M"
	Location       = "Asia/Chongqing"

	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Level int

var (
	RotationError = errors.New("rotation error")
	LogWriteError = errors.New("logwrite error")
	ReadDirError  = errors.New("readdir error")

	conversion = map[string]string{
		/*stdLongMonth      */ "B": "January",
		/*stdMonth          */ "b": "Jan",
		// stdNumMonth       */ "m": "1",
		/*stdZeroMonth      */ "m": "01",
		/*stdLongWeekDay    */ "A": "Monday",
		/*stdWeekDay        */ "a": "Mon",
		// stdDay            */ "d": "2",
		// stdUnderDay       */ "d": "_2",
		/*stdZeroDay        */ "d": "02",
		/*stdHour           */ "H": "15",
		// stdHour12         */ "I": "3",
		/*stdZeroHour12     */ "I": "03",
		// stdMinute         */ "M": "4",
		/*stdZeroMinute     */ "M": "04",
		// stdSecond         */ "S": "5",
		/*stdZeroSecond     */ "S": "05",
		/*stdLongYear       */ "Y": "2006",
		/*stdYear           */ "y": "06",
		/*stdPM             */ "p": "PM",
		// stdpm             */ "p": "pm",
		/*stdTZ             */ "Z": "MST",
		// stdISO8601TZ      */ "z": "Z0700",  // prints Z for UTC
		// stdISO8601ColonTZ */ "z": "Z07:00", // prints Z for UTC
		/*stdNumTZ          */ "z": "-0700", // always numeric
		// stdNumShortTZ     */ "b": "-07",    // always numeric
		// stdNumColonTZ     */ "b": "-07:00", // always numeric
	}

	fCache = &formatCache{}

	levelStrings = [...]string{
		"\x1b[36mDEBG\x1b[0m",
		"\x1b[34mINFO\x1b[0m",
		"\x1b[33mWARN\x1b[0m",
		"\x1b[31mEROR\x1b[0m",
		"\x1b[35mFATAL\x1b[0m"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "\x1b[37mUNKN\x1b[0m"
	}
	return levelStrings[int(l)]
}
