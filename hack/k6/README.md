# k6 Load Tests

This directory contains load tests for the [k6](https://k6.io/) load testing tool.

There are [different types](https://k6.io/docs/test-types/load-test-types/) of load tests, but the most common is a _stress test_.

The scenarios in this directory are designed to be run against a local instance of the admission controller.

The available scenarios are:

- **Smoke Test**: A simple test that creates a minor load on the admission controller to validate that the script works as expected under minimal load.
- **Stress Test**: A test that creates a significant load on the admission controller to validate that the script works as expected under heavy load.
- **Spike Test**: A test that creates a sudden spike in load **on** the admission controller to validate that the script works as expected under sudden load.
- **Soak Test**: A test that creates a sustained load on the admission controller to validate that the script works as expected under sustained load.
- **Spike Test**: A test that creates a sudden spike in load on the admission controller to validate that the script works as expected under sudden load.
- **Breakpoint Test**: A test that creates extreme load on the admission controller to evaluate the capacity limits of the admission controller.

## Running the Tests

To run the tests, you must have [k6](https://k6.io/docs/getting-started/installation) installed.

Also, a local instance of the admission controller should be running. You can achieve that using `make ssl` and `make run` in the root of the repository.

To run a test, use the following command:

```bash
k6 run hack/k6/<scenario>.js
```

To run all tests in order, use the following command:

```bash
make stress
```
