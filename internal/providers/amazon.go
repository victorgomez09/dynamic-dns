package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/victorgomez09/dynamic-dns/internal/models"
)

type AmazonProvider struct {
	url string
}

func NewAmazonProvider() *AmazonProvider {
	return &AmazonProvider{
		url: "https://checkip.amazonaws.com",
	}
}

func (p *AmazonProvider) Name() string {
	return "amazon"
}

func (p *AmazonProvider) GetIP(ctx context.Context, ipType models.IPType) (string, error) {
	if ipType == models.IPV6 {
		return "", fmt.Errorf("amazon provider does not support IPv6 yet") // Amazon does not support IPv6 yet.
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", p.url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := strings.TrimSpace(string(body))

	if ip == "" {
		return "", fmt.Errorf("empty response from amazon provider")
	}

	return ip, nil
}
