package providers

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/victorgomez09/dynamic-dns/internal/models"
)

type IpifyProvider struct {
	url4 string
	url6 string
}

func NewIpifyProvider() *IpifyProvider {
	return &IpifyProvider{
		url4: "https://api.ipify.org",
		url6: "https://api64.ipify.org",
	}
}

func (p *IpifyProvider) Name() string { return "ipify" }

func (p *IpifyProvider) GetIP(ctx context.Context, ipType models.IPType) (string, error) {
	url := p.url4
	if ipType == models.IPV6 {
		url = p.url6
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body)), nil
}
