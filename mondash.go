package main

import "time"

type mondashJSONDashboard struct {
	Metrics []mondashJSONDashboardMetric `json:"metrics"`
}

type mondashJSONDashboardMetric struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Value       float64   `json:"value"`
	LastUpdate  time.Time `json:"last_update"`
}
