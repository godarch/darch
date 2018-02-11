package reference

import (
	"testing"
)

type referenceTest struct {
	Input            string
	ShouldFail       bool
	ExpectedDomain   string
	ExpectedName     string
	ExpectedTag      string
	ExpectedFullName string
}

var (
	referenceTests = []referenceTest{
		{
			Input:            "base/image",
			ExpectedDomain:   "",
			ExpectedName:     "base/image",
			ExpectedTag:      "latest",
			ExpectedFullName: "base/base:latest",
		},
		{
			Input:            "base",
			ExpectedDomain:   "",
			ExpectedName:     "base",
			ExpectedTag:      "latest",
			ExpectedFullName: "base:latest",
		},
		{
			Input:            "base:tag",
			ExpectedDomain:   "",
			ExpectedName:     "base",
			ExpectedTag:      "tag",
			ExpectedFullName: "base:tag",
		},
		{
			Input:            "domain.com/base:tag",
			ExpectedDomain:   "domain.com",
			ExpectedName:     "domain.com/base",
			ExpectedTag:      "tag",
			ExpectedFullName: "domain.com/base:tag",
		},
		{
			Input:            "domain_com/base:tag",
			ExpectedDomain:   "",
			ExpectedName:     "domain_com/base",
			ExpectedTag:      "tag",
			ExpectedFullName: "domain_com/base:tag",
		},
		{
			Input:      "",
			ShouldFail: true,
		},
		{
			Input:      ":justtag",
			ShouldFail: true,
		},
		{
			Input:      "aa/asdf$$^/aa",
			ShouldFail: true,
		},
	}
)

func TestParsing(t *testing.T) {
	for _, referenceTest := range referenceTests {
		ref, err := ParseImage(referenceTest.Input)
		if referenceTest.ShouldFail {
			if err == nil {
				t.Fatalf("input %s should have failed", referenceTest.Input)
			}
		} else {
			if ref.Domain() != referenceTest.ExpectedDomain {
				t.Fatalf("expected domain \"%s\", got \"%s\"", referenceTest.ExpectedDomain, ref.Domain())
			}
			if ref.Name() != referenceTest.ExpectedName {
				t.Fatalf("expected name \"%s\", got \"%s\"", referenceTest.ExpectedName, ref.Name())
			}
			if ref.Tag() != referenceTest.ExpectedTag {
				t.Fatalf("expected tag \"%s\", got \"%s\"", referenceTest.ExpectedTag, ref.Tag())
			}
		}
	}
}

type referenceWithTagTest struct {
	Input          string
	WithTag        string
	ShouldFail     bool
	ExpcetedResult string
}

var (
	referenceWithTagTests = []referenceWithTagTest{
		{
			Input:          "base",
			WithTag:        "test",
			ExpcetedResult: "base:test",
		},
		{
			Input:      "base",
			WithTag:    "",
			ShouldFail: true,
		},
	}
)

func TestWithTag(t *testing.T) {
	for _, referenceTest := range referenceWithTagTests {
		ref, err := ParseImage(referenceTest.Input)
		if err != nil {
			t.Fatal(err)
		}
		ref, err = ref.WithTag(referenceTest.WithTag)
		if referenceTest.ShouldFail {
			if err == nil {
				t.Fatalf("input %s should have failed", referenceTest.Input)
			}
		} else {
			if ref.FullName() != referenceTest.ExpcetedResult {
				t.Fatalf("expected \"%s\", got \"%s\"", referenceTest.ExpcetedResult, ref.FullName())
			}
		}
	}
}

type referenceWithNameTest struct {
	Input          string
	WithName       string
	ShouldFail     bool
	ExpcetedResult string
}

var (
	referenceWithNameTests = []referenceWithNameTest{
		{
			Input:          "base",
			WithName:       "test",
			ExpcetedResult: "test:latest",
		},
		{
			Input:      "base",
			WithName:   "",
			ShouldFail: true,
		},
	}
)

func TestWithName(t *testing.T) {
	for _, referenceTest := range referenceWithNameTests {
		ref, err := ParseImage(referenceTest.Input)
		if err != nil {
			t.Fatal(err)
		}
		ref, err = ref.WithName(referenceTest.WithName)
		if referenceTest.ShouldFail {
			if err == nil {
				t.Fatalf("input %s should have failed", referenceTest.Input)
			}
		} else {
			if ref.FullName() != referenceTest.ExpcetedResult {
				t.Fatalf("expected \"%s\", got \"%s\"", referenceTest.ExpcetedResult, ref.FullName())
			}
		}
	}
}

type referenceWithDomainTest struct {
	Input          string
	WithDomain     string
	ShouldFail     bool
	ExpcetedResult string
}

var (
	referenceWithDomainTests = []referenceWithDomainTest{
		{
			Input:          "base",
			WithDomain:     "test.com",
			ExpcetedResult: "test.com/base:latest",
		},
		{
			Input:          "existing.com/base",
			WithDomain:     "new.com",
			ExpcetedResult: "new.com/base:latest",
		},
		{
			Input:          "existing.com/base",
			WithDomain:     "",
			ExpcetedResult: "base:latest",
		},
	}
)

func TestWithDomain(t *testing.T) {
	for _, referenceTest := range referenceWithDomainTests {
		ref, err := ParseImage(referenceTest.Input)
		if err != nil {
			t.Fatal(err)
		}
		ref, err = ref.WithDomain(referenceTest.WithDomain)
		if referenceTest.ShouldFail {
			if err == nil {
				t.Fatalf("input %s should have failed", referenceTest.Input)
			}
		} else {
			if ref.FullName() != referenceTest.ExpcetedResult {
				t.Fatalf("expected \"%s\", got \"%s\"", referenceTest.ExpcetedResult, ref.FullName())
			}
		}
	}
}
