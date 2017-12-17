
### Development Workflow

[![Go Report Card](https://goreportcard.com/badge/abhishekkr/dory)](https://goreportcard.com/report/abhishekkr/dory) [![Build Status](https://travis-ci.org/abhishekkr/dory.svg?branch=master)](https://travis-ci.org/abhishekkr/dory)

* How To Fetch Dependencies

```
./go-tasks deps
```


* How To Build

> following command will create a ./bin directory an create multi platform binaries there to run
> will also fetch dependencies if not done already
>
> default it runs at `:8080`, to run it at any other port (or IP::PORT), export environment variable `DORY_HTTP` with required listen-at value before running `dory` in that env scope.

```
./go-tasks build
```


* How To Use LocalAuth

> [curl example](w3assets/dory.sh)

---

building docker with custom local binary

```
docker build  -t abhishekkr/dory:alpha -f w3assets/Dockerfile .
```


---
