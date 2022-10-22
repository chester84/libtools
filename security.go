package libtools

import (
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

const (
	ProductDomain = ""
	DevDomain     = ""

	ProductH5Domain = ""
	DevH5Domain     = ""
)

// IsInternalIPV1 超简算法
func IsInternalIPV1(ip string) bool {
	if ip == "" {
		logs.Warning("[IsInternalIPV1] get empty input")
		return false
	}

	ipExp := strings.Split(ip, ".")
	if len(ipExp) != 4 {
		logs.Warning("[IsInternalIPV1] ip: %s address format is incorrect", ip)
		return false
	}

	if (ipExp[0] == "127" && ipExp[1] == "0") ||
		(ipExp[0] == "172" && (ipExp[1] == "31" || ipExp[1] == "16")) {
		return true
	}

	return false
}

func InternalApiDomain() string {
	if IsProductEnv() {
		return ProductDomain
	} else {
		return DevDomain
	}
}

func InternalH5Domain() string {
	if IsProductEnv() {
		return ProductH5Domain
	} else {
		return DevH5Domain
	}
}
