## Setup
### Start Infrastructure
```bash
docker-compose up -d
```
This launches Kafka (with UI at http://localhost:8040/), Reporter, and Redis.

### Install and tidy dependencies for both services:
```bash
make tidy
```
This will run `go mod tidy` in both the `forwarder` and `diff-calculator` folders.

## Running the Services
Use the provided Makefile for all main commands.
- Run the Forwarder:
    ```bash
    make run-forwarder
    ```

- Run the Diff-Calculator:
    ```bash
    make run-diff-calculator
    ```

- Run both services:
    ```bash
    make run-apps
    ```

## Testing
- Test Forwarder:
    ```bash
    make test-forwarder
    ```

- Test Diff-Calculator:
    ```bash
    make test-diff-calculator
    ```

- Test both:
    ```bash
    make test
    ```

## Q&A
*Q: What's one advantage and one disadvantage of joining the Forwarder and Diff-Calculator in the same application.*

A:
- Advantage: Lower deployment and operational overhead, and reduced latency (since you avoid Kafka as an intermediate queue).
- Disadvantage: Loss of separation of concerns and scalability. It's harder to scale, maintain, and deploy both functionalities indipendently.

*Q: Can you tell the differences? (between enpoint results)*

A: All endpoints return similar metrics (clicks, cost, data, impressions, installs), but some fields (especially install, clicks, cost, impressions) might apprear as either numbers or strings, requiring normalization.

*Q: How did you create the Kafka topic?*

A: Topics (`events` and `diffs`) are programmatically created by each service at startup. This ensures the required topic exists before writing any messages.

*Q: How do you recognize an existing entry in the datastore?*

A: Each event is keyed in the datastore by a combination of `partner` and `date`. If a previous record exists for this key, it is used for diff calculation, otherwise, the current event is treated as the initial value.

*Q: How do you create easily testable code?*

A: 
1. Business logic is separated from infrastructure code using interfaces.
2. All external dependencies (Kafka, Reporter, Redis) are abstracter behind interfaces and can be mocker in tests.
3. Core logic (such as diff calculation) is implemented as pure functions, which are trivial to unit test.
4. Each service is organized in a modular way with a clear separation between entrypoint, application logic, infra, and shared code.


