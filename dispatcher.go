// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2016 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package yangd

type Dispatcher interface {
	NewService(serviceName string) (Service, error)
}

type Service interface {
	GetRunning(path string) ([]byte, error)
	ValidateCandidate(candidate []byte) error
	SetRunning(candidate []byte) error

	GetState(path string) ([]byte, error)
}
