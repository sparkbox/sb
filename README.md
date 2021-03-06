# sb

A CLI for generating [SSH Certificates](https://engineering.fb.com/2016/09/12/security/scalable-and-secure-access-with-ssh/) via [Sign-in with Slack](https://api.slack.com/docs/sign-in-with-slack).

SSH Certificates are easier to manage than SSH keys primarily because Certificates can expire automatically.
This expiration means access to Sparkbox Slack is required to get SSH access to our various servers which strikes a good balance between security and maintenance overhead.

## Installation

### macOS

1. `brew tap sparkbox/brew`
1. `brew install sparkbox/brew/sb`

### Linux

1. Download `sb` from the [latest release](https://github.com/sparkbox/sb/releases)
1. Unzip the file
1. Move the `sb` binary to a location your `$PATH` understands: e.g. `mv sb /usr/local/bin/sb`
1. Start a fresh shell instance (new Terminal window)


## How to use

1. Login by running `sb login`. This should launch a Sign-in with Slack prompt in your browser. Paste the resultant ID and token back to `sb`.
1. Run `sb ssh` to generate a new, time limited SSH certificate.
1. Run `ssh-add -l` to verify your local `ssh-agent` has the cert by locating the `ECDSA-CERT` entry.
1. You can now SSH to any host that is configured to trust the Certificate Authority.

## How it works

![](https://sparkbox.github.io/sb/flow.png)

## Helpful Notes

* If you are using an Intel based machine, use the AMD64 file.
* If you are using an M1 MacBook, use the AMR64 file.
