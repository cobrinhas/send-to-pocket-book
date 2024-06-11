# proxy-server-worker

Worker that acts as a reverse proxy for proxy-server

---

This worker was bootstrapped using [worker brick](https://github.com/dart-pacotes/.brick) and configured with [wrangler](https://github.com/cloudflare/wrangler) CLI. You can install it via NPM: `npm install -g wrangler`

## Hooks

This repository is configured with client-side Git hooks that automatically format + lint the codebase before each push. You can install it by running the following command:

```bash
./hooks/INSTALL
```

## Development

To properly run this service, you will need to a set up a `.dev.vars` file. Start by creating a copy of the `.dev.vars.tpl` file and fill the variables with values appropriate for the execution context.

|  Variable Name   | Variable Description |
| :--------------: | :------------------- |
| `PROXY_SERVER_URL` | This is the URL of the proxy server that sends the e-mails.      |
| `PROXY_SERVER_SECURITY_TOKEN` | This is the security token set on the proxy-server to guard API usage.      |

Run the local server via `npm run start`

## Deploy

Setup worker environment variables with the following one-line:

```bash
IFS='='; ENV_FILE=.dev.vars; cat $ENV_FILE | while read line || [[ -n $line ]]; do read -ra envy <<< $line; wrangler secret put ${envy[0]} <<< ${envy[1]} ; done
```

Deploy to Cloudflare via `npm run deploy`

### Contact

This template was prepared by:

- João Freitas, @freitzzz
- Rute Santos, @rutesantos4

Contact us if you need help on your project!
