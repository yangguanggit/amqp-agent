package util

import (
	"amqp-agent/common/constant"
	"time"
)

/**
 * 获取时间
 */
func NowTime() time.Time {
	return time.Now()
}

/**
 * 获取本地时区
 */
func LocalTimeZone() *time.Location {
	return time.Local
}

/**
 * 使用本地时区解析时间
 */
func ParseTime(layout, datetime string) (time.Time, error) {
	return time.ParseInLocation(layout, datetime, time.Local)
}

/**
 * 获取当前日期时间
 */
func Now() string {
	return GetDateTime(0, constant.DateTimeLayout)
}

/**
 * 获取日期时间
 */
func GetDateTime(add int64, layout string) string {
	if add == 0 {
		return NowTime().Format(layout)
	}
	t := time.Unix(GetTimestamp(add), 0)
	return t.Format(layout)
}

/**
 * 获取时间戳
 */
func GetTimestamp(add int64) int64 {
	return NowTime().Unix() + add
}

/**
 * 根据时间戳获取日期时间
 */
func GetDateTimeByTimestamp(timestamp int64, layout string) string {
	t := time.Unix(timestamp, 0)
	return t.Format(layout)
}

/**
 * 根据日期时间获取时间戳
 */
func GetTimeStampByDateTime(datetime string, layout string) int64 {
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation(layout, datetime, loc)
	return t.Unix()
}
