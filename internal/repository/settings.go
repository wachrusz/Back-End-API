//go:build !exclude_swagger
// +build !exclude_swagger

// Package repository provides basic financial repository functionality.
package repository

type Settings struct {
	Subscriptions Subscription `json:"subscriptions"`
}
