// Copyright (c) 2017-2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2016 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package dbus

import (
	"fmt"

	"github.com/danos/yangd"
	"github.com/godbus/dbus"
)

type dbusDispatcher struct {
	conn *dbus.Conn
}

// Compile time check that the concrete type meets the interface
var _ yangd.Dispatcher = (*dbusDispatcher)(nil)

type dbusService struct {
	serviceName   string
	runningConfig dbus.BusObject
	runningState  dbus.BusObject
}

// Compile time check that the concrete type meets the interface
var _ yangd.Service = (*dbusService)(nil)

func (d *dbusDispatcher) dbusService(serviceName string) *dbusService {
	return &dbusService{
		serviceName: serviceName,
		runningConfig: d.conn.Object(serviceName,
			dbus.ObjectPath("/running")),
		runningState: d.conn.Object(serviceName,
			dbus.ObjectPath("/state")),
	}
}

func (s *dbusService) GetRunning(path string) ([]byte, error) {
	var result string
	err := s.runningConfig.Call("net.vyatta.vci.config.read.Get",
		0).Store(&result)
	return []byte(result), err
}

func (s *dbusService) ValidateCandidate(candidate []byte) error {
	err := s.runningConfig.Call("net.vyatta.vci.config.write.Check",
		0, string(candidate)).Store()
	return err
}

func (s *dbusService) SetRunning(candidate []byte) error {
	err := s.runningConfig.Call("net.vyatta.vci.config.write.Set",
		0, string(candidate)).Store()
	return err
}

func (s *dbusService) GetState(path string) ([]byte, error) {
	var result string
	err := s.runningState.Call("net.vyatta.vci.config.read.Get",
		0).Store(&result)
	if err != nil {
		fmt.Printf("FAILURE: %s\n", err.Error())
	}
	return []byte(result), err
}

func (d *dbusDispatcher) NewService(serviceName string) (yangd.Service, error) {
	service := d.dbusService(serviceName)
	return service, nil
}

func newDispatcher(conn *dbus.Conn, name string) (yangd.Dispatcher, error) {
	var err error

	if conn == nil {
		// If not specified, use system by default
		conn, err = dbus.Dial("unix:path=/var/run/vci/vci_bus_socket")
		if err != nil {
			return nil, err
		}
		if err = conn.Auth(nil); err != nil {
			conn.Close()
			return nil, err
		}
		if err = conn.Hello(); err != nil {
			conn.Close()
			return nil, err
		}
	}
	if name != "" {
		_, err := conn.RequestName(name, 0)
		if err != nil {
			return nil, err
		}
	}
	return &dbusDispatcher{
		conn: conn,
	}, nil

}

func NewDispatcherWithName(conn *dbus.Conn, name string) (yangd.Dispatcher, error) {
	return newDispatcher(conn, name)
}

func NewDispatcher(conn *dbus.Conn) (yangd.Dispatcher, error) {
	return newDispatcher(conn, "")
}
