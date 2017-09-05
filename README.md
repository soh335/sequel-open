[![Build Status](https://travis-ci.org/soh335/sequel-open.svg?branch=master)](https://travis-ci.org/soh335/sequel-open)

# sequel-open

my opening [Sequel Pro](http://www.sequelpro.com/) connection tool via command line.

## INSTALL

```
$ go get github.com/soh335/sequel-open
```

## USAGE

```
$ sequel-open -docker -host <container name> -ssh-password ... -ssh-user ... -user ... -password ...
```

if `-docker` option is given, overwrite `-ssh-host` form `DOCKER_HOST` and `-host` via `docker inspect` command.

## TODO

* support normal, socket type.

## KNOWN ISSUE

* When close seuqle window, alert window is shown. Because temporary connection file is already deleted.

## LICENSE

* MIT
