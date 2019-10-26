# w2sh

## Table of Contents
+ [About](#about)
+ [Getting Started](#getting_started)
+ [Usage](#usage)
+ [TODO](#todo)

## About <a name = "about"></a>

browse your cli app.

## Getting Started <a name = "getting_started"></a>
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

## TODO <a name = "todo"></a>

---

- [x] Create TODO doc
- [x] Fix page template reload issue
- [x] submit form
- [x] navigate subcommads
- [ ] submit subcommads
- [ ] recursive collect cmds
- [ ] cli frameworks as providers
- [ ] template update
- [ ] submit files
- [ ] intergrate urfave cli
