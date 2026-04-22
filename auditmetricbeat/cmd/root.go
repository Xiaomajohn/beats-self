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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/elastic/elastic-agent-libs/mapstr"

	// 使用本地auditbeat代码
	"github.com/elastic/beats/v7/auditmetricbeat/auditbeat/ab"
	"github.com/elastic/beats/v7/auditmetricbeat/auditbeat/core"
	"github.com/elastic/beats/v7/libbeat/cmd"
	"github.com/elastic/beats/v7/libbeat/cmd/instance"
	"github.com/elastic/beats/v7/libbeat/ecs"
	"github.com/elastic/beats/v7/libbeat/processors"
	"github.com/elastic/beats/v7/libbeat/publisher/processing"
	// 使用本地metricbeat代码
	"github.com/elastic/beats/v7/auditmetricbeat/metricbeat/beater"
	"github.com/elastic/beats/v7/auditmetricbeat/metricbeat/mb/module"

	// 注册本地模块
	_ "github.com/elastic/beats/v7/auditmetricbeat/auditbeat/include"
	_ "github.com/elastic/beats/v7/auditmetricbeat/metricbeat/include"
	_ "github.com/elastic/beats/v7/auditmetricbeat/metricbeat/include/fields"
)

const (
	// Name of the beat (auditmetricbeat).
	Name = "auditmetricbeat"
)

// RootCmd for running auditmetricbeat.
var RootCmd *cmd.BeatsRootCmd

// ShowCmd to display extra information.
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show modules information",
}

// withECSVersion is a modifier that adds ecs.version to events.
var withECSVersion = processing.WithFields(mapstr.M{
	"ecs": mapstr.M{
		"version": ecs.Version,
	},
})

// AuditMetricbeatSettings contains the default settings for auditmetricbeat
func AuditMetricbeatSettings(globals processors.PluginConfig) instance.Settings {
	runFlags := pflag.NewFlagSet(Name, pflag.ExitOnError)
	return instance.Settings{
		RunFlags:      runFlags,
		Name:          Name,
		HasDashboards: true,
		Processing:    processing.MakeDefaultSupport(true, globals, withECSVersion, processing.WithHost, processing.WithAgentMeta()),
		Initialize: []func(){
			// Initialize metricbeat monitoring modules
			func() { module.RegisterMonitoringModules("module") },
		},
	}
}

// Initialize initializes the entrypoint commands for auditmetricbeat
func Initialize(settings instance.Settings) *cmd.BeatsRootCmd {
	// Create beater with both auditbeat and metricbeat registries
	create := beater.CreatorWithRegistry(
		ab.Registry,
		beater.WithModuleOptions(
			module.WithEventModifier(core.AddDatasetToEvent),
		),
	)

	rootCmd := cmd.GenRootCmdWithSettings(create, settings)
	rootCmd.AddCommand(ShowCmd)

	return rootCmd
}

func init() {
	RootCmd = Initialize(AuditMetricbeatSettings(nil))
}
