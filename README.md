# What 

A little webapp that I made to read news from 'thehindu.com' by parsing RSS feed and scraping the articles.

This was made as a part of learning htmx and templ.


## Dependencies

Go, make, bunx, Ollama

## Installation

Uses tailwindcss for styling. I use bun ( bunx ) for tailwindcss. Feel free to change it in Makefile if you'd like to use npm or yarn.

Needs Ollama for news summarization feature to work.

To install dependencies run:

```shell
    make install
```

To run the app:

```shell
     make run
```

It tries to connect to Ollama at `http://localhost:111434` by default. You can change it by setting `OLLAMA_HOST` environment variable.

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
- Get rid of over-complicated directory structure


## Known Issues

- The theme switcher stops working after you navigate around. You have reload the page to make it work again. It maybe has to do with htmx and partial responses. I'll look into it later.
- Sometimes there's a double navigation bar. I think it's because of the way I'm using htmx. TFL
- The error pages need to be fixed as navigatging back makes navbar dissapear since they are just static html pages.


