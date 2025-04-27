# Rate Limiter Implementation

## Configuration and Environment Setup

### IP-Based Configuration

The configuration is stored in a file named `./env.json.example`:

```json
{
  "app": {
    "port": "8080"
  },
  "redis": {
    "db": 0,
    "host": "redis",
    "port": "6379"
  },
  "rate_limiter": {
    "by_ip": {
      "time_window": 1,
      "max_requests": 10,
      "blocked_duration": 60
    }
  }
}
```

Key parameters:
- `time_window`: Duration in **seconds** for tracking requests.
- `max_requests`: Maximum requests allowed within the `time_window`. For example, 10 requests per 1 second.
- `blocked_duration`: Duration in **seconds** that an IP is blocked after exceeding the limit.

## Launching the Application

`make run` to spin up applications.

### API Token Setup

To generate an API token, send the following HTTP request:

```http request
POST http://127.0.0.1:8080/api-key
Content-Type: application/json

{
  "time_window": 1,
  "max_requests": 10,
  "blocked_duration": 60
}
```

Or use cURL:

```shell
curl -s \
  -X POST \
  --header 'Content-Type: application/json' \
  --data '{
    "time_window": 1,
    "max_requests": 10,
    "blocked_duration": 60
  }' http://127.0.0.1:8080/api-key
```

Each token can have its own unique configuration.
You can use `.http` files on `api/` to test endpoints


## Redis Database Structure

The Redis database stores keys for rate limiting, visualized using [🔗 Redis Insight](https://redis.com/redis-enterprise/redis-insight/). The key types are:

### API Key Data

- `<apikey>`: Stores configuration for a specific API key. Example:
    - Key: `263b3bd80f0e7c9b4daceb4da8b2ef62d8c703b2f53bd924efdf79c3251a83ef`
    - Value:
      ```json
      {
        "max_requests": 10,
        "time_window": 1,
        "blocked_duration": 60
      }
      ```
- `rate:api-key_<apikey>`: Tracks request rates for the API key. Example:
    - Key: `rate:api-key_880d207159a7ac1a5a800eabbb310cf851e3c00cb5a2ff6e1ab9f38ce21bcc99`
    - Value:
      ```json
      {
        "max_requests": 10,
        "time_window_sec": 1,
        "requests": [
          1704787618,
          1704787619,
          1704787620,
          1704787621
        ]
      }
      ```
- `blocked:api-key_<apikey>`: Indicates a blocked API key with a set duration. Example:
    - Key: `blocked:api-key_880d207159a7ac1a5a800eabbb310cf851e3c00cb5a2ff6e1ab9f38ce21bcc99`
    - Value: `APIKeyBlocked`
    - Duration: Matches the `blocked_duration` from the API key config.

### IP Data

- `rate:ip_<ip>`: Tracks request rates for an IP. Example:
    - Key: `rate:ip_127.0.0.1`
    - Value:
      ```json
      {
        "max_requests": 10,
        "time_window_sec": 1,
        "requests": [
          1704787618,
          1704787619,
          1704787620,
          1704787621
        ]
      }
      ```
- `blocked:ip_<ip>`: Indicates a blocked IP with a set duration. Example:
    - Key: `blocked:ip_127.0.0.1`
    - Value: `IPBlocked`
    - Duration: Matches `rate_limiter.by_ip.blocked_duration` from `env.json`.

## Testing the Rate Limiter

To explore the testing CLI, run:

```shell
docker compose run --rm go-cli -h
```

Output:
```
Usage of ./cli-test:
  -k string
        API Key for the request
  -m string
        HTTP method to use (default "GET")
  -r int
        Maximum amount of requests to send (default 100)
  -t int
        Time in seconds of each request (default 1)
  -url string
        URL to test
```

### Testing with IP

Run the following command:

```shell
docker compose run --rm go-cli -url http://go-app:8080/hello-world -m GET -t 1 -r 10
```

### Testing with API Key

1. Create an API key by running:

```shell
curl -s \
  -X POST \
  --header 'Content-Type: application/json' \
  --data '{
    "time_window": 1,
    "max_requests": 10,
    "blocked_duration": 60
  }' http://127.0.0.1:8080/api-key
```

2. The response will include an API key, e.g.:
   ```json
   {
     "api_key": "c6f7363326f62f2483756447a963f2369a0dd5e90b7e8a36c32bc1a62ed38f51"
   }
   ```

3. Use the API key in the test command:

```shell
docker compose run \
  --rm \
  go-cli \
  -url http://go-app:8080/hello-world-key \
  -m GET \
  -t 1 \
  -r 10 \
  -k c6f7363326f62f2483756447a963f2369a0dd5e90b7e8a36c32bc1a62ed38f51
```