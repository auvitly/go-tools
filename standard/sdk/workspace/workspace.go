package workspace

import (
	"encoding/json"
	"maps"
	"sync"

	"github.com/auvitly/go-tools/collection/stderrs"
	"github.com/auvitly/go-tools/standard/collection/constraints"
)

type Stage interface {
	constraints.Integer | string
}

type Workspace[S Stage] struct {
	_mu          sync.RWMutex
	_stage       S
	_description string
	_values      map[string]any
}

type Value struct {
	Key   string
	Value any
}

type isWorkspace interface {
	mu() *sync.RWMutex
	values() map[string]any
}

func (w *Workspace[S]) mu() *sync.RWMutex      { return &w._mu }
func (w *Workspace[S]) values() map[string]any { return w._values }
func (w *Workspace[S]) Stage() S               { return w._stage }
func (w *Workspace[S]) Description() string    { return w._description }

func SetStage[S Stage](w *Workspace[S], stage S, message string) {
	w.mu().Lock()
	defer w.mu().Unlock()

	w._stage = stage
	w._description = message
}

func New[S Stage](stage S, message string) *Workspace[S] {
	return &Workspace[S]{
		_stage:       stage,
		_description: message,
		_values:      make(map[string]any),
	}
}

func StoreValue(w isWorkspace, key string, value any) *stderrs.Error {
	w.mu().Lock()
	defer w.mu().Unlock()

	raw, err := json.Marshal(value)
	if err != nil {
		return stderrs.Internal.SetMessage("not supported '%T' type as value", value)
	}

	var stored any

	err = json.Unmarshal(raw, &stored)
	if err != nil {
		return stderrs.Internal.SetMessage("not supported '%T' type as value", value)
	}

	w.values()[key] = stored

	return nil
}

func StoreValues(w isWorkspace, pairs ...Value) *stderrs.Error {
	w.mu().Lock()
	defer w.mu().Unlock()

	var parsed = make(map[string]any)

	for _, pair := range pairs {
		raw, err := json.Marshal(pair.Value)
		if err != nil {
			return stderrs.Internal.SetMessage("not supported '%T' type as value", pair.Value)
		}

		var stored any

		err = json.Unmarshal(raw, &stored)
		if err != nil {
			return stderrs.Internal.SetMessage("not supported '%T' type as value", pair.Value)
		}

		parsed[pair.Key] = stored
	}

	for key, value := range parsed {
		w.values()[key] = value
	}

	return nil
}

func LoadValue[V any](w isWorkspace, key string) (value V, stderr *stderrs.Error) {
	w.mu().RLock()
	defer w.mu().RUnlock()

	stored, ok := w.values()[key]
	if !ok {
		return value, stderrs.NotFound.SetMessage("not found value with '%s' key", key)
	}

	raw, err := json.Marshal(stored)
	if err != nil {
		return value, stderrs.Internal.SetMessage("not supported '%T' type as value", value)
	}

	err = json.Unmarshal(raw, &value)
	if err != nil {
		return value, stderrs.Internal.SetMessage("not supported '%T' type as value", value)
	}

	return value, nil
}

func ToMap[S Stage](w *Workspace[S]) map[string]any {
	w.mu().RLock()
	defer w.mu().RUnlock()

	return map[string]any{
		"stage":       w._stage,
		"description": w._description,
		"values":      maps.Clone(w._values),
	}
}

func FromMap[S Stage](data map[string]any) (*Workspace[S], *stderrs.Error) {
	type workspace[S Stage] struct {
		Stage       S              `json:"stage"`
		Description string         `json:"description"`
		Values      map[string]any `json:"values"`
	}

	var ws = new(workspace[S])

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, stderrs.Internal.SetMessage("cast from map error: %s", err.Error())
	}

	err = json.Unmarshal(raw, ws)
	if err != nil {
		return nil, stderrs.Internal.SetMessage("cast to workspace error: %s", err.Error())
	}

	return &Workspace[S]{
		_stage:       ws.Stage,
		_description: ws.Description,
		_values:      ws.Values,
	}, nil
}

func (w *Workspace[S]) MarshalJSON() ([]byte, error) {
	return json.Marshal(ToMap(w))
}

func (w *Workspace[S]) UnmarshalJSON(data []byte) error {
	var raw = make(map[string]any)

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	ws, stderr := FromMap[S](raw)
	if stderr != nil {
		return stderr
	}

	w._stage = ws._stage
	w._description = ws._description
	w._values = ws._values

	return nil
}
