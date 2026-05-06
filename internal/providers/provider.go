package providers

import (
	"fmt"

	"github.com/victorgomez09/dynamic-dns/internal/models"
)

func GetProviderByName(name string) (models.Provider, error) {
	switch name {
	case "amazon":
		return NewAmazonProvider(), nil
	case "cloudflare":
		return NewCFInternalProvider(), nil
	case "ipify":
		return NewIpifyProvider(), nil
	case "icanhazip":
		return NewICanHazIPProvider(), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
