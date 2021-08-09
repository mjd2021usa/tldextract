package tldextract

import (
	"net"
	"os"
	"strconv"
	"strings"
)

func SafeParseInt(str string, base int, bitSize int, defaultValue int64) int64 {
	result, err := strconv.ParseInt(str, base, bitSize)
	if err != nil {
		return defaultValue
	}
	return result
}

func GetEnvString(envVar string, defaultValue string) string {
	val := os.Getenv(envVar)
	val = strings.TrimSpace(val)
	if len(val) <= 0 {
		val = defaultValue
	}
	return val
}

func GetEnvInt64(envVar string, base int, bitSize int, defaultValue int64) int64 {
	val := os.Getenv(envVar)
	val = strings.TrimSpace(val)
	result, err := strconv.ParseInt(val, base, bitSize)
	if err != nil {
		return defaultValue
	}
	return result
}

// SubDomain - return sub-domain, domain
func SubDomain(domain string) (string, string) {
	splits := strings.Split(domain, ".")
	cnt := len(splits)
	if cnt == 1 {
		return "", domain
	}
	return strings.Join(splits[0:cnt-1], "."), splits[cnt-1]
}

func IsIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

func IsIPv6(ip net.IP) bool {
	return ip.To16() != nil
}
