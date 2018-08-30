# Go Vanity Import

Go Vanity Import includes web and backend services that allows you to set custom import paths for your Go packages.

## Install

```
$ go get github.com/RivenZoo/govanityimport
$ go install github.com/RivenZoo/govanityimport/web
```

## Deploy

```
# run backend rpc
$ $(GOPATH)/bin/govanityimport serve --config conf/stage/.govanityimport.yaml

# run web service, then bind domain to it's address
$ $(GOPATH)/bin/web --config conf/stage/.web.yaml
```

## TODO

* add test
* add manage tool