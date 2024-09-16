// Package crondescriptor converts cron expressions into human readable
// strings. The package includes four options for minor customization
// of output.
//
package crondescriptor

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInvalidDescriptionType = fmt.Errorf("invalid description type provided")
	ErrInvalidSegmentCase     = fmt.Errorf("invalid case reached; please investigate")
	ErrBlankExpression        = fmt.Errorf("expression cannot be blank")
	ErrInvalidFieldCount      = fmt.Errorf("at least five fields are required")
	ErrFieldCountExceeded     = fmt.Errorf("expression has too many fields; should not exceed 7 fields")
	ErrWeekStartIsZero        = fmt.Errorf("week start already 0; check DayOfWeekIndexZero option")
	ErrInvalidCharacters      = fmt.Errorf("invalid character(s). allowed value: 0-23. allowed special characters: '*' ',' '-'")
	ErrInvalidMinuteFormat    = fmt.Errorf("invalid minute format")
	ErrInvalidSecondsValue    = fmt.Errorf("invalid seconds value. only 0-59 allowed)")
	ErrInvalidDayOfWeekRange  = fmt.Errorf("day of week range can only be 0-6 or 1-7 depending on DayOfWeekIndexZero option")
)

// CronDaysShort contains the days of the week short form
var CronDaysShort = []string{"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"}

// CronDaysLong contains the days of the week long form
var CronDaysLong = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

// CronMonths are zero indexed with Jan
var CronMonths = []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}

// CronDescriptor contains the cron expression, the array equivalent, and other settings
type CronDescriptor struct {
	Expression      string
	expressionArray [7]string
	Options         Options
	Logger          *zap.Logger
	sugarLog        *zap.SugaredLogger
}

// Options for the CronDescriptor functions
type Options struct {
	CasingType          CasingType
	DayOfWeekIndexZero  bool
	Use24HourTimeFormat bool
	Verbose             bool
}

// NewCronDescriptor creates a new CronDescriptor object
func NewCronDescriptor(cronExpr string) (*CronDescriptor, error) {
	logLevelEnvVar := os.Getenv("CRON_DESCRIPTOR_LOG_LEVEL")

	cfg := zap.NewDevelopmentConfig()
	if logLevelEnvVar == "debug" {
		cfg.Level.SetLevel(zap.DebugLevel)
	} else {
		cfg.Level.SetLevel(zap.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	options := Options{
		CasingType:          CasingTypeSentence,
		DayOfWeekIndexZero:  true,
		Use24HourTimeFormat: false,
		Verbose:             false,
	}

	cd := &CronDescriptor{
		Expression: cronExpr,
		Logger:     logger,
		sugarLog:   logger.Sugar(),
		Options:    options,
	}

	err = cd.Parse(cronExpr)
	if err != nil {
		return nil, err
	}

	return cd, nil
}

// NewCronDescriptorWithOptions creates a new CronDescriptor object
func NewCronDescriptorWithOptions(cronExpr string, options Options) (*CronDescriptor, error) {
	logLevelEnvVar := os.Getenv("CRON_DESCRIPTOR_LOG_LEVEL")

	cfg := zap.NewDevelopmentConfig()
	if logLevelEnvVar == "debug" {
		cfg.Level.SetLevel(zap.DebugLevel)
	} else {
		cfg.Level.SetLevel(zap.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	cd := &CronDescriptor{
		Expression: cronExpr,
		Logger:     logger,
		sugarLog:   logger.Sugar(),
		Options:    options,
	}

	err = cd.Parse(cronExpr)
	if err != nil {
		return nil, err
	}

	return cd, nil
}

// Parse normalizes and breaks a cron expression into an array of length seven
func (cd *CronDescriptor) Parse(expression string) (err error) {
	parsedExpr := [7]string{"", "", "", "", "", "", ""}

	if expression == "" {
		return ErrBlankExpression
	}

	splitExpr := strings.Split(expression, " ")
	splitExprLen := len(splitExpr)

	switch {
	case splitExprLen < 5:
		return ErrInvalidFieldCount

	case splitExprLen == 5:
		cd.sugarLog.Debug("Cron expr len is 5 so shift array past seconds element")
		for i, field := range splitExpr {
			parsedExpr[i+1] = field
		}

		cd.sugarLog.Debugf("%s => %s", splitExpr, parsedExpr)

	case splitExprLen == 6:
		// If last element ends with 4 digits, a year element has been
		// supplied and no seconds element
		re := regexp.MustCompile(`\d{4}$`)
		if re.Find([]byte(splitExpr[5])) != nil {
			cd.sugarLog.Debugf("The last element (%s) appears to have four digits. Assuming a year.", splitExpr[5])
			for i, field := range splitExpr {
				parsedExpr[i+1] = field
			}
		} else {
			copy(parsedExpr[:], splitExpr[:6])
		}

	case splitExprLen == 7:
		copy(parsedExpr[:], splitExpr[:7])

	default:
		return ErrFieldCountExceeded
	}

	parsedExpr, err = cd.normalizeExpression(parsedExpr)
	if err != nil {
		return err
	}

	cd.expressionArray = parsedExpr
	return nil
}

func (cd *CronDescriptor) normalizeExpression(parsedExpr [7]string) ([7]string, error) {
	origExpr := parsedExpr

	// convert ? to * only for day of the month and day of the week
	parsedExpr[3] = strings.ReplaceAll(parsedExpr[3], "?", "*")
	parsedExpr[5] = strings.ReplaceAll(parsedExpr[5], "?", "*")

	// convert 0/, 1/ to */
	for i := 0; i < 3; i++ {
		if strings.HasPrefix(parsedExpr[i], "0/") {
			parsedExpr[i] = strings.ReplaceAll(parsedExpr[i], "0/", "*/") // seconds, minutes, hours
		}
	}

	// convert 1/ to */
	for i := 3; i < 7; i++ {
		if strings.HasPrefix(parsedExpr[i], "1/") {
			parsedExpr[i] = strings.ReplaceAll(parsedExpr[i], "1/", "*/") // day of the month, month, day of the week, ?
		}
	}

	// handle DayOfWeekStartIndexZero option where SUN=1 rather than SUN=0
	// if self._options.day_of_week_start_index_zero is False:
	// 	expression_parts[5] = self.decrease_days_of_week(expression_parts[5])
	cd.sugarLog.Debug(cd.Options.DayOfWeekIndexZero)
	if !cd.Options.DayOfWeekIndexZero {
		daysOfWeekDecreased, err := cd.decreaseDaysOfWeek(parsedExpr[5])
		if err != nil {
			return parsedExpr, err
		}
		parsedExpr[5] = *daysOfWeekDecreased
	}

	if parsedExpr[3] == "?" {
		parsedExpr[3] = "*"
	}

	// convert SUN-SAT format to 0-6 format
	for i, day := range CronDaysShort {
		parsedExpr[5] = strings.ReplaceAll(strings.ToUpper(parsedExpr[5]), day, strconv.Itoa(i))
	}

	// convert JAN-DEC format to 1-12 format
	for i, month := range CronMonths {
		parsedExpr[4] = strings.ReplaceAll(strings.ToUpper(parsedExpr[4]), month, strconv.Itoa(i+1))
	}

	// convert 0 second to (empty)
	if parsedExpr[0] == "0" {
		parsedExpr[0] = ""
	}

	for i := range parsedExpr {
		// convert all '*/1' to '*'
		if parsedExpr[i] == "*/1" {
			parsedExpr[i] = "*"
		}

		/*
			Convert Month,day of the week,Year step values with a starting value (i.e. not '*') to between expressions.
			This allows us to reuse the between expression handling for step values.
			For Example:
			- month part '3/2' will be converted to '3-12/2' (every 2 months between March and December)
			- day of the week part '3/2' will be converted to '3-6/2' (every 2 days between Tuesday and Saturday)
		*/

		slashFound := strings.Index(parsedExpr[i], "/")
		asterixFound := strings.Index(parsedExpr[i], "*")
		dashFound := strings.Index(parsedExpr[i], "-")
		commaFound := strings.Index(parsedExpr[i], ",")
		if slashFound != -1 && asterixFound == -1 && dashFound == -1 && commaFound == -1 {
			// Removed 9999 as a bounding year. This will result in years
			// with the word "starting in ..."
			choices := [7]string{"", "", "", "", "12", "6", ""}

			stepThroughRange := choices[i]

			if stepThroughRange != "" {
				fieldSplit := strings.Split(parsedExpr[i], "/")
				parsedExpr[i] = fmt.Sprintf("%s-%s/%s", fieldSplit[0], stepThroughRange, fieldSplit[1])
				cd.sugarLog.Debugf("Value %s change to %s", stepThroughRange, parsedExpr[i])
			}
		}
	}

	cd.sugarLog.Debugf("Expression %s normalized to %s", origExpr, parsedExpr)
	return parsedExpr, nil
}

func (cd *CronDescriptor) decreaseDaysOfWeek(weekDay string) (*string, error) {
	cd.sugarLog.Debugf("Decrease days of week for %s", weekDay)
	var dowChars []string

	if weekDay == "*" {
		return &weekDay, nil
	}

	for i, c := range weekDay {
		cd.sugarLog.Debugf("Start at %d for char %s", i, string(c))

		decrement := false
		if i == 0 {
			decrement = true
		} else if i > 0 {
			if string(weekDay[i-1]) != "#" && string(weekDay[i-1]) != "/" && string(weekDay[i]) != "-" {
				decrement = true
			}
		}

		if decrement {
			cd.sugarLog.Debugf("Decrementing string int %s", string(c))

			charToInt, err := strconv.ParseInt(string(c), 10, 64)
			if err != nil {
				return nil, err
			}

			cd.sugarLog.Debugf("Int conversion: %d", charToInt)
			if charToInt == 0 {
				return nil, ErrWeekStartIsZero
			}

			dowChars = append(dowChars, strconv.Itoa(int(charToInt)-1))
		} else {
			dowChars = append(dowChars, string(c))
		}
	}

	cd.sugarLog.Debugf("Joining %s into one string", dowChars)
	joinedChars := strings.Join(dowChars, "")
	finalString := &joinedChars

	cd.sugarLog.Debugf("Final string: %s", *finalString)
	return finalString, nil
}

// DescriptionTypeEnum for the type of description
type DescriptionTypeEnum int

const (
	// Full - provide full description
	Full DescriptionTypeEnum = iota
	// TimeOfDay - only the time of the day portion of the expression
	TimeOfDay
	// Seconds - only the seconds portion of the expression
	Seconds
	// Minutes - only the minutes portion of the expression
	Minutes
	// Hours - only the hours portion of the expression
	Hours
	// DayOfWeek - only the day of the week portion of the expression
	DayOfWeek
	// Month - only the month portion of the expression
	Month
	// DayOfMonth - only the day of the month portion of the expression
	DayOfMonth
	// Year - only the Year portion of the expression
	Year
)

// GetDescription returns the description of the expression based on the parameter provided
func (cd *CronDescriptor) GetDescription(descriptionType DescriptionTypeEnum) (*string, error) {
	var description *string
	var err error

	switch descriptionType {
	case Full:
		description, err = cd.getFullDescription()
	case TimeOfDay:
		description, err = cd.getTimeOfDayDescription()
	case Hours:
		description, err = cd.getHoursDescription()
	case Minutes:
		description, err = cd.getMinutesDescription()
	case Seconds:
		description, err = cd.getSecondsDescription()
	case DayOfMonth:
		description, err = cd.getDayOfMonthDescription()
	case Month:
		description, err = cd.getMonthDescription()
	case DayOfWeek:
		description, err = cd.getDayOfTheWeekDescription()
	case Year:
		description, err = cd.getYearDescription()
	default:
		return nil, ErrInvalidDescriptionType
	}

	if err != nil {
		return nil, err
	}

	return description, nil
}

// return description
func (cd *CronDescriptor) getFullDescription() (*string, error) {
	timeSegment, err := cd.getTimeOfDayDescription()
	if err != nil {
		return nil, err
	}
	cd.sugarLog.Debugf("timeSegment: %s", *timeSegment)

	dayOfMonthDesc, err := cd.getDayOfMonthDescription()
	if err != nil {
		return nil, err
	}
	cd.sugarLog.Debugf("dayOfMonthDesc: %s", *dayOfMonthDesc)

	dayOfWeekDesc, err := cd.getDayOfTheWeekDescription()
	if err != nil {
		return nil, err
	}
	cd.sugarLog.Debugf("dayOfWeekDesc: %s", *dayOfWeekDesc)

	monthDesc, err := cd.getMonthDescription()
	if err != nil {
		return nil, err
	}
	cd.sugarLog.Debugf("monthDesc: %s", *monthDesc)

	yearDesc, err := cd.getYearDescription()
	if err != nil {
		return nil, err
	}
	cd.sugarLog.Debugf("yearDesc: %s", *yearDesc)

	description := fmt.Sprintf("%s%s%s%s%s",
		*timeSegment, *dayOfMonthDesc, *dayOfWeekDesc, *monthDesc, *yearDesc)

	description = cd.transformVerbosity(description)
	description = cd.transformCase(description)

	return &description, nil
}

func (cd *CronDescriptor) getTimeOfDayDescription() (*string, error) {
	secondExpr := cd.expressionArray[0]
	minuteExpr := cd.expressionArray[1]
	hourExpr := cd.expressionArray[2]

	specialChars := []string{"/", "-", ",", "*"}
	description := []string{}

	switch {

	case !cd.contains(secondExpr, specialChars) &&
		!cd.contains(minuteExpr, specialChars) &&
		!cd.contains(hourExpr, specialChars):
		description = append(description, "At ")
		formattedTime, err := cd.formatTime(hourExpr, minuteExpr, secondExpr)
		if err != nil {
			return nil, err
		}
		description = append(description, *formattedTime)

	case cd.contains(minuteExpr, []string{"-"}) &&
		!cd.contains(minuteExpr, []string{","}) &&
		!cd.contains(hourExpr, specialChars):
		splitMinute := strings.Split(minuteExpr, "-")
		if len(splitMinute) > 2 {
			return nil, ErrInvalidMinuteFormat
		}
		descrMin0Format, err := cd.formatTime(hourExpr, splitMinute[0], secondExpr)
		if err != nil {
			return nil, err
		}
		descrMin1Format, err := cd.formatTime(hourExpr, splitMinute[1], secondExpr)
		if err != nil {
			return nil, err
		}
		descFormat := fmt.Sprintf("Every minute between %s and %s", *descrMin0Format, *descrMin1Format)
		description = append(description, descFormat)

	case cd.contains(hourExpr, []string{","}) &&
		!cd.contains(hourExpr, []string{"-"}) &&
		!cd.contains(minuteExpr, specialChars):
		splitHour := strings.Split(hourExpr, ",")
		description = append(description, "At")

		emptySecondsField := ""
		for i, hourPart := range splitHour {
			description = append(description, " ")
			formattedTime, err := cd.formatTime(hourPart, minuteExpr, emptySecondsField)
			if err != nil {
				return nil, err
			}
			description = append(description, *formattedTime)

			if i < len(splitHour)-2 {
				description = append(description, ",")
			}

			if i == len(splitHour)-2 {
				description = append(description, " and")
			}
		}

	default:
		secondsDesc, err := cd.getSecondsDescription()
		if err != nil {
			return nil, err
		}
		minutesDesc, err := cd.getMinutesDescription()
		if err != nil {
			return nil, err
		}
		hoursDesc, err := cd.getHoursDescription()
		if err != nil {
			return nil, err
		}

		description = append(description, *secondsDesc)
		if strings.Join(description, "") != "" {
			description = append(description, ", ")
		}

		description = append(description, *minutesDesc)
		if strings.Join(description, "") != "" {
			description = append(description, ", ")
		}

		description = append(description, *hoursDesc)
	}

	finalDescription := strings.Join(description, "")
	return &finalDescription, nil
}

func (cd *CronDescriptor) getSecondsDescription() (*string, error) {
	expr := cd.expressionArray[0]
	cd.sugarLog.Debugf("getSecondsDescription working with expression %s", expr)

	segDesc, err := cd.getSegmentDescription(
		expr, Seconds,
		"every second",
		func(s string) (*string, error) {
			secondsInt, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			if secondsInt < 0 || secondsInt > 59 {
				return nil, ErrInvalidSecondsValue
			}
			return &s, nil
		},
		func(s string) string { return fmt.Sprintf("every %s seconds", s) },
		func(s string, t string) string { return fmt.Sprintf(", seconds %s through %s past the minute", s, t) },
		func(s string) string {
			return "at %s seconds past the minute"
		})

	if err != nil {
		return nil, err
	}
	return segDesc, nil
}

func (cd *CronDescriptor) getMinutesDescription() (*string, error) {
	expr := cd.expressionArray[1]
	cd.sugarLog.Debugf("getMinutesDescription working with expression %s", expr)

	segDesc, err := cd.getSegmentDescription(
		expr, Minutes,
		"every minute",
		func(s string) (*string, error) {
			return &s, nil
		},
		func(s string) string { return fmt.Sprintf("every %s minutes", s) },
		func(s string, t string) string { return fmt.Sprintf(", minutes %s through %s past the hour", s, t) },
		func(s string) string {
			return "at %s minutes past the hour"
		})

	if err != nil {
		return nil, err
	}
	return segDesc, nil
}

func (cd *CronDescriptor) getHoursDescription() (*string, error) {
	expr := cd.expressionArray[2]
	cd.sugarLog.Debugf("getHoursDescription working with expression %s", expr)

	if strings.Contains(expr, ":") {
		return nil, ErrInvalidCharacters
	}

	segDesc, err := cd.getSegmentDescription(
		expr,
		Hours,
		"every hour",
		func(s string) (*string, error) {
			hour := s
			minute := "0"
			if strings.Contains(s, ":") {
				hourSplit := strings.Split(s, ":")
				hour = hourSplit[0]
				minute = hourSplit[1]
			}
			return cd.formatTime(hour, minute, "")
		},
		func(s string) string { return fmt.Sprintf("every %s hours", s) },
		func(s string, t string) string { return fmt.Sprintf("%s through %s", s, t) },
		func(s string) string { return "at %s" })

	if err != nil {
		return nil, err
	}
	return segDesc, nil
}

func (cd *CronDescriptor) getDayOfTheWeekDescription() (*string, error) {
	expr := cd.expressionArray[5]
	dayOfWeekDesc := ""
	// day of the month is specified and day of the week is * so to prevent contradiction like "on day 1 of the month, every day"
	// we will not specified a day of the week description.
	if expr == "*" && cd.expressionArray[3] != "*" {
		return &dayOfWeekDesc, nil
	}

	cd.sugarLog.Debugf("getDayOfTheWeekDescription working with expression %s", expr)
	segDesc, err := cd.getSegmentDescription(
		expr,
		DayOfWeek,
		", every day",
		func(s string) (*string, error) {
			expr := s
			cd.sugarLog.Debugf("Check: \"%s\" contains a #: %t", s, cd.contains(s, []string{"#"}))
			if cd.contains(s, []string{"#"}) {
				splitStr := strings.Split(s, "#")
				expr = splitStr[0]
			} else if cd.contains(s, []string{"L"}) {
				expr = strings.ReplaceAll(s, "L", "")
			}

			dayInt, err := strconv.Atoi(expr)
			if err != nil {
				return nil, err
			}

			if dayInt >= len(CronDaysLong) {
				return nil, ErrInvalidDayOfWeekRange
			}

			return &CronDaysLong[dayInt], nil
		},
		func(s string) string { return fmt.Sprintf(", every %s days of the week", s) },
		func(s string, t string) string { return fmt.Sprintf(", %s through %s", s, t) },
		func(s string) string {
			formatted := ""
			cd.sugarLog.Debugf("Check: \"%s\" contains a #: %t", s, cd.contains(s, []string{"#"}))
			if cd.contains(s, []string{"#"}) {
				hashIndex := strings.Index(s, "#")
				dayOfWeekOfMonth := string(s[hashIndex+1])
				dayOfWeekOfMonthNum, err := strconv.Atoi(dayOfWeekOfMonth)
				if err != nil {
					cd.Logger.Panic(err.Error())
				}

				dayOrderMap := map[int]string{
					1: "first",
					2: "second",
					3: "third",
					4: "fourth",
					5: "fifth",
				}

				dayOrder, ok := dayOrderMap[dayOfWeekOfMonthNum]
				if !ok {
					cd.Logger.Panic(fmt.Sprintf("No day of the week of the month mapping for: %s, %s", s, dayOfWeekOfMonth))
				}
				preformatted := fmt.Sprintf(", on the %s", dayOrder)
				formatted = preformatted + " %s of the month"

			} else if cd.contains(s, []string{"L"}) {
				formatted = ", on the last %s of the month"
			} else {
				formatted = ", only on %s"
			}
			return formatted
		})

	if err != nil {
		return nil, err
	}

	return segDesc, nil
}

func (cd *CronDescriptor) getMonthDescription() (*string, error) {
	expr := cd.expressionArray[4]
	cd.sugarLog.Debugf("getMonthDescription working with expression %s", expr)

	segDesc, err := cd.getSegmentDescription(
		expr,
		Month,
		"",
		func(s string) (*string, error) {
			monthInt, err := strconv.Atoi(s)
			if err != nil {
				monthInt = cd.indexOf(CronMonths, s)
				if monthInt == -1 {
					cd.sugarLog.Panicf("Invalid month value for month: %s", s)
				}
				monthInt++
			}
			month := time.Date(time.Now().Year(), time.Month(monthInt), 1, 0, 0, 0, 0, time.Now().Location()).Month()
			retStr := month.String()
			return &retStr, nil
		},
		func(s string) string { return fmt.Sprintf(", every %s months", s) },
		func(s string, t string) string { return fmt.Sprintf(", %s through %s", s, t) },
		func(s string) string { return ", only in %s" })

	if err != nil {
		return nil, err
	}

	return segDesc, nil
}

func (cd *CronDescriptor) getDayOfMonthDescription() (*string, error) {
	var description string
	var dayOfMonthDesc *string

	expr := cd.expressionArray[3]
	cd.sugarLog.Debugf("getDayOfMonthDescription working with expression %s", expr)
	expr = strings.ReplaceAll(expr, "?", "*")

	if expr == "L" {
		description = ", on the last day of the month"
	} else if expr == "LW" || expr == "WL" {
		description = ", on the last weekday of the month"
	} else {
		re := regexp.MustCompile(`(\d{1,2}W)|(W\d{1,2})`)
		if re.Match([]byte(expr)) {
			foundStr := string(re.Find([]byte(expr)))
			foundStr = strings.ReplaceAll(foundStr, "W", "")
			dayNum, err := strconv.Atoi(foundStr)
			if err != nil {
				return nil, err
			}

			dayString := ""
			if dayNum == 1 {
				dayString = "first weekday"
			} else {
				dayString = fmt.Sprintf("weekday nearest day %d", dayNum)
			}
			description = fmt.Sprintf(", on the %s of the month", dayString)

		} else {
			segDesc, err := cd.getSegmentDescription(
				expr,
				DayOfMonth,
				", every day",
				func(s string) (*string, error) { return &s, nil },
				func(s string) string {
					if s == "1" {
						return ", every day"
					}
					return fmt.Sprintf(", every %s days", s)
				},
				func(s string, t string) string { return fmt.Sprintf(", between day %s and %s of the month", s, t) },
				func(s string) string { return ", on day %s of the month" })

			if err != nil {
				return nil, err
			}

			description = *segDesc
		}
	}

	dayOfMonthDesc = &description
	return dayOfMonthDesc, nil
}

func (cd *CronDescriptor) getYearDescription() (*string, error) {
	var yearDescription *string

	yearField := cd.expressionArray[6]
	cd.sugarLog.Debugf("getYearDescription working with expression %s", yearField)

	formatYear := func(y string) (*string, error) {
		re := regexp.MustCompile(`^\d+$`)
		if re.Find([]byte(y)) != nil {
			yearInt, err := strconv.ParseInt(y, 10, 64)
			if err != nil {
				return nil, err
			}
			if yearInt < 1900 {
				convertedYear := strconv.Itoa(int(yearInt))
				return &convertedYear, nil
			}
			year := time.Date(int(yearInt), time.January, 1, 0, 0, 0, 0, time.Now().Location()).Year()
			cd.sugarLog.Debugf("getYearDescription: year converted to time.Date: %d, %s, %d", int(yearInt), strconv.Itoa(year), year)
			convertedYear := strconv.Itoa(year)
			return &convertedYear, nil
		}

		yearInt, err := strconv.ParseInt(y, 10, 64)
		if err != nil {
			return nil, err
		}
		convertedYear := strconv.Itoa(int(yearInt))
		return &convertedYear, nil
	}

	yearDescription, err := cd.getSegmentDescription(yearField,
		Year,
		"",
		func(s string) (*string, error) {
			formattedYear, err := formatYear(s)
			return formattedYear, err
		},
		func(s string) string { return fmt.Sprintf(", every %s years", s) },
		func(s string, t string) string { return fmt.Sprintf(", %s through %s", s, t) },
		func(s string) string { return ", in %s" })

	if err != nil {
		return nil, err
	}

	return yearDescription, nil
}

func (cd *CronDescriptor) indexOf(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func (cd *CronDescriptor) contains(s string, subStrs []string) bool {
	for _, substr := range subStrs {
		if strings.Contains(s, substr) {
			return true
		}
	}

	return false
}

func (cd *CronDescriptor) getSegmentDescription(expr string,
	descriptionType DescriptionTypeEnum,
	allDescription string,
	getSingleItemDesc func(string) (*string, error),
	getIntervalDescFormat func(string) string,
	getBetweenDescFormat func(string, string) string,
	getDescFormat func(string) string) (*string, error) {

	var description *string
	var err error
	test := ""
	description = &test

	cd.sugarLog.Debugf("Current expression: %s", expr)

	switch {
	case expr == "":
		*description = ""

	case expr == "*":
		*description = allDescription

	case !cd.contains(expr, []string{"/", "-", ","}):
		cd.sugarLog.Debugf("Expression %s does not contain \"/\", \"-\", \",\"", expr)
		descFormat := getDescFormat(expr)
		cd.sugarLog.Debugf("Description format: %s", descFormat)
		singleItemDesc, err := getSingleItemDesc(expr)
		if err != nil {
			return nil, err
		}
		cd.sugarLog.Debugf("Single item format: %s ", *singleItemDesc)
		*description = fmt.Sprintf(descFormat, *singleItemDesc)
		cd.sugarLog.Debugf("description: %s ", *description)

	case cd.contains(expr, []string{"/"}):
		cd.sugarLog.Debugf("Expression %s contains \"/\"", expr)

		segments := strings.Split(expr, "/")
		*description = getIntervalDescFormat(segments[1])
		cd.sugarLog.Debugf("Interval description format: %s", *description)

		if cd.contains(segments[0], []string{"-"}) {
			betweenSegDescr, err := cd.generateBetweenSegmentDescription(segments[0], getBetweenDescFormat, getSingleItemDesc)
			if err != nil {
				return nil, err
			}

			if !strings.HasPrefix(*betweenSegDescr, ", ") {
				temp := strings.Join([]string{*description, ",", *betweenSegDescr}, "")
				description = &temp
			}

			temp := fmt.Sprintf("%s%s", *description, *betweenSegDescr)
			description = &temp

		} else if !cd.contains(segments[0], []string{"*", ","}) {
			// rangeItemDesc := getDescFormat(segments[0])
			rangeItemDesc := strings.ReplaceAll(segments[0], ", ", "")

			var temp string
			if descriptionType == Seconds || descriptionType == Minutes {
				temp = fmt.Sprintf("%s, starting at %s", *description, rangeItemDesc)
			} else {
				temp = fmt.Sprintf("%s, starting in %s", *description, rangeItemDesc)
			}

			description = &temp
		}

	case cd.contains(expr, []string{","}):
		cd.sugarLog.Debugf("Case: \"%s\" contains \",\"", expr)

		segments := strings.Split(expr, ",")
		descriptionContent := ""
		for i, segment := range segments {
			if i > 0 && len(segments) > 2 {
				descriptionContent += ","

				if i < len(segments)-1 {
					descriptionContent += " "
				}
			}

			if i > 0 && len(segments) > 1 && (i == len(segments)-1 || len(segments) == 2) {
				descriptionContent += " and "
			}

			if cd.contains(expr, []string{"-"}) {
				betweenDescription, err := cd.generateBetweenSegmentDescription(
					segment, getBetweenDescFormat, getSingleItemDesc)
				if err != nil {
					return nil, err
				}

				descriptionContent += strings.ReplaceAll(*betweenDescription, ", ", "")
			} else {
				singleItemDesc, err := getSingleItemDesc(segment)
				if err != nil {
					return nil, err
				}
				descriptionContent += *singleItemDesc
			}
		}
		descFormat := getDescFormat(expr)
		*description = fmt.Sprintf(descFormat, descriptionContent)

	case cd.contains(expr, []string{"-"}):
		cd.sugarLog.Debugf("case contains(expr, []string{\"-\"}):")
		description, err = cd.generateBetweenSegmentDescription(expr, getBetweenDescFormat, getSingleItemDesc)
		if err != nil {
			return nil, err
		}

	default:
		return description, ErrInvalidSegmentCase
	}

	return description, nil
}

func (cd *CronDescriptor) generateBetweenSegmentDescription(betweenExpr string,
	getBetweenDescFormat func(string, string) string,
	getSingleItemDesc func(string) (*string, error)) (*string, error) {

	var description *string
	defaultDescription := ""
	description = &defaultDescription

	cd.sugarLog.Debugf("Split %s by -", betweenExpr)
	betweenSegments := strings.Split(betweenExpr, "-")
	betweenSeg1Desc, err := getSingleItemDesc(betweenSegments[0])
	if err != nil {
		return nil, err
	}

	betweenSeg2Desc, err := getSingleItemDesc(betweenSegments[1])
	if err != nil {
		return nil, err
	}
	replaced := strings.ReplaceAll(*betweenSeg2Desc, ":00", ":59")
	betweenSeg2Desc = &replaced

	cd.sugarLog.Debugf("Two segments: %s, %s", *betweenSeg1Desc, *betweenSeg2Desc)

	temp := getBetweenDescFormat(*betweenSeg1Desc, *betweenSeg2Desc)
	description = &temp

	return description, nil
}

func (cd *CronDescriptor) formatTime(hourExpr string, minuteExpr string, secondExpr string) (*string, error) {
	hour, err := strconv.Atoi(hourExpr)
	if err != nil {
		return nil, err
	}

	period := ""
	if !cd.Options.Use24HourTimeFormat {
		if hour >= 12 {
			period = " PM"
		} else {
			period = " AM"
		}
		if hour > 12 {
			hour -= 12
		}
	}

	second := ""
	if secondExpr != "" {
		second = fmt.Sprintf(":%02s", secondExpr)
	}

	formattedTime := fmt.Sprintf("%02d:%02s%s%s", hour, minuteExpr, second, period)
	return &formattedTime, nil
}

func (cd *CronDescriptor) transformVerbosity(description string) string {
	if !cd.Options.Verbose {
		description = strings.ReplaceAll(description, ", every minute", "")
		description = strings.ReplaceAll(description, ", every hour", "")
		description = strings.ReplaceAll(description, ", every day", "")
		description = strings.ReplaceAll(description, "at 0 minutes past the hour, ", "")
	}

	return description
}

// CasingType for the type of casing transformatin
type CasingType int

const (
	// CasingTypeSentence style casing
	CasingTypeSentence CasingType = iota
	// CasingTypeTitle style casing
	CasingTypeTitle
	// CasingTypeLower style casing
	CasingTypeLower
)

func (cd *CronDescriptor) transformCase(description string) string {
	switch cd.Options.CasingType {
	case CasingTypeSentence:
		description = strings.ToUpper(string(description[0])) + string(description[1:])
	case CasingTypeTitle:
		description = strings.Title(description)
	case CasingTypeLower:
		fallthrough
	default:
		description = strings.ToLower(description)
	}
	return description
}
