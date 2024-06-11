# proxy-server

Plain simple api to send documents to PocketBook via Send-To-PocketBook e-mail service

## Hooks

This repository is configured with client-side Git hooks which you need to install by running the following command:

```bash
./hooks/INSTALL
```

## Usage

To properly run this service, you will need to a setup a `.env` file. Start by creating a copy of the `.env.tpl` file and fill the variables with values appropriate for the execution context. Additionally, you will need to activate the GMail API on your account and create a consent screen. After that, download the `credentials.json` file and place it on this folder root.

Finally, you need to run the regen-token script which will create a credentials `token.json` file, enabling access to GMail API from the account you selected to send e-mails from.

```bash
go run cmd/regen-token/regen-token.go
```

Then, all you need to do is to run the service with the following command:

```bash
go run cmd/proxy-server/proxy-server.go
```

## Docker

To build the service image:

```bash
docker_tag=proxy-server:latest

docker build \
    -f ./deployments/Dockerfile \

    . -t $docker_tag
```

To run the service container:

```bash
export $(grep -v '^#' .env | xargs)

docker run \
    -p $server_port:$server_port \
    --mount "type=bind,src=$server_tls_crt_fp,dst=$server_tls_crt_fp" \
    --mount "type=bind,src=$server_tls_key_fp,dst=$server_tls_key_fp" \
    --mount "type=bind,src=$logging_fp,dst=$logging_fp" \
    --mount "type=bind,src=credentials.json,dst=credentials.json" \
    --mount "type=bind,src=token.json,dst=token.json" \
    --env-file .env \
    -t $docker_tag
```

### Contact

This template was prepared by:

- João Freitas, @freitzzz
- Rute Santos, @rutesantos4

Contact us for freelancing work!