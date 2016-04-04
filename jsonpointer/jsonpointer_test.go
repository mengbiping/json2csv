package jsonpointer

import (
	"reflect"
	"testing"
)

var testNewJSONPointerCases = []struct {
	pointer  string
	expected []Token
	err      string
}{
	{`/foo`, []Token{`foo`}, ``},
	{`/foo~0bar`, []Token{`foo~bar`}, ``},
	{`/foo~1bar`, []Token{`foo/bar`}, ``},
	{`/foo/bar`, []Token{`foo`, `bar`}, ``},
	{`/foo/0/bar`, []Token{`foo`, `0`, `bar`}, ``},
	{`/`, []Token{""}, ``},      // empty string key
	{`//`, []Token{"", ""}, ``}, // empty string key
	{``, []Token{}, ``},         // whole content (root)
	{`foo`, nil, `Invalid JSON Pointer "foo"`},
}

func TestNewJSONPointer(t *testing.T) {
	for caseIndex, testCase := range testNewJSONPointerCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		actual := []Token(pointer)
		if err != nil {
			if err.Error() != testCase.err {
				t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.err, err)
			}
		} else if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testLenCases = []struct {
	pointer  string
	expected int
}{
	{`/foo`, 1},
	{`/foo~0bar`, 1},
	{`/foo~1bar`, 1},
	{`/foo/bar`, 2},
	{`/foo/0/bar`, 3},
}

func TestLen(t *testing.T) {
	for caseIndex, testCase := range testLenCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Fatal(err)
		}
		actual := pointer.Len()
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testAppendCases = []struct {
	pointer  string
	token    string
	expected string
}{
	{`/foo`, `append`, `/foo/append`},
	{`/foo~0bar`, `append`, `/foo~0bar/append`},
	{`/foo~1bar`, `append`, `/foo~1bar/append`},
	{`/foo/bar`, `append`, `/foo/bar/append`},
	{`/foo/0/bar`, `append`, `/foo/0/bar/append`},
	{`/`, `append`, `//append`},
	{`//`, `append`, `///append`},
	{``, `append`, `/append`},
}

func TestAppend(t *testing.T) {
	for caseIndex, testCase := range testAppendCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Fatal(err)
		}
		pointer.Append(Token(testCase.token))
		actual := pointer.String()
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testPopCases = []struct {
	pointer  string
	removed  string
	expected string
}{
	{`/foo`, `foo`, ``},
	{`/foo~0bar`, `foo~bar`, ``},
	{`/foo~1bar`, `foo/bar`, ``},
	{`/foo/bar`, `bar`, `/foo`},
	{`/foo/0/bar`, `bar`, `/foo/0`},
	{`/`, ``, ``},
	{`//`, ``, `/`},
	{``, ``, ``},
}

func TestPop(t *testing.T) {
	for caseIndex, testCase := range testPopCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Fatal(err)
		}

		removed := pointer.Pop()
		if removed != Token(testCase.removed) {
			t.Errorf("%d: Expected removed %v, but %v", caseIndex, Token(testCase.removed), removed)
		}

		actual := pointer.String()
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

func TestClone(t *testing.T) {
	orig, err := NewJSONPointer("/foo/bar")
	pointer, err := NewJSONPointer("/foo/bar")
	if err != nil {
		t.Fatal(err)
	}

	cloned := pointer.Clone()
	if !reflect.DeepEqual(cloned, pointer) {
		t.Errorf("Expected %v, but %v", pointer, cloned)
	}

	cloned.AppendString("baz")
	if !reflect.DeepEqual(pointer, orig) {
		t.Errorf("Expected %v, but %v", orig, pointer)
	}
}

var testStringsCases = []struct {
	pointer  string
	expected []string
}{
	{`/foo`, []string{`foo`}},
	{`/foo~0bar`, []string{`foo~bar`}},
	{`/foo~1bar`, []string{`foo/bar`}},
	{`/foo/bar`, []string{`foo`, `bar`}},
	{`/foo/0/bar`, []string{`foo`, `0`, `bar`}},
	{`/`, []string{""}},      // empty string key
	{`//`, []string{"", ""}}, // empty string key
	{``, []string{}},         // whole content (root)
}

func TestStrings(t *testing.T) {
	for caseIndex, testCase := range testStringsCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Fatal(err)
		}
		actual := pointer.Strings()
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testEscapedStringsCases = []struct {
	pointer  string
	expected []string
}{
	{`/foo`, []string{`foo`}},
	{`/foo~0bar`, []string{`foo~0bar`}},
	{`/foo~1bar`, []string{`foo~1bar`}},
	{`/foo/bar`, []string{`foo`, `bar`}},
	{`/foo/0/bar`, []string{`foo`, `0`, `bar`}},
	{`/`, []string{""}},      // empty string key
	{`//`, []string{"", ""}}, // empty string key
	{``, []string{}},         // whole content (root)
}

func TestEscapedStrings(t *testing.T) {
	for caseIndex, testCase := range testEscapedStringsCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Fatal(err)
		}
		actual := pointer.EscapedStrings()
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testStringCases = []struct {
	pointer  string
	expected string
}{
	{`/foo`, `/foo`},
	{`/foo~0bar`, `/foo~0bar`},
	{`/foo~1bar`, `/foo~1bar`},
	{`/foo/bar`, `/foo/bar`},
	{`/foo/0/bar`, `/foo/0/bar`},
	{`/`, `/`},   // empty string key
	{`//`, `//`}, // empty string key
	{``, ``},     // whole content (root)
}

func TestString(t *testing.T) {
	for caseIndex, testCase := range testStringCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Fatal(err)
		}
		actual := pointer.String()
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testDotNotationCases = []struct {
	pointer         string
	expected        string
	expectedBracket string
}{
	{`/foo`, `foo`, `foo`},
	{`/foo~0bar`, `foo~bar`, `foo~bar`},
	{`/foo~1bar`, `foo/bar`, `foo/bar`},
	{`/foo/bar`, `foo.bar`, `foo.bar`},
	{`/foo/0/bar`, `foo.0.bar`, `foo[0].bar`},
	{`/`, ``, ``},    // empty string key
	{`//`, `.`, `.`}, // empty string key
	{``, ``, ``},     // whole content (root)
}

func TestDotNotation(t *testing.T) {
	for caseIndex, testCase := range testDotNotationCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Fatal(err)
		}
		actual := pointer.DotNotation(false)
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
		actual = pointer.DotNotation(true)
		if actual != testCase.expectedBracket {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expectedBracket, actual)
		}
	}
}