package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dropwhile/icbt/internal/util"
)

type JobList struct {
	jobs   []Job
	jobMap map[Job]bool
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

func (jl *JobList) LogValuer() slog.Value {
	names := make([]string, 0)
	for _, job := range jl.jobs {
		names = append(names, string(job))
	}
	return slog.StringValue(strings.Join(names, ","))
}

func (jl *JobList) Contains(job Job) bool {
	return jl.jobMap[job]
}
