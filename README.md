## dory

[![Go Report Card](https://goreportcard.com/badge/abhishekkr/dory)](https://goreportcard.com/report/abhishekkr/dory) [![Build Status](https://travis-ci.org/abhishekkr/dory.svg?branch=master)](https://travis-ci.org/abhishekkr/dory)

quick try using docker [abhishekkr/dory:1.0-alpine](https://hub.docker.com/r/abhishekkr/dory/)

```
docker run -it -p8080:8080 abhishekkr/dory:latest
```

Share your secret with a fish that have short term memory.

When you run `dory`, webserver by default will be available at [:8080](http://localhost:8080) and hosts [help document](http://localhost:8080/help) for quick overview.


Current Features:

> * added `/ping` api listing count of keys in cache and disk
>
> * blind list and purge of keys over `/admin/store/{cache,disk}` api using Admin Token; `just a reminder that value is not recoverable using Admin Token`
>
> * can also reach permanent (disk persisted, non-ttl) memories in Dory by providing `GET Param` as `persist=true`
>
> * provides local auth-store backed by `cache2go` and encrypted by AES crypto
>
> * allows POST, GET and DELETE for a `auth identifier path` (like `http://dory.local:8080/local-auth/:identifier`) to store, fetch and purge data
>
> * successful POST of data at `auth identifier path` returns reference `X-DORY-TOKEN` mapped with this `auth identifier path`, this token need to be sent as value of this header `X-DORY-TOKEN` for GET and DELETE.
>
> * created secret store have default TTL of 5minutes, custom TTL can be set as URL Param in POST request by value of `ttlsecond` in seconds
>
> * first GET of secret will purge it from store, unless GET Param `keep=true` is provided

---

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

![image of dory](w3assets/images/dory-1024px.jpg)

---

[MIT Licensed](./LICENSE)

---
