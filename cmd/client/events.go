// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/Masterminds/sprig/v3"
	"github.com/muesli/reflow/indent"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icanbringthat/internal/util"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

const eventTpl = `
{{- /* whitespace fix */ -}}
- ref_id: {{.GetRefId}}
  name: {{.GetName}}
  description: {{.GetDescription}}
  archived: {{.GetArchived}}
  when: {{.GetWhen.GetTs.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
  tz: {{.GetWhen.GetTz}}
  created: {{.GetCreated.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

type EventsListCmd struct {
	Archived bool `name:"archived" help:"show archived events"`
}

func (cmd *EventsListCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EventsListRequest_builder{
		Archived: &cmd.Archived,
	}.Build()
	resp, err := client.EventsList(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(template.New("eventTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventTpl))
	for _, event := range resp.Msg.GetEvents() {
		if err := t.Execute(os.Stdout, event); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	return nil
}

type EventsCreateCmd struct {
	Name        string    `name:"name" required:"" help:"event name"`
	Description string    `name:"description" required:"" help:"event description"`
	When        time.Time `name:"when" required:"" help:"event start time"`
	Tz          string    `name:"tz" required:"" help:"event timezone"`
}

func (cmd *EventsCreateCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EventCreateRequest_builder{
		Name:        cmd.Name,
		Description: cmd.Description,
		When: icbt.TimestampTZ_builder{
			Ts: timestamppb.New(cmd.When),
			Tz: cmd.Tz,
		}.Build(),
	}.Build()
	resp, err := client.EventCreate(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(template.New("eventTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventTpl))
	if err := t.Execute(os.Stdout, resp.Msg.GetEvent()); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EventsUpdateCmd struct {
	Name        *string    `name:"name" help:"event name"`
	Description *string    `name:"description" help:"event description"`
	When        *time.Time `name:"when" help:"event start time"`
	Tz          *string    `name:"tz" help:"event timezone"`
	RefID       string     `name:"ref-id" arg:"" required:""`
}

func (cmd *EventsUpdateCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EventUpdateRequest_builder{
		RefId: cmd.RefID,
	}.Build()
	if cmd.Name != nil {
		req.SetName(*cmd.Name)
	}
	if cmd.Description != nil {
		req.SetDescription(*cmd.Description)
	}
	if (cmd.When != nil && cmd.Tz == nil) ||
		(cmd.When == nil && cmd.Tz != nil) {
		return fmt.Errorf("either both or neither of when and tz are required")
	}
	if cmd.When != nil {
		req.SetWhen(
			icbt.TimestampTZ_builder{
				Ts: timestamppb.New(*cmd.When),
				Tz: *cmd.Tz,
			}.Build(),
		)
	}
	if cmd.Name == nil && cmd.Description == nil && cmd.When == nil {
		return fmt.Errorf("at least one field must be included to update anything")
	}

	if _, err := client.EventUpdate(meta.ctx, connect.NewRequest(req)); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}

type EventsDeleteCmd struct {
	RefID string `name:"ref-id" arg:"" required:""`
}

func (cmd *EventsDeleteCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EventDeleteRequest_builder{
		RefId: cmd.RefID,
	}.Build()
	if _, err := client.EventDelete(meta.ctx, connect.NewRequest(req)); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}

type EventsGetDetailsCmd struct {
	RefID string `name:"ref-id" arg:"" required:""`
}

func (cmd *EventsGetDetailsCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EventGetDetailsRequest_builder{
		RefId: cmd.RefID,
	}.Build()
	resp, err := client.EventGetDetails(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	outWriter := indent.NewWriterPipe(os.Stdout, 2, nil)

	fmt.Println("event:")
	t := util.Must(template.New("eventTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventTpl))
	if err := t.Execute(outWriter, resp.Msg.GetEvent()); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	fmt.Println("items:")
	items := resp.Msg.GetItems()
	if len(items) > 0 {
		t2 := util.Must(template.New("eventItemTpl").
			Funcs(sprig.FuncMap()).
			Parse(eventItemTpl))
		for _, item := range items {
			if err := t2.Execute(outWriter, item); err != nil {
				return fmt.Errorf("executing template: %w", err)
			}
		}
	}

	fmt.Println("earmarks:")
	earmarks := resp.Msg.GetEarmarks()
	if len(earmarks) > 0 {
		t2 := util.Must(template.New("earmarkTpl").
			Funcs(sprig.FuncMap()).
			Parse(earmarkTpl))
		for _, earmark := range earmarks {
			if err := t2.Execute(outWriter, earmark); err != nil {
				return fmt.Errorf("executing template: %w", err)
			}
		}
	}
	return nil
}

type EventsListItemsCmd struct {
	RefID string `name:"ref-id" arg:"" required:""`
}

func (cmd *EventsListItemsCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EventListItemsRequest_builder{
		RefId: cmd.RefID,
	}.Build()
	resp, err := client.EventListItems(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	outWriter := indent.NewWriterPipe(os.Stdout, 2, nil)

	fmt.Println("items:")
	items := resp.Msg.GetItems()
	if len(items) > 0 {
		t2 := util.Must(template.New("eventItemTpl").
			Funcs(sprig.FuncMap()).
			Parse(eventItemTpl))
		for _, item := range items {
			if err := t2.Execute(outWriter, item); err != nil {
				return fmt.Errorf("executing template: %w", err)
			}
		}
	}
	return nil
}

type EventsListEarmarksCmd struct {
	RefID string `name:"ref-id" arg:"" required:""`
}

func (cmd *EventsListEarmarksCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EventListEarmarksRequest_builder{
		RefId: cmd.RefID,
	}.Build()
	resp, err := client.EventListEarmarks(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	outWriter := indent.NewWriterPipe(os.Stdout, 2, nil)

	fmt.Println("earmarks:")
	earmarks := resp.Msg.GetEarmarks()
	if len(earmarks) > 0 {
		t2 := util.Must(template.New("earmarkTpl").
			Funcs(sprig.FuncMap()).
			Parse(earmarkTpl))
		for _, earmark := range earmarks {
			if err := t2.Execute(outWriter, earmark); err != nil {
				return fmt.Errorf("executing template: %w", err)
			}
		}
	}
	return nil
}
