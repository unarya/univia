#!/bin/bash
set -e
set -o pipefail

# --------------------------------------------------------------------
# UNIVIA - RELEASE TAG SCRIPT
# Author: Tiecont
# Usage : ./release.sh <stage> <base_version> <pre_number>
# Example: ./release.sh alpha v0.0.2 1  → v0.0.2-alpha.1
# --------------------------------------------------------------------

STAGE=$1
BASE_VERSION=$2
PRE_NUM=$3

# === Validation ===
if [ -z "$STAGE" ] || [ -z "$BASE_VERSION" ] || [ -z "$PRE_NUM" ]; then
    echo -e "\033[1;31m[ERROR]\033[0m Usage: $0 <stage> <base_version> <pre_number>"
    echo "Example: $0 alpha v0.0.2 1"
    exit 1
fi

if ! [[ "$BASE_VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "\033[1;31m[ERROR]\033[0m Invalid base version format: $BASE_VERSION"
    echo "Expected format: vMAJOR.MINOR.PATCH (e.g. v0.0.2)"
    exit 1
fi

if ! [[ "$PRE_NUM" =~ ^[0-9]+$ ]]; then
    echo -e "\033[1;31m[ERROR]\033[0m Invalid pre-release number: $PRE_NUM"
    exit 1
fi

# === Construct tag ===
TAG="${BASE_VERSION}-${STAGE}.${PRE_NUM}"

# === Utilities ===
log() {
    echo -e "\033[1;34m[INFO]\033[0m $1"
}
success() {
    echo -e "\033[1;32m[SUCCESS]\033[0m $1"
}

# === Main ===
log "Preparing release tag: ${TAG}"

# Commit current changes if any
if ! git diff --quiet || ! git diff --cached --quiet; then
    git add .
    git commit -m "release ${TAG}" || log "No new changes to commit."
else
    log "No changes to commit."
fi

# Delete old tag if exists (replace behavior)
if git rev-parse "${TAG}" >/dev/null 2>&1; then
    log "Tag ${TAG} already exists. Replacing..."
    git tag -d "${TAG}" >/dev/null 2>&1 || true
    git push origin ":refs/tags/${TAG}" >/dev/null 2>&1 || true
fi

# Create and push new tag
git tag "${TAG}"
git push origin "${TAG}"
success "Released ${TAG} successfully!"

# --------------------------------------------------------------------
# Update go.mod dependencies to new version
# --------------------------------------------------------------------
log "Updating internal go.mod dependencies to ${TAG} ..."

# Pattern: replace any line like
#   github.com/deva-labs/univia v0.0.2-alpha.3
#   → github.com/deva-labs/univia v0.0.2-alpha.4
for modfile in $(find cmd -type f -name "go.mod"); do
    if grep -q "github.com/deva-labs/univia" "$modfile"; then
        log "Updating $modfile ..."
        # Try replacing version
        if ! sed -i.bak -E "s#(github\.com/deva-labs/univia)[[:space:]]+v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z]+\.[0-9]+)?#\1 ${TAG}#g" "$modfile"; then
            error "Failed to update version in $modfile"
            exit 1
        fi
        rm -f "${modfile}.bak"

        # Run go mod tidy and catch any error explicitly
        if ! (cd "$(dirname "$modfile")" && go mod tidy >/dev/null 2>&1); then
            error "go mod tidy failed for $(dirname "$modfile")"
            exit 1
        fi
    fi
done

# Commit dependency updates
if ! git diff --quiet || ! git diff --cached --quiet; then
    git add cmd/**/go.mod cmd/**/go.sum
    git commit -m "chore: bump univia module to ${TAG}"
    git push origin HEAD
    success "Updated go.mod files to ${TAG}"
else
    log "No dependency changes to commit."
fi

# --------------------------------------------------------------------
echo "---------------------------------------------------"
echo "  • Stage        : ${STAGE}"
echo "  • Base Version : ${BASE_VERSION}"
echo "  • Tag Created  : ${TAG}"
echo "---------------------------------------------------"
