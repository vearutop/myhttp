# myhttp

[![Build Status](https://github.com/vearutop/myhttp/workflows/test-unit/badge.svg)](https://github.com/vearutop/myhttp/actions?query=branch%3Amaster+workflow%3Atest-unit)

Hashing HTTP fetcher. This is a demo project. See also [`plt`](https://github.com/vearutop/plt).

This repo has a depencency on `github.com/bool64/dev` for GitHub CI and Makefile features. 
If you would like to see version without any 3rd party dependencies, please check [`minimal`](https://github.com/vearutop/myhttp/tree/minimal) branch. 

## Installation

```
go get github.com/vearutop/myhttp
```

Or download binary from [releases](https://github.com/vearutop/myhttp/releases).

## Usage

Provide a list of URLs as arguments. Protocol prefix can be omitted and defaults to `http://`.

```
myhttp example.com https://google.com http://127.0.0.1:1234/foo
```

Resources are fetched concurrently with a default limit of 10 simultaneous requests. Limit can be controlled
with `-parallel` flag:

```
myhttp -parallel 2 example.com https://google.com http://127.0.0.1:1234/foo example2.com example3.com
```

### Flags

```
Usage of myhttp:
  -parallel int
        maximum number of concurrent requests (default 10)
```
