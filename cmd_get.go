package main

import (
	"encoding/json"
	"fmt"
	"os"
	// "sort"

	cli "github.com/jawher/mow.cli"
)

type Kind string

const (
	PROJECT  = "projects"
	TASK     = "tasks"
	POMODORO = "pomodoros"
)

func get(config *Config) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		cmd.Spec = "[OPTIONS] [KIND] [ID]"
		var (
			kind = cmd.StringArg("KIND", PROJECT, "project | task | pomodoro")
			id   = cmd.IntArg("ID", 0, "resource identifier")
			// ascend = cmd.BoolOpt("a ascend", false, "sort tasks in ascending order")
			// limit  = cmd.IntOpt("l limit", 0, "limit returned tasks")
			asJson = cmd.BoolOpt("json", false, "output result as JSON")
		)
		cmd.Action = func() {
			store, err := NewSQLiteStore(config.DBPath)
			maybe(err)
			defer store.Close()
			switch *kind {
			case PROJECT:
				project := &Project{
					ID: int64(*id),
				}
				maybe(store.With(func(s Store) error {
					return s.ReadProject(project)
				}))
				if *asJson {
					maybe(json.NewEncoder(os.Stdout).Encode(project))
					return
				}
				Tree(*project).Write(os.Stdout, 0, Tree(*project).MaxDepth() == 0)
			case TASK:
				var tasks []*Task
				maybe(store.With(func(s Store) error {
					if *id == 0 {
						result, err := s.ReadTasks(int64(*id))
						if err != nil {
							return err
						}
						tasks = result
						return err
					} else {
						task := &Task{
							ID: int64(*id),
						}
						err := s.ReadTask(task)
						if err != nil {
							return err
						}
						tasks = []*Task{task}
						return nil
					}
				}))
				if *asJson {
					maybe(json.NewEncoder(os.Stdout).Encode(tasks))
				}
			default:
				maybe(fmt.Errorf("unknown resource: %s", *kind))
			}
		}
	}
}