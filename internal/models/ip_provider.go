package models

import "context"

// IPType defines whether we are looking for IPv4 or IPv6.
type IPType int

const (
	IPV4 IPType = iota
	IPV6
)

// Provider is the interface implemented by all public IP detection services.
type Provider interface {
	Name() string
	GetIP(ctx context.Context, ipType IPType) (string, error)
}
