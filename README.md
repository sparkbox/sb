# sb

A CLI for generating [SSH Certificates](https://engineering.fb.com/2016/09/12/security/scalable-and-secure-access-with-ssh/) via [Sign-in with Slack](https://api.slack.com/docs/sign-in-with-slack).

SSH Certificates are easier to manage than SSH keys primarily because Certificates can expire automatically.
This expiration means access to Sparkbox Slack is required to get SSH access to our various servers which strikes a good balance between security and maintenance overhead.

## How to use

1. Download `sb` from the latest release for your architecture.
2. Login by running `sb login`. This should launch a Sign-in with Slack prompt in your browser. Paste the resultant ID and token back to `sb`.
3. You're now set to generate a SSH Certificate, run `sb ssh`.
4. Verify your local `ssh-agent` has the cert by running `ssh-add -l`and noting the `ECDSA-CERT` entry.
5. You can now SSH to any host that is configured to trust the Certificate Authority.

## How it works

![](https://sparkbox.github.io/sb/flow.png)
