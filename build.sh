#!/bin/sh

# GOOS=darwin GOARCH=arm64 go build -o arm64-macos-sb -a -ldflags="-X 'sb/cmd.AppVersion=$GIT_TAG'"

# GOOS=darwin GOARCH=amd64 go build -o amd64-macos-sb -a -ldflags="-X 'sb/cmd.AppVersion=$GIT_TAG'"

# GOOS=linux GOARCH=amd64 go build -o amd64-linux-sb -a -ldflags="-X 'sb/cmd.AppVersion=$GIT_TAG'"

# zip arm64-macos-sb.zip arm64-macos-sb

# zip amd64-macos-sb.zip amd64-macos-sb

# zip amd64-linux-sb.zip amd64-linux-sb

# curl -d "https://api.github.com/repos/sparkbox/sb/releases/{release_id}/assets"

cat "${GITHUB_EVENT_PATH}"
