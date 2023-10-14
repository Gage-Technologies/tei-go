# TEI (Text Embeddings Inference) Go Client

![Go Report Card](https://goreportcard.com/badge/github.com/Gage-Technologies/tei-go)
![License](https://img.shields.io/github/license/Gage-Technologies/tei-go)

## Overview

This client library is designed for Go developers who want to interact with the Text Embeddings Inference server. It handles the underlying HTTP requests and offers a simple API for embedding text.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Features](#features)
  - [Timeouts](#timeouts)
  - [Custom Headers](#custom-headers)
  - [Cookies](#cookies)

## Installation

To install the package, you can run:

```
go get github.com/gage-technologies/tei-go
```

## Usage

Here's a simple example to demonstrate how to use this library:

```go
import "github.com/gage-technologies/tei-go/tei"

func main() {
	client := tei.NewClient("http://localhost:8080", nil, nil, time.Second*30)
	res, err := client.Embed("Hi there!", false)
	if err != nil {
		panic(err)
	}
	fmt.Println("Embedding: ", res[0])
}
```

## Features

### Timeouts

The client allows you to specify a timeout for HTTP requests.

### Custom Headers

You can specify custom headers when creating a new client:

```go
headers := map[string]string{"Authorization": "Bearer token"}
client := tei.NewClient("http://localhost:8080", headers, nil, time.Second*30)
```

### Cookies

You can also specify cookies when creating a new client:

```go
cookies := map[string]string{"session_id": "abc123"}
client := tei.NewClient("http://localhost:8080", nil, cookies, time.Second*30)
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details
