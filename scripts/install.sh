#!/usr/bin/env sh
set -eu

repo_url="https://github.com/ESHAYAT102/lapip.git"
binary_name="lapip"
install_dir="${XDG_BIN_HOME:-${HOME}/.local/bin}"
clone_dir=""

cleanup() { if [ -n "${clone_dir:-}" ] && [ -d "$clone_dir" ]; then rm -rf "$clone_dir"; fi; }
trap cleanup EXIT INT TERM

command -v go >/dev/null 2>&1 || { echo "error: Go is required" >&2; exit 1; }
command -v git >/dev/null 2>&1 || { echo "error: git is required" >&2; exit 1; }

mkdir -p "$install_dir"
clone_dir="$(mktemp -d "${TMPDIR:-/tmp}/lapip.XXXXXX")"
echo "cloning $repo_url"
git clone --depth 1 "$repo_url" "$clone_dir"
echo "building $binary_name"
(cd "$clone_dir" && go build -o "$install_dir/$binary_name" .)
chmod +x "$install_dir/$binary_name"
echo "installed $binary_name to $install_dir/$binary_name"

case ":$PATH:" in
  *":$install_dir:"*) ;;
  *) echo "warning: $install_dir is not in PATH"; echo "add: export PATH=\"$install_dir:\$PATH\"" ;;
esac
