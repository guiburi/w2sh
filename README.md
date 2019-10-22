# w2sh

## Table of Contents
+ [About](#about)
+ [Getting Started](#getting_started)
+ [Usage](#usage)
+ [Contributing](../CONTRIBUTING.md)

## About <a name = "about"></a>

Generate a web handler for any cli app.

### Installing

```
$ go get github.com/guiburi/w2sh
```
### Supported Providers

* github.com/spf13/cobra


### Examples

See the examples folder for a working samples.

## Usage <a name = "usage"></a>

```
http.HandleFunc("/", w2sh.Handle(cmd.RootCmd))
```
