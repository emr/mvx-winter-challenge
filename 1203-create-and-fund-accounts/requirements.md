## Requirements
We need a script that generates three accounts in each shard on MultiversX chain and requests tokens from the faucet for them.

### Account creation
The script must generate exactly 3 unique accounts per shard. Total account count will be 9. 3 of them in shard0, 3 in shard1, 3 in shard2. After generating the wallet, you should check if it's in the desired shard, if not, create until it is and ignore the previous created ones.

Here's how to check an account's shard: https://blog.giantsvillage.com/how-to-generate-a-multiversx-wallet-in-a-specific-shard-tech-tutorial-458d565caed6

### Token request
Automate faucet request for the generated accounts via web-wallet.

### Usability
The script should be executed in a single command. Provide clear instruction on how to execute the script on README.md file.

## Resources
- Stick on official MultiversX documentation: https://docs.multiversx.com
- Check the source code of MultiversX Go SDK and examples:
  - https://github.com/multiversx/mx-sdk-go
  - https://pkg.go.dev/github.com/multiversx/mx-sdk-erdgo

## Instructions
- Create 3 wallets for each shard (shard0, shard1 and shard2) and generate their PEM files through the SDK and store them in a variable with their secret keys along with the PEM file path.
- For each wallet, do the following:
  - Check if the wallet has been created in the desired shard, if not, create until it is and ignore the previous created ones.
  - Open the page in the browser: https://testnet-wallet.multiversx.com/unlock/pem
  - Upload the PEM file to the file input in the page.
  - Click on "Access Wallet" button.
  - Navigate to https://testnet-wallet.multiversx.com/faucet
  - Wait for the "I'm not a robot" checkbox to be visible.
  - Wait for the "Request Tokens" button to be enabled and click on it when it is.
  - Close the browser tab.
  - Print success message.
  - Repeat the process for the rest of the wallets.
