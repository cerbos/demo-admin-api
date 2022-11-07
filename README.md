# Admin API demo

A demo implementation of a Golang backend with APIs for creating, updating and validating policies using the [Admin API](https://docs.cerbos.dev/cerbos/latest/api/admin_api.html). The demo also provides a sample client (built in React) which utilises the API.

The client provides basic policy creation, as well as an editing interface (via the [Monaco editor](https://microsoft.github.io/monaco-editor/)) which continually validates the policy and allows updates.

## Dependencies

- Docker for running the [Cerbos Policy Decision Point (PDP)](https://docs.cerbos.dev/cerbos/latest/installation/container.html), or both services via `docker-compose`
- Golang 1.19 (required if running manually)
- Node.js (required if running manually)

## Getting started

### Manually

1. Start up the Cerbos PDP instance docker container. This will be called by the Go app to check authorization.

```sh
cd cerbos
./start.sh
```

2. Build the React front end

```sh
# from project root
cd client
npm i
npm run build
```

3. Start the Go server

```sh
# from project root
go run main.go
```

4. Open the web app at `localhost:8090`

### Docker-compose

```sh
docker-compose up -d
```

## Audit logs

A separate endpoint is provided demoing the retrieval of audit logs from the PDP:

```sh
curl "http://localhost:8090/auditlog"
```

## Other considerations

- Similar APIs are provided by the AdminClient which allow management of Schemas.
- The Admin API is under heavy development, and might include breaking changes in future releases. This demo represents it's current state.
