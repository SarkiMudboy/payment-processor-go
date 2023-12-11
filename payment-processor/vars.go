package main

import "time"

var UserFile File
var TransactionFile File
var CardFile File

var Period = map[string]time.Duration{
	"week":  168 * time.Hour,
	"month": 672 * time.Hour,
	"year":  8064 * time.Hour,
}
