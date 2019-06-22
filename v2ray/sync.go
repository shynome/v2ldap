package v2ray

// SyncResponse of sync func
type SyncResponse struct {
	Delete []string
	Add    []string
}

// Sync ldap user to here
// TODO: sync
func (v2 V2ray) Sync(confirm bool) (resp SyncResponse, err error) {
	return
}
