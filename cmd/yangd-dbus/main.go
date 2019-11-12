// Copyright (c) 2017,2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2016 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/danos/yangd/dbus"
)

func handleError(operation string, err error) {
	if err != nil {
		fmt.Printf("Failed to %s: %s", operation, err)
		os.Exit(1)
	}
}

func main() {
	serviceBus := "net.vyatta.test.service.example"
	newConfig := []byte("{\"example\":{\"description\":\"New Config\",\"server\":{\"address\":\"10.0.0.1\",\"port\":900}}}")

	dispatch, err := dbus.NewDispatcher(nil)
	if err != nil {
		panic(err)
	}

	service, err := dispatch.NewService(serviceBus)
	if err != nil {
		panic(err)
	}

	fmt.Println("Original Config:")
	config, err := service.GetRunning("")
	handleError("get config", err)
	fmt.Println(string(config))

	fmt.Println("Validating Config:")
	err = service.ValidateCandidate(newConfig)
	handleError("check config", err)

	fmt.Println("Setting Config:")
	err = service.SetRunning(newConfig)
	handleError("set config", err)

	fmt.Println("New Config:")
	config, err = service.GetRunning("")
	handleError("get config", err)
	fmt.Println(string(config))
}
