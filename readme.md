# Rewrite Header
Rewrite Headers is a middleware plugin for [Traefik](https://traefik.io) which replaces HTTP response headers using regular expressions.

Based on: https://github.com/vincentinttsh/rewriteheaders

## Configuration

### Static

```yaml
pilot:
  token: "xxxx"

experimental:
  plugins:
    rewriteHeaders:
      modulename: "github.com/virtualzone/rewriteheaders"
      version: "v0.1.0"
```

### Dynamic

To configure the Rewrite Headers plugin, create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). 
The following example creates and uses the rewriteHeaders middleware plugin to modify the Location HTTP response header.

```yaml
http:
  routes:
    my-router:
      rule: "Host(`localhost`)"
      service: "my-service"
      middlewares : 
        - "rewriteHeaders"
  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1"
  middlewares:
    rewriteHeaders:
      plugin:
        rewriteHeaders:
          rewrites:
            - header: "Location"
              regex: "^http://(.+)$"
              replacement: "https://$1"
```
