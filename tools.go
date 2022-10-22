package libtools

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"golang.org/x/text/message"
)

const charset string = "abcdefghzkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ3456789" //随机因子

func GenerateRandomStr(length int) string {
	retStr := ""
	csLeng := len(charset)
	for i := 0; i < length; i++ {
		randNum := GenerateRandom(0, csLeng)
		retStr += string(charset[randNum])
	}

	return retStr
}

// 生成一个区间范围的随机数,左闭右开
func GenerateRandom(min, max int) int {
	if min >= max {
		return max
	}

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max - min)
	randNum += min

	return randNum
}

// 生成一个区间范围的随机数,左闭右开
func GenerateRandom64(min, max int64) int64 {
	if min >= max {
		return max
	}

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int63n(max - min)
	randNum += min

	return randNum
}

//! 手机验证在4-8位之间
func GenerateMobileCaptcha(length int) string {
	if length < 4 || length > 8 {
		return ""
	}

	minStr := "1" + strings.Repeat("0", length-1)
	maxStr := "1" + strings.Repeat("0", length)

	min, _ := strconv.Atoi(minStr)
	max, _ := strconv.Atoi(maxStr)

	captcha := GenerateRandom(min, max)
	return strconv.Itoa(captcha)
}

func GetCurrentEnv() string {
	runMode, _ := config.String("runmode")
	return runMode
}

func IsProductEnv() bool {
	return GetCurrentEnv() == "prod"
}

func EnvDisplay() string {
	var display string

	if IsProductEnv() {
		display = `生产`
	} else {
		display = `测试/开发`
	}

	return display
}

func DBDriver() string {
	dbType, _ := config.String("db_type")
	return dbType
}

func GetLocalUploadPrefix() string {
	uploadPrefix, _ := config.String("upload_prefix")
	return uploadPrefix
}

// CheckRequiredParameter 通用的检查必要参数的方法,只检测参数存在,不关心参数值
func CheckRequiredParameter(parameter map[string]interface{}, requiredParameter map[string]bool) bool {
	var requiredCheck int
	var rpCopy = make(map[string]bool)
	for rp, v := range requiredParameter {
		rpCopy[rp] = v
	}

	for k := range parameter {
		if requiredParameter[k] {
			requiredCheck++
			delete(rpCopy, k)
		}
	}

	if len(requiredParameter) != requiredCheck {
		var lostParam []string
		for l := range rpCopy {
			lostParam = append(lostParam, l)
		}
		logs.Error("request lost required parameter, parameter:", parameter, fmt.Sprintf("lostParam: [%s]", strings.Join(lostParam, ", ")))
		return false
	}

	return true
}

func ThreeElementExpression(status bool, exp1 interface{}, exp2 interface{}) (result interface{}) {
	if status {
		return exp1
	} else {
		return exp2
	}
}

func FullStack() string {
	var buf [2 << 11]byte
	runtime.Stack(buf[:], true)
	return string(buf[:])
}

func ClearOnSignal(handler func()) {
	signalChan := make(chan os.Signal, 1)

	// SIGINT  2  用户发送INTR字符(Ctrl+C)触发
	// SIGTERM 15 结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		handler()
		os.Exit(0)
	}()
}

// IntsSliceToWhereInString 将状态或者IDs集合转换为string
// interface{}支持所有int, int8 etc.: %d
func IntsSliceToWhereInString(intsSlice interface{}) (s string, err error) {
	sl, err := ToSlice(intsSlice)
	if err != nil {
		return
	}
	for _, i := range sl {
		s += fmt.Sprintf("%d,", i)
	}
	if len(s) > 0 {
		s = strings.TrimSuffix(s, ",")
	} else {
		err = fmt.Errorf("[IntsSliceToWhereInString] generate empty string, will occur sql error, with param %v", intsSlice)
	}
	return
}

// ToSlice 转化 泛型为 slice
func ToSlice(arr interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("[ToSlice] should be slice param, but %#v", arr)
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret, nil
}

//MobileFormat 去除空格，并且不能加拨0或62
func MobileFormat(mobile string) string {
	// 去除空格
	str := strings.Replace(mobile, " ", "", -1)
	if strings.HasPrefix(str, "08") {
		str = strings.Replace(str, "08", "8", 1)
	}
	if strings.HasPrefix(str, "628") {
		str = strings.Replace(str, "628", "8", 1)
	}
	return str
}

// NumberFormat 输出格式化数字，千分位以逗号分割
func NumberFormat(number interface{}) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprint(number)
}

// SliceInt64ToMap 输出格式化数字，千分位以逗号分割
func SliceInt64ToMap(s []int64) map[int64]interface{} {
	m := make(map[int64]interface{}, len(s))
	for _, v := range s {
		m[v] = nil
	}
	return m
}

// IsInMap template辅助方法 key 是否在map中
func IsInMap(m map[interface{}]interface{}, key interface{}) bool {
	//logs.Debug("m: %#v, key: %#v", m, key)
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

func IsInMapV2(key interface{}, m map[interface{}]interface{}) bool {
	//logs.Debug("m: %#v, key: %#v", m, key)
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

func DisplayByKey4Map(key interface{}, m map[interface{}]interface{}) (html string) {
	if v, ok := m[key]; ok {
		html = fmt.Sprintf(`%v`, v)
	}

	return
}

func MoneyDisplay(money int64) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

func MoneyDisplayInt64(money int64) int64 {
	return money / 100
}

func HumanMoney(money int64) string {
	str := MoneyDisplay(money)
	length := len(str)
	if length < 4 {
		return str
	}

	arr := strings.Split(str, ".") // 用小数点符号分割字符串,为数组接收
	length1 := len(arr[0])
	if length1 < 4 {
		return str
	}
	count := (length1 - 1) / 3

	for i := 0; i < count; i++ {
		arr[0] = arr[0][:length1-(i+1)*3] + "," + arr[0][length1-(i+1)*3:]
	}

	return strings.Join(arr, ".") // 将一系列字符串连接为一个字符串，之间用sep来分隔。
}

func CoinDisplay(money int64) string {
	return fmt.Sprintf("%d", (money)/100)
}

func CoinDisplayFloat64(money int64) float64 {
	return float64(money) / 100
}

func ScoreDisplay(score int) string {
	return fmt.Sprintf("%.1f", float64(score)/10)
}

func Float64TruncateWith2(num float64) (after float64, err error) {
	numStr := fmt.Sprintf(`%.02f`, num)
	after, err = Str2Float64(numStr)
	return
}

/*
func FeeRateTransform(origin float64) (after int64) {
	return int64(origin * types.FeeRateBase / 100)
}
*/

func FixCountNum(in int) int {
	var out = in
	if in < 0 {
		out = 0
	}

	return out
}

func AbsInt64(num int64) int64 {
	if num >= 0 {
		return num
	} else {
		return -num
	}
}

// MaxInt64 求最大值.
func MaxInt64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func MapSortKeyStr(params map[string]interface{}) string {
	paramLen := len(params)
	if paramLen <= 0 {
		logs.Warning("[MapSortKeyStr] params len is 0, params: %v", params)
		return ""
	}

	cntr := make([]string, paramLen)
	var i int = 0
	for k, _ := range params {
		cntr[i] = k
		i++
	}

	// 按字典序列排序
	sort.Strings(cntr)

	str := "" // 待签名字符串
	for i = 0; i < paramLen; i++ {
		key := cntr[i]
		str += fmt.Sprintf("%s=%s&", key, Stringify(params[key]))
	}
	//str += secret
	return str[0 : len(str)-1]
}

func GenOpenAvatar(openAvatar string) (avatar string) {
	if openAvatar != "" && strings.HasPrefix(openAvatar, "http") {
		avatar = openAvatar
	} else {
		avatar = `https://cdn-1302993108.cos.ap-guangzhou.myqcloud.com/img/default.png`
	}

	return
}

func StrLen(s string) int {
	return len(s)
}

func DivideRedPacket(count, amount int64) {
	//初始10个红包, 10000元钱, 单位是分
	//剩余金额
	remain := amount
	//验证红包算法的总金额,最后sum应该==amount
	sum := int64(0)
	//进行发红包
	for i := int64(0); i < count; i++ {
		x := DoubleAverage(count-i, remain)
		//金额减去
		remain -= x
		//发了多少钱
		sum += x
		//金额转成元
		fmt.Println(i+1, "=", float64(x)/float64(100))
	}
	fmt.Println()
	fmt.Printf("总和 %d分\n", sum)
}

//二倍均值算法,count剩余个数,amount剩余金额
func DoubleAverage(count, amount int64) int64 {
	//最小钱
	min := int64(1)

	if count == 1 {
		//返回剩余金额
		return amount
	}

	//计算最大可用金额,min最小是1分钱,减去的min,下面会加上,避免出现0分钱
	max := amount - min*count
	//计算最大可用平均值
	avg := max / count
	//二倍均值基础加上最小金额,防止0出现,作为上限
	avg2 := 2*avg + min
	//随机红包金额序列元素,把二倍均值作为随机的最大数
	rand.Seed(time.Now().UnixNano())
	//加min是为了避免出现0值,上面也减去了min
	x := rand.Int63n(avg2) + min
	return x
}
