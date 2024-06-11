# send-to-pocket-book

A universal browser extension (Chrome/Firefox/Edge/Chromium) to send documents/articles/web pages to your PocketBook. It downloads the document from the browser active tab and sends an e-mail to your `@pbsync.com` account.

## Setup

You can rollout the extension from source by cloning this repository and self-hosting the proxy-server. For that you will need to:

- Follow the setup instructions of `proxy-server` and deploy it in a machine
- Follow the setup instructions of `proxy-server-worker` to hide `proxy-server` behind Cloudflare
- Copy the deployed Cloudflare worker URL and paste in the `proxy-server-url` input element ([index.html](extension/src/index.html))