package factory

import (
	"errors"
	"fmt"
)

type ErrorFactory struct{}

func (*ErrorFactory) InvalidCredentials() error {
	return errors.New("invalid credentials")
}

func (*ErrorFactory) Malformatted(target string) error {
	return fmt.Errorf("malformatted %s", target)
}

func (*ErrorFactory) StoreSessionFailed() error {
	return errors.New("failed to store session data")
}

func (*ErrorFactory) TimeOut(target string) error { return fmt.Errorf("%s timeout", target) }

func (*ErrorFactory) NotFound(target string) error {
	return fmt.Errorf("%s not found", target)
}

func (*ErrorFactory) Unexpected() error {
	return errors.New("unexpected error")
}

func (*ErrorFactory) AlreadyExists(target string) error {
	return fmt.Errorf("%s already exists", target)
}

func (*ErrorFactory) NotEnabled(target string) error { return fmt.Errorf("%s is not enabled", target) }
