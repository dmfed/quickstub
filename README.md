# quickstub - quick http server stub

Quickstub is a HTTP server which can be configured in a very simple manner by yaml file .yaml.

The program is intended mainly for mocking some API responses when testing your applications. 

## quickstart
Install th binary
```bash
go install github.com/dmfed/quickstub/cmd/quickstub@latest
```
Generate the config
```bash
quickstub -sample > myconfig.yaml
```
Edit the config as required (examples are included in the sample) then launch your server.
**Note**: quickstub will fail to run with generated exmple config, since it has "file:" field for one of endpoints.
```bash
quickstub -conf myconfig.yaml
```

