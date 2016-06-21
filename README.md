# sequel-open

my opening [Sequel Pro](http://www.sequelpro.com/) connection tool via command line.

## INSTALL

```
$ go get github.com/soh335/sequel-oepn
```

## USAGE

```
$ sequel-open -docker -host <container name> -ssh-password ... -ssh-user ... -user ... -password ...
```

## KNOWN ISSUE

* When close seuqle window, alert window is shown. Because temporary connection file is already deleted.

## LICENSE

* MIT
