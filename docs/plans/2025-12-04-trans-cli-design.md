# trans-cli Design

Transmission RPC client for listing, filtering, and controlling torrents.

## Commands

```
trans-cli list                  # all torrents
trans-cli downloading           # only downloading
trans-cli seeding               # only seeding
trans-cli stopped               # only stopped
trans-cli show <id>             # single torrent details
trans-cli add <magnet|file>     # add torrent (magnet URI or .torrent file)
trans-cli start <id>            # resume torrent
trans-cli stop <id>             # pause torrent
```

## Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `TRANSMISSION_URL` | Yes | RPC endpoint (e.g., `http://localhost:9091`) |
| `TRANSMISSION_USER` | No | Username for HTTP Basic Auth |
| `TRANSMISSION_PASS` | No | Password for HTTP Basic Auth |

Flags (override env vars):
- `--url`
- `--user`
- `--pass`
- `--insecure` / `-k` (skip TLS verify)

RPC path defaults to `/transmission/rpc`.

## Output Format

### List Output

```yaml
torrents:
- id: 1
  name: Ubuntu 24.04 LTS
  status: seeding
  percentDone: 100
  totalSize: 4.2 GB
  uploadRatio: 2.35
  rateDownload: 0 B/s
  rateUpload: 125 KB/s
  eta: -
  tracker: torrent.ubuntu.com
  peersConnected: 12
```

### Show Output

```yaml
torrent:
  id: 1
  name: Ubuntu 24.04 LTS
  status: seeding
  percentDone: 100
  totalSize: 4.2 GB
  downloadedEver: 4.2 GB
  uploadedEver: 9.9 GB
  uploadRatio: 2.35
  rateDownload: 0 B/s
  rateUpload: 125 KB/s
  eta: -
  tracker: torrent.ubuntu.com
  peersConnected: 12
  addedDate: 2024-12-01 10:30:00
  doneDate: 2024-12-01 12:45:00
  downloadDir: /downloads
```

## Project Structure

```
trans/
├── cmd/trans-cli/
│   └── main.go          # cobra commands, flags
├── pkg/trans/
│   ├── client.go        # RPC client, session ID handling
│   ├── models.go        # torrent structs
│   ├── output.go        # YAML formatting
│   └── errors.go        # error types, exit codes
```

## RPC Implementation Notes

### Session ID Handling

Transmission uses CSRF protection via `X-Transmission-Session-Id` header. On HTTP 409 response:
1. Extract session ID from `X-Transmission-Session-Id` response header
2. Retry request with new session ID

### Torrent Status Mapping

RPC `status` field values:
- 0: stopped
- 1: queued to verify
- 2: verifying
- 3: queued to download
- 4: downloading
- 5: queued to seed
- 6: seeding

### Adding Torrents

Two methods supported:
- `filename`: magnet URI string
- `metainfo`: base64-encoded .torrent file contents

### Required RPC Fields

For list/show, request these fields via `torrent-get`:
- id, name, status
- percentDone, totalSize, sizeWhenDone
- uploadRatio, rateDownload, rateUpload
- eta, peersConnected
- trackers (array, extract first hostname)
- downloadedEver, uploadedEver (show only)
- addedDate, doneDate, downloadDir (show only)

## Updates to Existing Files

- `README.md` - add trans-cli section, add env vars to config table
- `install.sh` - include trans-cli in download/install
- `build.sh` - add trans build target
- `go.work` (if using) - add trans module

## Sources

- [Transmission RPC Spec](https://github.com/transmission/transmission/blob/main/docs/rpc-spec.md)
