# analytics-go e2e-cli

A small CLI tool used for end-to-end testing of the `analytics-go` SDK. It accepts a JSON description of event sequences, sends those events through the real SDK, and reports the result.

## Prerequisites

- Go 1.17+
- (For `run-e2e.sh`) Node.js 18+ and the `sdk-e2e-tests` repo checked out alongside this repo

## Building

```bash
cd e2e-cli
go build -o e2e-cli-bin ./...
```

## Usage

```bash
./e2e-cli-bin --input '<JSON>'
```

The CLI writes debug/log information to **stderr** and the JSON result to **stdout**.

Exit code is `0` on success and `1` on failure.

## Input JSON format

```json
{
  "writeKey": "YOUR_WRITE_KEY",
  "apiHost": "https://api.segment.io",
  "sequences": [
    {
      "delayMs": 0,
      "events": [
        {
          "type": "track",
          "event": "Button Clicked",
          "userId": "user-123",
          "properties": { "plan": "pro" }
        },
        {
          "type": "identify",
          "userId": "user-123",
          "traits": { "email": "user@example.com" }
        },
        {
          "type": "page",
          "userId": "user-123",
          "name": "Home"
        },
        {
          "type": "screen",
          "userId": "user-123",
          "name": "Main Screen"
        },
        {
          "type": "alias",
          "userId": "new-id",
          "previousId": "old-id"
        },
        {
          "type": "group",
          "userId": "user-123",
          "groupId": "group-456",
          "traits": { "name": "Acme Corp" }
        }
      ]
    }
  ],
  "config": {
    "flushAt": 15,
    "flushInterval": 1000,
    "maxRetries": 3,
    "timeout": 10
  }
}
```

### Top-level fields

| Field       | Type   | Description                                      |
|-------------|--------|--------------------------------------------------|
| `writeKey`  | string | Segment write key used to authenticate requests  |
| `apiHost`   | string | Full API endpoint URL (e.g. `https://api.segment.io`) |
| `sequences` | array  | List of event sequences (run in order)           |
| `config`    | object | SDK configuration overrides                      |

### Sequence fields

| Field     | Type  | Description                                                         |
|-----------|-------|---------------------------------------------------------------------|
| `delayMs` | int   | Milliseconds to wait before sending this sequence's events          |
| `events`  | array | List of events to enqueue                                           |

### Event fields

| Field         | Type   | Description                                                         |
|---------------|--------|---------------------------------------------------------------------|
| `type`        | string | Event type: `track`, `identify`, `page`, `screen`, `alias`, `group` |
| `event`       | string | Event name (required for `track`)                                   |
| `userId`      | string | User ID                                                             |
| `anonymousId` | string | Anonymous ID (used when `userId` is absent)                         |
| `messageId`   | string | Optional explicit message ID                                        |
| `previousId`  | string | Previous user ID (required for `alias`)                             |
| `groupId`     | string | Group ID (required for `group`)                                     |
| `name`        | string | Page or screen name                                                 |
| `timestamp`   | string | ISO 8601 timestamp (RFC3339), e.g. `2024-01-15T10:30:00Z`          |
| `properties`  | object | Event properties (for `track`, `page`, `screen`)                    |
| `traits`      | object | User or group traits (for `identify`, `group`)                      |
| `integrations`| object | Integration-specific settings                                       |

### Config fields

| Field           | Type | Description                              |
|-----------------|------|------------------------------------------|
| `flushAt`       | int  | Max events per batch (default: 250)      |
| `flushInterval` | int  | Flush interval in milliseconds           |
| `maxRetries`    | int  | (informational, not directly used by SDK)|
| `timeout`       | int  | Timeout in seconds (informational)       |

## Output JSON format

On success:

```json
{"success": true, "sentBatches": 1}
```

On failure:

```json
{"success": false, "sentBatches": 0, "error": "description of error"}
```

## Running E2E tests

Use the provided shell script to build the CLI and run the shared `sdk-e2e-tests` suite:

```bash
# From the e2e-cli directory, with sdk-e2e-tests checked out as a sibling of analytics-go:
./run-e2e.sh

# Override the e2e test directory:
E2E_TESTS_DIR=/path/to/sdk-e2e-tests ./run-e2e.sh

# Pass additional arguments to run-tests.sh:
./run-e2e.sh --suite basic
```

The script:
1. Builds the `e2e-cli-bin` binary using `go build`
2. Invokes `sdk-e2e-tests/scripts/run-tests.sh` with the appropriate `--sdk-dir` and `--cli` flags
