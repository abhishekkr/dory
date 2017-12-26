
## Dory

> Share your secret with a fish that have short term memory.

* [HowTo Use ~ quick try-out with it's docker](./usage)

* [Development Workflow ~ run a local copy and try local changes](./development)


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

![image of dory](https://github.com/abhishekkr/dory/raw/master/w3assets/images/dory-1024px.jpg)

[MIT Licensed](https://github.com/abhishekkr/dory/blob/master/LICENSE)

---
