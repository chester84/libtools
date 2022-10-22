package libtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/beego/beego/v2/core/logs"
	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 字串截取
func SubString(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}

// 字串截取
func SubStringByPos(s string, sPos, ePos int) string {
	runes := []rune(s)
	if ePos > len(runes) {
		ePos = len(runes)
	}

	return string(runes[sPos:ePos])
}

func Strim(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)

	return str
}

// StrReplace 在 origin 中搜索 search 组,替换成 replace
func StrReplace(origin string, search []string, replace string) (s string) {
	s = origin
	for _, find := range search {
		s = strings.Replace(s, find, replace, -1)
	}

	return
}

func TrimRealName(name string) string {
	// 将名字里面的标点符号替换成1个空格

	name = ReplaceInvalidRealName(name)
	reg := regexp.MustCompile(`[\pP]+?`)
	name = reg.ReplaceAllString(name, " ")

	// 将连续的空白替换成一个空格
	reg = regexp.MustCompile(`\s{2,}`)
	name = reg.ReplaceAllString(name, " ")

	name = strings.TrimSpace(name)

	return name
}

func ReplaceInvalidRealName(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z\\s]+")
	ret := reg.ReplaceAllString(name, "")
	return ret
}

// IsIndonesiaName 判断是否是合法的印尼名字
func IsIndonesiaName(name string) (valid bool) {
	var exp, _ = regexp.Compile(`^[a-zA-Z ]+$`)
	if exp.MatchString(name) {
		return true
	}
	return false
}

// IsNumber 判断是否都是数字
func IsNumber(str string) (valid bool) {
	var exp, _ = regexp.Compile(`^[0-9]+$`)
	if exp.MatchString(str) {
		return true
	}
	return false
}

func ContainNumber(str string) (valid bool) {
	var exp, _ = regexp.Compile(`[\d]`)
	if exp.MatchString(str) {
		return true
	}
	return false
}

/** 将数组默认的0，转为1 */
func IndexNumber(index int) int {
	return index + 1
}

func Str2Int64(str string) (int64, error) {
	number, err := strconv.ParseInt(str, 10, 64)
	return number, err
}

func Int642Str(number int64) string {
	return strconv.FormatInt(number, 10)
}

func Str2Int(str string) (int, error) {
	number, err := strconv.ParseInt(str, 10, 0)
	return int(number), err
}

func Int2Str(number int) string {
	return strconv.FormatInt(int64(number), 10)
}

func Float2Str(f float32) string {
	return Float642Str(float64(f))
}

func Float642Str(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Str2Float64(s string) (f float64, err error) {
	f, err = strconv.ParseFloat(s, 64)

	return
}

func Str2Float(s string) (f float32, err error) {
	f64, err := Str2Float64(s)
	f = float32(f64)

	return
}

func JsonEncode(d interface{}) (jsonStr string, err error) {
	bson, err := json.Marshal(d)
	jsonStr = string(bson)

	return
}

func Unicode(rs string) string {
	jsonStr := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			jsonStr += string(r)
		} else {
			jsonStr += "\\u" + strconv.FormatInt(int64(rint), 16)
		}
	}

	return jsonStr
}

func Escape(html string) string {
	return template.HTMLEscapeString(html)
}

func AddSlashes(str string) string {
	str = strings.Replace(str, `\`, `\\`, -1)
	str = strings.Replace(str, "'", `\'`, -1)
	str = strings.Replace(str, `"`, `\"`, -1)

	return str
}

func StripSlashes(str string) string {
	str = strings.Replace(str, `\'`, `'`, -1)
	str = strings.Replace(str, `\"`, `"`, -1)
	str = strings.Replace(str, `\\`, `\`, -1)

	return str
}

func RawUrlEncode(s string) (r string) {
	r = UrlEncode(s)
	r = strings.Replace(r, "+", "%20", -1)
	return
}

// 直接json.Marshal ，  会把 < > & 转成 unicode 编码
// JSONMarshal 解决直接json.Marshal 后单引号，双引号，< > & 符号的问题
func JSONMarshal(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)

	return buf.Bytes(), err
}

//将slice 转化成字符串
//[]int{1, 2, 3, 4, 5}  => 1,2,3,4,5 或
//[]string{"1", "2", "3", "4", "5"}  => 1,2,3,4,5
//"AAA bbb" 转为  AAA,bbb
//其他类型返回空字符串
func ArrayToString(a interface{}, delim string) (newStr string) {
	vtype := reflect.TypeOf(a).String()
	if vtype == "[]int" || vtype == "[]int64" {
		newStr = strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	} else if vtype == "[]string" {

		// 单独处理[]string类型是因为，我们要处理的字符串值可能有空格，但是也会被当成需要分格的对象，所以先替换处理，然后再替换回去
		specialChar := "^^"
		newSlice := make([]string, 0)
		for _, v := range a.([]string) {
			newSlice = append(newSlice, strings.Replace(v, " ", specialChar, -1))
		}
		fileds := strings.Fields(fmt.Sprint(newSlice))
		join := strings.Replace(strings.Join(fileds, delim), specialChar, " ", -1)
		newStr = strings.Trim(join, "[]")
	} else if vtype == "string" {
		newStr = strings.Replace(fmt.Sprint(a.(string)), " ", delim, -1)
	} else {
		newStr = ""
	}
	return
}

func GetIntKeysOfMap(mymap map[int]string) (keys []int) {
	keys = make([]int, 0, len(mymap))
	for k := range mymap {
		keys = append(keys, k)
	}
	return
}

// IsValidIndonesiaMobile 是否是印尼的有效电话号
// 08 开头, 10-13位数字, 2018.08,有13位手机号段了
// 2018.11 增加15位手机号判断
func IsValidIndonesiaMobile(mobile string) (yes bool, err error) {
	mobileLen := len(mobile)
	if mobileLen < 10 || mobileLen > 15 {
		err = fmt.Errorf("mobile length is invalid, mobile: %s", mobile)
		return
	}

	if "08" != SubString(mobile, 0, 2) {
		err = fmt.Errorf("mobile is invalid, mobile: %s", mobile)
		return
	}

	_, errNum := Str2Int64(mobile)
	if errNum != nil {
		err = fmt.Errorf("mobile has invalid char, mobile: %s, errNum: %v", mobile, errNum)
		return
	}

	yes = true

	return yes, nil
}

// ParseTableName 从SQL语句中解析主表名
func ParseTableName(sql string) (name string, err error) {
	re := regexp.MustCompile(`(?i).+FROM\s+(\S+)`)
	allMatch := re.FindAllStringSubmatch(sql, -1)

	if len(allMatch) == 0 {
		err = fmt.Errorf("parse sql has error.")
		return
	}

	tableName := strings.Replace(allMatch[0][1], "`", "", -1)
	names := strings.Split(tableName, ".")
	name = names[len(names)-1]

	return
}

// 手机号脱敏处理(只保留前两位和后四位，中间的每个字符都替换为"*")
// 例如：08123456789 改为 08*****6789, 0812345645678 改为 08*******5678
func MobileDesensitization(src string) (dst string) {
	length := len(src)

	if length >= 7 {
		prefix := src[0:3]
		var middle []string
		for i := 0; i < length-7; i++ {
			middle = append(middle, "*")
		}
		suffix := src[length-4 : length]

		dst = fmt.Sprintf("%s%s%s", prefix, strings.Join(middle, ""), suffix)
	}

	return
}

func RealNameMask(realName string) (maskStr string) {
	words := ([]rune)(realName)
	wordsLen := len(words)

	if wordsLen <= 1 {
		maskStr = realName
	} else if wordsLen == 2 {
		maskStr = fmt.Sprintf("%s*", string(words[0:1]))
	} else {
		var star string
		for i := 0; i < wordsLen-2; i++ {
			star += "*"
		}
		maskStr = fmt.Sprintf("%s%s%s", string(words[0:1]), star, string(words[wordsLen-1:]))
	}

	return
}

func ParseTargetList(str string) []string {
	list := make([]string, 0)
	if str == "" {
		return list
	}

	listStr := strings.Split(str, "\n")

	c := ","
	if strings.Contains(listStr[0], "\r") {
		c = "\r"
	} else if strings.Contains(listStr[0], ",") {
		c = ","
	}

	if len(listStr) == 1 {
		vec := strings.Split(listStr[0], c)

		list = append(list, strings.Trim(vec[0], " "))

		return list
	}

	for _, v := range listStr {
		vec := strings.Split(v, c)
		if len(vec) < 1 {
			continue
		}

		list = append(list, strings.Trim(vec[0], " "))
	}

	return list
}

func StrTenThousand2MoneyMul100(money string) int64 {
	price, _ := DecimalMoneyMul100(money)
	price *= 10000

	return price
}

// DecimalMoneyMul100 将带小数点的金额转换成数据库中的乘以100后的整数
func DecimalMoneyMul100(moneyStr string) (money int64, err error) {
	//moneyFloat64, err := Str2Float64(moneyStr)
	//if err != nil {
	//	return
	//}
	//
	//money = int64(moneyFloat64 * 100.0)

	moneyBig, err := decimal.NewFromString(moneyStr)
	if err != nil {
		logs.Warning("[DecimalMoneyMul100] convert to math big exception, str: %s, err: %v", moneyStr, err)
		return
	}

	moneyMul := moneyBig.Mul(decimal.NewFromInt(100))
	money = moneyMul.IntPart()

	return
}

func SecretKeyMask(secretKey string) (str string) {
	keyLen := len(secretKey)
	if keyLen > 9 {
		str = fmt.Sprintf(`%s***%s`, SubString(secretKey, 0, 3), SubString(secretKey, keyLen-3, 3))
	} else {
		str = secretKey
	}

	return
}

/**
  13:23
*/
func ConvertTime2Secs(str string) int {
	arr := strings.Split(str, ":")
	hours, _ := Str2Int(arr[0])
	mins, _ := Str2Int(arr[1])
	return hours*60*60 + mins*60
}

// Snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// Camel string, xx_yy to XxYy
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	flag, num := true, len(s)-1
	for i := 0; i <= num; i++ {
		d := s[i]
		if d == '_' {
			flag = true
			continue
		} else if flag {
			if d >= 'a' && d <= 'z' {
				d = d - 32
			}
			flag = false
		}
		data = append(data, d)
	}
	return string(data[:])
}

func SimpleNumStr2Num(s string) int {
	var ret int
	var base float64 = 1
	if strings.Contains(s, "w") || strings.Contains(s, "W") {
		base = 10000
		s = StrReplace(s, []string{"w", "W"}, "")
	}

	num, err := Str2Float64(s)
	if err != nil {
		logs.Warning("[SimpleNumStr2Num] str to num err: %v", err)
	}

	ret = int(num * base)

	return ret
}

func ExtractUrls(s string) ([]string, error) {
	var urlBox []string
	var err error

	findUrl := regexp.MustCompile(`(http\S*)`)
	allMatch := findUrl.FindAllStringSubmatch(s, -1)
	//log.Printf("[ExtractUrl] allMatch: %#v\n", allMatch)

	if len(allMatch) > 0 {
		for _, fu := range allMatch {
			urlBox = append(urlBox, fu[0])
		}
	} else {
		err = fmt.Errorf(`can not extract url from input: %s`, s)
	}

	return urlBox, err
}

func Gbk2Utf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 4位字符串版本 1.1234.1234.1234
func AppNumVersion(appVersion string) int64 {
	exp := strings.Split(appVersion, ".")
	if len(exp) != 4 {
		logs.Error("invalid app version, appVersion: %s", appVersion)
		return 0
	}

	numBox := make([]int64, 4)
	for i, subVer := range exp {
		num, _ := Str2Int64(subVer)

		if i > 0 {
			num = num % 10000
		}

		numBox[i] = num
	}

	verStr := fmt.Sprintf(`%d%04d%04d%04d`, numBox[0], numBox[1], numBox[2], numBox[3])
	numVersion, _ := Str2Int64(verStr)

	return numVersion
}

// 3 位版本号 1.2.3
func AppNumVersion3(appVersion string) int64 {
	exp := strings.Split(appVersion, ".")
	if len(exp) != 3 {
		logs.Error("invalid app version, appVersion: %s", appVersion)
		return 0
	}

	numBox := make([]int64, 4)
	for i, subVer := range exp {
		num, _ := Str2Int64(subVer)

		if i > 0 {
			num = num % 10000
		}

		numBox[i] = num
	}

	verStr := fmt.Sprintf(`%d%04d%04d`, numBox[0], numBox[1], numBox[2])
	numVersion, _ := Str2Int64(verStr)

	return numVersion
}

func MobileMask(mobile string) string {
	if len(mobile) != 11 {
		logs.Warning("input wrong mobile: %s", mobile)
		return ""
	}

	return fmt.Sprintf(`%s%s%s`, mobile[0:3], strings.Repeat(`*`, 5), mobile[8:])
}

func RegRemoveScript(in string) string {
	reg, _ := regexp.Compile(`\<(?i)script[\S\s]+?\</(?i)script\>`)
	in = reg.ReplaceAllString(in, "")
	return in
}

func NicknameMask(nickname string) string {
	var s = nickname

	rStr := []rune(nickname)
	rLen := len(rStr)

	if rLen > 0 && rLen < 2 {
		s = fmt.Sprintf(`%s*`, string(rStr[0:1]))
	} else if rLen > 2 && rLen <= 4 {
		s = fmt.Sprintf(`%s%s%s`, string(rStr[0:1]), strings.Repeat("*", rLen-2), string(rStr[rLen-1:]))
	} else if rLen > 4 {
		s = fmt.Sprintf(`%s%s%s`, string(rStr[0:1]), strings.Repeat("*", rLen-3), string(rStr[rLen-2:]))
	}

	return s
}

func TrimTags(s string) string {
	re := regexp.MustCompile(`#(\S+)`)
	out := strings.TrimSpace(re.ReplaceAllString(s, ""))
	return out
}

func StringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

func CompareEquals(val1, val2 float64) bool {
	value1 := fmt.Sprintf("%.4f", val1)
	value2 := fmt.Sprintf("%.4f", val2)

	return value1 == value2
}
