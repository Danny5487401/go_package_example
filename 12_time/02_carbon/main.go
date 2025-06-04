package main

import (
	"github.com/dromara/carbon/v2"
	"time"
)

func main() {
	// Carbon 和 time.Time 互转
	// 将标准 time.Time 转换成 Carbon
	carbon.CreateFromStdTime(time.Now())
	// 将 Carbon 转换成标准 time.Time
	carbon.Now().StdTime()

	// 时间边界
	// 本年开始时间
	carbon.Parse("2020-08-05 13:14:15").StartOfYear().ToDateTimeString() // 2020-01-01 00:00:00
	// 本年结束时间
	carbon.Parse("2020-08-05 13:14:15").EndOfYear().ToDateTimeString() // 2020-12-31 23:59:59

	// 本季度开始时间
	carbon.Parse("2020-08-05 13:14:15").StartOfQuarter().ToDateTimeString() // 2020-07-01 00:00:00
	// 本季度结束时间
	carbon.Parse("2020-08-05 13:14:15").EndOfQuarter().ToDateTimeString() // 2020-09-30 23:59:59

	// 本月开始时间
	carbon.Parse("2020-08-05 13:14:15").StartOfMonth().ToDateTimeString() // 2020-08-01 00:00:00
	// 本月结束时间
	carbon.Parse("2020-08-05 13:14:15").EndOfMonth().ToDateTimeString() // 2020-08-31 23:59:59

	// 本周开始时间
	carbon.Parse("2020-08-05 13:14:15").StartOfWeek().ToDateTimeString()                                // 2020-08-02 00:00:00
	carbon.Parse("2020-08-05 13:14:15").SetWeekStartsAt(carbon.Sunday).StartOfWeek().ToDateTimeString() // 2020-08-02 00:00:00
	carbon.Parse("2020-08-05 13:14:15").SetWeekStartsAt(carbon.Monday).StartOfWeek().ToDateTimeString() // 2020-08-03 00:00:00
	// 本周结束时间
	carbon.Parse("2020-08-05 13:14:15").EndOfWeek().ToDateTimeString()                                // 2020-08-08 23:59:59
	carbon.Parse("2020-08-05 13:14:15").SetWeekStartsAt(carbon.Sunday).EndOfWeek().ToDateTimeString() // 2020-08-08 23:59:59
	carbon.Parse("2020-08-05 13:14:15").SetWeekStartsAt(carbon.Monday).EndOfWeek().ToDateTimeString() // 2020-08-09 23:59:59

	// 本日开始时间
	carbon.Parse("2020-08-05 13:14:15").StartOfDay().ToDateTimeString() // 2020-08-05 00:00:00
	// 本日结束时间
	carbon.Parse("2020-08-05 13:14:15").EndOfDay().ToDateTimeString() // 2020-08-05 23:59:59

	// 本小时开始时间
	carbon.Parse("2020-08-05 13:14:15").StartOfHour().ToDateTimeString() // 2020-08-05 13:00:00
	// 本小时结束时间
	carbon.Parse("2020-08-05 13:14:15").EndOfHour().ToDateTimeString() // 2020-08-05 13:59:59

	// 本分钟开始时间
	carbon.Parse("2020-08-05 13:14:15").StartOfMinute().ToDateTimeString() // 2020-08-05 13:14:00
	// 本分钟结束时间
	carbon.Parse("2020-08-05 13:14:15").EndOfMinute().ToDateTimeString() // 2020-08-05 13:14:59

	// 本秒开始时间
	carbon.Parse("2020-08-05 13:14:15").StartOfSecond().ToString() // 2020-08-05 13:14:15 +0800 CST
	// 本秒结束时间
	carbon.Parse("2020-08-05 13:14:15").EndOfSecond().ToString() // 2020-08-05 13:14:15.999999999 +0800 CST

	// 时间获取
	// 获取本年总天数
	carbon.Parse("2019-08-05 13:14:15").DaysInYear() // 365
	carbon.Parse("2020-08-05 13:14:15").DaysInYear() // 366
	// 获取本月总天数
	carbon.Parse("2020-02-01 13:14:15").DaysInMonth() // 29
	carbon.Parse("2020-04-01 13:14:15").DaysInMonth() // 30
	carbon.Parse("2020-08-01 13:14:15").DaysInMonth() // 31
}
