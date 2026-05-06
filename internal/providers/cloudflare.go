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

type CFInternalProvider struct {
	url string
}

func NewCFInternalProvider() *CFInternalProvider {
	return &CFInternalProvider{url: "https://1.1.1.1/cdn-cgi/trace"}
}

func (p *CFInternalProvider) Name() string { return "cloudflare" }

func (p *CFInternalProvider) GetIP(ctx context.Context, ipType models.IPType) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", p.url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ip=") {
			return strings.TrimPrefix(line, "ip="), nil
		}
	}
	return "", fmt.Errorf("ip field not found in cloudflare trace")
}
