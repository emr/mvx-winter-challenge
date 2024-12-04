# MultiversX Account Generator

This script generates MultiversX testnet accounts and requests tokens from the faucet automatically.

## Prerequisites

- Go 1.23 or later
- Chrome/Chromium browser

## Installation

1. Clone the repository
   ```bash
   git clone https://github.com/emr/mvx-winter-challenge.git
   ```
2. Go to the project directory
   ```bash
   cd 1203-create-and-fund-accounts/
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage
Run the script:

```bash
go run cmd/main.go
```

The script will:
1. Generate 9 accounts (3 per shard)
2. Save their PEM files to the `accounts` directory
3. Open Chrome browser to use web-wallet
4. Request tokens from faucet for each account
5. You'll need to solve the CAPTCHA manually for each request
