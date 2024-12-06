## Requirements
Create a script that transfers 10000 units of each token from each account whose PEM files provided in the `accounts` directory to 1000 other accounts on MultiversX testnet.

### Usability
The script should be executed in a single command. Provide clear instruction on how to execute the script on README.md file.

## Resources
- Official MultiversX documentation: https://docs.multiversx.com
- Related part: https://docs.multiversx.com/tokens/fungible-tokens
- MultiversX Go SDK:
    - https://github.com/multiversx/mx-sdk-go
    - https://pkg.go.dev/github.com/multiversx/mx-sdk-erdgo
- SDK usage examples: https://github.com/multiversx/mx-sdk-go/tree/main/examples
- API documentation: https://testnet-api.multiversx.com

## Instructions
- Load accounts from PEM files in the `accounts` directory.
- For each account:
    - Get the identifiers of the tokens with name `WinterToken` owned by the account.
        ```
        curl https://testnet-api.multiversx.com/accounts/{address}/tokens?type=FungibleESDT&name=WinterToken
        ```
    - Generate random 1000 MultiversX addresses.
    - Transfer 10000 units of each token to the generated addresses.
- Print operation results in the console.
