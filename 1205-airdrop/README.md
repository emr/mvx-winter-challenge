# Winter Token Distributor

The script distributes WinterToken that owned by the accounts whose PEM files provided in the `accounts` directory to randomly generated addresses on testnet.

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
   cd 1205-airdrop/
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
- Load all PEM files from the `accounts` directory
- For each account, find all owned tokens named "WinterToken"
- For each token, distribute 10000 units to 1000 randomly generated addresses
