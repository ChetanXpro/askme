# askme CLI Tool

`askme` is a command-line interface (CLI) tool designed to manage and interact with a vector database. It enables the setup of environment secrets and performs operations such as adding PDFs to a vector DB for similarity searches and executing similarity searches like a large language model (LLM).

## Installation

You can install `askme` using the `go install` command:

```bash
go install github.com/chetanxpro/askme@latest
```

The binary will be installed in your $GOPATH/bin directory or $GOBIN if set.

Make sure to replace yourusername/askme with the actual path to your repository.


## Using go get (for versions before Go 1.16)

If you are using an older version of Go that does not support the go install command with versioning, you can use go get:
```
go get -u github.com/yourusername/askme
```

This will get the latest version of askme and install it.

Note that starting from Go 1.17, using go get to install executables is deprecated. go install is the recommended way to install binaries outside a module.

## Commands

askme comes with a set of commands that allow you to interact with your vector database in various ways:

- `askme setup`: Set up secrets like OpenAI API key, Pinecone API key, etc.
- `askme add`: Add a PDF into the vector DB to perform similarity searches later.
- `askme ask`: Perform a similarity search, akin to querying an LLM.


## Usage
