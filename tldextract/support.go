package tldextract

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

// CreateList -
func CreateList(timeout int64) map[string]struct{} {
	urls := []string{
		"https://publicsuffix.org/list/public_suffix_list.dat",
		"https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat",
	}
	uniqueList := make(map[string]struct{})
	for _, url := range urls {
		data, err := DownloadFile(url, timeout)
		if err != nil {
			continue
		}
		resp := NormalizeLines(string(data))
		for _, item := range resp {
			uniqueList[item] = struct{}{}
		}
	}
	return uniqueList
}

// CreateNewCacheFile - create new cache file from URLs
func CreateNewCacheFile(fqdn string, urls []string, timeout int64) (map[string]struct{}, error) {
	uniqueList := DownloadUrls2List(urls, timeout)
	keys := GetKeys(uniqueList)
	if len(keys) > 0 {
		buf := strings.Join(keys, "\n")
		err := WriteFile(fqdn, []byte(buf))
		return uniqueList, err
	}
	return nil, fmt.Errorf("no records found - skipping overwrite")
}

// LoadCache - Load cache file with Refresh and fail over options
func LoadCache(fqdn string, urls []string, refresh bool, timeout int64) (map[string]struct{}, error) {
	if refresh {
		// Refresh cache from URLs
		uList, err := CreateNewCacheFile(fqdn, urls, timeout)
		if err != nil || len(uList) <= 0 {
			// Fall back to cache file
			return LoadCacheFile(fqdn)
		}
		return uList, err
	}
	uList, err := LoadCacheFile(fqdn)
	if err != nil || len(uList) <= 0 {
		// try to refresh cache from URLs
		return CreateNewCacheFile(fqdn, urls, timeout)
	}
	return uList, err
}

// GetKeys - get list of keys from generic map as string array
func GetKeys(list map[string]struct{}) []string {
	keys := make([]string, 0, len(list))
	for key := range list {
		keys = append(keys, key)
	}
	return keys
}

// LoadCacheFile - read cache file and Normalize contents (remove comments, split lines, etc)
func LoadCacheFile(fqdn string) (map[string]struct{}, error) {
	data, err := ReadFile(fqdn)
	if err != nil {
		return nil, err
	}
	// Normalize buffer into string array
	resp := NormalizeLines(string(data))
	// Add lines to unique list
	uList := make(map[string]struct{})
	for _, item := range resp {
		uList[item] = struct{}{}
	}
	return uList, nil
}

// DownloadUrls2List - Download N number of URLs and merge the unique rows into a generic key map
func DownloadUrls2List(urls []string, timeout int64) map[string]struct{} {
	uList := make(map[string]struct{})
	for _, url := range urls {
		data, err := DownloadFile(url, timeout)
		if err != nil {
			// nothing good, move to next file
			continue
		}
		// Normalize buffer into string array
		resp := NormalizeLines(string(data))
		// Add lines to unique list
		for _, item := range resp {
			uList[item] = struct{}{}
		}
	}
	return uList
}

// RemoveNoiseLines - remove blank lines and comments, converting each kept line to lowercase
func RemoveNoiseLines(srcLines []string) []string {
	dstLines := []string{}
	for _, line := range srcLines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") {
			dstLines = append(dstLines, strings.ToLower(line))
		}
	}
	return dstLines
}

// NormalizeLines - split string into array of strings and remove noise lines
func NormalizeLines(buffer string) []string {
	lines := strings.Split(buffer, "\n")
	return RemoveNoiseLines(lines)
}

// ReadFile - Read file into byte array
func ReadFile(fqfn string) ([]byte, error) {
	return ioutil.ReadFile(fqfn)
}

// WriteFile - Write byte array into file
func WriteFile(fqfn string, buffer []byte) error {
	return ioutil.WriteFile(fqfn, buffer, fs.FileMode(0644))
}

// DownloadFile - Get body from URL
func DownloadFile(url string, timeout int64) ([]byte, error) {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return ioutil.ReadAll(resp.Body)
	}

	return []byte{}, fmt.Errorf("HTTP Status Code: %d returned", resp.StatusCode)
}
