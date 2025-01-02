package internal

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unsafe"

	cryptoRand "crypto/rand"
	"math/rand"
	randV2 "math/rand/v2"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandomString generates a random string of length n.
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
// This link provides a good explanation of how the random string is generated.
func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// RandomInt generates a random integer between min and max(inclusive).
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

// RandomInt64 generates a random int64 between min and max(inclusive).
func RandomInt64(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// RandomFloat generates a random float between min and max(inclusive).
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bcryptPassword), nil
}

// ChaCha8 function that generates a random byte slice using the ChaCha8 algorithm.
func ChaCha8() ([]byte, error) {
	seed := [32]byte{}
	chacha8Seed := randV2.NewChaCha8(seed)
	randomBytes := make([]byte, 32)

	n, err := chacha8Seed.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %v", err)
	}
	if n != 32 {
		return nil, fmt.Errorf("failed to generate random bytes: expected 32 bytes, got %d", n)
	}
	return randomBytes, nil
}

// IsValidEmail checks if the provided email is valid.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsStrongPassword checks if a password meets security requirements:
// - At least 8 characters long
// - Contains at least one uppercase letter
// - Contains at least one lowercase letter
// - Contains at least one number
// - Contains at least one special character
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// GenerateSecureToken generates a cryptographically secure random token of specified length.
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := cryptoRand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateBase64Token generates a base64-encoded secure random token.
func GenerateBase64Token(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := cryptoRand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Truncate returns a truncated string with optional suffix.
func Truncate(str string, maxLength int, suffix string) string {
	if len(str) <= maxLength {
		return str
	}
	return str[:maxLength-len(suffix)] + suffix
}

var (
	nonAlphanumericRegex = regexp.MustCompile("[^a-z0-9-]+")
	multipleHyphenRegex  = regexp.MustCompile("-+")
)

// Slugify converts a string to a URL-friendly slug.
func Slugify(str string) string {
	// Convert to lowercase and replace spaces with hyphens
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, " ", "-")

	// Remove all non-alphanumeric characters except hyphens
	str = nonAlphanumericRegex.ReplaceAllString(str, "")

	// Remove consecutive hyphens
	str = multipleHyphenRegex.ReplaceAllString(str, "-")

	// Remove leading and trailing hyphens
	str = strings.Trim(str, "-")

	return str
}

// TimeAgo returns a human-readable string representing time since the given time.
func TimeAgo(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, pluralize(minutes))
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, pluralize(hours))
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, pluralize(days))
	case duration < 365*24*time.Hour:
		months := int(duration.Hours() / 24 / 30)
		return fmt.Sprintf("%d month%s ago", months, pluralize(months))
	default:
		years := int(duration.Hours() / 24 / 365)
		return fmt.Sprintf("%d year%s ago", years, pluralize(years))
	}
}

// Helper function for TimeAgo
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

// DedupeSlice removes duplicate strings from a slice while preserving order.
func DedupeSlice[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, str := range slice {
		if _, ok := seen[str]; !ok {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}
	return result
}

// FormatFileSize formats a file size in bytes to a human-readable string.
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Dict function that creates a map from a list of key-value pairs.
// Acts as props to pass multiple values to a template.
// When registered by the template engine, it can be used as follows:
// {{ dict "key1" "value1" "key2" "value2" }}.
// This is useful because the template engine does not support passing multiple values directly
// to child components.
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}

	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
