package util

import (
	"testing"
)

type stringSlice []string

func (a stringSlice) hasValue(b string) bool {
	for _, v := range a {
		if b == v {
			return true
		}
	}
	return false
}

func TestQueryCreation(t *testing.T) {
	input := UrlQuery{
		"a": "a",
		"b": "b",
		"c": "c",
	}

	expectedOutput := stringSlice{
		"?a=a&b=b&c=c",
		"?a=a&c=c&b=b",
		"?b=b&a=a&c=c",
		"?b=b&c=c&a=a",
		"?c=c&a=a&b=b",
		"?c=c&b=b&a=a",
	}

	output := input.CreateQuery()

	if !expectedOutput.hasValue(output) {
		t.Errorf("Failed ! Expected %s got %s", expectedOutput, output)
		t.FailNow()
	}

	t.Log("Success !")
}
