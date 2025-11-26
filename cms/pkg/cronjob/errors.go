package cronjob

import "fmt"

// ErrJobAlreadyExists is returned when attempting to add a job with a name
// that already exists in the job registry.
type ErrJobAlreadyExists struct {
	JobName string
}

func (e *ErrJobAlreadyExists) Error() string {
	return fmt.Sprintf("job already exists: %s", e.JobName)
}

// NewErrJobAlreadyExists creates a new ErrJobAlreadyExists error.
func NewErrJobAlreadyExists(jobName string) error {
	return &ErrJobAlreadyExists{JobName: jobName}
}

// ErrJobNotFound is returned when attempting to execute a job that doesn't
// exist in the job registry.
type ErrJobNotFound struct {
	JobName string
}

func (e *ErrJobNotFound) Error() string {
	return fmt.Sprintf("job not found: %s", e.JobName)
}

// NewErrJobNotFound creates a new ErrJobNotFound error.
func NewErrJobNotFound(jobName string) error {
	return &ErrJobNotFound{JobName: jobName}
}

// ErrJobScheduleFailed is returned when the cron scheduler fails to add a job.
type ErrJobScheduleFailed struct {
	JobName string
	Cause   error
}

func (e *ErrJobScheduleFailed) Error() string {
	return fmt.Sprintf("failed to schedule job %s: %v", e.JobName, e.Cause)
}

func (e *ErrJobScheduleFailed) Unwrap() error {
	return e.Cause
}

// NewErrJobScheduleFailed creates a new ErrJobScheduleFailed error.
func NewErrJobScheduleFailed(jobName string, cause error) error {
	return &ErrJobScheduleFailed{JobName: jobName, Cause: cause}
}

// ErrQueuePushFailed is returned when pushing a job to the Redis queue fails.
type ErrQueuePushFailed struct {
	JobName string
	Cause   error
}

func (e *ErrQueuePushFailed) Error() string {
	return fmt.Sprintf("failed to push job %s to queue: %v", e.JobName, e.Cause)
}

func (e *ErrQueuePushFailed) Unwrap() error {
	return e.Cause
}

// NewErrQueuePushFailed creates a new ErrQueuePushFailed error.
func NewErrQueuePushFailed(jobName string, cause error) error {
	return &ErrQueuePushFailed{JobName: jobName, Cause: cause}
}

// ErrEmptyJobName is returned when a job name is empty or invalid.
type ErrEmptyJobName struct{}

func (e *ErrEmptyJobName) Error() string {
	return "job name cannot be empty"
}

// NewErrEmptyJobName creates a new ErrEmptyJobName error.
func NewErrEmptyJobName() error {
	return &ErrEmptyJobName{}
}
