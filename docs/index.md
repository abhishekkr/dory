
## Dory

> Share your secret with a fish that have short term memory.

* [HowTo Use ~ quick try-out with it's docker](/usage)

* [Development Workflow ~ run a local copy and try local changes](/development)

#### Current Features:

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

![image of dory](w3assets/images/dory-1024px.jpg)

[MIT Licensed](./LICENSE)

---
