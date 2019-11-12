// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2016-2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package client

type ServiceRead interface {
	Get(string) []byte
}

type ServiceWrite interface {
	Check([]byte) []error
	Set([]byte) []error
}

type RunningConfig interface {
	ServiceRead
	ServiceWrite
}

type ServiceHandle interface {
}

func RegisterService(namespace string, runningConfig RunningConfig) ServiceHandle {

	return dbusRegisterService(namespace, runningConfig)
}
