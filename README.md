# askme 

askme is a command-line interface (CLI) tool designed to manage and interact with PDF data. It can perform operations such as adding PDFs to a vector DB for similarity searches.

## Installation

### Using Precompiled Binaries

You can download the precompiled binaries directly from the latest release on GitHub.

For example, to download the binary for Linux amd64:

```
curl -LO https://github.com/chetanxpro/askme/releases/download/v0.0.1/askme-linux-amd64
```

Make the binary executable:
```
chmod +x askme-linux-amd64
```
Move the binary to your system's bin folder:
```
mv askme-linux-amd64 /usr/local/bin/askme
```
Make sure to use the correct URL for your release version and desired platform.

Note
- Make sure to use the correct URL for your release version and the desired platform

### Using go install

You can install `askme` using the `go install` command:

```bash
go install github.com/chetanxpro/askme@latest
```

The binary will be installed in your $GOPATH/bin directory or $GOBIN if set.



## Using go get (for versions before Go 1.16)

If you are using an older version of Go that does not support the go install command with versioning, you can use go get:
```
go get -u github.com/chetanxpro/askme
```

This will get the latest version of askme and install it.

Note that starting from Go 1.17, using go get to install executables is deprecated. go install is the recommended way to install binaries outside a module.

## Commands

askme comes with a set of commands that allow you to interact with your vector database in various ways:

- `askme setup`: Set up secrets like OpenAI API key, Pinecone API key, etc.
- `askme add`: Add a PDF into the vector DB to perform similarity searches later.
- `askme ask`: Perform a similarity search, and you will get better result with openai .


## Usage

### Setting Up Secrets
```askme setup```

Follow the on-screen prompts to enter your API keys and other necessary configuration details.

### Adding a PDF to Vector DB

```askme add```

This will prompt the user to provided PDF document path and then upsert it into vectorDB and making it available for similarity searches.

### Performing a Similarity Search

```askme ask```

This will prompt the user to decide for which pdf user want to search query to find similar documents within your vector database.


## Feature To-Do List

Here are some of the features that are planned for future releases:

- [ ] Add Qdrant DB support.
- [ ] Improve error handling and logging.
- [ ] Create a user-friendly setup wizard.
- [ ] Add support for llama 2 and some other llm


Feel free to suggest new features by [opening an issue](https://github.com/chetanxpro/askme/issues) on GitHub.


# License

askme is made available under the MIT License. See the LICENSE file for more information.
