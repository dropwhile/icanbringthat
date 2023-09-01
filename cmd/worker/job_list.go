// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"strings"

	"github.com/dropwhile/icanbringthat/internal/util"
)

type JobList struct {
	jobMap map[Job]bool
	jobs   []Job
}

func NewJobList() *JobList {
	return &JobList{
		jobs:   []Job{},
		jobMap: map[Job]bool{},
	}
}

func (jl *JobList) Add(jobs ...Job) {
	jl.jobs = append(jl.jobs, jobs...)
	jl.jobs = util.Uniq[Job](jl.jobs)
	for _, v := range jobs {
		jl.jobMap[v] = true
	}
}

func (jl *JobList) AddByName(jobnames ...string) error {
	for _, v := range jobnames {
		switch name := strings.ToLower(v); name {
		case "notifier":
			jl.Add(NotifierJob)
		case "archiver":
			jl.Add(ArchiverJob)
		case "all":
			jl.Add(NotifierJob, ArchiverJob)
		default:
			return fmt.Errorf("unknown job: %s", v)
		}
	}
	return nil
}

func (jl *JobList) String() string {
	names := make([]string, 0)
	for _, job := range jl.jobs {
		names = append(names, string(job))
	}
	return strings.Join(names, ",")
}

func (jl *JobList) Contains(job Job) bool {
	return jl.jobMap[job]
}
