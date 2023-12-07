// Copyright (C) 2020-2023 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package manageability

import (
	"strings"

	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/common"
	"github.com/test-network-function/cnf-certification-test/cnf-certification-test/identifiers"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	"github.com/test-network-function/cnf-certification-test/pkg/checksdb"
	"github.com/test-network-function/cnf-certification-test/pkg/provider"
	"github.com/test-network-function/cnf-certification-test/pkg/testhelper"
)

var (
	env provider.TestEnvironment

	beforeEachFn = func(check *checksdb.Check) error {
		check.LogInfo("Check %s: getting test environment.", check.ID)
		env = provider.GetTestEnvironment()
		return nil
	}

	skipIfNoContainersFn = func() (bool, string) {
		if len(env.Containers) == 0 {
			log.Warn("No containers to check...")
			return true, "There are no containers to check. Please check under test labels."
		}
		return false, ""
	}
)

func LoadChecks() {
	log.Debug("Entering %s suite", common.ManageabilityTestKey)

	checksGroup := checksdb.NewChecksGroup(common.ManageabilityTestKey).
		WithBeforeEachFn(beforeEachFn)

	testID, tags := identifiers.GetGinkgoTestIDAndLabels(identifiers.TestContainersImageTag)
	checksGroup.Add(checksdb.NewCheck(testID, tags).
		WithSkipCheckFn(skipIfNoContainersFn).
		WithCheckFn(func(c *checksdb.Check) error {
			testContainersImageTag(c, &env)
			return nil
		}))

	testID, tags = identifiers.GetGinkgoTestIDAndLabels(identifiers.TestContainerPortNameFormat)
	checksGroup.Add(checksdb.NewCheck(testID, tags).
		WithSkipCheckFn(skipIfNoContainersFn).
		WithCheckFn(func(c *checksdb.Check) error {
			testContainerPortNameFormat(c, &env)
			return nil
		}))
}

func testContainersImageTag(check *checksdb.Check, env *provider.TestEnvironment) {
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, cut := range env.Containers {
		check.LogDebug("check that %s image is tagged", cut)
		if cut.IsTagEmpty() {
			nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container is missing image tag(s)", false))
			check.LogDebug("Container %s is missing image tag(s)", cut.String())
		} else {
			compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "Container is tagged", true))
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}

// The name field in the ContainerPort section must be of the form <protocol>[-<suffix>] where <protocol> is one of the following,
// and the optional <suffix> can be chosen by the application. Allowed protocol names: grpc, grpc-web, http, http2, tcp, udp.
var allowedProtocolNames = map[string]bool{"grpc": true, "http": true, "http2": true, "tcp": true, "udp": true}

func containerPortNameFormatCheck(portName string) bool {
	res := strings.Split(portName, "-")
	return allowedProtocolNames[res[0]]
}

func testContainerPortNameFormat(check *checksdb.Check, env *provider.TestEnvironment) {
	for _, newProtocol := range env.ValidProtocolNames {
		allowedProtocolNames[newProtocol] = true
	}
	var compliantObjects []*testhelper.ReportObject
	var nonCompliantObjects []*testhelper.ReportObject
	for _, cut := range env.Containers {
		for _, port := range cut.Ports {
			if !containerPortNameFormatCheck(port.Name) {
				check.LogDebug("%s: ContainerPort %s does not follow the partner naming conventions", cut, port.Name)
				nonCompliantObjects = append(nonCompliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "ContainerPort does not follow the partner naming conventions", false).
					AddField(testhelper.ContainerPort, port.Name))
			} else {
				compliantObjects = append(compliantObjects, testhelper.NewContainerReportObject(cut.Namespace, cut.Podname, cut.Name, "ContainerPort follows the partner naming conventions", true).
					AddField(testhelper.ContainerPort, port.Name))
			}
		}
	}
	check.SetResult(compliantObjects, nonCompliantObjects)
}
