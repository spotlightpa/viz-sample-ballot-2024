## Spotlight PA District Finder

### Installation

- Install Go. See .go-version for minimum version.
- Install Hugo. See netlify.toml for minimum version.
- Run `yarn` to install JavaScript dependencies. See netlify.toml for minimum version.

## Usage

Get an API key for Google Maps. Export it as an environment variable:

```
export GEOLOCATOR_API_KEY=abcd1234
```

Launch geolocator API service:

```
go run ./cmd/geolocator -port 2064
```

In another terminal tab/window, launch hugo:

```
hugo server
```

Open district finder in a web browser: http://localhost:1313/

## Architecture

The homepage is the sample ballot lookup visualization. It is a Single Page Application. The HTML is built by Hugo and uses Alpine.JS to request data from geolocator, a serverless backend written in Go intended to be deployed on Netlify.

Geolocator relies on data generated using chunker (to split a map file into chunks) and extractor (to get demographic data out of a map file). See the old repo where they were used: https://github.com/spotlightpa/viz-redistricting-2020/ 
