// Copyright (c) 2017-2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2016-2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"os"

	"github.com/danos/yang/schema"
	"github.com/godbus/dbus"
	"github.com/jsouthworth/objtree"
)

type dbusServiceRead interface {
	Get(string) []byte
}

type dbusServiceWrite interface {
	Check([]byte) []error
	Set([]byte) []error
}

var supervisor *objtree.BusManager

func init() {
}

type dbusRunningConfig struct {
	supervisor *objtree.BusManager

	running RunningConfig
	schema  schema.Node
}

func (r *dbusRunningConfig) Get(path string) []byte {
	return r.running.Get(path)
}

func (r *dbusRunningConfig) Check(jsonTree []byte) []error {
	return r.running.Check(jsonTree)
}

func (r *dbusRunningConfig) Set(jsonTree []byte) []error {
	return r.running.Set(jsonTree)
}

func newDbusRunningConfig(
	supervisor *objtree.BusManager,
	client RunningConfig,
) *dbusRunningConfig {

	// TODO get schema

	return &dbusRunningConfig{supervisor, client, nil}
}

func vciConnectFn(
	hdlr dbus.Handler,
	sig dbus.SignalHandler,
) (*dbus.Conn, error) {
	return dbus.DialHandler("unix:path=/var/run/vci/vci_bus_socket", hdlr, sig)
}
func dbusRegisterService(namespace string, client RunningConfig) ServiceHandle {

	supervisor, err := objtree.NewAnonymousBusManager(vciConnectFn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	runningConfig := newDbusRunningConfig(supervisor, client)
	obj := supervisor.NewObject("/running", runningConfig)
	err = obj.Implements("net.vyatta.vci.config.read", (*dbusServiceRead)(nil))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = obj.Implements("net.vyatta.vci.config.write", (*dbusServiceWrite)(nil))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return obj
}
