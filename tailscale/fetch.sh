#!/bin/sh
# Fetch the pinned Tailscale bundle into build/tailscale.
#
# The pinned version lives in tailscale/VERSION (single source of truth); keep
# it in sync with the SHA-256 below and with tailscale/README.md. Downloads the
# static arm64 tarball and verifies its SHA-256 before extracting the binaries.
set -eu

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VERSION="$(cat "$SCRIPT_DIR/VERSION")"
ARCH=arm64
BASE_URL="https://pkgs.tailscale.com/stable"
TARBALL="tailscale_${VERSION}_${ARCH}.tgz"

# SHA-256 of the above tarball. This is specific to the pinned VERSION.
EXPECTED_SHA256="d6bed2720310b328b861f943d7b4016c5358a9bad60575eb51e8143b79068484"

OUT_DIR="$SCRIPT_DIR/../build/tailscale"
mkdir -p "$OUT_DIR"

# sha256check prints the SHA-256 of a file, using sha256sum on GNU/Linux and
# shasum on macOS/BSD.
sha256check() {
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$1" | awk '{print $1}'
    else
        shasum -a 256 "$1" | awk '{print $1}'
    fi
}

cd "$OUT_DIR"

if [ ! -f "$TARBALL" ]; then
    echo "Downloading $BASE_URL/$TARBALL"
    curl -fsSL -o "$TARBALL" "$BASE_URL/$TARBALL"
fi

echo "Verifying checksum..."
ACTUAL="$(sha256check "$TARBALL")"
if [ "$ACTUAL" != "$EXPECTED_SHA256" ]; then
    echo "checksum mismatch for $TARBALL" >&2
    echo "expected: $EXPECTED_SHA256" >&2
    echo "actual:   $ACTUAL" >&2
    exit 1
fi

echo "Extracting..."
tar -xzf "$TARBALL"

# The tarball extracts into a top-level directory named after the release.
EXTRACT_DIR="tailscale_${VERSION}_${ARCH}"
if [ ! -x "$EXTRACT_DIR/tailscaled" ] || [ ! -x "$EXTRACT_DIR/tailscale" ]; then
    echo "extracted tree missing expected binaries" >&2
    exit 1
fi

echo "Done: $OUT_DIR"
ls -l "$OUT_DIR/$EXTRACT_DIR/tailscaled" "$OUT_DIR/$EXTRACT_DIR/tailscale"