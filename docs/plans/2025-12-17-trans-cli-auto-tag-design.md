# trans-cli: Auto-tag torrents with "cli"

## Summary

Auto-add label `"cli"` to torrents added via `trans-cli add`.

## Implementation

### 1. models.go - Add Labels field

```go
type TorrentAddArgs struct {
    Filename string   `json:"filename,omitempty"`
    Metainfo string   `json:"metainfo,omitempty"`
    Labels   []string `json:"labels,omitempty"`
}
```

### 2. client.go - Set labels in add methods

Both `AddTorrentMagnet` and `AddTorrentFile` set `Labels: []string{"cli"}`.

## Requirements

- Transmission 3.00+ (RPC v16)
- Older versions ignore the field silently
