package libtools

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

const (
	SecondAHour      int64 = 3600
	MillsSecondAHour       = SecondAHour * 1000
	SecondADay       int64 = 86400
	MillsSecondADay        = SecondADay * 1000
	MillsSecondAYear       = MillsSecondADay * 365
)

func GetDateFormat(timestamp int64, format string) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(timestamp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format(format)
}

func GetDate(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}

	tm := time.Unix(timestamp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("2006-01-02")
}

/** 获取时间计数 */
func GetTime(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(timestamp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("15:04:05")
}

/** 获取一个月的周期时间(毫秒) */
func GetMonthRange(timestamp int64) (begin, end int64) {
	tm := time.Unix(GetDateParse(GetDate(timestamp)), 0)
	bTime := tm.AddDate(0, 0, -tm.Day())
	eTime := tm.AddDate(0, 1, 0)
	return GetUnixMillisByTime(bTime), GetUnixMillisByTime(eTime)
}

func GetDateMH(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}

	tm := time.Unix(timestamp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("2006-01-02 15:04")
}

// 格式化毫秒时间
func MDateMH(timestamp int64) string {
	return GetDateMH(timestamp / 1000)
}

func GetDateMHS(timestamp int64) string {
	if timestamp <= 0 {
		return "-"
	}

	tm := time.Unix(timestamp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("2006-01-02 15:04:05")
}

func RFC3339TimeTransfer(datetime string) int64 {

	timeLayout := "2006-01-02T15:04:05Z" //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区

	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix() * 1000 //转化为时间戳 类型是int64

	return timestamp
}

func RFC3339TransferDate(str string) string {
	t1, _ := time.Parse(time.RFC3339, str)
	return UnixMsec2Date(t1.UnixMilli(), "Y-m-d")
}

func MDateMHSLocalDate(timestamp int64) string {
	tmp := timestamp / 1000

	if tmp <= 0 {
		return "-"
	}

	tm := time.Unix(tmp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("20060102")
}

func MDateMHSLocalDateAllNum(timestamp int64) string {
	tmp := timestamp / 1000

	if tmp <= 0 {
		return "-"
	}

	tm := time.Unix(tmp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("20060102150405")
}

func LocalYearMonth(timestamp int64) string {
	tmp := timestamp / 1000

	if tmp <= 0 {
		return "-"
	}

	tm := time.Unix(tmp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("200601")
}

func DateMHSZ(timestamp int64) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(timestamp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("2006-01-02")
}

// 毫秒
func MDateUTC(timestamp int64) string {
	return DateMHSZ(timestamp / 1000)
}

/*
*   从一种时间格式转为另一种 ，或者转为时间戳
*	@param timestr 即将处理的时间字符串
*	@param fromFormat 当前时间格式  Mon, 02 Jan 2006  MST
*	@param toFormat 目标时间格式   	2006-01-02 15:04:05
*	@param fromFormat 当前时间格式
*	@param unixtime 为真返回时间戳，否则正常转换时间格式
*	@return string []byte
 */
func TimeStrFormat(timestr, fromFormat, toFormat string, unixtime bool) interface{} {
	timeparse, _ := time.Parse(fromFormat, timestr)
	timestsmp := timeparse.Unix()
	if unixtime {
		return timestsmp
	}
	tm := time.Unix(timestsmp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format(toFormat)

}

// GetDateParse 用于跑批, 或者需要以 UTC时区为基准的时间解析
func GetDateParse(dates string) int64 {
	if "" == dates {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02", dates, loc)
	return parse.Unix()
}

// GetDateParse 用于跑批, 或者需要以 UTC时区为基准的时间解析
func GetDateParses(dates string) int64 {
	if "" == dates {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04:05", dates, loc)
	return parse.Unix()
}

// Str2TimeByLayout 使用layout将时间字符串转unix时间戳(毫秒)
func Str2TimeByLayout(layout, timeStr string) int64 {
	if "" == timeStr {
		return 0
	}

	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation(layout, timeStr, loc)
	return parse.UnixNano() / 1000000
}

// DateParseYMDHMS 解析 YYYY-MM-DD HH:MM:SS 格式的时间串为Unix时间戳
func DateParseYMDHMS(dates string) int64 {
	if "" == dates {
		return 0
	}

	local, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04:05", dates, local)

	return parse.Unix()
}

// 毫秒,输出北京时间
func MDateMHSBeijing(timestamp int64) string {
	tmp := timestamp / 1000

	if tmp <= 0 {
		return "-"
	}

	tm := time.Unix(tmp, 0)
	local, _ := time.LoadLocation("Asia/Shanghai")
	return tm.In(local).Format("2006-01-02 15:04:05")
}

// ParseDateRangeToDayRange 将时间范围字符串解析成毫秒时间戳
// 默认日期分隔符 " - "
// start, end, err
func ParseDateRangeToDayRange(dateRange string) (start, end int, err error) {
	splitSep := " - "
	start, end, err = ParseDateRangeToDayRangeWithSep(dateRange, splitSep)
	return
}

// PareseDateRangeToDayRangeWithSep 将时间范围字符串解析成毫秒时间戳
// start, end, err
func ParseDateRangeToDayRangeWithSep(dateRange string, splitSep string) (int, int, error) {
	if len(dateRange) == 0 {
		// 后台正常逻辑, 因此不记录log, 只是返回err, 便于处理
		return 0, 0, errors.New("Empty date range, just ignore it")
	}

	tr := strings.Split(dateRange, splitSep)
	if (len(tr)) != 2 {
		err := fmt.Errorf("[PareseDateRangeToMillsecondWithCustomSep][wrong date range format], (%s) cantnot split to 2 date by (%s)",
			dateRange, splitSep)
		logs.Error(err)
		return 0, 0, err
	}

	start, _ := strconv.Atoi(strings.Replace(tr[0], "-", "", -1))
	end, _ := strconv.Atoi(strings.Replace(tr[1], "-", "", -1))

	if start <= 0 || end <= 0 {
		err := fmt.Errorf("[PareseDateRangeToMillsecondWithCustomSep][wrong date range format], (%s) cantnot split to 2 format date like 2006-01-02",
			dateRange)
		logs.Error(err)
		return 0, 0, err
	}

	return start, end, nil
}

// 取当前系统时间的毫秒
func GetUnixMillis() int64 {
	return GetUnixMillisByTime(time.Now())
}

func GetUnixMillisByTime(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

func TimeNow() int64 {
	return time.Now().Unix()
}

func NaturalDay(offset int64) (um int64) {
	t := time.Now()
	date := GetDate(t.Unix())
	baseUm := GetDateParse(date) * 1000
	offsetUm := MillsSecondADay * offset

	um = baseUm + offsetUm

	return
}

/*
*
基于指定时间的偏移量
*/
func BaseDayOffset(baseDay int64, offset int64) (um int64) {
	date := GetDate(baseDay / 1000)
	baseUm := GetDateParse(date) * 1000
	offsetUm := MillsSecondADay * offset
	um = baseUm + offsetUm
	return
}

func GetDateRange(begin, end int64) int64 {
	return (end - begin) / SecondADay
}

func GetDateRangeMillis(begin, end int64) int64 {
	return (end - begin) / MillsSecondADay
}

// 返回的单位是秒
func GetMonth(timetag int64) int64 {
	dateStr := GetDateFormat(timetag/1000, "2006-01-02")
	dateStr = dateStr[0:len(dateStr)-2] + "01"

	return GetDateParse(dateStr)
}

// 毫秒,输出本地时间
func MDateMHS(timestamp int64) string {
	tmp := timestamp / 1000

	if tmp <= 0 {
		return "-"
	}

	tm := time.Unix(tmp, 0)
	local, _ := time.LoadLocation("Local")
	return tm.In(local).Format("2006-01-02 15:04:05 MST")
}

// GetDateParseBackend 所有后台使用
func GetDateParseBackend(dates string) int64 {
	if "" == dates {
		return 0
	}

	local, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02", dates, local)

	return parse.Unix()
}

/** 获取一天的0点0分0秒 */
func GetDateTimeByBegin(t int64) int64 {
	tm := time.Unix(t/1000, 0)
	local, _ := time.LoadLocation("Local")
	var begin = time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, local)
	return begin.Unix()
}

/** 获取一天的固定时间的毫秒数 h 24*/
func GetHourDateTime(t int64, h int) int64 {
	tm := time.Unix(t/1000, 0)
	local, _ := time.LoadLocation("Local")
	var begin = time.Date(tm.Year(), tm.Month(), tm.Day(), h, 0, 0, 0, local)
	return begin.UnixNano() / 1000000
}

/** 获取过去时中最近的5分数 */
func GetDateTimeBy5step(t int64) int64 {
	tm := time.Unix(t/1000, 0)
	local, _ := time.LoadLocation("Local")
	var begin = time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute()-tm.Minute()%5, 0, 0, local)
	return begin.Unix()
}

func GetDateTimeParseBackend(dates string) int64 {
	if "" == dates {
		return 0
	}

	local, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04:05", dates, local)

	return parse.Unix()
}

func Default7DaysTimeRange() string {
	last7days := NaturalDay(-7)
	return fmt.Sprintf(`%s - %s`, DateMHSZ(last7days/1000), DateMHSZ(GetUnixMillis()/1000))
}

func DefaultTodayTimeRange() string {
	now := GetUnixMillis()
	return fmt.Sprintf(`%s - %s`, DateMHSZ(now/1000), DateMHSZ(now/1000))
}

func DefaultYesterdayTimeRange() string {
	now := NaturalDay(-1)
	return fmt.Sprintf(`%s - %s`, DateMHSZ(now/1000), DateMHSZ(now/1000))
}

func DefaultTodayMHS() string {
	now := GetUnixMillis()
	return MDateMHSLocalDate(now)
}

func DefaultToday() string {
	now := GetUnixMillis()
	return DateMHSZ(now / 1000)
}

func DefaultYesterday() string {
	now := NaturalDay(-1)
	return DateMHSZ(now / 1000)
}

func GetTimeByTodaySecs(secs int) string {
	today := DateMHSZ(TimeNow())
	todayUnix := GetDateParse(today)
	t := todayUnix + int64(secs)
	todayTime := GetDateMHS(t)
	todayTimeArr := strings.Split(todayTime, " ")
	return todayTimeArr[1]
}

func HumanUnixMillis(t int64) (display string) {
	t = t / 1000

	var second int64 = 1
	var minute = 60 * second
	var oneHour = minute * 60
	var oneDay = oneHour * 24
	var oneWeek = oneDay * 7
	var oneMonth = oneDay * 30
	var oneYear = oneDay * 365

	var box []string
	if t >= oneYear {
		y := t / oneYear
		box = append(box, fmt.Sprintf(`%d year(s)`, y))
		t -= y * oneYear
	}
	if t >= oneMonth {
		m := t / oneMonth
		box = append(box, fmt.Sprintf(`%d month(s)`, m))
		t -= m * oneMonth
	}
	if t >= oneWeek {
		w := t / oneWeek
		box = append(box, fmt.Sprintf(`%d week(s)`, w))
		t -= w * oneWeek
	}
	if t >= oneHour {
		h := t / oneHour
		box = append(box, fmt.Sprintf(`%d hour(s)`, h))
		t -= h * oneHour
	}
	if t >= minute {
		m := t / minute
		box = append(box, fmt.Sprintf(`%d minute(s)`, m))
		t -= m * minute
	}

	if t > 0 {
		box = append(box, fmt.Sprintf(`%d second(s)`, t))
	}

	if len(box) > 0 {
		display = strings.Join(box, ", ")
	}

	return
}

func HumanUnixMillisV2(t int64) (display string) {
	t = t / 1000

	var second int64 = 1
	var minute = 60 * second
	var oneHour = minute * 60
	var oneDay = oneHour * 24
	var oneWeek = oneDay * 7
	var oneMonth = oneDay * 30
	var oneYear = oneDay * 365

	var box []string
	if t >= oneYear {
		y := t / oneYear
		box = append(box, fmt.Sprintf(`%d year(s)`, y))
		t -= y * oneYear
	}
	if t >= oneMonth {
		m := t / oneMonth
		box = append(box, fmt.Sprintf(`%d month(s)`, m))
		t -= m * oneMonth
	}
	if t >= oneWeek {
		w := t / oneWeek
		box = append(box, fmt.Sprintf(`%d week(s)`, w))
		t -= w * oneWeek
	}
	if t >= oneHour {
		h := t / oneHour
		box = append(box, fmt.Sprintf(`%02d`, h))
		t -= h * oneHour
	} else {
		box = append(box, "00")
	}
	if t >= minute {
		m := t / minute
		box = append(box, fmt.Sprintf(`%02d`, m))
		t -= m * minute
	} else {
		box = append(box, "00")
	}

	if t > 0 {
		box = append(box, fmt.Sprintf(`%02d`, t))
	} else {
		box = append(box, `00`)
	}

	if len(box) > 0 {
		display = strings.Join(box, ":")
	}

	return
}

func CalculateAgeByBirthday(birthday string) int {
	exp := strings.Split(birthday, "-")
	if len(exp) < 1 {
		return 0
	}

	year, _ := Str2Int(exp[0])
	age := time.Now().Year() - year
	if age < 0 {
		age = 0
	}
	return age
}

// 针对 golang 的时间函数库难记难用,封装以下两个函数,采用共识标识符来简化原始库的使用 {{{
// millisecond <-> msec
// see: https://www.php.net/manual/zh/function.date.php
// 采用类 linux 时间格式
// 仅取以下值:
// 日: d, D, l, j
// 月: m, M, n
// 年:  Y, y
// 时间: a, H, i, s
// 时区: e
var (
	find = []string{
		`a`, `M`, `n`, // 需要优先替换,否则出现误替换
		`d`, `D`, `l`, `j`,
		`m`,
		`Y`, `y`,
		`H`, `i`, `s`,
		`e`,
	}

	replace = []string{
		`3:04PM`, `Jan`, `1`,
		`02`, `Mon`, `Monday`, `2`,
		`01`,
		`2006`, `06`,
		`15`, `04`, `05`,
		`MST`,
	}
)

func UnixMsec2Date(um int64, layout string) string {
	timestamp := um / 1000
	if timestamp <= 0 {
		return `-`
	}

	tm := time.Unix(timestamp, 0)
	local, _ := time.LoadLocation("Local")

	for i, f := range find {
		layout = strings.Replace(layout, f, replace[i], -1)
	}

	//logs.Debug("[UnixMsec2Date] layout: %s", layout)
	return tm.In(local).Format(layout)
}

func Date2UnixMsec(dateStr, layout string) int64 {
	if "" == dateStr {
		return 0
	}

	for i, f := range find {
		layout = strings.Replace(layout, f, replace[i], -1)
	}

	loc, _ := time.LoadLocation("Local")
	parse, err := time.ParseInLocation(layout, dateStr, loc)
	if err != nil {
		logs.Error("[Date2UnixMsec] parse layout get exception, layout: %s, err: %v", layout, err)
		return 0
	}

	return parse.UnixNano() / 1000000
}

func Int642Time(t int64) time.Time {
	return time.Unix(t/1000, 0)
}

func GetMonthLastDay(t time.Time) int {
	lastTime := t.AddDate(0, 1, -1)
	_, _, d := lastTime.Date()
	return d
}

// }}}

func ExcelConvertToFormatDay(excelDaysString string) string {
	// 2006-01-02 距离 1900-01-01的天数
	baseDiffDay := 38719 //在网上工具计算的天数需要加2天，什么原因没弄清楚
	curDiffDay := excelDaysString
	b, _ := strconv.Atoi(curDiffDay)
	// 获取excel的日期距离2006-01-02的天数
	realDiffDay := b - baseDiffDay
	//fmt.Println("realDiffDay:",realDiffDay)
	// 距离2006-01-02 秒数
	realDiffSecond := realDiffDay * 24 * 3600
	//fmt.Println("realDiffSecond:",realDiffSecond)
	// 2006-01-02 15:04:05距离1970-01-01 08:00:00的秒数 网上工具可查出
	baseOriginSecond := 1136185445
	resultTime := time.Unix(int64(baseOriginSecond+realDiffSecond), 0).Format("2006-01-02")
	return resultTime
}

/**
 * @Description 获得当前月的初始和结束日期
 **/
func GetMonthDay() (string, string) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	f := firstOfMonth.Unix()
	l := lastOfMonth.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

/**
 * @Description 获得当前周的初始和结束日期
 **/
func GetWeekDay() (string, string) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	lastoffset := int(time.Saturday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if lastoffset == 6 {
		lastoffset = -1
	}

	firstOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	lastOfWeeK := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, lastoffset+1)
	f := firstOfWeek.Unix()
	l := lastOfWeeK.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

/**
 * @Description //获得当前季度的初始和结束日期
 * @return
 **/
func GetQuarterDay() (string, string) {
	year := time.Now().Format("2006")
	month := int(time.Now().Month())
	var firstOfQuarter string
	var lastOfQuarter string
	if month >= 1 && month <= 3 {
		//1月1号
		firstOfQuarter = year + "-01-01 00:00:00"
		lastOfQuarter = year + "-03-31 23:59:59"
	} else if month >= 4 && month <= 6 {
		firstOfQuarter = year + "-04-01 00:00:00"
		lastOfQuarter = year + "-06-30 23:59:59"
	} else if month >= 7 && month <= 9 {
		firstOfQuarter = year + "-07-01 00:00:00"
		lastOfQuarter = year + "-09-30 23:59:59"
	} else {
		firstOfQuarter = year + "-10-01 00:00:00"
		lastOfQuarter = year + "-12-31 23:59:59"
	}
	return firstOfQuarter, lastOfQuarter
}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
func GetBetweenDates(sdate, edate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(sdate) {
		timeFormatTpl = timeFormatTpl[0:len(sdate)]
	}
	date, err := time.Parse(timeFormatTpl, sdate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, edate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}
