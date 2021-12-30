package common

// Global Variables
const (
	// Report Level
	OK    = 0
	INFO  = 1
	DEBUG = 2
	ERROR = 3
)

var ReportingLevel = OK

var levelMap = map[int]string{
	OK:    "[   OK   ]   ",
	INFO:  "[  INFO  ]   ",
	DEBUG: "[  DEBUG  ]   ",
	ERROR: "[  ERROR  ]   ",
}
