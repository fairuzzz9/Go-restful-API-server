package utilities

import (
	"testing"
)

type rot18TestCase struct {
	name     string
	input    string
	expected string
}

var rot18TestCases = []rot18TestCase{
	{name: `test abc`, input: "abc", expected: `nop`},
	{name: `test abc123`, input: "abc123", expected: `nop678`},
	{name: `test %3Ds2#$sswws`, input: "%3Ds2#$sswws", expected: `%8Qf7#$ffjjf`},
	{name: `test tab white space `, input: "\t!@ mm0klm;;la$$#", expected: "	!@ zz5xyz;;yn$$#"},
}

func assertRot18(t *testing.T, input, expected string) {
	actual := Rot18(input)

	t.Log(actual)

	if actual != expected {
		t.Fatalf("error! expected string is \"%s\", but got \"%s\" \n", expected, actual)
	}
}

func TestRot18(t *testing.T) {

	for _, c := range rot18TestCases {
		t.Run(c.name, func(t *testing.T) {
			assertRot18(t, c.input, c.expected)
		})
	}
}
