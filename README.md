# Rocket Template

# Usage

## Setup

1. Go
```
go mod tidy
```

2. Web Dependencies

```
go run cmd/web/build/main.go
```

3. Datastar Pro + Rocket

Grab the `datastar-pro-rocket.js` file and drop it into `/web/resources/static/datastar/`

## Development Mode

```shell
task live
```

OR

```shell
air -build.cmd "go build -tags=dev -o tmp/bin/main ./cmd/web" -build.bin "tmp/bin/main" -misc.clean_on_exit true -build.include_ext "go,templ"

# watch and rebuild web assets + hotreload
go run cmd/web/build/main.go -watch

# watch and rebuild templ components
go tool templ generate -watch
````