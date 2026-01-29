# Rocket Map Component

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

Grab the `datastar-pro-rocket.js` file from the pro repo and drop it alongside the datastar bundle in `/web/resources/static/datastar/`

## Development Mode

```shell
go tool task live
```

OR

```shell
air -build.cmd "go build -tags=dev -o tmp/bin/main ./cmd/web" -build.bin "tmp/bin/main" -misc.clean_on_exit true -build.include_ext "go,templ"

# watch and rebuild web assets + hotreload
go run cmd/web/build/main.go -watch

# watch and rebuild templ components
go tool templ generate -watch
```

# Data Sources

- [land.json](https://github.com/martynafford/natural-earth-geojson/blob/master/110m/physical/ne_110m_land.json)
- all-members.json -> scraped from Discord API on 2026-01-28
