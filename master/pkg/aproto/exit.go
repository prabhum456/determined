package aproto

import (
	"fmt"

	"github.com/pkg/errors"
)

// ContainerFailure holds the reason why a container did not complete successfully.
type ContainerFailure struct {
	FailureType FailureType
	ErrMsg      string
	ExitCode    *ExitCode
}

func (c ContainerFailure) Error() string {
	if c.ExitCode == nil {
		return fmt.Sprintf("%s: %s", c.FailureType, c.ErrMsg)
	}
	return fmt.Sprintf("%s: %s (exit code %d)", c.FailureType, c.ErrMsg, *c.ExitCode)
}

// ExitCode is the process exit code of the container.
type ExitCode int

const (
	// SuccessExitCode is the 0 zero value exit code.
	SuccessExitCode = 0
)

// ContainerError returns a container failure wrapping the provided error. If the error is nil,
// a stack trace is provided instead.
func ContainerError(failureType FailureType, err error) ContainerStopped {
	if err == nil {
		return ContainerStopped{
			Failure: &ContainerFailure{
				FailureType: failureType,
				ErrMsg:      errors.WithStack(errors.Errorf("unknown error occurred")).Error(),
			},
		}
	}
	return ContainerStopped{
		Failure: &ContainerFailure{
			FailureType: failureType,
			ErrMsg:      err.Error(),
		},
	}
}

// ContainerExited returns a container failure with the encoded exit code. If the exit code is a
// the zero value, no failure is returned.
func ContainerExited(code ExitCode) ContainerStopped {
	if code == SuccessExitCode {
		return ContainerStopped{}
	}
	return ContainerStopped{
		&ContainerFailure{
			FailureType: ContainerFailed,
			ErrMsg:      errors.Errorf("%s: %d", ContainerFailed, code).Error(),
			ExitCode:    &code,
		},
	}
}

// FailureType denotes the type of failure that resulted in the container stopping.
// Each FailureType must be handled by ./internal/task/allocation.go.
type FailureType string

const (
	// ContainerFailed denotes that the container ran but failed with a non-zero exit code.
	ContainerFailed = FailureType("container failed with non-zero exit code")

	// TaskAborted denotes that the task was canceled before it was started.
	TaskAborted = FailureType("task was aborted before the task was started")

	// TaskError denotes that the task failed without an associated exit code.
	TaskError = FailureType("task failed without an associated exit code")

	// AgentFailed denotes that the agent failed while the container was running.
	AgentFailed = FailureType("agent failed while the container was running")

	// AgentError denotes that the agent failed to launch the container.
	AgentError = FailureType("agent failed to launch the container")
)

// IsRestartableSystemError checks if the error is caused by the system and
// shouldn't count against `max_restarts`.
func IsRestartableSystemError(err error) bool {
	switch contErr := err.(type) {
	case ContainerFailure:
		switch contErr.FailureType {
		case ContainerFailed, TaskError:
			return false
		// Questionable, could be considered failures, but for now we don't.
		case AgentError, AgentFailed:
			return true
		// Definitely not a failure.
		case TaskAborted:
			return true
		default:
			return false
		}
	default:
		return false
	}
}