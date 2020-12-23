package job_status

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
