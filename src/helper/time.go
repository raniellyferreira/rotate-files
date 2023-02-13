package helper

import "time"

// mustParseTime é uma função auxiliar para converter uma string no formato "YYYY-MM-DD HH:MM:SS" para time.Time
func MustParseTime(str string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		panic(err)
	}
	return t
}
