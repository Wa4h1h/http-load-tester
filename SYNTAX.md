# Syntax
```yaml
concurrency: integer
timeout: integer
base: string
iterations: integer
schema:
  requests:
    - name: string
      url: string
      timeout: integer
      method: string
      headers:
        - string:string
      body: "{\"key\":\"value\"}"
```
##### Global settings
`concurrency`: Number of concurrent iterations. (default: 1)\
`timeout`: Global timeout: Applied to requests that don't have their own timeout. (default 1s)\
`base`: Base URL for all schema requests. Each request URL is appended to the base URL\
`iterations`: Number of loops to do (default: 1)


##### Schema settings
`schema`: Wrapper over requests: List of requests\
`requests`:\
&nbsp;&nbsp;`name`: request name\
&nbsp;&nbsp;`url`: Url to call\
&nbsp;&nbsp;`timeout (Optional)`: Local timeout. Only this request must be executed within the timeout\
&nbsp;&nbsp;`method`: Http method (default GET)\
&nbsp;&nbsp;`headers`: list key:value pairs\
&nbsp;&nbsp;`body`: json string body


