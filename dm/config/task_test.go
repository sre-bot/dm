// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	. "github.com/pingcap/check"
	"io/ioutil"
	"path"
)

func (t *testConfig) TestInvalidTaskConfig(c *C) {
	var errorTaskConfig1 = `---
name: test
task-mode: all
is-sharding: true
meta-schema: "dm_meta"
remove-meta: false
enable-heartbeat: true
timezone: "Asia/Shanghai"
ignore-checking-items: ["all"]

target-database:
  host: "127.0.0.1"
  port: 4000
  user: "root"
  password: ""

mysql-instances:
  - source-id: "mysql-replica-01"
    server-id: 101 
    black-white-list:  "instance"
    route-rules: ["sharding-route-rules-table", "sharding-route-rules-schema"]
    column-mapping-rules: ["instance-1"]
    mydumper-config-name: "global"
    loader-config-name: "global"
    syncer-config-name: "global"
`
	var errorTaskConfig2 = `---
name: test
name: test1
task-mode: all
is-sharding: true
meta-schema: "dm_meta"
remove-meta: false
enable-heartbeat: true
timezone: "Asia/Shanghai"
ignore-checking-items: ["all"]

target-database:
  host: "127.0.0.1"
  port: 4000
  user: "root"
  password: ""

mysql-instances:
  - source-id: "mysql-replica-01"
    black-white-list:  "instance"
    route-rules: ["sharding-route-rules-table", "sharding-route-rules-schema"]
    column-mapping-rules: ["instance-1"]
    mydumper-config-name: "global"
    loader-config-name: "global"
    syncer-config-name: "global"
`
	taskConfig := NewTaskConfig()
	err := taskConfig.Decode(errorTaskConfig1)
	// field server-id is not a member of TaskConfig
	c.Check(err, NotNil)
	c.Assert(err, ErrorMatches, "*line 19: field server-id not found in type config.MySQLInstance*")

	err = taskConfig.Decode(errorTaskConfig2)
	// field name duplicate
	c.Check(err, NotNil)
	c.Assert(err, ErrorMatches, "*line 3: field name already set in type config.TaskConfig*")

	filepath := path.Join(c.MkDir(), "test_invalid_task.yaml")
	configContent := []byte(`---
aaa: xxx
name: test
task-mode: all
is-sharding: true
meta-schema: "dm_meta"
remove-meta: false
enable-heartbeat: true
timezone: "Asia/Shanghai"
ignore-checking-items: ["all"]
`)
	err = ioutil.WriteFile(filepath, configContent, 0644)
	err = taskConfig.DecodeFile(filepath)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "*line 2: field aaa not found in type config.TaskConfig*")

	filepath = path.Join(c.MkDir(), "test_invalid_task.yaml")
	configContent = []byte(`---
name: test
task-mode: all
task-mode: all
is-sharding: true
meta-schema: "dm_meta"
remove-meta: false
enable-heartbeat: true
timezone: "Asia/Shanghai"
ignore-checking-items: ["all"]
`)
	err = ioutil.WriteFile(filepath, configContent, 0644)
	err = taskConfig.DecodeFile(filepath)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "*line 4: field task-mode already set in type config.TaskConfig*")
}
