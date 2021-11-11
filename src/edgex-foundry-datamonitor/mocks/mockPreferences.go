package mocks

import "github.com/stretchr/testify/mock"

type MockPreferences struct {
	mock.Mock
}

// The panics here are normal, they are here to signal the developer that these methods are not implemented
// Only the ones I needed are implemented and these are necessary to make sure that the mock implements the mocked interface

// Bool looks up a boolean value for the key
func (p *MockPreferences) Bool(key string) bool {
	panic("not implemented") // TODO: Implement
}

// BoolWithFallback looks up a boolean value and returns the given fallback if not found
func (p *MockPreferences) BoolWithFallback(key string, fallback bool) bool {
	panic("not implemented") // TODO: Implement
}

// SetBool saves a boolean value for the given key
func (p *MockPreferences) SetBool(key string, value bool) {
	panic("not implemented") // TODO: Implement
}

// Float looks up a float64 value for the key
func (p *MockPreferences) Float(key string) float64 {
	panic("not implemented") // TODO: Implement
}

// FloatWithFallback looks up a float64 value and returns the given fallback if not found
func (p *MockPreferences) FloatWithFallback(key string, fallback float64) float64 {
	panic("not implemented") // TODO: Implement
}

// SetFloat saves a float64 value for the given key
func (p *MockPreferences) SetFloat(key string, value float64) {
	panic("not implemented") // TODO: Implement
}

// Int looks up an integer value for the key
func (p *MockPreferences) Int(key string) int {
	panic("not implemented") // TODO: Implement
}

// IntWithFallback looks up an integer value and returns the given fallback if not found
func (p *MockPreferences) IntWithFallback(key string, fallback int) int {
	args := p.Called(key, fallback)
	return args.Int(0)
}

// SetInt saves an integer value for the given key
func (p *MockPreferences) SetInt(key string, value int) {
	panic("not implemented") // TODO: Implement
}

// String looks up a string value for the key
func (p *MockPreferences) String(key string) string {
	panic("not implemented") // TODO: Implement
}

// StringWithFallback looks up a string value and returns the given fallback if not found
func (p *MockPreferences) StringWithFallback(key string, fallback string) string {
	args := p.Called(key, fallback)
	return args.String(0)
}

// SetString saves a string value for the given key
func (p *MockPreferences) SetString(key string, value string) {
	panic("not implemented") // TODO: Implement
}

// RemoveValue removes a value for the given key (not currently supported on iOS)
func (p *MockPreferences) RemoveValue(key string) {
	panic("not implemented") // TODO: Implement
}

// AddChangeListener allows code to be notified when some preferences change. This will fire on any update.
func (p *MockPreferences) AddChangeListener(_ func()) {
	panic("not implemented") // TODO: Implement
}

// Preferences describes the ways that an app can save and load user preferences
type Preferences interface {
	// Bool looks up a boolean value for the key
	Bool(key string) bool
	// BoolWithFallback looks up a boolean value and returns the given fallback if not found
	BoolWithFallback(key string, fallback bool) bool
	// SetBool saves a boolean value for the given key
	SetBool(key string, value bool)

	// Float looks up a float64 value for the key
	Float(key string) float64
	// FloatWithFallback looks up a float64 value and returns the given fallback if not found
	FloatWithFallback(key string, fallback float64) float64
	// SetFloat saves a float64 value for the given key
	SetFloat(key string, value float64)

	// Int looks up an integer value for the key
	Int(key string) int
	// IntWithFallback looks up an integer value and returns the given fallback if not found
	IntWithFallback(key string, fallback int) int
	// SetInt saves an integer value for the given key
	SetInt(key string, value int)

	// String looks up a string value for the key
	String(key string) string
	// StringWithFallback looks up a string value and returns the given fallback if not found
	StringWithFallback(key, fallback string) string
	// SetString saves a string value for the given key
	SetString(key string, value string)

	// RemoveValue removes a value for the given key (not currently supported on iOS)
	RemoveValue(key string)

	// AddChangeListener allows code to be notified when some preferences change. This will fire on any update.
	AddChangeListener(func())
}
