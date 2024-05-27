package utils

import (
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	numberRegexMatcher *regexp.Regexp = regexp.MustCompile(`\d`)
	intStrMatcher      *regexp.Regexp = regexp.MustCompile(`^[\+-]?\d+$`)
	urlMatcher         *regexp.Regexp = regexp.MustCompile(`^((ftp|http|https?):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(([a-zA-Z0-9]+([-\.][a-zA-Z0-9]+)*)|((www\.)?))?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?))(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`)
	dnsMatcher         *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9\-]+[\.]?)*[a-zA-Z0-9]$`)
	emailMatcher       *regexp.Regexp = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	base64Matcher      *regexp.Regexp = regexp.MustCompile(`^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`)
	base64URLMatcher   *regexp.Regexp = regexp.MustCompile(`^([A-Za-z0-9_-]{4})*([A-Za-z0-9_-]{2}(==)?|[A-Za-z0-9_-]{3}=?)?$`)
)

// IsBase64 检查字符串是否为 Base64 字符串
func IsBase64(base64 string) bool {
	return base64Matcher.MatchString(base64)
}

// IsBase64URL 检查给定字符串是否是有效的 URL 安全 Base64 编码字符串
func IsBase64URL(v string) bool {
	return base64URLMatcher.MatchString(v)
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func IsIp(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	return ip != nil
}

func IsIpV4(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return strings.Contains(ipstr, ".")
}

func IsIpV6(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return strings.Contains(ipstr, ":")
}

func IsPort(str string) bool {
	if i, err := strconv.ParseInt(str, 10, 64); err == nil && i > 0 && i < 65536 {
		return true
	}
	return false
}

func IsUrl(str string) bool {
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return urlMatcher.MatchString(str)
}

func IsDns(dns string) bool {
	return dnsMatcher.MatchString(dns)
}

// IsEmail 检查字符串是否是电子邮件地址
func IsEmail(email string) bool {
	return emailMatcher.MatchString(email)
}

func IsJWT(v string) bool {
	_jwt := strings.Split(v, ".")
	if len(_jwt) != 3 {
		return false
	}

	for _, s := range _jwt {
		if !IsBase64URL(s) {
			return false
		}
	}

	return true
}

// IsInt 检查字符串是否可以转换为数字
func IsInt(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		return true
	}
	return false
}

func IsFloat(v any) bool {
	switch v.(type) {
	case float32, float64:
		return true
	}
	return false
}

func IsNumber(v any) bool {
	return IsInt(v) || IsFloat(v)
}

func ContainNumber(input string) bool {
	return numberRegexMatcher.MatchString(input)
}

// IsIntStr 检查字符串是否可以转换为整数
func IsIntStr(str string) bool {
	return intStrMatcher.MatchString(str)
}

func IsFloatStr(str string) bool {
	_, e := strconv.ParseFloat(str, 64)
	return e == nil
}

func IsNumberStr(s string) bool {
	return IsIntStr(s) || IsFloatStr(s)
}

// IsRegexMatch 检查字符串是否与正则表达式匹配
func IsRegexMatch(str, regex string) bool {
	reg := regexp.MustCompile(regex)
	return reg.MatchString(str)
}
