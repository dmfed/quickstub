# quickstub - quick http server stub

Quickstub is a HTTP server which can be configured in a very simple manner with a .yaml file.

The program is intended mainly for mocking some API responses when testing your applications.

In addition to primary configuration quickstab may be reconfigured on the fly using the "magic endpoint" (see below).

## Quickstart
Install the binary
```bash
go install github.com/dmfed/quickstub/cmd/quickstub@latest
```
Generate the config
```bash
quickstub -sample > myconfig.yaml
```
Edit the config as required (examples are included in the sample) then launch your server.

```bash
quickstub -conf myconfig.yaml
```

## Configuring the server

Let's take a look at the sample configuration file.
```yaml 
version: 2
listen_addr: ':8080'
magic_endpoint: '/magic'
endpoints:
  # simple response with text body
  'GET /hello':
    code: 200
    body: "hello world"

  # responds with 201 code and text body
  'POST /hello':
    code: 201
    body: "created"

  # responds with 200 code a header and JSON body
  'GET /hello':
    code: 200
    headers:
        'Content-Type': 'application/json'
    body: '{"hello": "world"}'

  # redirects to /hello endpoint on the same host
  'GET /redirect':
    code: 301
    headers: 
      'Connection': 'keep-alive'
      'Location': '/hello'

  # responds with 400 code 
  'GET /badrequest':
    code: 400
    body: 'These are not the droids you are looking for!'
```
The first part configures the server itself.
```yaml
version: 2
listen_addr: ':8080'
magic_endpoint: '/magic'
```
**version** is the version of API.

**listen_addr** tells the server the hostname and port to listen on. It may take the following forms: "172.0.0.1:8080" or just ":8080". This field is reauired and must not be empty.

**magic_endpoint** is the path of endpoint where server accepts reconfigure requests (see below).

The **endpoints** part of the config is a configuration of responses of the server.

### Plain text response
```yaml
'GET /hello':
    code: 200
    body: 'hello world'
```
The above tells quickstub to repond with 200 and "hello world" text on endpoint "/hello".

### JSON response
```yaml 
'GET /hello':
    code: 200
    headers:
        'Content-Type': 'application/json'
    body: '{"hello": "world"}'
```
This will instruct quickstub to repon with JSON in the body.

### Response with file.
```yaml 
'GET /config':
    code: 200
    headers: 
      'Content-Type': 'application/octet-stream'
      'Content-Disposition': 'attachment; filename=myconfig.yaml'
    body: '@myconfig.yaml'
```
This tells the server to accept GET requests to "/config" endpoint and respond with provided headers and contents of file "myconfig.yaml" in response body. 

Note the "@" character. It tells qucikstub that the remaining part of the string is a local path to look for file. It may take form of "@/home/myuser/somefile" etc. If you want a string starting with "@" in response body, just escape it like this: "\\@myresponse body" for double-quoted strings and like this '\@myresponse body' for single-quoted string. (Note that yaml distinguishes between double and single quotes for strings). 

**Note:** the actual contents of the file is preloaded into RAM on server start (or reconfiguration) so be midfull of what files to use. This program is just an http stub it is not intended to actually server huge files etc. So JSON is OK while ISO image is not. 

## Reconfigure endpoint or "magic endpoint"
If "magic_endpoint" is not empty in the config file then the server will listen to GET and PATCH requests on this endpoint. 

```yaml
version: 2
listen_addr: ':8080'
magic_endpoint: '/magic'
```
In the above example doing
```bash 
curl -X GET localhost:8080/magic
```
will return the current configuration of the server.

To reconfigure the server do.
```bash
curl -X POST --data-binary '@my_new_config.yaml' localhost:8080/magic
```
This will result in pushing your file "my_new_config.yaml" to the server. The server will validate the new config config and respond with 201 OK. Then quickstub will be restarted with new parameters.

If any validation errors occur quickstub will respond with 400 and error text in response body.
