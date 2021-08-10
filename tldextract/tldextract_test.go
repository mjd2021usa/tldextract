package tldextract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New_empty_file(t *testing.T) {
	assert := assert.New(t)

	actual, err := New("", true)

	assert.Nil(err, "Error nil")
	assert.NotNil(actual, "Result not nil")
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

func Test_Extract(t *testing.T) {
	assert := assert.New(t)

	tld, err := New("", true)
	assert.Nil(err, "Error nil")

	testCases := []struct {
		Url            string
		ExpectedResult Result
		ExpectedError  error
		Description    string
	}{
		{
			Url:            "",
			ExpectedResult: Result{},
			ExpectedError:  nil,
			Description:    "empty string",
		},
		{
			Url: "users@myhost.com",
			ExpectedResult: Result{
				SubDomain: "",
				Domain:    "myhost",
				Tld:       "com",
			},
			ExpectedError: nil,
			Description:   "user@ address",
		},
		{
			Url: "mailto:users@myhost.com",
			ExpectedResult: Result{
				SubDomain: "",
				Domain:    "myhost",
				Tld:       "com",
			},
			ExpectedError: nil,
			Description:   "email address",
		},
		{
			Url: "myhost.com:999",
			ExpectedResult: Result{
				SubDomain: "",
				Domain:    "myhost",
				Tld:       "com",
			},
			ExpectedError: nil,
			Description:   "host:port",
		},
		{
			Url: "myhost.com",
			ExpectedResult: Result{
				SubDomain: "",
				Domain:    "myhost",
				Tld:       "com",
			},
			ExpectedError: nil,
			Description:   "basic host",
		},
		{
			Url: "https://user:pass@foo.myhost.com:999/some/path?param1=value1&param2=value2",
			ExpectedResult: Result{
				SubDomain: "foo",
				Domain:    "myhost",
				Tld:       "com",
			},
			ExpectedError: nil,
			Description:   "Full URL with subdomain",
		},
	}

	for _, tc := range testCases {
		actualResult := tld.Extract(tc.Url)

		//assert.Equal(tc.ExpectedError, actualErr, "Error are rquil")
		assert.Equal(tc.ExpectedResult, *actualResult, "Result are equal")
	}
}