package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Luzifer/rconfig"
	"github.com/opsgenie/opsgenie-go-sdk/alerts"
	"github.com/opsgenie/opsgenie-go-sdk/client"
)

var (
	version = "dev"
	cfg     = struct {
		MondashURL     string   `flag:"mondash-url" default:"https://mondash.org/" description:"URL of the mondash installation"`
		MondashBoard   string   `flag:"mondash-board" description:"Board ID of the mondash board to monitor"`
		OpsGenieAPIKey string   `flag:"opsgenie-key" description:"API-Key for OpsGenie API integration"`
		AlertStatus    []string `flag:"alert,a" default:"Critical,Unknown" description:"List of status to trigger alerts for"`
		OKStatus       []string `flag:"ok,o" default:"OK" description:"List of status to resolve alerts for"`
	}{}
)

func validateConfig() error {
	if len(cfg.OKStatus) == 0 {
		return errors.New("Need at least one OK status")
	}

	if len(cfg.AlertStatus) == 0 {
		return errors.New("Need at least one Alert status")
	}

	if len(cfg.OpsGenieAPIKey) == 0 {
		return errors.New("Need OpsGenie API key")
	}

	if len(cfg.MondashBoard) == 0 {
		return errors.New("Need Mondash board ID")
	}

	return nil
}

func main() {
	rconfig.Parse(&cfg)
	if err := validateConfig(); err != nil {
		log.Printf("An error occurred while reading the configuration: %s", err)
		os.Exit(1)
	}

	cli := new(client.OpsGenieClient)
	cli.SetAPIKey(cfg.OpsGenieAPIKey)

	alertCli, cliErr := cli.Alert()
	if cliErr != nil {
		log.Fatalf("Unable to open API connection: %s", cliErr)
	}
	v
	boardURL := fmt.Sprintf("%s/%s.json", strings.TrimRight(cfg.MondashURL, "/"), cfg.MondashBoard)
	res, err := http.Get(boardURL)

	if err != nil {
		log.Fatalf("Unable to fetch dashboard JSON: %s", err)
	}
	defer res.Body.Close()

	board := mondashJSONDashboard{}
	if err := json.NewDecoder(res.Body).Decode(&board); err != nil {
		log.Fatalf("Unable to read dashboard JSON: %s", err)
	}

	var alertTitles []string
	var okTitles []string

	for _, metric := range board.Metrics {
		for _, a := range cfg.AlertStatus {
			if a == metric.Status {
				alertTitles = append(alertTitles, fmt.Sprintf("%s (%s)", metric.Title, metric.Status))
			}
		}

		for _, o := range cfg.OKStatus {
			if o == metric.Status {
				okTitles = append(okTitles, metric.Title)
			}
		}
	}

	rs, err := alertCli.List(alerts.ListAlertsRequest{
		Status: "open",
	})
	if err != nil {
		log.Fatalf("Unable to fetch open alerts: %s", err)
	}

	var existingAlertPresent bool
	for _, alert := range rs.Alerts {
		if alert.Alias == fmt.Sprintf("mondash_%s", cfg.MondashBoard) {
			existingAlertPresent = true
			break
		}
	}

	if len(alertTitles) > 0 && !existingAlertPresent {
		res, err := alertCli.Create(alerts.CreateAlertRequest{
			Message:     fmt.Sprintf("Mondash board %s has metrics in alert status", cfg.MondashBoard),
			Alias:       fmt.Sprintf("mondash_%s", cfg.MondashBoard),
			Description: fmt.Sprintf("Checks in alert status:\n - %s", strings.Join(alertTitles, "\n - ")),
			Source:      fmt.Sprintf("mondash-opsgenie %s", version),
		})
		if err != nil {
			log.Fatalf("Unable to open alert: %s", err)
		}

		log.Printf("Created alert %s", res.AlertID)
	} else if len(okTitles) == len(board.Metrics) && existingAlertPresent {
		res, err := alertCli.Close(alerts.CloseAlertRequest{
			Alias: fmt.Sprintf("mondash_%s", cfg.MondashBoard),
		})
		if err != nil {
			log.Fatalf("Unable to close alert: %s", err)
		}

		log.Printf("Closed alert: %s", res.Status)
	} else {
		log.Printf("Nothing to do for me. (existingAlertPresent = %v)", existingAlertPresent)
	}
}
