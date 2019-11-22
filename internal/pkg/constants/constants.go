package constants

import "time"

const (
	DictNumeric      = "0123456789"
	DictUppercase    = "abcdefghijklmnopqrstuvwxyz"
	DictLowercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DictLetter       = DictLowercase + DictUppercase
	DicLetterNumeric = DictLetter + DictNumeric
)

var (
	MaxTime = time.Unix(1<<31-1, 0)
)
