# 🚀 Cloudflare DDNS Multi-Provider

An ultra-lightweight, resilient, cloud-native Dynamic DNS (DDNS) client written in **Go**. It is designed for modern **Kubernetes** infrastructure and **Raspberry Pi (ARM64)** devices.

---

## ✨ Features

*   **🛡️ Multi-Provider Resilience:** Queries a prioritized list of services (Amazon, Cloudflare, ICanHazIP, and more). If one fails, the system automatically fails over to the next one.
*   **🌐 Native Dual-Stack Support:** Simultaneous detection and update support for **IPv4 (A)** and **IPv6 (AAAA)** records.
*   **🔍 Zone ID Auto-Discovery:** Automatically resolves the Cloudflare zone ID from the configured domain name, reducing manual setup.
*   **⚡ Zero-Footprint (Distroless):** Statically compiled binary ready to run on a minimal `scratch` or `alpine` image (~12MB), minimizing OS attack surface.
*   **🔄 Exponential Backoff:** Smarter error handling to avoid hitting Cloudflare API limits too aggressively.
*   **🧪 Dry-Run Mode:** Lets you verify the changes that would be applied without modifying production records.

---

## 🛠️ Supported Public IP Providers

The system uses an internal plugin architecture to discover the public IP. These providers are currently implemented:

| Provider | Endpoint | Support |
| :--- | :--- | :--- |
| **Amazon** | `https://checkip.amazonaws.com` | IPv4 |
| **Cloudflare** | `https://1.1.1.1/cdn-cgi/trace` | IPv4/IPv6 |
| **ICanHazIP** | `https://ipv4.icanhazip.com` / `https://ipv6.icanhazip.com` | IPv4/IPv6 |
| **Ipify** | `https://api.ipify.org` / `https://api64.ipify.org` | IPv4/IPv6 |

Valid values for `IP_PROVIDERS` are: `amazon`, `cloudflare`, `icanhazip`, `ipify`.

---

## ⚙️ Configuration (Environment Variables)

Configuration is handled entirely through environment variables, which makes it easy to use with **Kubernetes Secrets**:

| Variable | Description | Example Value |
| :--- | :--- | :--- |
| `CF_API_TOKEN` | API token with DNS edit permissions | `qJ8...` |
| `CF_RECORD_NAMES` | Comma-separated records to update (A or AAAA) | `dns.esmosolutions.es` |
| `IP_PROVIDERS` | Ordered list of IP providers | `amazon,cloudflare,icanhazip,ipify` |
| `CHECK_INTERVAL` | How often to check | `10m` |
| `ENABLE_IPV6` | Enable AAAA record updates | `true` |
| `DRY_RUN` | Simulate changes without applying them | `false` |

---

## 🏗️ Project Structure
```text
.
├── cmd/main.go        # Main loop and signal handling (graceful shutdown)
├── internal/
│   ├── config/        # Environment variable parser
│   ├── cloudflare/    # API client and auto-discovery logic
│   └── providers/     # Public IP detection logic (Amazon, Cloudflare, etc.)
├── Dockerfile         # Multi-stage ARM64 build
└── deployment.yaml    # Kubernetes manifest