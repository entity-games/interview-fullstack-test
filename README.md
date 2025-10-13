# Entity Games: Technical Interview Test Project

## Setup

### Prerequisites
1. Install [Go](https://go.dev/doc/install)
2. Install [Go Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
3. Install [Docker](https://docs.docker.com/engine/install/)
4. Install [Node.js & npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)

### Ensure symlinks are present
* static/js -> dist/js
* static/css -> dist/css

### Build and run
1. Install go dependencies with `go mod download`
2. Install npm dependencies with `npm install`
3. Perform database migrations with `make migrate`
4. Start the stack and watch for changes with `make start_all`

## Access
The application should be available at http://localhost:4444

The test user is `fixture_user` with password `fixturepass`