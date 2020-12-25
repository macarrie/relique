package job_status

import "fmt"

const (
	Pending = iota
	Active
	Success
	Incomplete
	Error
)

type JobStatus struct {
	Status uint8
}

func New(status uint8) JobStatus {
	return JobStatus{
		Status: status,
	}
}

func (s *JobStatus) String() string {
	switch s.Status {
	case Pending:
		return "pending"
	case Active:
		return "active"
	case Success:
		return "success"
	case Incomplete:
		return "incomplete"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

func FromString(val string) (JobStatus, error) {
	s := JobStatus{}
	switch val {
	case "pending":
		s.Status = Pending
	case "active":
		s.Status = Active
	case "success":
		s.Status = Success
	case "incomplete":
		s.Status = Incomplete
	case "error":
		s.Status = Error
	default:
		return s, fmt.Errorf("unknown variant '%s'", val)
	}

	return s, nil
}
