#!/usr/bin/env bash
set -euo pipefail

rm -f coverage.out
profiles=()
for pkg in $(go list ./...); do
  f="coverage_$(echo $pkg | tr '/' '_').out"
  echo "[cover] testing $pkg -> $f"
  go test -covermode=atomic -coverprofile="$f" "$pkg"
  profiles+=("$f")
done
echo "[cover] merging ${#profiles[@]} profiles"
go run ./tools/covermerge "${profiles[@]}" > coverage.out
echo "[cover] total:"
go tool cover -func coverage.out | tail -n 1
