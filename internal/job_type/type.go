package job_type

const (
	_ = iota
	Unknown
	Backup
	Restore
)

type JobType struct {
	Type uint8
}

func New(t uint8) JobType {
	return JobType{Type: t}
}

func (t *JobType) String() string {
	switch t.Type {
	case Backup:
		return "backup"
	case Restore:
		return "restore"
	default:
		return "unknown"
	}
}

func FromString(val string) JobType {
	t := JobType{}

	switch val {
	case "backup":
		t.Type = Backup
	case "restore":
		t.Type = Restore
	default:
		t.Type = Unknown
	}

	return t
}

func (t *JobType) UnmarshalText(b []byte) error {
	tmp := FromString(string(b))

	*t = tmp

	return nil
}

func (t JobType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
