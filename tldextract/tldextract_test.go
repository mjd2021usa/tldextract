package tldextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_empty_file(t *testing.T) {
	assert := assert.New(t)

	actual, err := New("", true)

	assert.Nil(actual, "Result should be nil")
	assert.NotNil(err, "Error should not be nil")
}

func Test_New_missing_cache_file(t *testing.T) {
	assert := assert.New(t)

	actual, err := New("i.do.not.exist.cache", true)

	assert.Nil(err, "Error nil")
	assert.NotNil(actual, "Result not nil")
}

func Test_New_good_cache_file(t *testing.T) {
	assert := assert.New(t)

	actual, err := New("../test/tld.cache", true)

	assert.Nil(err, "Error nil")
	assert.NotNil(actual, "Result not nil")
}

func assertResult(t *testing.T, url string, expected *Result, returned *Result, description string) {
	if (expected.Flag == returned.Flag) && (expected.Domain == returned.Domain) && (expected.SubDomain == returned.SubDomain) && (expected.Tld == returned.Tld) {
		return
	}
	t.Errorf("%s - %s;expected:%+v;returned:%+v", description, url, expected, returned)
}

func Test_Extract(t *testing.T) {
	assert := assert.New(t)

	tld, err := New("../test/tld.cache", true)
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
			Url:            "https://user:pass@foo.myhost.com:999/some/path?param1=value1&param2=value2",
			ExpectedResult: Result{Flag: Domain, SubDomain: "foo", Domain: "myhost", Tld: "com"},
			ExpectedError:  nil,
			Description:    "Full URL with subdomain",
		},
	}

	for _, tc := range testCases {
		actualResult := tld.Extract(tc.Url)

		//assert.Equal(tc.ExpectedError, actualErr, "Error are equal")
		assertResult(t, tc.Url, &tc.ExpectedResult, actualResult, tc.Description)
	}
}
