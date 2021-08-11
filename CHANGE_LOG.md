# Change Log

## [v1.1.5](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.5) (2021-08-11)

Describe releases of the project.

- describe for releases:
  - features;
  - change log.

## [v1.1.4](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.4) (2021-08-05)

Add the API description in the [AsyncAPI](https://www.asyncapi.com/) format.

- API description in the [AsyncAPI](https://www.asyncapi.com/) format:
  - describe the API:
    - for consuming of the solutions;
    - for producing of the solution results;
  - add the Docker Compose configuration:
    - for generating of the API description;
    - for displaying of the API description.

## [v1.1.3](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.3) (2021-08-04)

Install the [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) tool in the Travis CI configuration and add the general integration tests.

- rename the used environment variables;
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) tool:
  - specify the version of the tool in the Docker configuration;
  - install the tool in the Travis CI configuration;
- add the general integration tests;
- fix and complete the badges in the README.md file.

## [v1.1.2](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.2) (2021-06-26)

Cover with unit tests the loading of the allowed imports from an outer configuration, and add the script for generating such configuration.

- solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package:
  - loading of the allowed imports from an outer configuration:
    - cover the code with unit tests;
- add the script for generating the allowed import configuration.

### Features

- interaction with queues:
  - common properties:
    - automatic declaring of the used queues;
    - passing of a message data in JSON;
  - operations:
    - consuming of the solutions:
      - concurrent handling;
      - once requeue the solution on failure;
    - producing of the solution results:
      - generating of the custom message ID;
- solution runners:
  - solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package:
    - loading of the allowed imports from an outer configuration;
    - canceling by a timeout:
      - code compiling;
      - code execution;
- tools:
  - script for generating the allowed import configuration.

## [v1.1.1](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.1) (2021-06-21)

Generating of the custom message ID on producing of the solution results, and canceling the code compiling by a timeout in the solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package.

- dependencies:
  - update the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package;
  - use the [github.com/thewizardplusplus/go-rabbitmq-utils](https://github.com/thewizardplusplus/go-rabbitmq-utils) package;
- producing of the solution results:
  - generating of the custom message ID;
  - cover the code with unit tests;
- solution runners:
  - remove the dummy solution runner;
  - solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package:
    - canceling the code compiling by a timeout;
    - cover the code with unit tests.

### Features

- interaction with queues:
  - common properties:
    - automatic declaring of the used queues;
    - passing of a message data in JSON;
  - operations:
    - consuming of the solutions:
      - concurrent handling;
      - once requeue the solution on failure;
    - producing of the solution results:
      - generating of the custom message ID;
- solution runners:
  - solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package:
    - loading of the allowed imports from an outer configuration;
    - canceling by a timeout:
      - code compiling;
      - code execution.

## [v1.1](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1) (2021-04-06)

Canceling the code execution by a timeout in the solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package.

- solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package:
  - canceling the code execution by a timeout.

### Features

- interaction with queues:
  - common properties:
    - automatic declaring of the used queues;
    - passing of a message data in JSON;
  - operations:
    - consuming of the solutions:
      - concurrent handling;
      - once requeue the solution on failure;
    - producing of the solution results;
- solution runners:
  - dummy solution runner (returns the current timestamp as the result);
  - solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package:
    - loading of the allowed imports from an outer configuration;
    - canceling the code execution by a timeout.

## [v1.0](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.0) (2021-03-29)

Major version. Implement the solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package.

- implement the solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package:
  - loading of the allowed imports from an outer configuration;
- add the examples:
  - of the request;
  - of the allowed import configuration;
- add the integration Docker Compose configuration.

### Features

- interaction with queues:
  - common properties:
    - automatic declaring of the used queues;
    - passing of a message data in JSON;
  - operations:
    - consuming of the solutions:
      - concurrent handling;
      - once requeue the solution on failure;
    - producing of the solution results;
- solution runners:
  - dummy solution runner (returns the current timestamp as the result);
  - solution runner based on the [github.com/thewizardplusplus/go-code-runner](https://github.com/thewizardplusplus/go-code-runner) package:
    - loading of the allowed imports from an outer configuration.

## [v1.0-alpha.1](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.0-alpha.1) (2021-03-27)

Second alpha of the major version. Update the models and requeue the solution on failure once only.

- update the models:
  - of the solution;
  - of the solution result;
- requeue the solution on failure once only;
- add the debug logging to the dummy solution runner.

### Features

- interaction with queues:
  - common properties:
    - automatic declaring of the used queues;
    - passing of a message data in JSON;
  - operations:
    - consuming of the solutions:
      - concurrent handling;
      - once requeue the solution on failure;
    - producing of the solution results;
- solution runners:
  - dummy solution runner (returns the current timestamp as the result).

## [v1.0-alpha](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.0-alpha) (2021-03-26)

Alpha of the major version. Implementing of the consuming of the solutions, the producing of the solution results, and the dummy solution runner.

### Features

- interaction with queues:
  - common properties:
    - automatic declaring of the used queues;
    - passing of a message data in JSON;
  - operations:
    - consuming of the solutions:
      - concurrent handling;
      - requeue the solution on failure;
    - producing of the solution results;
- solution runners:
  - dummy solution runner (returns the current timestamp as the result).
