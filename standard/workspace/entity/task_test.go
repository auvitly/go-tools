package entity_test

import (
	"testing"
	"time"

	"github.com/auvitly/go-tools/standard/workspace/entity"
	"github.com/auvitly/go-tools/test/lab"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCastTask(t *testing.T) {
	type (
		Code     int
		Mode     int
		Status   int
		Data     = map[string]any
		External = entity.SpecifiedTask[Code, Mode, Status, Data, Data, Data]
	)

	var data = Data{
		"int":    float64(1),
		"string": "string",
		"slice":  []any{"string"},
	}

	var in = &External{
		ID:           uuid.New(),
		ParentTaskID: nil,
		Type:         0,
		Mode:         0,
		Args:         data,
		StateData:    data,
		Result:       &data,
		StatusCode:   lab.Pointer[Status](0),
		CreatedTS:    time.Now(),
		UpdatedTS:    time.Now(),
		DoneTS:       lab.Pointer(time.Now()),
		SessionID:    lab.Pointer(uuid.New()),
		AssignTS:     lab.Pointer(time.Now()),
		Labels: map[string]string{
			"string": "string",
		},
	}

	mid, stderr := entity.MarshalTask(in)
	if stderr != nil {
		panic(stderr)
	}

	out, stderr := entity.UnmarshalTask[Code, Mode, Status, Data, Data, Data](mid)
	if stderr != nil {
		panic(stderr)
	}

	require.EqualValues(t, in, out)
}
