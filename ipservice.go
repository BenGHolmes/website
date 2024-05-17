package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/robfig/cron/v3"
)

type IpService interface {
	Start()
}

type ipService struct {
	cloudflare *cloudflare.API
}

func New() (IpService, error) {
	cloudflare, err := cloudflare.New(
		os.Getenv("CLOUDFLARE_API_KEY"),
		os.Getenv("CLOUDFLARE_API_EMAIL"),
	)
	if err != nil {
		return nil, err
	}

	return &ipService{cloudflare}, nil
}

func (ips *ipService) Start() {
	c := cron.New()
	c.AddFunc("*/5 * * * * *", ips.verifyCorrectIPAddress)
}

func (ips *ipService) verifyCorrectIPAddress() {
	oldIp, err := getPreviousPublicIP()
	if err != nil {
		// TODO: logger.Error
	}

	newIp, err := getCurrentPublicIP()
	if err != nil {
		// TODO: logger.Error
	}

	// Return if the IP matches
	if oldIp == newIp {
		// TODO: logger.Info
	}

	err = ips.updateIPAddress(*oldIp, *newIp)
	if err != nil {
		// TODO: logger.Warn
	}
}

func (ips *ipService) updateIPAddress(oldIP, newIP string) error {
	ctx := context.Background()
	zones, err := ips.cloudflare.ListZones(ctx)
	if err != nil {
		return err
	}

	// For each zone, look for outdated A records
	for _, zone := range zones {
		// Ignoring info since I'm assuming we won't have multiple pages for a single result
		outdated, _, err := ips.cloudflare.ListDNSRecords(
			ctx,
			cloudflare.ZoneIdentifier(zone.ID),
			cloudflare.ListDNSRecordsParams{
				Type:    "A",
				Content: oldIP,
			},
		)
		if err != nil {
			return err
		}

		// Update all outdated records
		for _, record := range outdated {
			_, err := ips.cloudflare.UpdateDNSRecord(
				ctx,
				cloudflare.ZoneIdentifier(zone.ID),
				cloudflare.UpdateDNSRecordParams{
					ID:      record.ID,
					Type:    record.Type,
					Name:    record.Name,
					Content: newIP,
					TTL:     record.TTL,
					Proxied: record.Proxied,
				},
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getPreviousPublicIP() (*string, error) {
	data, err := os.ReadFile(".public-ip")
	if err != nil {
		return nil, err
	}

	ip := string(data)

	return &ip, nil
}

func getCurrentPublicIP() (*string, error) {
	resp, err := http.Get("ipv4.icanhazip.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error with status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ip := string(data)
	return &ip, nil
}
