#!/usr/bin/env bash
set -euo pipefail

# This script updates pkg/version/version.go with the new version,
# commits the change, and pushes it to the appropriate branch.
#
# Usage: ./scripts/bump_version.sh <new_version> [branch]
#
# If the branch argument is omitted, the script will attempt to extract
# a prerelease branch name from the version string (e.g. for "v1.0.0-alpha.1",
# it will use "alpha"). If no prerelease segment is found, it defaults to "main".

if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <new_version> [branch]"
  exit 1
fi

NEW_VERSION="$1"

if [ "$#" -ge 2 ]; then
  BRANCH="$2"
else
  if [[ "$NEW_VERSION" =~ -([a-zA-Z]+)\. ]]; then
    BRANCH="${BASH_REMATCH[1]}"
  else
    BRANCH="main"
  fi
fi

echo "New version: $NEW_VERSION"
echo "Branch to update: $BRANCH"

# Fetch remote branches to ensure we have the latest.
git fetch origin

# Checkout the branch; if it doesn't exist locally, create it from origin.
if git rev-parse --verify "$BRANCH" >/dev/null 2>&1; then
  git checkout "$BRANCH"
else
  git checkout -B "$BRANCH" "origin/$BRANCH"
fi

# Stash any unstaged changes if present.
if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Working directory is not clean. Stashing changes..."
  git stash push --include-untracked -m "Auto-stash before bump version"
  STASHED=true
fi

# Rebase the branch onto the latest remote changes.
git pull --rebase origin "$BRANCH"

# Update pkg/version/version.go in place.
sed -i "s/var Version = \".*\"/var Version = \"$NEW_VERSION\"/" pkg/version/version.go

# Stage the change.
git add pkg/version/version.go

# Check if there are any staged changes.
if git diff --cached --quiet; then
  echo "No changes to commit in pkg/version/version.go."
  exit 0
fi

# Configure Git for committing.
git config user.email "release-bot@sailfin.io"
git config user.name "Release Bot"

# Commit the change.
git commit -m "Bump version to $NEW_VERSION [skip ci]"

# Push the commit.
git push origin "$BRANCH"

# ... after pushing, restore stashed changes if any.
if [ "${STASHED:-false}" = true ]; then
  echo "Restoring stashed changes..."
  git stash pop || echo "Warning: There were conflicts when restoring stashed changes."
fi
