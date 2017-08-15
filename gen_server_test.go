package otp

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

import (
	"log"
	"testing"
	"time"
)

type NameServer struct {
	Col []string
}

func (server *NameServer) Init() {
}

func (server *NameServer) Handle(request Request) Responce {
	switch request.Name {
	case "SetName":
		server.Col = append(server.Col, request.Data.(string))
		return Responce{Data: len(server.Col) - 1}
	case "GetName":
		return Responce{Data: server.Col[request.Data.(int)]}
	default:
		return Responce{}
	}
}

func (server *NameServer) HandleTimer(t time.Time) {
}

func (server *NameServer) Terminate() {
}

func NewNameServer() Module {
	s := new(NameServer)
	s.Col = make([]string, 0)
	return s
}

func TestGenServer(t *testing.T) {
	gs := NewGenServer()
	gs.Start(NewNameServer())
	defer gs.Stop()

	//req := Request{Name: "SetName", Data: "Sergey Penkovsky", Resp: make(chan Responce)}
	resp := gs.Rpc("SetName", "Sergey Penkovsky")
	res := <-resp
	log.Println("RES SetName: ", res)
	close(resp)

	resp = gs.Rpc("GetName", res.Data.(int))
	res = <-resp
	log.Println("RES GetName: ", res)
	close(resp)
	time.Sleep(5000)
}
