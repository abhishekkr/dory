
> [home](./index) [develop](./development)

### HowTo Use

Since v1.1 Dory binary can be used in dual mode, server and client.

quick try using docker [abhishekkr/dory:1.2-alpine](https://hub.docker.com/r/abhishekkr/dory/)

#### Dory Server

```
docker run -it -p8080:8080 abhishekkr/dory:latest
```

When you run `dory`, webserver by default will be available at [:8080](http://localhost:8080) and hosts [help document](http://localhost:8080/help) for quick overview.

* While running Dory server, you can get it's port change by configuring environment variable `DORY_HTTP`. Setting it's value to `:9999` will make it listen at port 9999; value `127.0.0.1:8080` will force it to listen at only `127.0.0.1` on port `8080`.

* Memberlist/Gossip binding/advertising will by default bind to port `7946`. They can be overridden by environment variables like `DORY_MEMBERS_BIND=7947` and `DORY_MEMBERS_ADVERTISE=7947`.

* Dory server provides admin api accessible by HTTP Header `X-DORY-ADMIN-TOKEN: <admin-token>`.
> This `admin-token` need to be configured as environment variable `DORY_ADMIN_TOKEN` with value of more than 256 characters. Else it is not usable.
>
> Admin tasks of List and Purge of all keys.
>
> Just a reminder values are not decipherable using this,as that's only possible using tokens... though they can be deleted in encrypted state itself.


**Sample docker run command with http port changed to 8000 and admin-token configured**

```
docker run \
  -e DORY_HTTP=:8080 \
  -e DORY_ADMIN_TOKEN=some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  -it -p8000:8080 \
  abhishekkr/dory:latest
```

---
---

#### Dory Client

* Ping is kind of health-check for Dory server, returns count of keys in Cache and Disk,
> it's default task as well so no need to provide that flag

```
dory-linux-amd64 -mode client \
  -task ping \
  -url http://127.0.0.1:8080
```

> above `ping` is equivalent to curl call below

```
curl -skL -X GET 'http://127.0.0.1:8080/ping'
```

---


**Interacting with Cache (Store in Memory with Data Expiry)**

* to publish a secret on Dory with default expiry of 300seconds
> below would return a token, that is used as value in later, assume returned token is _BmwkoOsB6KeMIUbqW8BC0u5vfDgdsr06_

```
dory-linux-amd64 -mode client \
  -task set \
  -url http://127.0.0.1:8080 \
  -key some-pass -val "what's in the name"
```

> above `set secret to cache` is equivalent to curl call below

```
curl -skL -X POST \
  --data "what's in the name"
  'http://127.0.0.1:8080/local-cache/some-pass'
```

---

* to publish a secret with expiry of 1hour

```
dory-linux-amd64 -mode client \
  -task set \
  -url http://127.0.0.1:8080 \
  -ttl 3600 \
  -key some-pass -val "what's in the name"
```

> above `set secret to cache with ttl of 1hr` is equivalent to curl call below

```
curl -skL -X POST \
  --data "what's in the name"
  'http://127.0.0.1:8080/local-cache/some-pass?ttl=3600'
```

---

* to publish a secret with data to be read from a file, as it's a blob (like private-key, credential files, image, anything)

```
dory-linux-amd64 -mode client \
  -task set \
  -url http://127.0.0.1:8080 \
  -key some-pass -val-from secret-store.log
```

> above `set secret to cache from file blob` is equivalent to curl call below

```
curl -skL -X POST \
  --data @secret-store.log
  'http://127.0.0.1:8080/local-cache/some-pass?ttl=3600'
```

---


* fetch secret published to cache using key using it's token returned while publish

```
dory-linux-amd64 -mode client \
  -task get \
  -url http://127.0.0.1:8080 \
  -key what -token BmwkoOsB6KeMIUbqW8BC0u5vfDgdsr06
```

> above `get secret from cache using token` is equivalent to curl call below

```
curl -skL -X GET \
  -H 'X-DORY-TOKEN: BmwkoOsB6KeMIUbqW8BC0u5vfDgdsr06' \
  'http://127.0.0.1:8080/local-cache/some-pass'
```

---


* delete secret published to a key using it's token returned while publish

```
dory-linux-amd64 -mode client \
  -task del \
  -url http://127.0.0.1:8080 \
  -key what -token BmwkoOsB6KeMIUbqW8BC0u5vfDgdsr06
```

> above `delete secret from cache using token` is equivalent to curl call below

```
curl -skL -X DELETE \
  -H 'X-DORY-TOKEN: BmwkoOsB6KeMIUbqW8BC0u5vfDgdsr06' \
  'http://127.0.0.1:8080/local-cache/some-pass'
```

---


**Interacting with non-expiry disk stored secrets**

> same commands as above would work, require just an extra flag '-persist true'

* example like `set` a key

```
dory-linux-amd64 -mode client \
  -persist true \
  -task set \
  -url http://127.0.0.1:8080 \
  -key some-pass -val "what's in the name"
```

> above `set secret on disk` is equivalent to curl call below, call at api `/local-disk` instead of `/local-cache`

```
curl -skL -X GET \
  --data "what's in the name"
  'http://127.0.0.1:8080/local-disk/some-pass'
```

---


**Interacting with Admin API**

* get list of all current keys in cache

```
## from cache
dory-linux-amd64 -mode client \
  -task list \
  -url http://127.0.0.1:8080 \
  -token some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

## from disk
dory-linux-amd64 -mode client \
  -persist true \
  -task list \
  -url http://127.0.0.1:8080 \
  -token some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

> above `get list of keys using admin-token` is equivalent to curl call below

```
## from cache
curl -skL -X GET \
  -H 'X-DORY-ADMIN-TOKEN: some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  'http://127.0.0.1:8080/admin/store/cache'

## from disk
curl -skL -X GET \
  -H 'X-DORY-ADMIN-TOKEN: some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  'http://127.0.0.1:8080/admin/store/disk'
```

---


* purge all current keys

```
## from cache
dory-linux-amd64 -mode client \
  -task purge \
  -url http://127.0.0.1:8080 \
  -token some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

## from disk
dory-linux-amd64 -mode client \
  -persist true \
  -task purge \
  -url http://127.0.0.1:8080 \
  -token some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

> above `purge all keys using admin-token` is equivalent to curl call below

```
## from cache
curl -skL -X DELETE \
  -H 'X-DORY-ADMIN-TOKEN: some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  'http://127.0.0.1:8080/admin/store/cache'

## from disk
curl -skL -X DELETE \
  -H 'X-DORY-ADMIN-TOKEN: some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  'http://127.0.0.1:8080/admin/store/disk'
```

---


* purge one key

```
## from cache
dory-linux-amd64 -mode client \
  -task purge \
  -url http://127.0.0.1:8080 \
  -token some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

## from disk
dory-linux-amd64 -mode client \
  -persist true \
  -task purge \
  -url http://127.0.0.1:8080 \
  -token some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

> above `purge all keys using admin-token` is equivalent to curl call below

```
## from cache
curl -skL -X DELETE \
  -H 'X-DORY-ADMIN-TOKEN: some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  'http://127.0.0.1:8080/admin/store/cache/some-pass'

## from disk
curl -skL -X DELETE \
  -H 'X-DORY-ADMIN-TOKEN: some-token-more-than-256-chars-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' \
  'http://127.0.0.1:8080/admin/store/disk/some-pass'
```


---
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
    	the kind of action dory client need to perform, supports {set,get,del,list,purge,purge-one,ping}; defaults ping (default "ping")
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
---
