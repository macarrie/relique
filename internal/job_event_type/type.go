package job_event_type

const (
	_ = iota
	JobStart
	JobRetry
	DiffReference
	FallbackToFull
	JobSetup
	ClientPing
	ClientPingError
	SyncStart
	SyncEnd
	SyncError
	ImageCreated
	ImageError
	JobEnd
	JobError
	Unknown
)

var names = map[uint8]string{
	JobStart:        "job_start",
	JobRetry:        "job_retry",
	DiffReference:   "diff_reference",
	FallbackToFull:  "fallback_to_full",
	JobSetup:        "job_setup",
	ClientPing:      "client_ping",
	ClientPingError: "client_ping_error",
	SyncStart:       "sync_start",
	SyncEnd:         "sync_end",
	SyncError:       "sync_error",
	ImageCreated:    "image_created",
	ImageError:      "image_error",
	JobEnd:          "job_end",
	JobError:        "job_error",
}

type JobEventType struct {
	Type uint8 `json:"type"`
}

func New(t uint8) JobEventType {
	return JobEventType{
		Type: t,
	}
}

func (t *JobEventType) String() string {
	if name, ok := names[t.Type]; ok {
		return name
	}

	return "unknown"
}

func FromString(val string) JobEventType {
	for t, name := range names {
		if name == val {
			return JobEventType{Type: t}
		}
	}

	return JobEventType{Type: Unknown}
}

func (t *JobEventType) UnmarshalText(b []byte) error {
	tmp := FromString(string(b))

	*t = tmp

	return nil
}

func (t JobEventType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
