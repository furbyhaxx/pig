#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
set -euo pipefail

: "${REPOSITORY:?}" "${GH_TOKEN:?}" "${PIG_UPDATE_SIGNING_KEY:?}"
[[ "$REPOSITORY" = furbyhaxx/pig ]] || { echo 'Refusing publication outside furbyhaxx/pig' >&2; exit 1; }
sha=$(git rev-parse HEAD)
# API failures must not be mistaken for an absent release.
releases=$(gh api "repos/$REPOSITORY/releases?per_page=100" --paginate --slurp)
if jq -e --arg sha "$sha" 'any(.[][]; .draft == false and (.body // "" | contains("Source: " + $sha)))' <<< "$releases" >/dev/null; then
  echo "Source $sha is already released."
  exit 0
fi
base=$(awk -F '"' '/^const PigVersion = / {print $2}' internal/coding/pigversion/pigversion.go)
[[ "$base" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "Invalid upstream PiG version: $base" >&2; exit 1; }
# Include draft/failed publication tags so a version is never reused.
git fetch origin --tags
revision=1
while IFS= read -r tag; do
  number=${tag#"v$base-c."}
  [[ "$number" =~ ^[1-9][0-9]*$ ]] || continue
  if ((number >= revision)); then revision=$((number + 1)); fi
done < <(git tag -l "v$base-c.*")
version="$base-c.$revision"
make set-version VERSION="$version"
go test ./automation/release/set-version
go test ./internal/codingagent -run 'Test(IsNewerPackageVersion|UpdateSourceURL|FetchUpdateManifest|CheckForBinaryUpdate)'

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
umask 077
printf '%s\n' "$PIG_UPDATE_SIGNING_KEY" > "$work/signing.pem"
unset PIG_UPDATE_SIGNING_KEY
openssl pkey -in "$work/signing.pem" -pubout -out "$work/update-trust.pem"
cmp "$work/update-trust.pem" automation/release/personal-update-trust.pem
trust=$(base64 < "$work/update-trust.pem" | tr -d '\n')
update_url="https://github.com/$REPOSITORY/releases/latest/download/update.json"
flags="-s -w -X github.com/MichaelKinsy/PiG/coding/cli.Build=$sha -X github.com/MichaelKinsy/PiG/internal/codingagent.DefaultUpdateURL=$update_url -X github.com/MichaelKinsy/PiG/internal/codingagent.DefaultUpdateTrustRoot=$trust"
mkdir -p "$work/out"
for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64; do
  os=${target%/*}
  arch=${target#*/}
  name="pig-$version-$os-$arch"
  mkdir "$work/$name"
  binary=pig
  if [[ "$os" = windows ]]; then binary=pig.exe; fi
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags "$flags" -o "$work/$name/$binary" ./cmd/pig
  if [[ "$target" = linux/amd64 ]]; then "$work/$name/pig" --version; fi
  cp LICENSE NOTICE THIRD_PARTY_NOTICES.md "$work/$name/"
  cp -R LICENSES "$work/$name/"
  cp "$work/update-trust.pem" "$work/$name/update-trust.pem"
  if [[ "$os" = windows ]]; then
    (cd "$work" && zip -qr "out/$name.zip" "$name")
  else
    tar -C "$work" -czf "$work/out/$name.tar.gz" "$name"
  fi
done
# Android requires the upstream NDK build path, never a Linux binary.
name="pig-$version-android-arm64"
mkdir "$work/$name"
GOOS=android GOARCH=arm64 automation/release/android-binary.sh build -trimpath -ldflags "$flags" -o "$work/$name/pig" ./cmd/pig
automation/release/android-binary.sh check "$work/$name/pig"
cp LICENSE NOTICE THIRD_PARTY_NOTICES.md "$work/$name/"
cp -R LICENSES "$work/$name/"
cp "$work/update-trust.pem" "$work/$name/update-trust.pem"
tar -C "$work" -czf "$work/out/$name.tar.gz" "$name"
cp docs/site/public/install.sh "$work/out/install.sh"
cp "$work/update-trust.pem" "$work/out/update-trust.pem"
(cd "$work/out" && sha256sum ./*.tar.gz ./*.zip install.sh update-trust.pem > SHA256SUMS)
python3 automation/release/gen-update-manifest.py --version "$version" \
  --base-url "https://github.com/$REPOSITORY/releases/download/v$version" \
  --sha256sums "$work/out/SHA256SUMS" > "$work/out/update.json"
bash automation/release/sign-update-manifest.sh "$work/out/update.json" "$work/signing.pem" \
  "$work/update-trust.pem" > "$work/out/update.json.sig"
# Stage all assets on a draft so clients never observe a partial release.
gh release create "v$version" --repo "$REPOSITORY" --target "$sha" --draft \
  --title "PiG $version" --notes "Personal PiG build. Upstream PiG: $base. Source: $sha

Configuration and executable names are unchanged. Linux amd64 smoke-tested; other targets cross-built, not natively tested."
gh release upload "v$version" "$work/out/"* --repo "$REPOSITORY"
gh release edit "v$version" --repo "$REPOSITORY" --draft=false --prerelease=false --latest
