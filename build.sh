#!/bin/sh

# GITHUB_EVENT_PATH documented here:
# https://docs.github.com/en/actions/reference/environment-variables#default-environment-variables
GIT_TAG=$(jq .release.tag_name < "${GITHUB_EVENT_PATH}" | sed -e 's/"//g')
UPLOAD_URL=$(jq .release.upload_url < "${GITHUB_EVENT_PATH}" | sed -e 's/"//g' | cut -d "{" -f 1)
RELEASES="arm64-darwin-sb amd64-linux-sb amd64-darwin-sb"

upload_file() {
    NAME=$1

    zip "${NAME}.zip" "${NAME}"
    curl -H "Accept: application/vnd.github.v3+json" \
         -H "Authorization: Bearer ${GITHUB_TOKEN}" \
         -H "Content-Type: application/zip" \
         --data-binary "@${NAME}.zip" \
         "${UPLOAD_URL}?name=${NAME}.zip"
}

for PLATFORM in ${RELEASES}; do
    GOOS=$(echo "${PLATFORM}" | cut -d - -f 2) \
    GOARCH=$(echo "${PLATFORM}" | cut -d - -f 1) \
    go build -o "${PLATFORM}" -a -ldflags="-X 'sb/cmd.AppVersion=${GIT_TAG}'"

    if [ "${UPLOAD_URL}" != null ]; then
        upload_file "${PLATFORM}"
    fi
done
