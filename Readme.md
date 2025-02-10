# XM Company API  
This API provides full CRUD operations for a company.  

## Build  
Several build modes are available:  

### Go Binary  
You can create the executable with `make build` or `go build -a -mod=vendor -ldflags '-linkmode external -w -s -extldflags "-static"' -o ./xm-company`.  

### Docker Image  
You can build a Docker image for the API with `docker build -t xm-company:<your tag> .`  

**NOTE:** For all the build methods mentioned, you can refer to the `Makefile`, which contains basic build operations.  

## Running the API  
This API can be run in several ways.  

### Locally  
You can run the project locally after building the binary with `XM_COMPANY_CONFIG=./config ./<binary name>` or directly with `go run main.go`. In both cases, the provided configuration file (`./config/config.yaml`) is used.  

### Docker Image  
You can use the Docker image built earlier with `docker run -e XM_COMPANY_CONFIG='.' xm-company:<your tag>`  
The default configuration file is used. If you want to modify the configuration, refer to the `Configuration` section.  

### Docker Compose  
You can use the provided Docker Compose files to deploy the API and its services in `dev` or `staging` environments.  

- For `dev`: `make deploy_dev`  
- For `staging`: `make deploy_staging`  

**NOTE:** For all the deployment methods mentioned, you can refer to the `Makefile`, which contains basic deployment operations.  

## Configuration  
This API is configured using a configuration file. The **directory** containing this configuration file must be specified in the environment variable `XM_COMPANY_CONFIG`.  

An example is provided in `./config/config.yaml`.  

The format is as follows:  
```yaml
datasources:
  mongo:
    db: <name of the database>
    host: <host of mongo>
    port: <port of mongo>
    username: <username>
    password: <password>
    timeout: <number of seconds before connection timeout>

router:
  addr: <host of the API>
  port: <port of the API. If you are using Docker for build and/or deployment, use the port specified in the Dockerfile (8080)>
  gin_mode: "<gin mode: debug|release>"
  shutdown_timeout: <number of seconds before shutdown timeout>
```  

### Environment Variable Overrides  
Each configuration element can be overridden by an environment variable.  

The variable names are dynamic and follow the YAML file hierarchy.  

**Example:**  

To override the value of `datasources.mongo.db`, the environment variable is `XM_COMPANY_DATASOURCES_MONGO_DB`.  
To override the value of `router.gin_mode`, the environment variable is `XM_COMPANY_ROUTER_GIN_MODE`.  

## Project Structure  
The project architecture follows the principles of screaming architecture. Its goal is to isolate each feature individually.  

The main file reads the configuration, launches the application, and is ready to handle interruption signals.  

The `app` package uses the configuration and starts the API.  

Each `features` package is self-contained and should only include elements related to its specific feature.  

The project structure is as follows:  
```yaml
config:
    config.yaml
deployment:
    docker-compose-<env>.yaml
documentation:
    swagger.yaml
    postman-collection.json
internal:
    app:
    config:
    feature:
        <feature name>:
            <feature>.go
            controller.go
            model.go
            repository.go
    interfaces:
        datasources:
            <datasource connector>
        event:
            <event connector>
        http:
            <feature>_handler.go
.gitignore
.golangci-lint.yaml
Dockerfile
go.mod
go.sum
main.go
Makefile
Readme.md
```

## Integration testing

You can use the Postman Flow to run integration tests : [Postman Flow](https://aloe-corp.postman.co/workspace/My-Workspace~c23b3208-facf-4932-8db8-dd4c47714ed9/flow/67a9dd828d573d003172fec7)

In order to run this flow, please deploy locally the API using the `staging` docker-compose.