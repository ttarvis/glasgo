# Glasgo Static Analysis Tool

A static analysis tool intended to check for potential security issues.  New tests will be added soon.  Special thanks to NCC Group Plc.

## Project

This is a static analysis tool written in Go for Go code.  It will find security and some correctness issues that may have a 
security implication.

## Compiling

To compile the tool, be sure to have the Go compiler first.

1. Use `Go build` for a local binary
2. Use `Go install` to compile and install in Go Path

## Using the tool

For now, all tests are run.

```
Glasgo directory1, directory2
```

or

```
Glasgo file1.go, file2.go
```

`Note:` The tool does not run on both directories and individual files

## Architecture

tbd

## Tests

* `error` - errors ignored
* `closer` - no file.Close() method called in function with file.Open()
* `insecureCrypto` - insecure cryptographic primitives
* `insecureRand` - insecurely generated random numbers
* `intToStr` - integer to string conversion without calling strconv
* `readAll` - ioutil.ReadAll called
* `textTemp` - checks if HTTP methods and template/text are in use
* `hardcoded` - looks for hardcoded credentials
* `bind` - checks if listener bound to all interfaces
* 'TLSConfig' - checks for insecure TLS configuration

## Design Choices

see the wiki

## Updates

Initial wave of tests have been uploaded and checked on test data

More tests to come

## to do

* add tests
* document tests
* document design choices

