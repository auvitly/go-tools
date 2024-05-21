package models

import "google.golang.org/grpc/codes"

type EmbedError struct {
	Value any           `json:"value,omitempty"`
	Embed []*EmbedError `json:"embed,omitempty"`
}

type Error struct {
	Code    string         `json:"code,omitempty"`
	Message string         `json:"message,omitempty"`
	Embed   *EmbedError    `json:"embed,omitempty"`
	Wraps   []string       `json:"wraps,omitempty"`
	Fields  map[string]any `json:"fields,omitempty"`
	Codes   struct {
		GRPC codes.Code `json:"grpc"`
		HTTP int        `json:"http"`
	} `json:"codes"`
}
