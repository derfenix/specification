package specification

import (
	"github.com/derfenix/specification/internal"
)

type Specification[T any] interface {
	Validate(*T) error
	Create(*T) error
}

type andSpecification[T any] struct {
	specs []Specification[T]
}

// NewAndSpecification holds a slice of specifications and run each of them for Create and Validate methods.
// Fails if *any* specification failed.
func NewAndSpecification[T any](specs ...Specification[T]) Specification[T] {
	return &andSpecification[T]{specs: specs}
}

// Validate returns error if Validate method of any specification returns error.
func (s *andSpecification[T]) Validate(item *T) error {
	for _, spec := range s.specs {
		if err := spec.Validate(item); err != nil {
			return err
		}
	}

	return nil
}

// Create returns error if Create method of any specification returns error.
func (s *andSpecification[T]) Create(item *T) error {
	for _, spec := range s.specs {
		if err := spec.Create(item); err != nil {
			return err
		}
	}

	return nil
}

// OrSpecification holds a slice of specifications and run each of them for Validate methods.
// Fails if *all* specifications failed.
func OrSpecification[T any](specs ...Specification[T]) Specification[T] {
	return &orSpecification[T]{specs: specs}
}

type orSpecification[T any] struct {
	specs []Specification[T]
}

// Validate returns error if Validate method of *all* specification returns error.
func (o *orSpecification[T]) Validate(item *T) error {
	errs := internal.NewMultierr().WithCap(len(o.specs))

	for _, spec := range o.specs {
		if err := spec.Validate(item); err != nil {
			errs.Append(err)
		} else {
			return nil
		}
	}

	return errs
}

// Create is a no-op method.
func (o *orSpecification[T]) Create(*T) error {
	return nil
}
