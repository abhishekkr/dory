## dory

Share your secret with a fish that have short term memory.

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

---

[developer's documentation](https://abhishekkr.github.io/dory/development)

---

![image of dory](w3assets/images/dory-1024px.jpg)

---

[MIT Licensed](./LICENSE)

---
