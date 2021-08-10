package tldextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_empty_file(t *testing.T) {
	assert := assert.New(t)

	actual, err := New("", false)

	assert.Nil(actual, "Result should be nil")
	assert.NotNil(err, "Error should not be nil")
}

func Test_New_missing_cache_file(t *testing.T) {
	assert := assert.New(t)

	actual, err := New("i.do.not.exist.cache", false)

	assert.Nil(err, "Error nil")
	assert.NotNil(actual, "Result not nil")
}

func Test_New_good_cache_file(t *testing.T) {
	assert := assert.New(t)

	actual, err := New("../test/tld.cache", false)

	assert.Nil(err, "Error nil")
	assert.NotNil(actual, "Result not nil")
}

func assertResult(t *testing.T, url string, expected *Result, actual *Result, description string) {
	if (expected.Flag == actual.Flag) &&
		(expected.SubDomain == actual.SubDomain) &&
		(expected.Domain == actual.Domain) &&
		(expected.Tld == actual.Tld) {
		return
	}
	t.Errorf("%s - %s - expected:%+v, actual:%+v", description, url, expected, actual)
}

func Test_Extract(t *testing.T) {
	assert := assert.New(t)

	tld, err := New("../test/tld.cache", false)
	assert.Nil(err, "Error nil")

	testCases := []struct {
		Url            string
		ExpectedResult Result
		ExpectedError  error
		Description    string
	}{
		{
			Url:            "",
			ExpectedResult: Result{Flag: Malformed, SubDomain: "", Domain: "", Tld: ""},
			ExpectedError:  nil,
			Description:    "empty string",
		},
		{
			Url:            "users@myhost.com",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "myhost", Tld: "com"},
			ExpectedError:  nil,
			Description:    "user@ address",
		},
		{
			Url:            "mailto:users@myhost.com",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "myhost", Tld: "com"},
			ExpectedError:  nil,
			Description:    "email address",
		},
		{
			Url:            "myhost.com:999",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "myhost", Tld: "com"},
			ExpectedError:  nil,
			Description:    "host:port",
		},
		{
			Url:            "myhost.com",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "myhost", Tld: "com"},
			ExpectedError:  nil,
			Description:    "basic host",
		},
		{
			Url:            "255.255.myhost.com",
			ExpectedResult: Result{Flag: Domain, SubDomain: "255.255", Domain: "myhost", Tld: "com"},
			ExpectedError:  nil,
			Description:    "basic host with numerit subdomains",
		},
		{
			Url:            "https://user:pass@foo.myhost.com:999/some/path?param1=value1&param2=value2",
			ExpectedResult: Result{Flag: Domain, SubDomain: "foo", Domain: "myhost", Tld: "com"},
			ExpectedError:  nil,
			Description:    "Full URL with subdomain",
		},
		{
			Url:            "http://www.duckduckgo.com",
			ExpectedResult: Result{Flag: Domain, SubDomain: "www", Domain: "duckduckgo", Tld: "com"},
			ExpectedError:  nil,
			Description:    "Full URL with subdomain",
		},
		{
			Url:            "http://duckduckgo.co.uk/path?param1=value1&param2=value2&param3=value3&param4=value4&src=https%3A%2F%2Fwww.yahoo.com%2F",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "duckduckgo", Tld: "co.uk"},
			ExpectedError:  nil,
			Description:    "Full HTTP URL with no subdomain",
		},
		{
			Url:            "http://big.long.sub.domain.duckduckgo.co.uk/path?param1=value1&param2=value2&param3=value3&param4=value4&src=https%3A%2F%2Fwww.yahoo.com%2F",
			ExpectedResult: Result{Flag: Domain, SubDomain: "big.long.sub.domain", Domain: "duckduckgo", Tld: "co.uk"},
			ExpectedError:  nil,
			Description:    "Full HTTP URL with subdomain",
		},
		{
			Url:            "https://duckduckgo.co.uk/path?param1=value1&param2=value2&param3=value3&param4=value4&src=https%3A%2F%2Fwww.yahoo.com%2F",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "duckduckgo", Tld: "co.uk"},
			ExpectedError:  nil,
			Description:    "Full HTTPS URL with no subdomain",
		},
		{
			Url:            "ftp://peterparker:multipass@mail.duckduckgo.co.uk:666/path?param1=value1&param2=value2&param3=value3&param4=value4&src=https%3A%2F%2Fwww.yahoo.com%2F",
			ExpectedResult: Result{Flag: Domain, SubDomain: "mail", Domain: "duckduckgo", Tld: "co.uk"},
			ExpectedError:  nil,
			Description:    "Full ftp URL with subdomain",
		},
		{
			Url:            "git+ssh://www.github.com/",
			ExpectedResult: Result{Flag: Domain, SubDomain: "www", Domain: "github", Tld: "com"},
			ExpectedError:  nil,
			Description:    "Full git+ssh URL with subdomain",
		},
		{
			Url:            "git+ssh://www.!github.com/",
			ExpectedResult: Result{Flag: Malformed, SubDomain: "", Domain: "", Tld: ""},
			ExpectedError:  nil,
			Description:    "Full git+ssh URL with bad domain",
		},
		{
			Url:            "ssh://server.domain.com/",
			ExpectedResult: Result{Flag: Domain, SubDomain: "server", Domain: "domain", Tld: "com"},
			ExpectedError:  nil,
			Description:    "Full ssh URL with subdomain",
		},
		{
			Url:            "//server.domain.com/path",
			ExpectedResult: Result{Flag: Domain, SubDomain: "server", Domain: "domain", Tld: "com"},
			ExpectedError:  nil,
			Description:    "Missing protocol URL with subdomain",
		},
		{
			Url:            "server.domain.com/path",
			ExpectedResult: Result{Flag: Domain, SubDomain: "server", Domain: "domain", Tld: "com"},
			ExpectedError:  nil,
			Description:    "Full ssh URL with subdomain",
		},
		{
			Url:            "10.10.10.10",
			ExpectedResult: Result{Flag: IPv4, SubDomain: "", Domain: "10.10.10.10", Tld: ""},
			ExpectedError:  nil,
			Description:    "Basic IPv4 Address",
		},
		{
			Url:            "http://10.10.10.10",
			ExpectedResult: Result{Flag: IPv4, SubDomain: "", Domain: "10.10.10.10", Tld: ""},
			ExpectedError:  nil,
			Description:    "Basic IPv4 Address URL",
		},
		{
			Url:            "http://10.10.10.256",
			ExpectedResult: Result{Flag: Malformed, SubDomain: "", Domain: "", Tld: ""},
			ExpectedError:  nil,
			Description:    "Basic IPv4 Address URL with bad IP",
		},
		/*
			{
				Url:            "http://2001:0db8:0000:0000:0000:ff00:0042:8329",
				ExpectedResult: Result{Flag: IPv6, SubDomain: "", Domain: "2001:0db8:0000:0000:0000:ff00:0042:8329", Tld: ""},
				ExpectedError:  nil,
				Description:    "Basic IPv6 Address URL",
			},
			{
				Url:            "http://2001:db8:0:0:0:ff00:42:8329",
				ExpectedResult: Result{Flag: IPv6, SubDomain: "", Domain: "2001:db8:0:0:0:ff00:42:8329", Tld: ""},
				ExpectedError:  nil,
				Description:    "Basic IPv6 Address URL",
			},
			{
				Url:            "http://2001:db8::ff00:42:8329",
				ExpectedResult: Result{Flag: IPv6, SubDomain: "", Domain: "2001:db8::ff00:42:8329", Tld: ""},
				ExpectedError:  nil,
				Description:    "Basic IPv6 Address URL",
			},
			{
				Url:            "http://::ffff:192.0.2.128",
				ExpectedResult: Result{Flag: IPv6, SubDomain: "", Domain: "::ffff:192.0.2.128", Tld: ""},
				ExpectedError:  nil,
				Description:    "Basic IPv6 Address URL",
			},
			{
				Url:            "http://::192.0.2.128",
				ExpectedResult: Result{Flag: IPv6, SubDomain: "", Domain: "::192.0.2.128", Tld: ""},
				ExpectedError:  nil,
				Description:    "Basic IPv6 Address URL",
			},
		*/
		{
			Url:            "http://godaddy.godaddy",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "godaddy", Tld: "godaddy"},
			ExpectedError:  nil,
			Description:    "Basic URL",
		},
		{
			Url:            "http://godaddy.godaddy.godaddy",
			ExpectedResult: Result{Flag: Domain, SubDomain: "godaddy", Domain: "godaddy", Tld: "godaddy"},
			ExpectedError:  nil,
			Description:    "Basic URL with subdomain",
		},
		{
			Url:            "http://godaddy.godaddy.co.uk",
			ExpectedResult: Result{Flag: Domain, SubDomain: "godaddy", Domain: "godaddy", Tld: "co.uk"},
			ExpectedError:  nil,
			Description:    "Basic URL with subdomain",
		},
		{
			Url:            "http://godaddy",
			ExpectedResult: Result{Flag: Malformed, SubDomain: "", Domain: "", Tld: ""},
			ExpectedError:  nil,
			Description:    "Basic URL with TLD only",
		},
		{
			Url:            "http://godaddy.cannon-fodder",
			ExpectedResult: Result{Flag: Malformed, SubDomain: "", Domain: "", Tld: ""},
			ExpectedError:  nil,
			Description:    "Basic URL with bad TLD",
		},
		{
			Url:            "http://godaddy.godaddy.cannon-fodder",
			ExpectedResult: Result{Flag: Malformed, SubDomain: "", Domain: "", Tld: ""},
			ExpectedError:  nil,
			Description:    "Basic URL with subdomainand bad TLD",
		},
		{
			Url:            "http://domainer.个人.hk",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "domainer", Tld: "个人.hk"},
			ExpectedError:  nil,
			Description:    "Basic URL with partial encoded punycode TLD",
		},
		{
			Url:            "http://domainer.公司.香港",
			ExpectedResult: Result{Flag: Domain, SubDomain: "", Domain: "domainer", Tld: "公司.香港"},
			ExpectedError:  nil,
			Description:    "Basic URL with fully encoded punycode TLD",
		},
	}

	for _, tc := range testCases {
		actualResult := tld.Extract(tc.Url)

		//assert.Equal(tc.ExpectedError, actualErr, "Error are equal")
		assertResult(t, tc.Url, &tc.ExpectedResult, actualResult, tc.Description)
	}
}
