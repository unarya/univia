#!/bin/bash

module=$1
version=$2

git tag ${module}/v${version}
git push origin ${module}/v${version}

echo "✅ Tagged ${module} version v${version}"
