#!/bin/sh

# GITHUB_EVENT_PATH documented here:
# https://docs.github.com/en/actions/reference/environment-variables#default-environment-variables
GIT_TAG=$(jq .release.tag_name < "${GITHUB_EVENT_PATH}" | sed -e 's/"//g')
UPLOAD_URL=$(jq .release.upload_url < "${GITHUB_EVENT_PATH}" | sed -e 's/"//g' | cut -d "{" -f 1)
CERT_FILE="${HOME}/developer_id_certificate.p12"
RELEASES="arm64-darwin-sb amd64-linux-sb amd64-darwin-sb"

upload_file() {
    NAME=$1

    if [ "${NAME}" = "amd64-linux-sb" ]; then
        zip "${NAME}.zip" "${NAME}"
        NAME="${NAME}.zip"
        CONTENT_TYPE="Content-Type: application/zip"
    else
        NAME="${NAME}.dmg"
        CONTENT_TYPE="Content-Type: application/octet-stream"
    fi

    curl -H "Accept: application/vnd.github.v3+json" \
         -H "Authorization: Bearer ${GITHUB_TOKEN}" \
         -H "${CONTENT_TYPE}" \
         --data-binary "@${NAME}" \
         "${UPLOAD_URL}?name=${NAME}"
}

setup_keychain() {
  echo "${APPLE_DEVELOPER_CERTIFICATE_P12_BASE64}" | base64 --decode > "${CERT_FILE}"
  EPHEMERAL_KEYCHAIN="ci-ephemeral-keychain"
  EPHEMERAL_KEYCHAIN_PASSWORD="$(openssl rand -base64 100)"
  security create-keychain -p "${EPHEMERAL_KEYCHAIN_PASSWORD}" "${EPHEMERAL_KEYCHAIN}"
  EPHEMERAL_KEYCHAIN_FULL_PATH="${HOME}/Library/Keychains/${EPHEMERAL_KEYCHAIN}-db"
  security import "${CERT_FILE}" -k "${EPHEMERAL_KEYCHAIN_FULL_PATH}" -P "${APPLE_DEVELOPER_CERTIFICATE_PASSWORD}" -T "$(command -v codesign)"
  security set-key-partition-list -S "apple-tool:,apple:" -s -k "${EPHEMERAL_KEYCHAIN_PASSWORD}" "${EPHEMERAL_KEYCHAIN_FULL_PATH}"
  security default-keychain -d "user" -s "${EPHEMERAL_KEYCHAIN_FULL_PATH}"
}

sign() {
  PLATFORM=$1
  ./gon -log-json -log-level=info "./${PLATFORM}-gon-config.json"
}

setup_keychain
curl -LO "https://github.com/mitchellh/gon/releases/download/v0.2.3/gon_macos.zip"
unzip ./gon_macos.zip

for PLATFORM in ${RELEASES}; do
    GOOS=$(echo "${PLATFORM}" | cut -d - -f 2) \
    GOARCH=$(echo "${PLATFORM}" | cut -d - -f 1) \
    go build -o "${PLATFORM}" -a -ldflags="-X 'sb/cmd.AppVersion=${GIT_TAG}'"

    if [ "${PLATFORM}" != "amd64-linux-sb" ]; then
        sign "${PLATFORM}"
    fi

    if [ "${UPLOAD_URL}" != null ]; then
        upload_file "${PLATFORM}"
    fi
done
