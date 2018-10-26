# goyang
YANG parser and compiler for Go programs.

The yang package (pkg/yang) is used to convert a YANG schema into either an
in memory abstract syntax trees (ast) or more fully resolved, in memory, "Entry"
trees.  An Entry tree consists only of Entry structures and has had
augmentation, imports, and includes all applied.

goyang is a sample program that uses the yang (pkg/yang) package.

goyang uses the yang package to create an in-memory tree representation of
schemas defined in YANG and then dumps out the contents in several forms.
The forms include:

*  tree - a simple tree representation
*  proto - something "protobuf like" as a proof of concept
*  types - list understood types extracted from the schema

The yang package, and the goyang program, are not complete and are a work in
progress.

### Getting started

To build goyang, ensure you have go language tools installed
(available at [golang.org](golang.org/dl)) and that the `GOPATH`
environment variable is set to your Go workspace.

1. `go get github.com/openconfig/goyang`
    * This will download goyang code and dependencies into the src
subdirectory in your workspace.

2. `cd <workspace>/src/github.com/openconfig/goyang`

3. `go build`

   * This will build the goyang binary and place it in the bin
subdirectory in your workspace.

### Contributing to goyang

goyang is still a work-in-progress and we welcome contributions.  Please see
the `CONTRIBUTING` file for information about how to contribute to the codebase.

### Disclaimer

This is not an official Google product.
