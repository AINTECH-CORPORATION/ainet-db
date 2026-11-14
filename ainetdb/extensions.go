package ainetdb

// ExtensionDescriptor describes an integrated SOTA capability in AINET DB.
type ExtensionDescriptor struct {
	ID       string
	Upstream string
	Category string
	Enabled  bool
}

// DefaultExtensions returns the SOTA extensions registered for AINET DB.
func DefaultExtensions() []ExtensionDescriptor {
	return []ExtensionDescriptor{
		{ID: "vector-vec0", Upstream: "asg017/sqlite-vec", Category: "vector-search", Enabled: true},
		{ID: "wal-stream", Upstream: "benbjohnson/litestream", Category: "replication", Enabled: true},
		{ID: "db-explorer", Upstream: "simonw/datasette", Category: "explorer", Enabled: true},
		{ID: "edge-branch", Upstream: "tursodatabase/libsql", Category: "edge-sync", Enabled: true},
	}
}
