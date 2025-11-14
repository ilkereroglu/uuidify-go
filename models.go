package uuidify

// UUIDResponse represents a response returning a single UUID value.
type UUIDResponse struct {
	UUID string `json:"uuid"`
}

// UUIDListResponse represents a response returning multiple UUID values.
type UUIDListResponse struct {
	UUIDs []string `json:"uuids"`
}

// ULIDResponse represents a response returning a single ULID value.
type ULIDResponse struct {
	ULID string `json:"ulid"`
}

// ULIDListResponse represents a response returning multiple ULID values.
type ULIDListResponse struct {
	ULIDs []string `json:"ulids"`
}
