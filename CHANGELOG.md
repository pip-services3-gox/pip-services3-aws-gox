# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> AWS specific components for Golang Changelog

## <a name="1.1.2"></a> 1.1.2 (2023-03-03)

### Features
- Now Register is not required for rewriting in Lambda container for inheritance
## <a name="1.1.0 - 1.1.1"></a> 1.1.0 - 1.1.1 (2023-03-02)

### Breaking changes
* Renamed descriptors for services:
    - "\*:service:lambda\*:1.0" -> "\*:service:awslambda\*:1.0"
    - "\*:service:commandable-lambda\*:1.0" -> "\*:service:commandable-awslambda\*:1.0"

* Renamed package container -> containers 

### Features
- Updated dependencies


## <a name="1.0.0"></a> 1.0.0 (2022-10-06)

### Features
* **lambda** AWS Lambda service and client
* **build** factories for constructing module components
* **clients** client components for working with Lambda AWS
* **services**  components for creating services
* **connect** components of installation and connection settings
* **containers**  components for creating containers
* **count** components of working with counters (metrics)
* **log** logging components with saving data


