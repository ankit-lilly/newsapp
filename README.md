# NewsMaxxing

## What

Playing around with HTMX, templ, Echo v5, and Go.

This is a simple news reader that parses the RSS feeds from various portals and
then allows reader to read the news by scraping the content from the original
news site.

Uses Ollama for summary and chat. Chat allows user to ask questions about the
contents of the article.

## Running

The only local runtime requirement is Go. The install step downloads the
standalone Tailwind CSS executable, DaisyUI plugin, and streaming-markdown
browser module; Bun, Node.js, and npm are not required.

Templ is pinned in `go.mod` as a tool dependency and runs through
`go tool templ`, so no global Templ CLI installation is needed.

```bash
make install
make build

./newsmaxxing
```

For development with Go and CSS hot reload:

```bash
make run
```

It tries to connect to Ollama at `http://localhost:11434`.

If you have a different host for Ollama then you can set the environment
variable `OLLAMA_HOST` to the correct host.

```bash
OLLAMA_HOST=http://yourhost:11434 ./newsmaxxing
```

Otherwise you could just download prebuilt binary from the releases.

If you are on Mac then it'll complain about it being from an unidentified
developer and ask you to move it to trash.

The workaround is to run the following command in the terminal:

```bash
xattr -d com.apple.quarantine /path/to/binary
```

Or you can just right click on the binary and then click on `Open` and then
click on `Open` again.

Otherwise you can go to `System Preferences -> Security & Privacy -> General`
and then click on `Open Anyway`.

## Env variables


The app uses environment variables to for configuration. You can take a look at sample
`.env` file in the root directory.

Note:

You may get too many requests error from time to time due to rate limit on the
various websites. You'll just have wait or click on some other menu.

It's using pre-load to load everything at once and then cache it. So, it should
be faster after the first load.
