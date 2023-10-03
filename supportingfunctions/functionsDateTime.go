package supportingfunctions

import "time"

func leapYears(date time.Time) (leaps int) {

	// returns year, month,
	// date of a time object
	y, m, _ := date.Date()

	if m <= 2 {
		y--
	}
	leaps = y/4 + y/400 - y/100
	return leaps
}

// The function calculates the
// difference between two dates and times
// and returns the days, hours, minutes,
// seconds between two dates

func GetDifference(a, b time.Time) (days, hours, minutes, seconds int) {
	if a.After(b) {
		a, b = b, a
	}

	// month-wise days
	monthDays := [12]int{31, 28, 31, 30, 31,
		30, 31, 31, 30, 31, 30, 31}

	// extracting years, months,
	// days of two dates
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()

	// extracting hours, minutes,
	// seconds of two times
	h1, min1, s1 := a.Clock()
	h2, min2, s2 := b.Clock()

	// totalDays since the
	// beginning = year*365 + number_of_days
	totalDays1 := y1*365 + d1

	// adding days of the months
	// before the current month
	for i := 0; i < (int)(m1)-1; i++ {
		totalDays1 += monthDays[i]
	}

	// counting leap years since
	// beginning to the year "a"
	// and adding that many extra
	// days to the totaldays
	totalDays1 += leapYears(a)

	// Similar procedure for second date
	totalDays2 := y2*365 + d2

	for i := 0; i < (int)(m2)-1; i++ {
		totalDays2 += monthDays[i]
	}

	totalDays2 += leapYears(b)

	// Number of days between two days
	days = totalDays2 - totalDays1

	// calculating hour, minutes,
	// seconds differences
	hours = h2 - h1
	minutes = min2 - min1
	seconds = s2 - s1

	// if seconds difference goes below 0,
	// add 60 and decrement number of minutes
	if seconds < 0 {
		seconds += 60
		minutes--
	}

	// performing similar operations
	// on minutes and hours
	if minutes < 0 {
		minutes += 60
		hours--
	}

	// performing similar operations
	// on hours and days
	if hours < 0 {
		hours += 24
		days--
	}

	return days, hours, minutes, seconds

}
