<h1 align="center">mcp-x</h1>

<p align="center">
  <b>An MCP server that gives an LLM the X (Twitter) API: read and search posts, manage users and lists, upload media, and publish — as the account whose keys it runs on.</b>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-8bc34a?style=for-the-badge" alt="License MIT"></a>
  <img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/MCP-Server-6E56CF?style=for-the-badge&logo=anthropic&logoColor=white" alt="MCP">
  <img src="https://img.shields.io/badge/Transport-stdio_%7C_http-26A5E4?style=for-the-badge" alt="Transport">
  <img src="https://img.shields.io/badge/X_API-v2-000000?style=for-the-badge&logo=x&logoColor=white" alt="X API v2">
  <img src="https://img.shields.io/badge/Auth-OAuth_1.0a-eab308?style=for-the-badge" alt="OAuth 1.0a">
  <img src="https://img.shields.io/badge/Paid_API-pay--per--use-e11d48?style=for-the-badge" alt="Paid API">
</p>

<p align="center">
  <a href="https://varmcp.com/mcp-x">Website</a> ·
  <a href="#the-x-api-costs-money">Costs</a> ·
  <a href="#credentials">Credentials</a> ·
  <a href="#tools">Tools</a> ·
  <a href="#quick-start">Quick start</a> ·
  <a href="#configuration">Configuration</a> ·
  <a href="#architecture">Architecture</a> ·
  <a href="CONTRIBUTING.md">Contributing</a>
</p>

---

## What it is

`mcp-x` is a [Model Context Protocol](https://modelcontextprotocol.io) server written in Go. It exposes the **X API v2** to any MCP-compatible client (Claude Desktop, IDE agents, custom LLM apps) as 42 tools covering posts, users, lists and media.

It authenticates with **OAuth 1.0a user context**, which means every call acts *as a real X account* — the one the four keys belong to. `x_post_create` publishes publicly. `x_user_follow` really follows. `x_post_delete` is irreversible. This is not a sandbox, and it is not free: see [The X API costs money](#the-x-api-costs-money) before you wire it into an agent.

Both transports the MCP SDK supports are available and expose the identical tool set:

- **stdio** — the client launches the binary and talks over stdin/stdout (the default, ideal for desktop clients).
- **http** — a long-running streamable HTTP server (useful for remote/shared deployments).

---

## The X API costs money

> [!WARNING]
> **There is no free tier any more.** X retired the Free/Basic/Pro subscription tiers for new developers and moved to **pay-per-use credits**: you buy credits upfront in the [Developer Console](https://developer.x.com) and every request deducts from the balance in real time. Legacy Basic ($200/mo) and Pro ($5,000/mo) subscriptions survive only for accounts that already had them; Enterprise starts around $42,000/mo. A new developer account today gets pay-per-use and nothing else.

Rates at the time of writing ([official pricing](https://docs.x.com/x-api/getting-started/pricing) — always re-check the Console, they have changed several times in 2026):

| Operation | Price |
| :--- | :--- |
| Post read | **$0.005** per post returned |
| Owned read (your own posts, bookmarks, followers, likes, lists) | **$0.001** per resource |
| User read | **$0.010** per user returned |
| Likes / mutes / blocks read | **$0.001** per resource |
| Followers / following read | **$0.010** per resource |
| **Publishing a post** | **$0.015** per request |
| **Publishing a post containing a URL** | **$0.200** per request |
| Like / repost and other interactions | **$0.015** per request |
| List and bookmark writes | $0.005–$0.010 per request |

Two things follow from this, and both are baked into the server:

- **Reads are charged per *resource returned*, not per request.** `max_results: 100` on `x_posts_search` costs twenty times what `max_results: 5` costs for the same query. Every read tool's description tells the model to ask for the smallest `max_results` that answers the question, and the batching tools (`x_posts_lookup`, `x_users_lookup`) tell it to batch rather than loop.
- **`x_posts_count` does not consume the post-read budget.** It returns match counts bucketed by minute/hour/day for the same query syntax. Size a topic with `x_posts_count` first, then pay for `x_posts_search`.

X also **deduplicates**: the same resource fetched twice inside a 24-hour UTC window is billed once. And pay-per-use is capped at 3 million post reads per billing cycle — past that, only Enterprise.

When the money runs out the API answers with a distinct error, and the server maps it to a message that explicitly tells the model **not to retry** — see [Errors](#errors).

---

## Credentials

The server needs **four** OAuth 1.0a values, all from one X app:

```
X_API_KEY
X_API_KEY_SECRET
X_ACCESS_TOKEN
X_ACCESS_TOKEN_SECRET
```

The first pair identifies the *app*; the second pair identifies the *account acting through it*. `mcp-x` uses `AuthenMethodOAuth1UserContext`, not app-only bearer auth, because every write endpoint and every "me" endpoint (`x_users_me`, `x_posts_home`, `x_bookmarks_list`, mentions) requires a user context. There is no bearer-token mode.

### Getting them — and the one step everybody gets wrong

1. Go to [developer.x.com](https://developer.x.com) → your project → your app.
2. Open **User authentication settings** and set **App permissions** to **Read and Write**. Do this **first**.
3. *Only then* go to **Keys and tokens** and generate the **Access Token and Secret**.

> [!IMPORTANT]
> **An access token permanently keeps the permissions the app had at the moment it was generated.** If you created the token while the app was Read-only and then flipped the app to Read and Write, the token is *still* read-only. Nothing about the app settings page will tell you this. Every write will fail with X's `oauth1-permissions` problem type, forever, until you go back to **Keys and tokens** and **regenerate** the Access Token and Secret.
>
> This is the single most common setup failure with the X API, which is why the server checks for it at startup and refuses to start with:
>
> ```
> the access token is read-only: set the app permissions to Read and Write
> in the X Developer Console, then regenerate the Access Token and Secret
> ```
>
> Regenerating the *API Key/Secret* is not the fix. Regenerate the **Access Token and Secret**.

### Startup verification

Before registering a single tool, the server calls `GET /2/users/me` once (`client.Bootstrap`) with a 15-second deadline. This does three jobs:

- proves the four keys are valid — bad keys fail the process, not the first tool call;
- catches the read-only-token trap above;
- caches the key owner's numeric user id, because every write endpoint is `POST /2/users/:id/...` and looking the id up per write would be another billed request.

A failure here is fatal by design. A server that starts and then fails every call is worse than one that does not start.

The cost of that choice is worth stating plainly: **there is no way to try this server without a funded X developer account.** No credentials means no startup, which means no tool list — `/mcp` and `claude mcp list` will show a failed connection and nothing else. To confirm an install short of that, run the binary with `-version`, and read the [Tools](#tools) section for what it would have exposed.

### Keys are secrets

The four values are credentials for a live account with write access. Keep them in a file only you can read, pass it with `-env`, and never commit it — `.env` is gitignored, `.env.example` is the template. When running under an MCP client, the client's `env` block works too; it takes precedence over the `.env` file.

---

## Tools

42 tools in four groups. Every tool carries MCP annotations: `readOnlyHint` on reads, `destructiveHint` on anything irreversible (`x_post_delete`, `x_list_delete`, unlike, unrepost, unfollow, member removal). Each returns a structured JSON payload matching its output schema; the SDK mirrors the same JSON into the text content block for clients that do not read `structuredContent`.

Every tool accepts an optional `timeout_ms`, clamped into the group's `[MIN, MAX]` window (see [Limits](#limits)).

### Posts — reading

| Tool | Description |
| :--- | :--- |
| **`x_posts_search`** | Searches recent posts for several queries **in parallel**. Recent search reaches back **7 days only**. |
| **`x_posts_count`** | Counts matches per query bucketed by minute/hour/day. **Does not spend post reads** — use it to size a topic before searching. |
| **`x_posts_lookup`** | Fetches up to **100 posts by id** in one call. |
| **`x_posts_by_user`** | Recent posts for several usernames, fetched in parallel. |
| **`x_posts_mentions`** | Posts mentioning the key owner. |
| **`x_posts_home`** | The key owner's home timeline. |
| **`x_posts_quotes`** | Posts quoting a given post. |
| **`x_posts_liked`** | Posts the key owner liked. |
| **`x_bookmarks_list`** | The key owner's bookmarks. |
| **`x_post_liked_by`** | Users who liked a given post. |
| **`x_post_reposted_by`** | Users who reposted a given post. |

#### `x_posts_search`

| Parameter | Type | Default | Notes |
| :--- | :--- | :--- | :--- |
| `queries` | `[]string` | — | **Required.** Run in parallel, capped at `POSTS_MAX_QUERIES` (`5`). ≤ 512 chars each. |
| `max_results` | `int` | `10` | Posts **per query**, 1..100. Every one is billed. |
| `sort_order` | `string` | — | `recency` (newest first) or `relevancy` (best match). |
| `days` | `int` | `7` | How far back, 1..7. The API cannot go further. |
| `include_retweets` | `bool` | `false` | When false the server appends `-is:retweet` to every query. |
| `timeout_ms` | `int64` | `15000` | Whole-call timeout, clamped to `[2000, 60000]`. |

Query operators go **inside the query string**:

| Operator | Meaning |
| :--- | :--- |
| *space* | AND |
| `OR` | explicit OR |
| `-term` | NOT |
| `( )` | grouping |
| `from:user` / `to:user` | by author / by recipient |
| `@user` / `#tag` | mentions / hashtags |
| `"exact phrase"` | exact phrase |
| `lang:en` | language |
| `is:retweet` `is:reply` `is:quote` `is:verified` | post kind |
| `has:media` `has:images` `has:videos` `has:links` | attachments |
| `url:example.com` | linked domain |
| `conversation_id:123` | one thread |

#### `x_posts_count`

| Parameter | Type | Default | Notes |
| :--- | :--- | :--- | :--- |
| `queries` | `[]string` | — | **Required.** Same operators as `x_posts_search`. |
| `granularity` | `string` | `hour` | `minute`, `hour` or `day`. |
| `days` | `int` | `7` | 1..7, same window. |
| `timeout_ms` | `int64` | `15000` | Clamped to `[2000, 60000]`. |

#### `x_posts_lookup`

| Parameter | Type | Default | Notes |
| :--- | :--- | :--- | :--- |
| `ids` | `[]string` | — | **Required.** 1..100 post ids. Batch them; do not call once per id. |

#### `x_posts_by_user`

| Parameter | Type | Default | Notes |
| :--- | :--- | :--- | :--- |
| `usernames` | `[]string` | — | **Required.** Without the `@`. Capped at `POSTS_MAX_USERNAMES` (`5`), fetched in parallel. |
| `max_results` | `int` | `10` | Posts per user, 1..100. |
| `exclude` | `[]string` | — | `replies` and/or `retweets`. |
| `days` | `int` | `7` | 1..7. |
| `timeout_ms` | `int64` | `15000` | Clamped to `[2000, 60000]`. |

#### Timeline tools

`x_posts_mentions`, `x_posts_home`, `x_posts_liked` and `x_bookmarks_list` share one shape: `max_results` (1..100), `pagination_token`, `timeout_ms`. `x_posts_quotes`, `x_post_liked_by` and `x_post_reposted_by` add a required `id`.

Pagination is explicit and manual: a response carries `next_token`, and you pass it back as `pagination_token` for the next page. The server never walks pages on its own — each page is billed, so that decision stays with the caller.

### Posts — writing

| Tool | Annotation | Description |
| :--- | :--- | :--- |
| **`x_post_create`** | write | Publishes a post. **Public and billed** ($0.015, or $0.20 with a link). |
| **`x_post_delete`** | destructive | Deletes one of your posts. Irreversible. |
| **`x_post_like`** / **`x_post_unlike`** | write / destructive | Like is public. |
| **`x_post_repost`** / **`x_post_unrepost`** | write / destructive | Repost is public. |

#### `x_post_create`

| Parameter | Type | Notes |
| :--- | :--- | :--- |
| `text` | `string` | ≤ 280 characters. Required unless `media_ids` is set. |
| `reply_to_id` | `string` | Reply to this post. |
| `quote_id` | `string` | Quote this post. |
| `media_ids` | `[]string` | 1..4 ids from `x_media_upload`. |
| `poll` | `object` | `options` (2..4) + `duration_minutes` (5..10080). **Cannot be combined with `media_ids`.** |
| `reply_settings` | `string` | `following`, `mentionedUsers`, `subscribers` or `verified`. |

Publishing the same text twice in a row is rejected by X as a duplicate — the server surfaces that as a distinct message rather than a generic failure.

### Users

| Tool | Annotation | Description |
| :--- | :--- | :--- |
| **`x_users_lookup`** | read | Up to **100 usernames or 100 ids** in one call — one or the other, not both. Unresolved names come back in `not_found`. |
| **`x_users_me`** | read | The profile behind the keys. |
| **`x_user_followers`** / **`x_user_following`** | read | One page at a time; `max_results` up to 1000. |
| **`x_user_follow`** / **`x_user_unfollow`** | write / destructive | Public. |
| **`x_user_mute`** / **`x_user_unmute`** | write | A mute is silent: the other account is not told and can still see you. |
| **`x_users_muted`** / **`x_users_blocked`** | read | The key owner's muted / blocked accounts. |

Blocking and unblocking are deliberately **not** exposed. Reading the block list is.

### Lists

| Tool | Annotation | Description |
| :--- | :--- | :--- |
| **`x_list_create`** | write | Name ≤ 25 chars, description ≤ 100. Private lists are owner-only. |
| **`x_list_update`** | write | Send only the fields you want changed. |
| **`x_list_delete`** | destructive | Irreversible: list, membership and followers all go. |
| **`x_list_get`** | read | One list with member and follower counts. |
| **`x_lists_owned`** | read | Lists owned by a username, or by the key owner when omitted. |
| **`x_list_posts`** | read | Recent posts from the list's members. Billed per post. |
| **`x_list_members`** | read | Accounts in the list, paginated. |
| **`x_list_member_add`** / **`x_list_member_remove`** | write / destructive | By username. |
| **`x_list_follow`** / **`x_list_unfollow`** | write / destructive | Public lists only. |
| **`x_list_pin`** / **`x_list_unpin`** / **`x_lists_pinned`** | write / write / read | Pinned lists of the key owner. |

### Media

#### `x_media_upload`

Uploads an image, GIF or video with the chunked INIT → APPEND → FINALIZE flow and returns a `media_id` for `x_post_create`.

| Parameter | Type | Notes |
| :--- | :--- | :--- |
| `path` | `string` | Local file. **Must sit inside `MEDIA_ROOT`.** |
| `base64` | `string` | File contents instead of a path; requires `mime`. |
| `mime` | `string` | e.g. `image/jpeg`, `video/mp4`. Required with `base64`. |
| `category` | `string` | **Required.** `tweet_image` (≤ 5 MB), `tweet_gif` (≤ 15 MB), `tweet_video` (≤ 512 MB). |
| `timeout_ms` | `int64` | Default `120000`, clamped to `[5000, 600000]` — covers upload *and* X's transcoding. |

Two things to know:

- **`MEDIA_ROOT` is a hard boundary.** It is required configuration with no default, and the server refuses any path resolving outside it, symlinks included. Without it an LLM with this tool could read any file on the host and post it. Choose a dedicated directory.
- **Video is transcoded asynchronously by X.** The server polls FINALIZE status every `MEDIA_POLL_INTERVAL_MS`, but a slow encode can still return `status: "pending"`. The `media_id` is valid but not yet postable — retry `x_post_create` shortly. A `media_id` expires after a while and is meant for a single post.

### Partial results

The fan-out tools — `x_posts_search`, `x_posts_count`, `x_posts_by_user` — run their inputs in parallel and return **one entry per query/username, each with its own `status`**, so one bad input does not lose the results that worked. `x_posts_lookup` does the same per id (`not found` sits next to the posts that resolved), and `x_users_lookup` collects unresolved names in `not_found`.

A call fails outright only when the input is rejected before any work starts, or when **every** item in it fails.

### Errors

Failures come back as a tool result with `isError: true` and a plain-text message, not as a JSON-RPC error — the model reads the message and can correct the call itself. The messages are written for a model rather than a log reader, so the ones that must not be retried say so explicitly.

| Message | Meaning |
| :--- | :--- |
| `the access token is read-only: …regenerate the Access Token and Secret` | The [read-only-token trap](#getting-them--and-the-one-step-everybody-gets-wrong). |
| `the monthly usage cap for this project is reached; …do not retry` | Usage cap hit. Nothing works until it resets. |
| `the project has no X API credits left; …do not retry` | Credits exhausted. They do not reset — top up in the Console. |
| `rate limit exceeded; the window resets in about 15 minutes` | Back off; do not hammer. |
| `X rejected the credentials; check the four keys` | Bad or revoked keys. |
| `not authorized for this resource` | Protected, deleted, or somebody else's. |
| `X rejected this as a duplicate of a recent post` | Change the text. |
| `user not found` / `post not found` / `list not found` | Self-explanatory. |
| `X failed to process the uploaded media; re-encode the file` | Transcoding failed. |
| `media exceeds the size limit for its category` | 5 MB / 15 MB / 512 MB. |
| `the path is outside the directory this server is allowed to read` | `MEDIA_ROOT` violation. |
| `post text exceeds 280 characters` | — |
| `a post needs text unless it carries media` | — |
| `a poll needs 2 to 4 options and a duration between 5 and 10080 minutes` | — |
| `a post carries at most 4 media items` | — |
| `too many ids in one call; the limit is 100` | — |
| `invalid username: 1 to 15 letters, digits or underscores` | Validated locally, before spending a request. |
| `every query failed; the X API may be unreachable` | All items failed. |
| `the call timed out; retry with a larger timeout_ms or fewer inputs` | — |
| `the X API is unavailable` | 5xx or unclassified upstream failure. |

Classification happens in [`adapter/x/apierr`](app/internal/adapter/x/apierr/apierr.go), and it is deliberately defensive. X answers errors with `Content-Type: application/problem+json`, which the underlying `gotwi` client does not recognise as JSON — so the entire problem body lands verbatim in a message field instead of being decoded into typed fields. The mapper therefore concatenates every source the response can carry and matches problem types (`oauth1-permissions`, `usage-capped`, `credits-depleted`, `rate-limit-exceeded`, …) against that haystack, falling back to HTTP status codes when nothing matches.

### Known limitations

- **Recent search only.** 7 days back, full stop. Historical/archive search is a different (Enterprise) product and is not wired up.
- **Blocking is not exposed.** `x_users_blocked` reads the list; there is no block/unblock tool.
- **No streaming endpoints.** Filtered/sampled stream is not implemented.
- **No DMs.**
- **One account per process.** The keys are process-wide; the key owner is resolved once at startup. Serving several accounts means several processes.
- **Pagination is manual.** Tools return `next_token`; nothing walks pages automatically, because every page is billed.

---

## Quick start

### Install

Pick whichever fits — all of them give the identical server.

**Container** (no Go toolchain needed):

```bash
docker pull ghcr.io/role1776/mcp-x:0.1.1     # or :latest to track the newest release
```

Pin an explicit version in anything you rely on. `:latest` moves on every stable release, and a release can add or change tool behaviour; [`server.json`](server.json) pins the same version the MCP Registry advertises.

**Prebuilt binary** — grab the archive for your platform from the [latest release](https://github.com/Role1776/mcp-x/releases/latest), unpack it, and put `mcp-x` on your `PATH`.

**MCP Bundle** — for clients that install `.mcpb` files, download `mcp-x_<version>_<os>_<arch>.mcpb` from the [latest release](https://github.com/Role1776/mcp-x/releases/latest) and open it with your client. The bundle carries the compiled binary, so it needs neither Docker nor Go, and the client prompts for the five required values on install. Pick the file matching your OS *and* CPU architecture: a bundle holds one native binary.

> [!NOTE]
> Bundles built before v0.1.1 declared no configuration fields, so the client never asked for the credentials and the server exited on startup every time. Use v0.1.1 or newer.

**From source** (needs Go 1.26+):

```bash
git clone https://github.com/Role1776/mcp-x
cd mcp-x
make build          # -> bin/mcp-x
```

`go install` also works, with one caveat worth knowing before you type it:

```bash
go install github.com/Role1776/mcp-x/app/cmd/mcp-x@latest
```

The Go module lives in `app/`, so the release tags (`v0.1.0`) do not name it — `@v0.1.0` fails outright and `@latest` resolves to a pseudo-version of the newest commit on `main`. A `go install` build also reports its version as `dev`, because the version is stamped in by the release pipeline and not by the module. For a build that matches a release exactly, use the prebuilt binary, the container, or `make build` from a checked-out tag.

### Configure

From a clone, or from an unpacked release archive (both ship the template):

```bash
cp .env.example .env
$EDITOR .env        # fill in the four X_* keys and MEDIA_ROOT
```

Installing from the container or an `.mcpb` bundle gives you no checkout to copy from — take the template from [`.env.example`](.env.example), or skip the file entirely and pass the five values in your client's `env` block (below).

Five values are required; everything else has a default. `MEDIA_ROOT` must be an **absolute path to a directory that already exists** — the server checks this at startup and refuses to start otherwise, rather than failing at the first upload.

### Run

```bash
./bin/mcp-x -env /absolute/path/to/.env
```

| Flag | Meaning |
| :--- | :--- |
| `-version` | Print the version the binary reports to MCP clients, and exit. Answerable without credentials, unlike the handshake. |
| `-env` | Path to a `.env` file. **There is no implicit lookup** — under stdio the working directory is chosen by the MCP client, so a relative default would be unpredictable. If the flag is omitted, or the file does not exist, the server falls back to the ambient environment and says so on `stderr`; the five required values must still be set somewhere or startup fails. |

A successful start logs the authenticated account:

```
INFO authenticated with the X API op=app.Run user_id=1234567890
INFO MCP server has started op=app.Run transport=stdio
```

### Connecting an MCP client (stdio)

**Claude Code** — one command, credentials inline (`--` separates the server's own command from the flags above it):

```bash
claude mcp add mcp-x \
  -e X_API_KEY=... \
  -e X_API_KEY_SECRET=... \
  -e X_ACCESS_TOKEN=... \
  -e X_ACCESS_TOKEN_SECRET=... \
  -e MEDIA_ROOT=/absolute/path/to/media \
  -- /absolute/path/to/mcp-x
```

Or point it at a `.env` file instead: `-- /absolute/path/to/mcp-x -env /absolute/path/to/.env`. Check the result with `claude mcp list`, or `/mcp` inside a session.

**Claude Desktop** and other clients that take a JSON config — point them at the built binary:

```json
{
  "mcpServers": {
    "x": {
      "command": "/absolute/path/to/mcp-x",
      "args": ["-env", "/absolute/path/to/.env"]
    }
  }
}
```

Or skip the file and pass the credentials in the `env` block — variables already in the environment win over the `.env` file, so a client's `env` block always takes effect:

```json
{
  "mcpServers": {
    "x": {
      "command": "/absolute/path/to/mcp-x",
      "env": {
        "X_API_KEY": "...",
        "X_API_KEY_SECRET": "...",
        "X_ACCESS_TOKEN": "...",
        "X_ACCESS_TOKEN_SECRET": "...",
        "MEDIA_ROOT": "/absolute/path/to/media"
      }
    }
  }
}
```

### Connecting an MCP client (container)

Run the image on stdio. Docker needs each variable named on the command line with `-e` for it to reach the process, and `MEDIA_ROOT` only makes sense if the directory is mounted:

```json
{
  "mcpServers": {
    "x": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "X_API_KEY",
        "-e", "X_API_KEY_SECRET",
        "-e", "X_ACCESS_TOKEN",
        "-e", "X_ACCESS_TOKEN_SECRET",
        "-e", "MEDIA_ROOT",
        "-v", "/host/media:/media:ro",
        "ghcr.io/role1776/mcp-x:0.1.1"
      ],
      "env": {
        "X_API_KEY": "...",
        "X_API_KEY_SECRET": "...",
        "X_ACCESS_TOKEN": "...",
        "X_ACCESS_TOKEN_SECRET": "...",
        "MEDIA_ROOT": "/media"
      }
    }
  }
}
```

`-i` is required — without it the container gets no stdin and the client sees the server die immediately. Clients that install *from the [MCP Registry](https://registry.modelcontextprotocol.io) itself* build this invocation themselves and prompt for the variables declared in [`server.json`](server.json) — that is a property of the client, not of every client that can talk to this server; Claude Code's CLI, for one, does not install from the registry, so use `claude mcp add` above.

### When the client says the connection closed

Under stdio the server's `stderr` belongs to the client, and most clients discard it. So a configuration problem that the server explains perfectly well in one line reaches you as nothing but:

```
✘ Failed to connect — -32000: Connection closed
```

**Run the binary by hand with the same environment to see the real reason:**

```bash
/absolute/path/to/mcp-x -env /absolute/path/to/.env
```

It prints exactly what is wrong — which variables are missing and where to get them, or that the keys were rejected, or that the access token is read-only — and exits. Almost every failed install is one of those three.

### Running over HTTP

Set `MCP_TRANSPORT=http` and the server listens on `SERVER_PORT` at `MCP_PATH` (default `http://localhost:8080/mcp`).

> [!CAUTION]
> The HTTP transport has **no authentication of its own**. Anyone who can reach the endpoint can publish, delete and follow as your account, on your credits. Bind it to localhost or put it behind an authenticating reverse proxy — never expose it to the open internet.

---

## Configuration

Everything is configured through **environment variables**, and each value is validated before startup: a missing required key, or a non-numeric or non-positive number, is a startup error naming the variable you actually set (`X_API_KEY`, not `APIKey`). Relationships *between* limits are not checked at startup — see [Limits](#limits). Variables already present in the environment win over a `.env` file.

See [`.env.example`](.env.example) for the full list at its default values, ready to copy to `.env`.

### Required

| Env | Notes |
| :--- | :--- |
| `X_API_KEY` | OAuth 1.0a consumer key. |
| `X_API_KEY_SECRET` | OAuth 1.0a consumer secret. |
| `X_ACCESS_TOKEN` | User access token — [generate it *after* setting Read and Write](#credentials). |
| `X_ACCESS_TOKEN_SECRET` | User access token secret. |
| `MEDIA_ROOT` | Absolute path to the only directory `x_media_upload` may read from. Must exist and be a directory; both are checked at startup. No default, on purpose. |

Missing any of them aborts startup with the list of what is missing plus a pointer to the Developer Console.

### MCP server

| Env | Default | Notes |
| :--- | :--- | :--- |
| `MCP_TRANSPORT` | `stdio` | `stdio` or `http`. |
| `MCP_NAME` | `mcp-x` | Server name advertised to clients. |
| `MCP_PATH` | `/mcp` | HTTP route (http transport only). |

The version advertised to clients is not configurable: it is stamped into the binary at build time from the git tag.

### HTTP server (http transport only)

| Env | Default |
| :--- | :--- |
| `SERVER_PORT` | `8080` |
| `SERVER_READ_TIMEOUT` | `60s` |
| `SERVER_WRITE_TIMEOUT` | `60s` |

### HTTP client

| Env | Default | Notes |
| :--- | :--- | :--- |
| `MAX_IDLE_CONNS_PER_HOST` | `100` | Connection pooling toward the X API. |
| `X_DEBUG` | `false` | Log raw `gotwi` requests. Verbose; useful when an X error makes no sense. |

### Logging

| Env | Default | Notes |
| :--- | :--- | :--- |
| `LOG_MODE` | `local` | `local` → text handler at debug level; `prod` → JSON handler at info level. Logs go to `stderr` (they must: `stdout` carries the MCP protocol). |

### Limits

Each group has its own limits, so a slow media upload cannot be capped by a search timeout.

**Posts**

| Env | Default | Notes |
| :--- | :--- | :--- |
| `POSTS_MAX_QUERIES` | `5` | Queries per `x_posts_search` / `x_posts_count` call. |
| `POSTS_MAX_USERNAMES` | `5` | Usernames per `x_posts_by_user` call. |
| `POSTS_MAX_IDS` | `100` | Ids per `x_posts_lookup`; the API's own ceiling. |
| `POSTS_DEFAULT_RESULTS` | `10` | |
| `POSTS_MAX_RESULTS` | `100` | |
| `POSTS_DEFAULT_DAYS` | `7` | |
| `POSTS_MAX_DAYS` | `7` | Recent search cannot look back further. |
| `POSTS_DEFAULT_TIMEOUT_MS` | `15000` | |
| `POSTS_MIN_TIMEOUT_MS` | `2000` | |
| `POSTS_MAX_TIMEOUT_MS` | `60000` | |

**Users**

| Env | Default |
| :--- | :--- |
| `USERS_MAX_USERNAMES` | `100` |
| `USERS_MAX_IDS` | `100` |
| `USERS_DEFAULT_RESULTS` | `100` |
| `USERS_MAX_RESULTS` | `1000` |
| `USERS_DEFAULT_TIMEOUT_MS` | `15000` |
| `USERS_MIN_TIMEOUT_MS` | `2000` |
| `USERS_MAX_TIMEOUT_MS` | `60000` |

**Lists**

| Env | Default |
| :--- | :--- |
| `LISTS_DEFAULT_RESULTS` | `25` |
| `LISTS_MAX_RESULTS` | `100` |
| `LISTS_DEFAULT_TIMEOUT_MS` | `15000` |
| `LISTS_MIN_TIMEOUT_MS` | `2000` |
| `LISTS_MAX_TIMEOUT_MS` | `60000` |

**Media**

| Env | Default | Notes |
| :--- | :--- | :--- |
| `MEDIA_CHUNK_BYTES` | `4194304` | 4 MB. The APPEND endpoint caps a segment at 5 MB. |
| `MEDIA_POLL_INTERVAL_MS` | `2000` | How often FINALIZE status is polled while X transcodes. |
| `MEDIA_DEFAULT_TIMEOUT_MS` | `120000` | |
| `MEDIA_MIN_TIMEOUT_MS` | `5000` | |
| `MEDIA_MAX_TIMEOUT_MS` | `600000` | 10 minutes, for large video. |

Each value is checked on its own — it must be greater than zero, and the ones the API itself bounds (`*_MAX_IDS`, `POSTS_MAX_RESULTS`, `LISTS_MAX_RESULTS`, `POSTS_MAX_DAYS`) are additionally capped at the API's ceiling so a typo cannot produce a request X will reject. The `DEFAULT_*`, `MIN_*` and `MAX_*` triples are **not** cross-checked against each other at startup. An inconsistent set does not stop the server; it is reconciled per request instead:

- a value the caller omits, or passes as zero or negative, falls back to the matching `DEFAULT_*`;
- the result is then clamped into `[MIN_*, MAX_*]`, so a `DEFAULT_*` larger than its `MAX_*` simply yields `MAX_*`;
- if `MIN_*` exceeds `MAX_*`, the maximum wins.

The effective limit is therefore always within the configured maximum, and misconfiguration degrades to a working server rather than a failed start. The trade-off is that it degrades **silently**: `POSTS_MAX_RESULTS=1` instead of `10` produces no warning, only quietly smaller responses. Worth double-checking these values when results look truncated.

---

## Architecture

The project follows a clean, layered structure. Dependencies point inward toward the domain, and each layer talks to the next through interfaces.

```
app/                       the Go module: sources plus its build files
                           (Dockerfile, .dockerignore, .goreleaser.yaml)

cmd/mcp-x/main.go          entry point: parse flags, load config, run app

internal/
  app/                     wiring + lifecycle (X client, bootstrap, run, graceful shutdown)
  config/                  config loading (.env → env vars → validate → env-name error messages)
  domain/x/                core types and errors
    posts/ users/ lists/ media/   value objects: PostID, UserID, Username, Draft, Query, Upload…
    format.go              shared snowflake-id and username validation
    errors.go              the sentinel error set the whole app maps onto
  dto/x/                   request/response shapes for the MCP tools, with jsonschema + validate tags
  transport/mcp/           MCP layer
    router/                registers every tool group on the MCP server
    x/posts|users|lists|media/   tool definitions, handlers and error → tool-result mapping
  usecase/x/               business logic: validation, parallelism, timeouts, limit resolution
  adapter/x/               X API wiring
    client/                gotwi client + startup Bootstrap (key check, key-owner id)
    apierr/                X problem responses → domain errors
    posts|users|lists|media/     endpoint calls and mappers to DTOs
  pkg/                     reusable building blocks (mcpserver, server, logger, validator, parallel)
```

Request flow for a tool call:

```
MCP client → transport/mcp/x/... (handler) → usecase/x/... → adapter/x/... → gotwi → X API
                    ↑ maps errors                 ↑ validates, resolves limits,
                      to isError text               fans out, applies timeouts
```

Two boundaries are worth calling out:

- **The domain layer refuses to construct invalid values.** `NewPostID`, `NewUserID`, `NewUsername`, `NewDraft` and friends validate on construction, so a malformed id or a 300-character post is rejected locally — before it costs a request. Ids are checked against the snowflake format and usernames against X's 1–15 `[A-Za-z0-9_]` rule in a single shared `format` package.
- **Every layer speaks the same error vocabulary.** Adapters convert X's problem responses into the sentinels in `domain/x/errors.go`; the transport layer is the only place that turns a sentinel into human text. That is why the same failure reads identically no matter which of the 42 tools produced it.

---

## Development

Everything Go lives in `app/`, so either use the makefile from the repository root or pass `-C app` to the toolchain:

```bash
make build          # compile the binary -> bin/mcp-x
make test           # go test -v ./...
make cover          # total coverage percentage
make cover-html     # coverage report in the browser
make version        # the version that would be stamped into the binary

go -C app build ./...
go -C app test ./...
go -C app vet ./...
```

Mocks are generated with `mockgen` from the `//go:generate` directives next to each `usecase` interface:

```bash
go -C app generate ./...
```

**New code is expected to come with tests.** See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the full pull-request guidelines.

## License

Released under the [MIT License](LICENSE).
