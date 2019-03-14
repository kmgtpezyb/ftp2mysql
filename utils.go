package main

import (
	"time"
)

func YesDay(ti time.Time, years, months, days int) string {
	
	newT := ti.AddDate(years,months,days)
	return newT.Format("20060102")
}
