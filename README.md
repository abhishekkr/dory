## dory

Share your secret with a fish that have short term memory.

When you run `dory`, webserver by default will be available at [:8080](http://localhost:8080) and hosts [help document](http://localhost:8080/help) for quick overview.


Current Features:

* provides local auth-store backed by `cache2go` and encrypted by AES crypto

* allows POST, GET and DELETE for a `auth identifier path` (like `http://dory.local:8080/local-auth/:identifier`) to store, fetch and purge data

* successful POST of data at `auth identifier path` returns reference `X-DORY-TOKEN` mapped with this `auth identifier path`, this token need to be sent as value of this header `X-DORY-TOKEN` for GET and DELETE.

* created secret store have default TTL of 5minutes, custom TTL can be set as URL Param in POST request by value of `ttlsecond` in seconds

* first GET of secret will purge it from store, unless GET Param `keep=true` is provided

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

GO_MAIN_FILE=dory.go ./go-tasks build
```


* How To Use LocalAuth

> [curl example](w3assets/dory.sh)

---

![image of dory](w3assets/images/dory-1024px.jpg)

---

[MIT Licensed](./LICENSE)

---
