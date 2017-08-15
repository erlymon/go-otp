package otp

import (
	"time"
)

/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

type State interface{}

type Request struct {
	Name string
	Data interface{}
	Resp chan Responce
}

type Responce struct {
	Data  interface{}
	Error error
}

type Module interface {
	Init()
	Handle(request Request) Responce
	HandleTimer(t time.Time)
	Terminate()
}

type GenServer struct {
	mod Module
	Q   chan struct{}
	R   chan Request
	T   <-chan time.Time
}

func NewGenServer() *GenServer {
	return &GenServer{
		Q: make(chan struct{}),
		R: make(chan Request),
	}
}

func NewGenServerWithTimer(c <-chan time.Time) *GenServer {
	return &GenServer{
		Q: make(chan struct{}),
		R: make(chan Request),
		T: c,
	}
}

func (server *GenServer) Start(mod Module) {
	go func() {
		mod.Init()
		defer mod.Terminate()
		server.loop(mod)
	}()
}

func (server *GenServer) Stop() {
	close(server.Q)
}

func (server *GenServer) Rpc(name string, data interface{}) chan Responce {
	responce := make(chan Responce)
	server.R <- Request{Name: name, Data: data, Resp: responce}
	return responce
}

func (server *GenServer) loop(mod Module) {

	for {
		select {
		case request := <-server.R:
			responce := mod.Handle(request)
			request.Resp <- responce
		case t := <-server.T:
			mod.HandleTimer(t)
		case <-server.Q:
			close(server.R)
			return
		}
	}
}
