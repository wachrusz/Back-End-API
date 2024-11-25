//go:build !exclude_swagger
// +build !exclude_swagger

// Package repository provides basic financial repository functionality.
package repository

import "github.com/wachrusz/Back-End-API/internal/repository/models"

type Settings struct {
	Subscriptions models.Subscription `json:"subscriptions"`
}
