package utils

import "context"

// Checker is a function that we can give to the api endpoint
// If it returns false or an error, the server will be set as "down"
type Checker func(ctx context.Context) (bool, error)

type Checkable interface {
	ReadyCheck() Checker
}
