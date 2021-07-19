# Minecraft Manager

So far, it's meant to be a Discord bot to:
* start/stop an azure vm (running Minecraft) through a user account (ie. without AAD access)
* switch between Minecraft worlds running on the server

# Dev Setup
Follow this to set up a dev environment:

## Prerequisites
Recent versions of:
* Azure CLI
* Go

## Config

Fill values in `config.yaml` or set them as environment variables.

To find the environment variables names - uppercase the yaml keys, join them with `_` and add a `MC_` prefix.

For example the environment variable for
```yaml
bot:
    token:
```
would be `MC_BOT_TOKEN`

## Run (dev)

Run the bot by running:
```shell
make run
```

## Build
Build a binary by running:
```shell
make
```
This creates `bin/mcbot`