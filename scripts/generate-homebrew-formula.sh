#!/usr/bin/env bash

set -euo pipefail

repo=""
tag=""
sha256=""

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --repo)
      repo="$2"
      shift 2
      ;;
    --tag)
      tag="$2"
      shift 2
      ;;
    --sha256)
      sha256="$2"
      shift 2
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

if [[ -z "${repo}" || -z "${tag}" || -z "${sha256}" ]]; then
  echo "usage: $0 --repo OWNER/REPO --tag vX.Y.Z --sha256 SHA256" >&2
  exit 1
fi

cat <<EOF
class Codexclaw < Formula
  desc "Lightweight Codex task daemon"
  homepage "https://github.com/${repo}"
  url "https://github.com/${repo}/archive/refs/tags/${tag}.tar.gz"
  sha256 "${sha256}"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"codexclaw"), "./cmd/codexclaw"
  end

  test do
    assert_match "Lightweight Codex task daemon", shell_output("#{bin}/codexclaw --help")
  end
end
EOF
