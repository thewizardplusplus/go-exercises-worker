# Change Log

## [v1.1.4](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.4) (2021-08-05)

## [v1.1.3](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.3) (2021-08-04)

## [v1.1.2](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.2) (2021-06-26)

## [v1.1.1](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1.1) (2021-06-21)

## [v1.1](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.1) (2021-04-06)

## [v1.0](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.0) (2021-03-29)

## [v1.0-alpha.1](https://github.com/thewizardplusplus/go-exercises-worker/tree/v1.0-alpha.1) (2021-03-27)

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
