// Copyright (c) 2006-2019, xiaobo
//
// This is free software, licensed under the GNU General Public License v3.
// See /LICENSE for more information.
//

package main

import (
	"bytes"
	"strings"
)

type RawHeader struct {
	Content string
	Lines   []string
}

func NewRawHeader(content string) *RawHeader {
	ret := &RawHeader{Content: content}
	ret.Lines = strings.Split(content, "\n")

	return ret
}

type HeaderHandler interface {
	GetExt() string
	Execute(rh *RawHeader) string
}

var HeaderHandlers = []HeaderHandler{
	&GoHeaderHandler{Base{Ext: ".go"}},
	&ExHeaderHandler{Base{Ext: ".ex"}},
	&ExHeaderHandler{Base{Ext: ".rb"}},
	&ExHeaderHandler{Base{Ext: ".py"}},
	&ErlHeaderHandler{Base{Ext: ".erl"}},
	&LuaHeaderHandler{Base{Ext: ".lua"}},
	&CSSHeaderHandler{Base{Ext: ".js"}},
	&CSSHeaderHandler{Base{Ext: ".java"}},
	&CSSHeaderHandler{Base{Ext: ".css"}},
}

func GetHandler(ext string) HeaderHandler {
	for _, handler := range HeaderHandlers {
		if ext == handler.GetExt() {
			return handler
		}
	}

	return nil
}

//////// Base Handler ////////

type Base struct {
	Ext string
}

func (base *Base) GetExt() string {
	return base.Ext
}

//////// Handlers ////////

//// Go ////
type GoHeaderHandler struct {
	Base
}

func (handler *GoHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	for _, line := range rh.Lines {
		if "\r" == line || "\n" == line {
			buffer.WriteString("//\n")
		} else {
			buffer.WriteString("// " + line + "\n")
		}
	}

	return buffer.String()
}

//// Elixir Ruby Python ////
type ExHeaderHandler struct {
	Base
}

func (handler *ExHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	for _, line := range rh.Lines {
		if "\r" == line || "\n" == line {
			buffer.WriteString("#\n")
		} else {
			buffer.WriteString("# " + line + "\n")
		}
	}

	return buffer.String()
}

//// Erlang ////
type ErlHeaderHandler struct {
	Base
}

func (handler *ErlHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	for _, line := range rh.Lines {
		if "\r" == line || "\n" == line {
			buffer.WriteString("%%\n")
		} else {
			buffer.WriteString("%% " + line + "\n")
		}
	}

	return buffer.String()
}

//// Lua ////
type LuaHeaderHandler struct {
	Base
}

func (handler *LuaHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	for _, line := range rh.Lines {
		if "\r" == line || "\n" == line {
			buffer.WriteString("--\n")
		} else {
			buffer.WriteString("-- " + line + "\n")
		}
	}

	return buffer.String()
}

//// Java JavaScript CSS ////
type CSSHeaderHandler struct {
	Base
}

func (handler *CSSHeaderHandler) Execute(rh *RawHeader) string {
	var buffer bytes.Buffer

	buffer.WriteString("/*\n")
	for _, line := range rh.Lines {
		if "\r" == line || "\n" == line {
			buffer.WriteString(" *\n")
		} else {
			buffer.WriteString(" * " + line + "\n")
		}
	}
	buffer.WriteString(" */\n")

	return buffer.String()
}
