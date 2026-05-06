package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/victorgomez09/dynamic-dns/internal/cloudflare"
	"github.com/victorgomez09/dynamic-dns/internal/config"
	"github.com/victorgomez09/dynamic-dns/internal/models"
	"github.com/victorgomez09/dynamic-dns/internal/providers"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("🚀 Starting Cloudflare DDNS")
	log.Printf("📡 Records to monitor: %s", cfg.RecordNames)
	log.Printf("⏱️ Check interval: %v", cfg.CheckInterval)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	run(ctx, cfg)

	ticker := time.NewTicker(cfg.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			run(ctx, cfg)
		case <-ctx.Done():
			log.Println("Shutting down gracefully...")
			return
		}
	}
}

func run(ctx context.Context, cfg *config.Config) {
	var currentIP string
	var providerUsed models.Provider

	log.Println("🔍 Starting IP verification round...")

	// 1. Try each provider defined in the configuration.
	for _, name := range cfg.Providers {
		p, err := providers.GetProviderByName(name)
		if err != nil {
			log.Printf("⚠️  Configuration: unknown provider '%s', skipping...", name)
			continue
		}

		// Try to fetch the IP with a short timeout per provider so a slow one
		// does not block the whole cycle.
		childCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
		ip, err := p.GetIP(childCtx, models.IPV4)
		cancel() // Release resources immediately.

		if err != nil {
			log.Printf("❌ [%s] Failed: %v", p.Name(), err)
			continue // Move on to the next provider.
		}

		if ip != "" {
			currentIP = ip
			providerUsed = p
			break // Success. Stop iterating providers.
		}
	}

	// 2. If no provider returned an IP, abort this cycle.
	if currentIP == "" {
		log.Println("🚨 CRITICAL ERROR: No IP provider responded. Retrying on the next interval.")
		return
	}

	log.Printf("✨ [%s] IP detected successfully: %s", providerUsed.Name(), currentIP)

	// 3. Proceed with the Cloudflare update.
	updateCloudflare(ctx, cfg, currentIP)
}

func updateCloudflare(ctx context.Context, cfg *config.Config, ip string) {
	cfClient, err := cloudflare.NewClient(cfg.ApiToken)
	if err != nil {
		log.Printf("❌ Error connecting to the Cloudflare API: %v", err)
		return
	}

	for _, record := range cfg.RecordNames {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		err := cfClient.UpdateDNSRecord(ctx, record, ip, cfg.DryRun)
		if err != nil {
			log.Printf("❌ [%s] Error updating record: %v", record, err)
		}
	}
}
