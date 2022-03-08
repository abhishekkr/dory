## dory

Share your secret with a fish that have short term memory.

> * this is a secret sharing service for masses, where you don't need to be authenticated at service to store and share secret
>
> * anyone with access to service can upload a secret and share the token with people they wanna share it
>
> * if accessed without an explicit retention parameter, the secret gets purged on first fetch
>
> * if stored in `cache` mode, it self expires after a TTL if not accessed for that duration
>
> * even service admin can't decipher a secret posted by any user

[![Go Report Card](https://goreportcard.com/badge/abhishekkr/dory)](https://goreportcard.com/report/abhishekkr/dory) [![Build Status](https://travis-ci.org/abhishekkr/dory.svg?branch=master)](https://travis-ci.org/abhishekkr/dory)

quick try using docker [abhishekkr/dory:1.2-alpine](https://hub.docker.com/r/abhishekkr/dory/)

[quick usage guide for dory server and client](https://abhishekkr.github.io/dory/usage)

When you run `dory`, webserver by default will be available at [:8080](http://localhost:8080) and hosts [help document](http://localhost:8080/help) for quick overview.


Current Features:

* local-auth has 2 kind of datastores, non-persistent `cache` and persistent `disk`

> * both of these stores purge their entry on any fetch by default, unless asked to `keep=true` it further
>
> * `cache` store has extra expiry attached to it, for keys to self-delete if not accessed for TTL (default 300seconds)
>
> * `disk` store persists secrets which are resumable even after service restart
>
> * all secrets are stored post AES encryption with per secret unique token
>
> * every secret can only be deciphered by it's token which is returned as response of posting the secret, so if token is lost... so is data in secret
>
> * token is also required to delete the secret by normal users


* admin tasks

> * listing of all the keys against which secrets are stored
>
> * purge one key, if by mistake wrong or undesired secret has been shared
>
> * purge all keys, for cleanup or in times of threat
>
> `just a reminder that value is not recoverable using Admin Token`


* `/ping` api listing count of keys in cache and disk


### Using

* Managing simple secrets

```
## adding expirable secret file, life of 180sec if not read
DORY_KEY=$(curl -skL -v -X POST -H "Content-Type: multipart/form-data" \
        -d@secret.json  \
        "http://127.0.0.1:8080/local-auth/mysecret?ttlsecond=180")

## fetching secret, which expires on read
curl -skL -v --request GET -o secret.json  \
        --header "X-DORY-TOKEN: ${DORY_KEY}" \
        "http://127.0.0.1:8080/local-auth/mysecret"
```


* Support for big files which need multipart form-data

```
SECRET_FILE="mysecret.store"

## adding expirable secret file, life of 300sec if not read
DORY_KEY=$(curl -skL -v -X POST -H "Content-Type: multipart/form-data" \
        -F "form=@${SECRET_FILE}" \
        "http://127.0.0.1:8080/local-auth/${SECRET_FILE}?ttlsecond=300&file-field=form")

## fetching secret, which expires on read
curl -skL -v --request GET -o secret.store  \
        --header "X-DORY-TOKEN: ${DORY_KEY}" \
        "http://127.0.0.1:8080/local-auth/${SECRET_FILE}"
```

---

[developer's documentation](https://abhishekkr.github.io/dory/development)

---

![image of dory](w3assets/images/dory-1024px.jpg)

---

[MIT Licensed](./LICENSE)

---
