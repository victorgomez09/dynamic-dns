package providers

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/victorgomez09/dynamic-dns/internal/models"
)

type ICanHazIPProvider struct {
	url4 string
	url6 string
}

func NewICanHazIPProvider() *ICanHazIPProvider {
	return &ICanHazIPProvider{
		url4: "https://ipv4.icanhazip.com",
		url6: "https://ipv6.icanhazip.com",
	}
}

func (p *ICanHazIPProvider) Name() string { return "icanhazip" }

func (p *ICanHazIPProvider) GetIP(ctx context.Context, ipType models.IPType) (string, error) {
	url := p.url4
	if ipType == models.IPV6 {
		url = p.url6
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	// Some of these services behave better with a simple User-Agent.
	req.Header.Set("User-Agent", "Cloudflare-DDNS-Go/2026")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body)), nil
}
