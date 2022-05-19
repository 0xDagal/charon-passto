# Charon Passto

## What is it ?

Charon Passto is a Traefik plugin providing a middleware and it is part of the
Charon Project.

### What is the Charon Project

The Charon Project aim to protect your SPF from DOS/DDOS attacks.

### What is this middleware doing

It permit Traefik to interact with the ES database in the Charon Project.

## Configuration

To configure this plugin you should add its configuration to the Traefik dynamic
configuration as explained here. The following snippet shows how to configure
this plugin with the File provider in TOML and YAML:

Static:

```toml
[pilot]
  token="xxx"

[experimental.plugins.cache]
  modulename = "github.com/0xDagal/charon-passto"
  version = "v0.1.18"
```

```yaml
pilot:
  token: "xxx"

experimental:
  plugins:
    charon-passto:
      moduleName: "github.com/0xDagal/charon-passto"
      version: "v0.1.8"
```

Dynamic:

```toml
[http.middlewares]
  [http.middlewares.my-passto.plugin.charon-passto]
    es-address = "http://your-elastic-sever:port"
```

```yaml
http:
  middlewares:
    my-passto:
      plugin:
        charon-passto:
          es-address: 'http://your-elastic-sever:port'

```

### Options

#### Elastic Search Address (`es-address`)

The address to the elastic server passto will sends the information to.
