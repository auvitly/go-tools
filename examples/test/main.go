package main

import (
	"github.com/auvitly/go-tools/recovery"
	"log/slog"
	"runtime/debug"
)

type modelDB struct {
	String string
	PtrInt *int
}

type modelAPI struct {
	String string
	Int    int
}

func log(values ...any) func(any) {
	return func(msg any) {
		slog.Error("we obtain panic",
			slog.Any("panic", msg),
			slog.Any("stack", string(debug.Stack())),
			slog.Any("values", values),
		)
	}
}

func convert(item modelDB) (result *modelAPI) {
	defer recovery.WithAsyncHandlers(log(item)).Do()

	return &modelAPI{
		String: item.String,
		Int:    *item.PtrInt,
	}
}

func main() {
	var (
		records = []modelDB{
			{
				String: "valid",
				PtrInt: new(int),
			},
			{
				String: "not valid",
				PtrInt: nil,
			},
		}
		results []modelAPI
	)

	for _, record := range records {
		if result := convert(record); result != nil {
			results = append(results, *result)
		}
	}

	slog.Info("Our results", "results", results)
}
