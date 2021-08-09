package tldextract

import (
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"
)

// RemoveNoiseLines - remove blankc lines and comments, with lines converter to lowercase
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

// NormalizeLines - split buffer into lines and remove noise
func NormalizeLines(buffer string) []string {
	lines := strings.Split(buffer, "\n")
	return RemoveNoiseLines(lines)
}

// CreateList -
func CreateList() map[string]struct{} {
	urls := []string{
		"https://publicsuffix.org/list/public_suffix_list.dat",
		"https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat",
	}
	uniqueList := make(map[string]struct{})
	for _, url := range urls {
		data, err := DownloadFile(url)
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

func GetKeys(list map[string]struct{}) []string {
	keys := make([]string, 0, len(list))
	for key := range list {
		keys = append(keys, key)
	}
	return keys
}

func DownloadUrls2List(urls []string) map[string]struct{} {
	uniqueList := make(map[string]struct{})
	for _, url := range urls {
		data, err := DownloadFile(url)
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

func CreateNewCacheFile(fqdn string, urls []string) (map[string]struct{}, error) {
	uniqueList := DownloadUrls2List(urls)
	keys := GetKeys(uniqueList)
	var buf string
	if len(keys) > 0 {
		buf = strings.Join(keys, "\n")
	}
	err := WriteFile(fqdn, []byte(buf))
	return uniqueList, err
}

func LoadCache(fqdn string, urls []string, refresh bool) (map[string]struct{}, error) {
	if refresh {
		return CreateNewCacheFile(fqdn, urls)
	} else {
		data, err := ReadFile(fqdn)
		if err != nil {
			return CreateNewCacheFile(fqdn, urls)
		}
		resp := NormalizeLines(string(data))
		uniqueList := make(map[string]struct{})
		for _, item := range resp {
			uniqueList[item] = struct{}{}
		}
		return uniqueList, nil
	}
}

// ReadFile -
func ReadFile(fqfn string) ([]byte, error) {
	return ioutil.ReadFile(fqfn)
}

// WriteFile -
func WriteFile(fqfn string, buffer []byte) error {
	return ioutil.WriteFile(fqfn, buffer, fs.FileMode(0644))
}

// DownloadFile -
func DownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
