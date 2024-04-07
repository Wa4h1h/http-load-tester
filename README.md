# http-load-tester
CLI tool that can make a high volume of concurrent HTTP(S) request against an website/HTTP based API

### install
```bash
go install github.com/Wa4h1h/http-load-tester@latest
```

OR

### generate build
```bash
make
```


## Usage
```
http-load-tester --help

Usage: http-load-tester [<command>] [options]
Commands:
bulk    perform http load test on different urls from a file

Use http-load-tester <command> -h or --help for more information about a command.
Options:
  -H value
        http headers
  -b string
        http request body
  -c int
        number of concurrent requests (default 1)
  -i int
        number of request iterations (default 1)
  -m string
        http method (default "GET")
  -timeout float
        number of seconds before a http request times out (default 0)
  -url string
        url to call

```
The cli tool provides two ways to load test an HTTP-API. 
Send a request exactly to one endpoint:
```bash
http-load-tester -url=http://eu.httpbin.org/get -i=10 -c=5
```

You can also provide multiple header
```bash
http-load-tester -url=http://eu.httpbin.org/get -H=k1:v1 -H=k2:v2 -i=10 -c=10
```
For now only json body is supported.\
POST request example:
```bash
http-load-tester -url="http://eu.httpbin.org/post" -m="POST" -b="{key:value}" -H=k1:v1 -H=k2:v2 -i=10 -c=10
```

Bulk usage:
```bash
http-load-tester bulk -h
Options:
  -f string
        path to yaml file containing the urls configurations
```
```bash
http-load-tester bulk -f=request.yaml
```

## Example urls configuration
For more details see SYNTAX.md
```yaml
concurrency: 4
timeout: 5
iterations: 10
schema:
  requests:
    - name: test
      url: http://eu.httpbin.org/get
      method: GET
      headers:
        - X-Forward-For:10.3.4.5
    - name: post
      url: http://eu.httpbin.org/post
      method: POST
      body: >
            {
              "key":"value"
            }
```
## Output example
```bash
http://eu.httpbin.org/post 200 OK 228ms
http://eu.httpbin.org/post 200 OK 229ms
http://eu.httpbin.org/post 200 OK 231ms
http://eu.httpbin.org/post 200 OK 232ms
http://eu.httpbin.org/post 200 OK 235ms
http://eu.httpbin.org/post 200 OK 237ms
http://eu.httpbin.org/post 200 OK 240ms
http://eu.httpbin.org/post 200 OK 242ms
http://eu.httpbin.org/post 200 OK 340ms
http://eu.httpbin.org/post 200 OK 380ms

Concurrency: 10
Total time: 0.382s
Total sent requests: 10
Received responses:
  ..............1xx: 0
  ..............2xx: 10
  ..............3xx: 0
  ..............4xx: 0
  ..............5xx: 0
  ..............Timed out: 0
Total requests failed to send: 0
Request per second: 10.000
(Min, Max, Avg) Request time: 228ms, 380ms, 259.40ms
```

## Note
Please be careful not to load test a website that you don’t own/have permission to do so - it will look like, and could become, a denial of service attack!