package specification_test

import (
	"context"
	"fmt"

	"github.com/derfenix/specification"
)

// Domain model

type Account struct {
	ID     string
	Active bool
}

// Domain specification

var (
	ErrValidationError = fmt.Errorf("validation error")
	ErrNoIDGenerator   = fmt.Errorf("no id generator specified")
)

type activeAccountSpecification struct {
	value bool
}

func ActiveAccountSpecification(forcedValue bool) *activeAccountSpecification {
	return &activeAccountSpecification{value: forcedValue}
}

func (i *activeAccountSpecification) Validate(*Account) error {
	return nil
}

func (i *activeAccountSpecification) Create(account *Account) error {
	account.Active = i.value

	return nil
}

func MustIDSpecification(idGenerator func() string) *mustIDSpecification {
	return &mustIDSpecification{idGenerator: idGenerator}
}

type mustIDSpecification struct {
	idGenerator func() string
}

func (m *mustIDSpecification) Validate(account *Account) error {
	if account.ID == "" {
		return fmt.Errorf("%w: id can't be empty", ErrValidationError)
	}

	return nil
}

func (m *mustIDSpecification) Create(account *Account) error {
	if account.ID != "" {
		return nil
	}

	if m.idGenerator == nil {
		return ErrNoIDGenerator
	}

	account.ID = m.idGenerator()

	return nil
}

// Domain service

var defaultSpecification = specification.NewAndSpecification[Account](
	MustIDSpecification(
		func() string {
			return "fixed_id"
		},
	),
	ActiveAccountSpecification(false),
)

type Repository interface {
	Save(ctx context.Context, account *Account) error
}

type Service struct {
	repo Repository
}

func (s *Service) Create(ctx context.Context, account *Account, spec specification.Specification[Account]) error {
	if spec == nil {
		spec = defaultSpecification
	}

	if err := spec.Create(account); err != nil {
		return fmt.Errorf("specification create: %w", err)
	}

	if err := spec.Validate(account); err != nil {
		return fmt.Errorf("specification validate: %w", err)
	}

	if err := s.repo.Save(ctx, account); err != nil {
		return fmt.Errorf("save account: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, account *Account, spec specification.Specification[Account]) error {
	if spec == nil {
		spec = defaultSpecification
	}

	if err := spec.Validate(account); err != nil {
		return fmt.Errorf("specification validate: %w", err)
	}

	if err := s.repo.Save(ctx, account); err != nil {
		return fmt.Errorf("sace account: %w", err)
	}

	return nil
}
