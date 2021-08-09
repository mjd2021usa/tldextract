package main

import (
	"fmt"

	"github.com/mjd2021usa/tldextract"
)

func main() {
	tlde, err := tldextract.New("tld.cache", true)
	if err != nil {
		panic(fmt.Sprintf("tldextract.New() error: %s", err))
	}

	result := tlde.Extract("my.foo.bar.com")

	fmt.Printf("Result.Tld:%s\n", result.Tld)
	fmt.Printf("Result.Domain:%s\n", result.Domain)
	fmt.Printf("Result.SubDomain:%s\n", result.SubDomain)
}
