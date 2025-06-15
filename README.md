# webhook-proxy

A simple webhook proxy written in Go.

## Usage

1. Create a `config.yaml` file with the following structure:

```yaml
services:
  - name: service1
    urls:
      - https://www.postb.in/b/1749991958464-2974628300871
```

2. Run the proxy with the following command:

```bash
go run main.go
```

3. Send a POST request to the proxy with the following headers:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"message": "Hello, world!"}' http://localhost:8080/webhook
```

4. The request will be forwarded to the specified URLs. (Service name will be stripped from the request path.)


## Configuration

The proxy can be configured using a YAML file. The following configuration options are available:

- `services`: An array of service objects, each representing a target service.
  - `name`: The name of the service.
  - `urls`: An array of URLs to forward requests to.

## License

This project is licensed under the MIT License.