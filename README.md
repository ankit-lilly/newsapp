Playing around with HTMX and Templ by building a webapp that lets you read news by parsing RSS feed
and scraping the articles.

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



