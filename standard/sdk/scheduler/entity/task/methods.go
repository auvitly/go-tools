package task

import (
	"time"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/google/uuid"
)

type NewParams struct {
	ID        uuid.UUID
	ParentID  *uuid.UUID
	Type      string
	Mode      Mode
	PendingTS *time.Time
	Arguments map[string]any
	Labels    map[string]string
}

func New(params NewParams) (*Task, *stderrs.Error) {
	var ts = time.Now()

	switch {
	case params.ID == uuid.Nil:
		return nil, stderrs.InvalidArgument.SetMessage("specify task id")
	case params.ParentID != nil && *params.ParentID == uuid.Nil:
		return nil, stderrs.InvalidArgument.SetMessage("specify task parent id")
	case len(params.Type) == 0:
		return nil, stderrs.InvalidArgument.SetMessage("specify task type")
	case len(params.Mode) == 0:
		return nil, stderrs.InvalidArgument.SetMessage("specify task mode")
	default:
		return &Task{
			ID:       params.ID,
			ParentID: params.ParentID,

			Type:   params.Type,
			Mode:   params.Mode,
			Status: StatusCreated,

			Arguments: params.Arguments,
			Data:      nil,
			Results:   nil,

			CreatedTS: ts,
			UpdatedTS: ts,
			PendingTS: func() time.Time {
				if params.PendingTS != nil {
					return *params.PendingTS
				}

				return ts
			}(),

			WorkerSessionID: nil,
			WorkerAssignTS:  nil,
			WorkerLabels:    params.Labels,
		}, nil
	}
}
