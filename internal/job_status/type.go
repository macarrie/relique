package job_status

const (
	_ = iota
	Pending
	Active
	Success
	Incomplete
	Error
	Unknown
)

type JobStatus struct {
	Status uint8 `json:"status"`
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

func FromString(val string) JobStatus {
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
		s.Status = Unknown
	}

	return s
}

func (t *JobStatus) UnmarshalText(b []byte) error {
	tmp := FromString(string(b))

	*t = tmp

	return nil
}

func (t JobStatus) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
