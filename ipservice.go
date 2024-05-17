package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
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
	fmt.Println("INFO: Starting IP Service")
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-ticker.C
			ips.verifyCorrectIPAddress()
		}
	}()
}

func (ips *ipService) verifyCorrectIPAddress() {
	fmt.Println("INFO: verifying correct IP address")

	oldIp, err := getPreviousPublicIP()
	if err != nil {
		fmt.Printf("ERROR: couldn't get previous public IP: %v\n", err)
		return
	}

	newIp, err := getCurrentPublicIP()
	if err != nil {
		fmt.Printf("ERROR: couldn't get current public IP: %v\n", err)
		return
	}

	// Return if the IP matches
	if *oldIp == *newIp {
		fmt.Println("INFO: IP has not changed")
		return
	}

	err = ips.updateIPAddress(*oldIp, *newIp)
	if err != nil {
		fmt.Printf("ERROR: couldn't update IP: %v\n", err)
		return
	}

	os.WriteFile(".public-ip", []byte(*newIp), 0644)
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
			fmt.Printf("INFO: Updating record %s\n", record.ID)

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
	resp, err := http.Get("https://ipv4.icanhazip.com")
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
