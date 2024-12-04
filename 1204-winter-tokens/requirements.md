## Requirements
Create a script that issues three tokens (fungible) on the MultiversX testnet for each account whose PEM files provided in the `accounts` directory, with a supply of 100 million tokens each.

### Token issuing
Ticker should be `WINTER`. It should have 8 decimals and 100 million supply.

### Usability
The script should be executed in a single command. Provide clear instruction on how to execute the script on README.md file.

## Resources
You must inspect all the resources below:
- Official MultiversX documentation: https://docs.multiversx.com
- Related part: https://docs.multiversx.com/tokens/fungible-tokens
- MultiversX Go SDK:
    - https://github.com/multiversx/mx-sdk-go
    - https://pkg.go.dev/github.com/multiversx/mx-sdk-erdgo
- SDK usage examples: https://github.com/multiversx/mx-sdk-go/tree/main/examples

## Instructions
- Load accounts from PEM files in the `accounts` directory.
- Issue 3 new tokens for each account.
- Transfer 100 million tokens to the issuer's account.
- Print the transaction hashes.
