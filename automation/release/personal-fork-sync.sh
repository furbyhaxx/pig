#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
# Merges upstream PiG main into the checked-out fork main and pushes it (.project/specs/0001-personal-distribution.md D2).
set -euo pipefail

: "${GITHUB_OUTPUT:?}"
upstream_url=${PIG_UPSTREAM_URL:-https://github.com/MichaelKinsy/PiG.git}
git config user.name 'github-actions[bot]'
git config user.email '41898282+github-actions[bot]@users.noreply.github.com'
# .gitattributes routes fork-owned files to this driver; it keeps the fork side.
git config merge.ours.driver true
git remote add upstream "$upstream_url" 2>/dev/null || git remote set-url upstream "$upstream_url"
git fetch --no-tags upstream main

if git merge-base --is-ancestor upstream/main HEAD; then
  echo "Upstream main is already merged."
else
  base=$(git merge-base HEAD upstream/main)
  superseded=$(git diff --name-only "$base" upstream/main -- ':(glob)**/AGENTS.md' ':(exclude,glob)**/testdata/**')
  message="chore: merge PiG upstream"
  if [[ -n "$superseded" ]]; then
    list=$(sed 's/^/- /' <<< "$superseded")
    message+=$'\n\nUpstream changed these AGENTS.md files; the fork versions were kept. Review them for rules to port:\n'"$list"
    echo "::warning title=Upstream AGENTS.md changes::Review superseded upstream AGENTS.md changes listed in the merge commit."
    printf '### Upstream AGENTS.md changes to review\n\n%s\n' "$list" >> "${GITHUB_STEP_SUMMARY:-/dev/null}"
  fi
  if ! git merge --no-ff --no-edit -m "$message" upstream/main; then
    git diff --name-only --diff-filter=U | sed 's/^/::error title=Upstream merge conflict::/'
    git merge --abort
    echo "Resolve the upstream merge locally and push it; nothing was published." >&2
    exit 1
  fi
  git push origin HEAD:main
fi
printf 'sha=%s\n' "$(git rev-parse HEAD)" >> "$GITHUB_OUTPUT"
