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

# Ensure the latest remote branches are fetched.
git fetch origin

# Check if the branch exists locally. If not, create it from the remote.
if git rev-parse --verify "$BRANCH" >/dev/null 2>&1; then
  git checkout "$BRANCH"
else
  git checkout -B "$BRANCH" "origin/$BRANCH"
fi

# Update version.go in place without creating a backup.
sed -i "s/var Version = \"dev\"/var Version = \"$NEW_VERSION\"/" pkg/version/version.go

# Configure Git for committing.
git config user.email "release-bot@sailfin.io"
git config user.name "Release Bot"

# Stage the change and commit.
git add pkg/version/version.go
git commit -m "Bump version to $NEW_VERSION [skip ci]"

# Push the commit to the target branch.
git push origin "$BRANCH"
