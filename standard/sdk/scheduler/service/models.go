package service

import (
	"cmp"
	"time"

	"github.com/google/uuid"
)

type CreateTaskParams[T, M cmp.Ordered] struct {
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	Args         []byte
	Labels       map[string]string
}

type ReceiveTaskParams[T cmp.Ordered] struct {
	WorkerID uuid.UUID
	Type     T
	Version  string
	Labels   map[string]string
}

type ReportStateParams[S cmp.Ordered] struct {
	TaskID      uuid.UUID
	WorkerID    uuid.UUID
	SessionID   uuid.UUID
	ReportState ReportState
}

type ReporGetParams struct {
	TaskID uuid.UUID
}

type ReportState interface {
	implReportState()
}

type SetStateDone[S cmp.Ordered] struct {
	StatusCode S
	Result     []byte
}

type SetStateInWork struct {
	StateData map[string]any
}

type SetStatePutOff struct {
	StateData    map[string]any
	CatchLaterAT time.Time
}

func (SetStateDone[S]) implReportState() {}
func (SetStateInWork) implReportState()  {}
func (SetStatePutOff) implReportState()  {}
