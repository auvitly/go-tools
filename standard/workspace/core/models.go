package core

import (
	"cmp"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateTaskParams[T, M cmp.Ordered] struct {
	ParentTaskID *uuid.UUID
	Type         T
	Mode         M
	Args         json.RawMessage
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
	Result     json.RawMessage
}

type SetStateInWork struct {
	StateData json.RawMessage
}

type SetStatePutOff struct {
	StateData    json.RawMessage
	CatchLaterAT time.Time
}

func (SetStateDone[S]) implReportState() {}
func (SetStateInWork) implReportState()  {}
func (SetStatePutOff) implReportState()  {}
