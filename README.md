# GoEnValue - a simple templates processor from environment variables.

This tool very useful if you need to configure your containerised software.
You can do it by passing environment variables into container on your Docker container start.

How it works
---
You need to:
 1. create templates from your configuration files
 2. add in your entrypoint script state with execution of goenvalue
 3. build image with your software, templates and goenvalue
 4. run container with environment variables

Example
---
**php.ini.tpl**
```ini
  ...
  memory_limit = {{ .MEMORY_LIMIT }}
  ...
```

**Dockerfile**
```Dockerfile
  FROM php:7
  ENV PHP_MEMORY_LIMIT="64M"
  COPY php.ini.tpl /etc/php/php.ini.tpl
  COPY entrypoint.sh /entrypoint.sh
  ...
  ENTRYPOINT /entrypoint.sh
```

**entrypoint.sh**
```bash
  #!/bin/bash
  goenvalue -p PHP -i /etc/php/php.ini.tpl
  exec $@
```

Test it:
```bash
  $ docker run --rm -it -e PHP_MEMORY_LIMIT=128M my/php:7 php -i | grep memory_limit
  memory_limit => 128M => 128M
```
