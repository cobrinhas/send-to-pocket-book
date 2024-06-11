# send-to-pocket-book

A universal browser extension (Chrome/Firefox/Edge/Chromium) to send documents/articles/web pages to your PocketBook. It downloads the document from the browser active tab and sends an e-mail to your `@pbsync.com` account.

## Why?

I've recently acquired a PocketBook e-book reader ([PocketBook Verse HD]()), and despite these devices being known as cheap e-readers, they have a neat drop-in feature called **Send-To-PockeBook**. It allows users to remotely send a document/e-book to their devices by sending an e-mail to their `@pbsync.com` account with the document in attachement. The document is then downloaded and stored locally on the device.

The value proposition of this extension is allowing users to use this functionality directly from their active browser tab, without leaving the tab that has the document they would like to read on their PocketBook.

## How does it look like?

tbd.

## How do I install it?

tbd.

## Setup

You can rollout the extension from source by cloning this repository and self-hosting the proxy-server. For that you will need to:

- Follow the setup instructions of `proxy-server` and deploy it in a machine
- Follow the setup instructions of `proxy-server-worker` to hide `proxy-server` behind Cloudflare
- Copy the deployed Cloudflare worker URL and paste in the `proxy-server-url` input element ([index.html](extension/src/index.html))