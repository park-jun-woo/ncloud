# NCloud SDK for Go

Welcome to the **NCloud SDK for Go**! This SDK provides a convenient way to interact with the NCloud API using the Go programming language.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Examples](#examples)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Features

- Simplified integration with NCloud services
- Lightweight and easy-to-use SDK
- Fully compatible with the Go programming language

## Installation

To install the NCloud SDK for Go, use the following `go get` command:

```bash
go get github.com/park-jun-woo/ncloud-sdk-go
```

Ensure that you have Go installed on your system. If not, download and install it from [golang.org](https://golang.org/).

## Getting Started

Follow these steps to get started with the NCloud SDK for Go:

1. Import the SDK in your Go application:

    ```go
    import "github.com/park-jun-woo/ncloud-sdk-go"
    ```

2. Initialize the SDK by providing your NCloud credentials:

    ```go
    config := ncloud.Config{
        AccessKey: "your-access-key",
        SecretKey: "your-secret-key",
    }
    client := ncloud.NewClient(config)
    ```

3. Start using the SDK to interact with NCloud services:

    ```go
    result, err := client.SomeService.SomeAction()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Println(result)
    ```

## Examples

### Example: Using GlobalDNS

```go
package main

import (
    "fmt"
    "log"

    "github.com/park-jun-woo/ncloud-sdk-go/services"
    "github.com/park-jun-woo/ncloud-sdk-go/services/Networking/GlobalDNS"
)

func main() {
    access := &services.Access{
        AccessKey: "your-access-key",
        SecretKey: "your-secret-key",
    }

    // Example: Retrieve or Create a Domain
    domainName := "example.com"
    domain, err := GlobalDNS.GetDomain(access, domainName, true)
    if err != nil {
        log.Fatalf("Failed to get or create domain: %v", err)
    }
    fmt.Printf("Domain: %v\n", domain)

    // Example: Set a DNS Record
    recordType := "A"
    recordContent := "192.0.2.1"
    recordTtl := 300
    domain, record, err := GlobalDNS.SetRecord(access, domainName, recordType, recordContent, recordTtl, true)
    if err != nil {
        log.Fatalf("Failed to set record: %v", err)
    }
    fmt.Printf("Record: %v\n", record)
}
```

## Documentation

For detailed API documentation and usage instructions, please refer to the [official documentation](https://github.com/park-jun-woo/ncloud-sdk-go/wiki).

## Contributing

We welcome contributions to improve this SDK! To contribute:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes with clear and descriptive messages.
4. Open a pull request to the `main` branch of this repository.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Thank you for using the **NCloud SDK for Go**! If you encounter any issues or have questions, feel free to open an issue in this repository.
