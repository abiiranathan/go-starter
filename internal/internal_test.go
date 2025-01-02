package internal_test

import (
	"testing"

	"github.com/abiiranathan/go-starter/internal"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Parallel()

	password := "password"
	hashedPassword, err := internal.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Fatal("expected hashed password to have length greater than 0")
	}

	// verify that the hashed password is valid
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Fatalf("expected hashed password to be valid: %v", err)
	}
}

func TestRandomString(t *testing.T) {
	t.Parallel()

	n := 10
	randomString := internal.RandomString(n)
	if len(randomString) != n {
		t.Fatalf("expected random string to have length %d, got %d", n, len(randomString))
	}

	for i := 5; i < 20; i++ {
		randomString := internal.RandomString(i)
		if len(randomString) != i {
			t.Fatalf("expected random string to have length %d, got %d", i, len(randomString))
		}

		t.Logf("random string: %s", randomString)
	}
}

func TestRandomInt(t *testing.T) {
	t.Parallel()

	min := 5
	max := 10
	randomInt := internal.RandomInt(min, max)
	if randomInt < min || randomInt > max {
		t.Fatalf("expected random int to be between %d and %d, got %d", min, max, randomInt)
	}

	for i := 0; i < 10; i++ {
		randomInt := internal.RandomInt(0, 100)
		if randomInt < 0 || randomInt > 100 {
			t.Fatalf("expected random int to be between 0 and 100, got %d", randomInt)
		}

		t.Logf("random int: %d", randomInt)
	}
}

func TestRandomInt64(t *testing.T) {
	t.Parallel()

	min := int64(5)
	max := int64(10)
	randomInt := internal.RandomInt64(min, max)
	if randomInt < min || randomInt > max {
		t.Fatalf("expected random int to be between %d and %d, got %d", min, max, randomInt)
	}

	for i := 0; i < 10; i++ {
		randomInt := internal.RandomInt64(0, 100)
		if randomInt < 0 || randomInt > 100 {
			t.Fatalf("expected random int to be between 0 and 100, got %d", randomInt)
		}

		t.Logf("random int: %d", randomInt)
	}
}

func TestRandomFloat(t *testing.T) {
	t.Parallel()

	min := 5.0
	max := 10.0
	randomFloat := internal.RandomFloat(min, max)
	if randomFloat < min || randomFloat > max {
		t.Fatalf("expected random float to be between %f and %f, got %f", min, max, randomFloat)
	}

	for i := 0; i < 10; i++ {
		randomFloat := internal.RandomFloat(0, 100)
		if randomFloat < 0 || randomFloat > 100 {
			t.Fatalf("expected random float to be between 0 and 100, got %f", randomFloat)
		}

		t.Logf("random float: %f", randomFloat)
	}
}

func TestSlugify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple Text", "Hello World", "hello-world"},
		{"With Special Characters", "Hello, World!", "hello-world"},
		{"With Spaces and Tabs", "Hello\t World   ", "hello-world"},
		{"Non-ASCII Characters", "¡Hola Mundo!", "hola-mundo"},
		{"Empty String", "", ""},
		{"Already Slugified", "already-slugified", "already-slugified"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.Slugify(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid Email", "test@example.com", true},
		{"Valid Email", "Test <test@example.com>", true},
		{"Invalid Email - No Domain", "test@", false},
		{"Invalid Email - No Username", "@example.com", false},
		{"Invalid Email - No @ Symbol", "testexample.com", false},
		{"Empty String", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
