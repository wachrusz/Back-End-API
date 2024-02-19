//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

type Settings struct {
	Subscriptions Subscription `json:"subscriptions"`
}
