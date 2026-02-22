package models

// To ensure later that enums are validated
type Enum interface {
    IsValid() bool
}
