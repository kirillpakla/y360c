package helper

import (
	"strings"
	"testing"
	"unicode"
)

func TestGenerateTemporaryPasswordShape(t *testing.T) {
	first, err := GenerateTemporaryPassword()
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateTemporaryPassword()
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != temporaryPasswordLength || len(second) != temporaryPasswordLength {
		t.Fatalf("lengths = %d %d", len(first), len(second))
	}
	if first == second {
		t.Fatal("two generated passwords were identical")
	}
	for _, password := range []string{first, second} {
		var upper, lower, digit bool
		for _, r := range password {
			switch {
			case unicode.IsUpper(r):
				upper = true
			case unicode.IsLower(r):
				lower = true
			case unicode.IsDigit(r):
				digit = true
			default:
				t.Fatalf("password contains %q", r)
			}
		}
		if !upper || !lower || !digit {
			t.Fatalf("password classes upper=%v lower=%v digit=%v", upper, lower, digit)
		}
		if strings.ContainsAny(password, "0O1Il") {
			t.Fatal("password contains an excluded ambiguous character")
		}
	}
}
