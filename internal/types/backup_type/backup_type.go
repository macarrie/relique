package backup_type

import (
	"fmt"
)

const (
	Diff = iota
	Full
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
	default:
		return "unknown"
	}
}

func FromString(val string) (BackupType, error) {
	t := BackupType{}
	switch val {
	case "diff":
		t.Type = Diff
	case "full":
		t.Type = Full
	default:
		return t, fmt.Errorf("unknown variant '%s'", val)
	}

	return t, nil
}
