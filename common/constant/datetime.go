package constant

const (
	DateTimeLayout = "2006-01-02 15:04:05"
	DateLayout     = "2006-01-02"
	TimeLayout     = "15:04:05"

	DD_MM_YYYY_UL = "02/01/2006"
	DD_MM_YYYY_HL = "02-01-2006"

	TimeSecond     = 1
	TimeFiveSecond = 5 * TimeSecond
	TimeMinute     = 60 * TimeSecond
	TimeHour       = 60 * TimeMinute
	TimeDay        = 24 * TimeHour
	TimeWeek       = 7 * TimeDay
	Time30Days     = 30 * TimeDay
)
