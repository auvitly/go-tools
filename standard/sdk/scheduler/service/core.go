package service

import (
	"cmp"
	"context"
	"time"

	"github.com/auvitly/go-tools/standard/sdk/scheduler/storage"
	"github.com/auvitly/go-tools/stderrs"
)

type Service[T, M, S cmp.Ordered] struct {
	dependencies Dependencies[T, M, S]
	config       Config
}

func New[T, M, S cmp.Ordered](
	ctx context.Context,
	dependencies Dependencies[T, M, S],
	config Config,
) (
	*Service[T, M, S], *stderrs.Error,
) {
	var c = &Service[T, M, S]{
		dependencies: dependencies,
		config:       config,
	}

	if c.config.PullingInterval == 0 {
		return nil, stderrs.FailedPrecondition.SetMessage("pulling interval is 0")
	}

	go c.start(ctx)

	return c, nil
}

func (c *Service[T, M, S]) start(ctx context.Context) {
	var ticker = time.NewTicker(c.config.PullingInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tasks, stderr := c.dependencies.TaskStorage.List(ctx, storage.TaskListParams{
				OnlyAssigned: true,
			})
			if stderr != nil {
				continue
			}

			for _, task := range tasks {
				if time.Since(task.UpdatedTS) > c.config.TaskDowntime {
					_, stderr = c.dependencies.TaskStorage.Update(ctx, storage.TaskUpdateParams[S]{
						TaskID:       task.ID,
						StatusCode:   task.StatusCode,
						StateData:    task.StateData,
						Result:       task.Result,
						UpdatedTS:    time.Now(),
						CatchLaterTS: nil,
						DoneTS:       nil,
						SessionID:    nil,
						AssignTS:     nil,
					})
					if stderr != nil {
						continue
					}

					stderr = c.dependencies.SessionStorage.Done(ctx, storage.SessionDoneParams{
						SessionID: *task.SessionID,
					})
					if stderr != nil {
						continue
					}
				}
			}
		}
	}
}
