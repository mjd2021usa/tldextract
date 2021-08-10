package main

import (
	"fmt"

	"github.com/mjd2021usa/tldextract"
)

func main() {
	cacheFileName := "tld.cache"
	tlde, err := tldextract.New(cacheFileName, false)
	if err != nil {
		panic(fmt.Sprintf("tldextract.New() error: %s", err))
	}

	uList, err := tldextract.LoadCacheFile(cacheFileName)
	if err != nil {
		panic(fmt.Errorf("error loading cache file '%s'", cacheFileName))
	}

	cnt := int(0)
	skipCnt := int(0)
	domainRoot := "fruitloops"
	subDomainRoot := "mr.ozzy"

	for key, _ := range uList {
		exceptionRule := key[0] == '!'
		if exceptionRule {
			//key = key[1:]
			// for now skip exceptionRules
			skipCnt++
			continue
		}
		astrickRule := key[0] == '*'
		if astrickRule {
			//key = key[1:]
			// for now skip astrickRule
			skipCnt++
			continue
		}

		domain := fmt.Sprintf("%s.%s", domainRoot, key)
		subDomain := fmt.Sprintf("%s.%s.%s", subDomainRoot, domainRoot, key)

		result1 := tlde.Extract(domain)
		if result1.Tld != key {
			fmt.Printf("result1.Tld Mismatch key:%s Tld:%s\n", key, result1.Tld)
		}
		if result1.Domain != domainRoot {
			fmt.Printf("result1.Domain Mismatch: key:%s Domain:%s domainRoot:%s\n", key, result1.Domain, domainRoot)
		}
		if result1.SubDomain != "" {
			fmt.Printf("result1.SubDomain Mismatch: key:%s SubDomain:%s subDomainRoot:%s\n", key, result1.SubDomain, subDomainRoot)
		}

		result2 := tlde.Extract(subDomain)
		if result2.Tld != key {
			fmt.Printf("result2.Tld Mismatch key:%s Tld:%s\n", key, result2.Tld)
		}
		if result2.Domain != domainRoot {
			fmt.Printf("result2.Domain Mismatch: key:%s Domain:%s domainRoot:%s\n", key, result2.Domain, domainRoot)
		}
		if result2.SubDomain != subDomainRoot {
			fmt.Printf("result2.SubDomain Mismatch: key:%s SubDomain:%s subDomainRoot:%s\n", key, result2.SubDomain, subDomainRoot)
		}
		cnt++
	}
	fmt.Printf("Cache Rows Tested: %d, skip count: %d\n", cnt, skipCnt)
}
