# drift-watch

Lightweight daemon that detects configuration drift between running containers and their source manifests.

## Installation

```bash
go install github.com/yourorg/drift-watch@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/drift-watch.git && cd drift-watch && make build
```

## Usage

Point `drift-watch` at your manifests directory and let it run alongside your containers:

```bash
drift-watch --manifests ./deploy/manifests --interval 30s
```

Example output when drift is detected:

```
[DRIFT] container: api-server
  expected image: myapp:v1.2.0
  running image:  myapp:v1.1.9
  env mismatch:   LOG_LEVEL (expected=info, got=debug)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--manifests` | `./manifests` | Path to source manifest files |
| `--interval` | `60s` | How often to check for drift |
| `--output` | `text` | Output format: `text` or `json` |
| `--alert-webhook` | `` | Optional webhook URL for drift alerts |

### JSON Output

```bash
drift-watch --manifests ./deploy --output json | jq .
```

```json
{
  "container": "api-server",
  "drifted": true,
  "fields": ["image", "env.LOG_LEVEL"]
}
```

## Requirements

- Go 1.21+
- Docker or a compatible container runtime

## License

MIT © yourorg