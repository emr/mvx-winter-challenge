# MultiversX Winter Token Challenge

This script issues three ESDT tokens on MultiversX testnet for each account whose PEM files provided in the `accounts` directory.

## Prerequisites

- Go 1.23 or later
- PEM files in the `accounts` directory

## Installation

1. Clone the repository
   ```bash
   git clone https://github.com/emr/mvx-winter-challenge.git
   ```
2. Go to the project directory
   ```bash
   cd 1204-winter-tokens/
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

Run the script:

```bash
go run main.go
```

The script will:
1. Load all PEM files from the `accounts` directory
2. For each account, issue 3 tokens with:
   - Name: WinterToken
   - Ticker: WINTER
   - Supply: 100 million
   - Decimals: 8
3. Print the transaction hashes.

Preview:
![preview](preview.mp4)
