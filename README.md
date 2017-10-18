## dory

Share your secret with a fish that have short term memory.

When you run `dory`, webserver by default will be available at [:8080](http://localhost:8080) and hosts [help document](http://localhost:8080/help) for quick overview.

To run it at any other port (or IP::PORT), export environment variable `DORY_HTTP` with required listen-at value before running `dory` in that env scope.

---

* How To Fetch Dependencies

```
./go-tasks deps
```


* How To Build

> following command will create a ./bin directory an create multi platform binaries there to run
> will also fetch dependencies if not done already

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
