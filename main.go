package main

import (
	"math/rand/v2"
	"time"
)

type Date struct {
	day   int
	month int
	year  int
}

var monthDays map[int][]int

func main() {
	monthDays = map[int][]int{
		1:  {31},
		2:  {28, 29},
		3:  {31},
		4:  {30},
		5:  {31},
		6:  {30},
		7:  {31},
		8:  {31},
		9:  {30},
		10: {31},
		11: {30},
		12: {31},
	}
}

func checkWeekday(date Date, guess string) bool {
	t := time.Date(date.year, time.Month(date.month), date.day, 0, 0, 0, 0, time.UTC)
	return t.Weekday().String() == guess
}

func newDate(startYear int, endYear int, monthDays map[int][]int) Date {
	// choose year
	year := rand.IntN(endYear+1-startYear) + startYear
	leapYear := isLeapYear(year)
	// choose month
	month := rand.IntN(12) + 1
	// choose day
	var day int
	if leapYear && month == 2 {
		day = rand.IntN(monthDays[month][1]) + 1
	} else {
		day = rand.IntN(monthDays[month][0]) + 1
	}

	return Date{day, month, year}
}

func isLeapYear(year int) bool {
	if year%100 == 0 && year%400 != 0 {
		return false
	}
	return year%4 == 0
}
