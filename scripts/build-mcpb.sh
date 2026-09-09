#!/usr/bin/env bash
# Packs each binary GoReleaser produced into an MCP Bundle (.mcpb).
#
# An .mcpb is a zip holding a manifest.json plus the server payload; the
# "binary" server type means the bundle carries a pre-compiled executable and
# needs no runtime on the user's machine. One bundle per OS/arch, because a
# native binary is not portable across them.
#
# Usage: scripts/build-mcpb.sh <version>
# Reads:  dist/artifacts.json (written by `goreleaser release`)
# Writes: dist/mcpb/*.mcpb
set -euo pipefail

version="${1:?usage: build-mcpb.sh <version>}"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
dist="$root/dist"
out="$dist/mcpb"

[ -f "$dist/artifacts.json" ] || { echo "no $dist/artifacts.json — run goreleaser first" >&2; exit 1; }

rm -rf "$out"
mkdir -p "$out"

# goos -> the platform identifier the MCPB manifest schema expects
mcpb_platform() {
  case "$1" in
    linux)   echo linux ;;
    darwin)  echo darwin ;;
    windows) echo win32 ;;
    *)       echo "unsupported goos: $1" >&2; return 1 ;;
  esac
}

count=0
while IFS=$'\t' read -r path goos goarch; do
  binary="$(basename "$path")"
  platform="$(mcpb_platform "$goos")"
  stage="$(mktemp -d)"

  mkdir -p "$stage/server"
  cp "$root/$path" "$stage/server/$binary"
  chmod +x "$stage/server/$binary"

  # ${__dirname} is expanded by the host app to the unpacked bundle directory.
  jq -n \
    --arg version "$version" \
    --arg binary "$binary" \
    --arg platform "$platform" \
    '{
      manifest_version: "0.3",
      name: "io.github.Role1776/mcp-x",
      display_name: "X (Twitter)",
      version: $version,
      description: "Search, read and publish X posts and manage users, lists and media via the X API v2.",
      author: { name: "Role1776", url: "https://github.com/Role1776" },
      homepage: "https://github.com/Role1776/mcp-x",
      repository: { type: "git", url: "https://github.com/Role1776/mcp-x" },
      license: "MIT",
      server: {
        type: "binary",
        entry_point: ("server/" + $binary),
        mcp_config: {
          command: ("${__dirname}/server/" + $binary),
          args: [],
          # Substituted by the host app from the answers to user_config below.
          # Without this the server would start with no credentials and exit.
          env: {
            X_API_KEY: "${user_config.x_api_key}",
            X_API_KEY_SECRET: "${user_config.x_api_key_secret}",
            X_ACCESS_TOKEN: "${user_config.x_access_token}",
            X_ACCESS_TOKEN_SECRET: "${user_config.x_access_token_secret}",
            MEDIA_ROOT: "${user_config.media_root}"
          }
        }
      },
      # The five variables the server refuses to start without. Everything else
      # has a default and is left to the .env file or the ambient environment.
      user_config: {
        x_api_key: {
          type: "string", title: "API key", sensitive: true, required: true,
          description: "OAuth 1.0a consumer key. developer.x.com -> your project -> your app -> Keys and tokens."
        },
        x_api_key_secret: {
          type: "string", title: "API key secret", sensitive: true, required: true,
          description: "OAuth 1.0a consumer secret, issued together with the API key."
        },
        x_access_token: {
          type: "string", title: "Access token", sensitive: true, required: true,
          description: "User access token. Set the app permissions to Read and Write BEFORE generating it: a token keeps the permissions it was created with, and a read-only one fails every write until it is regenerated."
        },
        x_access_token_secret: {
          type: "string", title: "Access token secret", sensitive: true, required: true,
          description: "User access token secret, issued together with the access token."
        },
        media_root: {
          type: "directory", title: "Media directory", required: true,
          description: "The only directory x_media_upload may read files from. Anything outside it is refused, symlinks included."
        }
      },
      compatibility: { platforms: [$platform] }
    }' > "$stage/manifest.json"

  bundle="$out/mcp-x_${version}_${goos}_${goarch}.mcpb"
  (cd "$stage" && zip -q -r -X "$bundle" manifest.json server)
  rm -rf "$stage"

  echo "packed $(basename "$bundle")"
  count=$((count + 1))
done < <(jq -r '.[] | select(.type == "Binary") | [.path, .goos, .goarch] | @tsv' "$dist/artifacts.json")

[ "$count" -gt 0 ] || { echo "no binaries found in artifacts.json" >&2; exit 1; }
