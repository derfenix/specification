package internal

import (
	"errors"
	"strings"
)

const defaultSeparator = "; "

type Multierr struct {
	separator string
	errors    []error
}

func NewMultierr() *Multierr {
	return &Multierr{}
}

func (m *Multierr) WithSeparator(s string) *Multierr {
	m.separator = s

	return m
}

func (m *Multierr) WithCap(c int) *Multierr {
	m.errors = make([]error, 0, c)

	return m
}

func (m *Multierr) Error() string {
	if len(m.errors) == 0 {
		return ""
	}

	errStrings := make([]string, len(m.errors))
	for idx := range errStrings {
		errStrings[idx] = m.errors[idx].Error()
	}

	separator := m.separator

	if separator == "" {
		separator = defaultSeparator
	}

	return strings.Join(errStrings, separator)
}

func (m *Multierr) Append(err error) {
	m.errors = append(m.errors, err)
}

func (m *Multierr) Errors() []error {
	if len(m.errors) == 0 {
		return nil
	}

	return m.errors
}

func (m *Multierr) Is(target error) bool {
	for _, err := range m.errors {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (m *Multierr) As(target interface{}) bool {
	for _, err := range m.Errors() {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}
