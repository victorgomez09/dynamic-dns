package cloudflare

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

type Client struct {
	api *cloudflare.API
}

func NewClient(token string) (*Client, error) {
	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return nil, err
	}
	log.Printf("✅ Connected to the Cloudflare API")

	return &Client{api: api}, nil
}

func (c *Client) UpdateDNSRecord(ctx context.Context, recordName string, newIp string, dryRun bool) error {
	// Extract the base domain to resolve the zone ID.
	parts := strings.Split(recordName, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid record name: %s", recordName)
	}

	zoneName := strings.Join(parts[len(parts)-2:], ".")
	zoneID, err := c.api.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("resolve zone %s: %w", zoneName, err)
	}

	// Look up the DNS record within that zone.
	records, _, err := c.api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Name: recordName,
		Type: "A",
	})
	if err != nil {
		return fmt.Errorf("list DNS records for zone %s and record %s: %w", zoneName, recordName, err)
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found for %s", recordName)
	}

	record := records[0]
	// Skip the update when the current value already matches.
	if record.Content == newIp {
		log.Printf("☁️  [Cloudflare] %s is already up to date with IP %s", recordName, newIp)
		return nil
	}
	if dryRun {
		log.Printf("🧪 [DryRun] Would update %s: %s -> %s", recordName, record.Content, newIp)
		return nil
	}

	log.Printf("🔄 [Cloudflare] Updating %s: %s -> %s", recordName, record.Content, newIp)
	_, err = c.api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
		ID:      record.ID,
		Content: newIp,
	})
	if err != nil {
		return fmt.Errorf("update DNS record %s in zone %s: %w", recordName, zoneName, err)
	}
	log.Printf("✅ [Cloudflare] %s updated successfully to %s", recordName, newIp)

	return nil
}
