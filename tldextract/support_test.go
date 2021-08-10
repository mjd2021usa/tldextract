package tldextract

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SafeParseInt(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		IntString    string
		DefaultValue int64
		Expected     int64
		Description  string
	}{
		{
			IntString:    "",
			DefaultValue: int64(999),
			Expected:     int64(999),
			Description:  "empty string",
		},
		{
			IntString:    "-1",
			DefaultValue: int64(999),
			Expected:     int64(-1),
			Description:  "negative one",
		},
		{
			IntString:    "0",
			DefaultValue: int64(999),
			Expected:     int64(0),
			Description:  "zero",
		},
		{
			IntString:    "+1",
			DefaultValue: int64(999),
			Expected:     int64(1),
			Description:  "positive one",
		},
		{
			IntString:    "bad",
			DefaultValue: int64(999),
			Expected:     int64(999),
			Description:  "bad value, default",
		},
	}

	for _, tc := range testCases {
		actual := SafeParseInt(tc.IntString, 10, 64, tc.DefaultValue)

		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_GetEnvString(t *testing.T) {
	assert := assert.New(t)
	envValue := "envValue"
	testCases := []struct {
		EnvKey       string
		EnvValue     *string
		DefaultValue string
		Expected     string
		Description  string
	}{
		{
			EnvKey:       "",
			EnvValue:     nil,
			DefaultValue: "default",
			Expected:     "default",
			Description:  "empty key",
		},
		{
			EnvKey:       "exampleKey",
			EnvValue:     nil,
			DefaultValue: "default",
			Expected:     "default",
			Description:  "nil env",
		},
		{
			EnvKey:       "exampleKey",
			EnvValue:     &envValue,
			DefaultValue: "default",
			Expected:     envValue,
			Description:  "good env",
		},
	}

	for _, tc := range testCases {
		if tc.EnvValue == nil {
			os.Unsetenv(tc.EnvKey)
		} else {
			os.Setenv(tc.EnvKey, *tc.EnvValue)
		}

		actual := GetEnvString(tc.EnvKey, tc.DefaultValue)

		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_GetEnvInt64(t *testing.T) {
	assert := assert.New(t)
	envBadInt64Value := "bad"
	envGoodInt64Value := "77777777777"
	testCases := []struct {
		EnvKey       string
		EnvValue     *string
		DefaultValue int64
		Expected     int64
		Description  string
	}{
		{
			EnvKey:       "",
			EnvValue:     nil,
			DefaultValue: int64(777),
			Expected:     int64(777),
			Description:  "empty key",
		},
		{
			EnvKey:       "exampleKey",
			EnvValue:     nil,
			DefaultValue: int64(777),
			Expected:     int64(777),
			Description:  "nil env",
		},
		{
			EnvKey:       "exampleKey",
			EnvValue:     &envBadInt64Value,
			DefaultValue: int64(777),
			Expected:     int64(777),
			Description:  "good env, bad value",
		},
		{
			EnvKey:       "exampleKey",
			EnvValue:     &envGoodInt64Value,
			DefaultValue: int64(777),
			Expected:     int64(77777777777),
			Description:  "good env, good value",
		},
	}

	for _, tc := range testCases {
		if tc.EnvValue == nil {
			os.Unsetenv(tc.EnvKey)
		} else {
			os.Setenv(tc.EnvKey, *tc.EnvValue)
		}

		actual := GetEnvInt64(tc.EnvKey, 10, 64, tc.DefaultValue)

		assert.Equal(tc.Expected, actual, tc.Description)
	}
}

func Test_SubDomain(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		Domain            string
		ExpectedSubDomain string
		ExpectedDomain    string
		Description       string
	}{
		{
			Domain:            "",
			ExpectedSubDomain: "",
			ExpectedDomain:    "",
			Description:       "empty value",
		},
		{
			Domain:            "example",
			ExpectedSubDomain: "",
			ExpectedDomain:    "example",
			Description:       "domain only",
		},
		{
			Domain:            "sub.example",
			ExpectedSubDomain: "sub",
			ExpectedDomain:    "example",
			Description:       "domain with single level sub domain",
		},
		{
			Domain:            "sub.sub.sub.example",
			ExpectedSubDomain: "sub.sub.sub",
			ExpectedDomain:    "example",
			Description:       "domain with multi level sub domain",
		},
	}

	for _, tc := range testCases {

		actualSubDomain, actualDomain := SubDomain(tc.Domain)

		assert.Equal(tc.ExpectedDomain, actualDomain, tc.Description)
		assert.Equal(tc.ExpectedSubDomain, actualSubDomain, tc.Description)
	}
}

func Test_CreateList(t *testing.T) {
	assert := assert.New(t)

	actual := CreateList(10)

	assert.NotNil(actual, "Not Nil")
	assert.NotEqual(0, len(actual), "Length is not 0")
}
