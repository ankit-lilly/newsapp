A simple webapp using HTMX and Go that zips through RSS feeds, serving up fresh news without those annoying ads or pesky paywalls. 


Needs Ollama for news summarization feature to work.

To install dependencies run:

```shell
    make install
```

To run in development mode:

```shell
    OLLAMA_HOST="http://localhost:111434' make run
```

To build the executable:

```shell
    make build

    OLLAMA_HOST="http://localhost:111434' ./newsapp
```



## TODO

- Fix error handling with better error messages
- Get rid of over-complicated ridiculous directory structure



