// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build mage
// +build mage

package main

import (
	"github.com/elastic/beats/v7/dev-tools/mage"
)

func init() {
	// Set custom beat name
	mage.SetBuildArgs("auditmetricbeat")
}

// Build builds the Beat binary.
func Build() error {
	return mage.Build(mage.DefaultBuildArgs())
}

// GolangcrossBuild cross-builds the Beat binary.
func GolangcrossBuild() error {
	return mage.GolangCrossBuild(mage.DefaultCrossBuildArgs())
}

// Clean cleans the build environment.
func Clean() error {
	return mage.Clean()
}

// Test runs unit tests.
func Test() error {
	return mage.GoTest(mage.DefaultGoTestArgs())
}

// Coverage runs tests with coverage.
func Coverage() error {
	return mage.GoTestUnitCoverage()
}

// Update updates the generated files.
func Update() error {
	return mage.Update()
}

// Fields regenerates fields.yml.
func Fields() error {
	return mage.Fields()
}

// Config regenerates config files.
func Config() error {
	return mage.Config()
}

// Package packages the Beat for distribution.
func Package() error {
	return mage.Package(mage.DefaultPackageArgs())
}

// CrossBuild cross-compiles the beat for all platforms.
func CrossBuild() error {
	return mage.CrossBuild()
}
