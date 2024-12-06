## Requirements
Create a script that queries the MultiversX testnet blockchain to retrieve and display the list of transactions of each account defined constantly in memory.

### Blockchain querying
- Use public API https://testnet-api.multiversx.com
- Have a retry mechanism for failed requests.
- Use pagination with a page size of 200.

### Results
- Process each account concurrently.
- Fetch results in pages and save them into the file as they come.
- Create csv files directory under `results` directory for each account with their address as the name.

### Usability
The script should be executed in a single command. Provide clear instruction on how to execute the script on README.md file.

## Resources
- Official MultiversX documentation: https://docs.multiversx.com
- MultiversX Go SDK:
    - https://github.com/multiversx/mx-sdk-go
    - https://pkg.go.dev/github.com/multiversx/mx-sdk-erdgo
- SDK usage examples: https://github.com/multiversx/mx-sdk-go/tree/main/examples
- API documentation: https://testnet-api.multiversx.com

## Instructions
Help me implement @requirements.md in `1206-query-transactions` directory by sticking on @MultiversX, @MvX API Doc, and @MvX Go SDK API Ref as well as examples in @MvX Go SDK Source Code. Try to use best practices of idiomatic Go, use separation of concerns and respect SOLID principles while keeping the code simple. Don't overcomplicate it to apply every possible principle. Write code that is easily readable and testable. Write unit tests with only essential cases, don't include every possible scenario that will never happen in runtime. You won't complete everything by yourself. You have to ask me questions rather than making assumptions when something's not very clear on your side.
