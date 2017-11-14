

v0.3

* can also reach permanent (disk persisted, non-ttl) memories in Dory by providing `GET Param` as `persist=true`, gets saved using `diskv`


v0.2

* moved CacheTable based local-auth to DataStore interface, so other store backends can be added


v0.1

* provides local auth-store backed by `cache2go` and encrypted by AES crypto

* allows POST, GET and DELETE for a `auth identifier path` (like `http://dory.local:8080/local-auth/:identifier`) to store, fetch and purge data

* successful POST of data at `auth identifier path` returns reference `X-DORY-TOKEN` mapped with this `auth identifier path`, this token need to be sent as value of this header `X-DORY-TOKEN` for GET and DELETE.

* created secret store have default TTL of 5minutes, custom TTL can be set as URL Param in POST request by value of `ttlsecond` in seconds

* first GET of secret will purge it from store, unless GET Param `keep=true` is provided
