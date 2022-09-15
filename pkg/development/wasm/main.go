//go:build wasm
// +build wasm

package main

import (
	"context"
	"encoding/json"
	"syscall/js"

	"github.com/Permify/permify/internal/commands"
	"github.com/Permify/permify/pkg/development"
	`github.com/Permify/permify/pkg/dsl/schema`
	`github.com/Permify/permify/pkg/errors`
	`github.com/Permify/permify/pkg/tuple`
)

var dev *development.Development

// check -
func check() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := &development.CheckQuery{}
		mErr := json.Unmarshal([]byte(string(args[0].String())), params)
		if mErr != nil {
			return js.ValueOf([]interface{}{false, mErr.Error()})
		}
		var err errors.Error
		var result commands.CheckResponse
		result, err = development.Check(context.Background(), dev.P, params.Subject, params.Action, params.Entity, string(args[1].String()))
		if err != nil {
			return js.ValueOf([]interface{}{false, err.Error()})
		}
		return js.ValueOf([]interface{}{result.Can, nil})
	})
}

// writeSchema -
func writeSchema() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var err errors.Error
		var version string
		version, err = development.WriteSchema(context.Background(), dev.M, string(args[0].String()))
		if err != nil {
			return js.ValueOf([]interface{}{"", err.Error()})
		}
		return js.ValueOf([]interface{}{version, nil})
	})
}

// writeTuple -
func writeTuple() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var t = &tuple.Tuple{}
		mErr := json.Unmarshal([]byte(string(args[0].String())), t)
		if mErr != nil {
			return js.ValueOf([]interface{}{mErr.Error()})
		}
		var err errors.Error
		err = development.WriteTuple(context.Background(), dev.R, *t, string(args[1].String()))
		if err != nil {
			return js.ValueOf([]interface{}{err.Error()})
		}
		return js.ValueOf([]interface{}{nil})
	})
}

// readSchema -
func readSchema() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var sch schema.Schema
		var err errors.Error
		sch, err = development.ReadSchema(context.Background(), dev.M, string(args[0].String()))
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		pretty, err := json.MarshalIndent(sch, "", "  ")
		if err != nil {
			return "", err
		}
		return js.ValueOf([]interface{}{string(pretty), nil})
	})
}

func main() {
	ch := make(chan struct{}, 0)
	dev = development.NewDevelopment()

	js.Global().Set("check", check())
	js.Global().Set("writeSchema", writeSchema())
	js.Global().Set("writeTuple", writeTuple())
	js.Global().Set("ReadSchema", readSchema())
	// js.Global().Set("DeleteTuple", deleteTuple())
	<-ch
}