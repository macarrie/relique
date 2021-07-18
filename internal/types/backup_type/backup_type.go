package backup_type

const (
	Unknown = iota
	Diff
	Full
	Restore // For display purposes
	CumulativeDiff
)

type BackupType struct {
	Type uint8
}

func (t *BackupType) String() string {
	switch t.Type {
	case Diff:
		return "diff"
	case Full:
		return "full"
	case Restore:
		return "restore"
	case CumulativeDiff:
		return "cumulative_diff"
	default:
		return "unknown"
	}
}

func FromString(val string) BackupType {
	t := BackupType{}

	switch val {
	case "diff":
		t.Type = Diff
	case "full":
		t.Type = Full
	case "restore":
		t.Type = Restore
	case "cumulative_diff":
		t.Type = CumulativeDiff
	default:
		t.Type = Unknown
	}

	return t
}

func (t *BackupType) UnmarshalText(b []byte) error {
	tmp := FromString(string(b))

	*t = tmp

	return nil
}

func (t BackupType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
