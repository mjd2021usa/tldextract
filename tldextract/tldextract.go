package tldextract

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	DefaultCacheTimeout int64 = 10

	// Use raw strings to avoid having to quote the backslashes.
	DomainRegexText = `^[a-z0-9-\p{Han}]{1,63}$`
	//IP4RegexText    = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	SchemeRegexText = `^([a-z0-9\+\-\.]+:)?//`

	Malformed = iota
	Domain
	Ip4
	Ip6
)

var (
	DefaultTldUrls = []string{
		"https://publicsuffix.org/list/public_suffix_list.dat",
		"https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat",
	}

	// Compile the expression once, preferably at init time.
	domainRegex = regexp.MustCompile(DomainRegexText)
	//ip4Regex    = regexp.MustCompile(IP4RegexText)
	schemeRegex = regexp.MustCompile(SchemeRegexText)
)

type Result struct {
	Flag      int
	SubDomain string
	Domain    string
	Tld       string
}

type TldNode struct {
	ExceptRule bool
	ValidTld   bool
	matches    map[string]*TldNode
}

type TLDExtract struct {
	//CacheTimeout int64
	CacheFile string
	TldNodes  *TldNode
	Debug     bool
}

func New(fqdn string, debug bool) (*TLDExtract, error) {
	// Load Unique Cache List
	//timeout := GetEnvInt64("TLDEXTRACT_CACHE_TIMEOUT", 10, 64, DefaultCacheTimeout)
	urlsString := GetEnvString("TLDEXTRACT_URLS", strings.Join(DefaultTldUrls, ","))
	urls := strings.Split(urlsString, ",")
	cache, err := LoadCache(fqdn, urls, true)
	if err != nil {
		panic(fmt.Sprintf("Cache error: %s", err))
	}

	// Load Unique Cache List into TldNode structure
	newEmptyMap := make(map[string]*TldNode)
	tldNodes := &TldNode{ExceptRule: false, ValidTld: false, matches: newEmptyMap}
	for key, _ := range cache {
		exceptionRule := key[0] == '!'
		if exceptionRule {
			key = key[1:]
		}
		parts := strings.Split(key, ".")
		addTldRule(tldNodes, parts, exceptionRule)
	}

	tld := TLDExtract{
		CacheFile: fqdn,
		//CacheTimeout: timeout,
		Debug:    debug,
		TldNodes: tldNodes,
	}
	return &tld, nil
}

func addTldRule(rootNode *TldNode, parts []string, ex bool) {
	numParts := len(parts)
	current := rootNode
	for idx := numParts - 1; idx >= 0; idx-- {
		lab := parts[idx]
		match, found := current.matches[lab]
		if !found {
			except := ex
			valid := !ex && idx == 0
			newEmptyMap := make(map[string]*TldNode)
			current.matches[lab] = &TldNode{ExceptRule: except, ValidTld: valid, matches: newEmptyMap}
			match = current.matches[lab]
		} else if !ex && idx == 0 {
			match.ValidTld = true
		}

		current = match
	}
}

func (tlde *TLDExtract) Extract(urlString string) *Result {
	data := strings.ToLower(urlString)

	data = schemeRegex.ReplaceAllString(data, "")
	atIdx := strings.Index(data, "@")
	if atIdx != -1 {
		data = data[atIdx+1:]
	}

	index := strings.IndexFunc(data, func(r rune) bool {
		switch r {
		case '&', '/', '?', ':', '#':
			return true
		}
		return false
	})
	if index != -1 {
		data = data[0:index]
	}

	if tlde.Debug {
		fmt.Printf("%s;%s\n", data, urlString)
	}

	return tlde.extract(data)
}

func (tlde *TLDExtract) extract(url string) *Result {
	domain, tld := tlde.extractTld(url)
	if tld == "" {
		ip := net.ParseIP(url)
		if ip != nil {
			if IsIPv4(ip) {
				return &Result{Flag: Ip4, Domain: url}
			} else if IsIPv6(ip) {
				return &Result{Flag: Ip6, Domain: url}
			}
			return &Result{Flag: Malformed, Domain: url}
		}
		return &Result{Flag: Malformed}
	}
	subDomain, domain := SubDomain(domain)
	if domainRegex.MatchString(domain) {
		return &Result{Flag: Domain, Domain: domain, SubDomain: subDomain, Tld: tld}
	}
	return &Result{Flag: Malformed}
}

func (tlde *TLDExtract) extractTld(url string) (domain, tld string) {
	spl := strings.Split(url, ".")
	tldIndex, validTld := tlde.getTldIndex(spl)
	if validTld {
		domain = strings.Join(spl[:tldIndex], ".")
		tld = strings.Join(spl[tldIndex:], ".")
	} else {
		domain = url
	}
	return
}

func (tlde *TLDExtract) getTldIndex(labels []string) (int, bool) {
	current := tlde.TldNodes
	parentValid := false
	for idx := len(labels) - 1; idx >= 0; idx-- {
		lab := labels[idx]
		node, foundLabel := current.matches[lab]
		_, foundAsterisk := current.matches["*"]

		switch {
		case foundLabel && !node.ExceptRule:
			parentValid = node.ValidTld
			current = node
		// Found an exception rule
		case foundLabel:
			fallthrough
		case parentValid:
			return idx + 1, true
		case foundAsterisk:
			parentValid = true
		default:
			return -1, false
		}
	}
	return -1, false
}
