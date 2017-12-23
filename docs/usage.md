
> [home](./index) [develop](./development)

### HowTo Use

Since v1.1 Dory binary can be used in dual mode, server and client.

quick try using docker [abhishekkr/dory:0.3-alpine](https://hub.docker.com/r/abhishekkr/dory/)

#### Dory Server

```
docker run -it -p8080:8080 abhishekkr/dory:latest
```

When you run `dory`, webserver by default will be available at [:8080](http://localhost:8080) and hosts [help document](http://localhost:8080/help) for quick overview.

* While running Dory server, you can get it's port change by configuring environment variable `DORY_HTTP`. Setting it's value to `:9999` will make it listen at port 9999; value `127.0.0.1:8000` will force it to listen at only `127.0.0.1` on port `8000`.

* Dory server provides admin api accessible by HTTP Header `X-DORY-ADMIN-TOKEN: <admin-token>`.
> This `admin-token` need to be configured as environment variable `DORY_ADMIN_TOKEN` with value of more than 256 characters. Else it is not usable.
> admin tasks of List and Purge of all keys
> just a reminder values are not decipherable using this,as that's only possible using tokens... though they can be deleted in encrypted state itself


**Sample docker run command with http port changed to 8000 and admin-token configured**

```
docker run \
  -e DORY_HTTP=:8000 \
  -e DORY_ADMIN_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  -it -p8000:8000 \
  abhishekkr/dory:latest
```

---

#### Dory Client

* Ping is kind of health-check for Dory server, returns count of keys in Cache and Disk,
> it's default task as well so no need to provide that flag

```
dory-linux-amd64 -mode client -task ping -url http://127.0.0.1:8080
```


**Interacting with Cache (Store in Memory with Data Expiry)**

* to publish a secret on Dory with default expiry of 300seconds
> below would return a token, that is used as value in later, assume returned token is `BmwkoOsB6KeMIUbqW8BC0u5vfDgdsr06`

```
dory-linux-amd64 -mode client -task set -url http://127.0.0.1:8080 -key some-pass -val "what's in the name"

```

* to publish a secret with expiry of 1hour

```
dory-linux-amd64 -mode client -task set -url http://127.0.0.1:8080 -ttl 3600 -key some-pass -val "what's in the name"
```

* to publish a secret with data to be read from a file, as it's a blob (like private-key, credential files, image, anything)

```
dory-linux-amd64 -mode client -task set -url http://127.0.0.1:8080 -ttl 3600 -key some-pass -val-from secret-store.log
```

* fetch secret published to cache using key using it's token returned while publish

```
dory-linux-amd64 -mode client -task get -url http://127.0.0.1:8080 -key what -token BmwkoOsB6KeMIUbqW8BC0u5vfDgdsr06
```

* delete secret published to a key using it's token returned while publish

```
dory-linux-amd64 -mode client -task del -url http://127.0.0.1:8080 -key what -token BmwkoOsB6KeMIUbqW8BC0u5vfDgdsr06
```


**Interacting with non-expiry disk stored secrets**

> same commands as above would work, require just an extra flag '-persist true'

* example like `set` a key

```
dory-linux-amd64 -mode client -task set -url http://127.0.0.1:8080 -persist true -key some-pass -val "what's in the name"
```


**Interacting with Admin API**

* get list of all current keys in cache

```
## from cache
dory-linux-amd64 -mode client -url http://127.0.0.1:8080 -task list \
  -token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

## from disk
dory-linux-amd64 -mode client -url http://127.0.0.1:8080 -task list -persist true \
  -token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

* purge all current keys

```
## from cache
dory-linux-amd64 -mode client -url http://127.0.0.1:8080 -task purge \
  -token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

## from disk
dory-linux-amd64 -mode client -url http://127.0.0.1:8080 -task purge -persist true \
  -token xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

---

#### Current supported flags

```
dory-linux-amd64:
  -keep
    	to decide if to purge key post read or not, defaults as false for purge on read
  -key string
    	key name to be provided to dory
  -mode string
    	run mode, allowed modes are client and server, defaults server (default "server")
  -persist
    	to decide datastore as cache or disk, defaults as false for cache
  -task string
    	the kind of action dory client need to perform, supports {set,get,del,list,purge,ping}; defaults ping (default "ping")
  -token string
    	token for secret, required when trying to Get or Delete a key
  -ttl int
    	ttl for key, if it's set task for cache datastore; defaults 300 sec (default 300)
  -url string
    	url for dory server to be talked to
  -val string
    	value to be provided to dory, required when trying to Post or Delete a key
  -val-from string
    	value from a file to be provided to dory, required when trying to Post or Delete a key

```

---
