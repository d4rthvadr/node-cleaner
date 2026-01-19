#!/usr/bin/env bash
set -euo pipefail

# Build script for DepoCleaner
# Cross-compiles binaries and generates checksums.

APP_NAME="depo-cleaner"
DIST_DIR="dist"
LD_FLAGS="-s -w"

# Build matrix (OS ARCH)
MATRIX=(
  "darwin arm64"
  "darwin amd64"
  "linux amd64"
  "linux arm64"
)

mkdir -p "${DIST_DIR}"

build_target() {
  local os="$1" arch="$2"
  local out="${DIST_DIR}/${APP_NAME}_${os}_${arch}"
  echo "Building ${out}..."
  GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -ldflags "${LD_FLAGS}" -o "${out}" ./
}

checksum_file() {
  local file="$1"
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "${file}" > "${file}.sha256"
  elif command -v sha256sum >/dev/null 2>&1; then
    sha256sum "${file}" > "${file}.sha256"
  else
    echo "Warning: no sha256 program found; skipping checksum for ${file}" >&2
  fi
}

# Ensure go.mod exists
if [[ ! -f go.mod ]]; then
  echo "Error: go.mod not found in current directory" >&2
  exit 1
fi

# Tidy modules (optional)
if command -v go >/dev/null 2>&1; then
  go mod tidy
fi

for entry in "${MATRIX[@]}"; do
  os=$(echo "$entry" | awk '{print $1}')
  arch=$(echo "$entry" | awk '{print $2}')
  build_target "$os" "$arch"
  checksum_file "${DIST_DIR}/${APP_NAME}_${os}_${arch}"
done

# Copy README and LICENSE to dist for convenience
for f in README.md LICENSE; do
  if [[ -f "$f" ]]; then
    cp "$f" "${DIST_DIR}/" || true
  fi
done

echo "Build complete. Binaries in ${DIST_DIR}/"
