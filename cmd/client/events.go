package main

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/muesli/reflow/indent"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dropwhile/icbt/rpc/icbt"
)

const eventTpl = `
{{- /* whitespace fix */ -}}
- ref_id: {{.RefId}}
  name: {{.Name}}
  description: {{.Description}}
  archived: {{.Archived}}
  when: {{.When.Ts.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
  tz: {{.When.Tz}}
  created: {{.Created.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

type EventsListCmd struct {
	Archived bool `name:"archived" help:"show archived events"`
}

func (cmd *EventsListCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.ListEventsRequest{
		Archived: &cmd.Archived,
	}
	resp, err := client.ListEvents(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(template.New("eventTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventTpl))
	for _, event := range resp.Events {
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
	req := &icbt.CreateEventRequest{
		Name:        cmd.Name,
		Description: cmd.Description,
		When: &icbt.TimestampTZ{
			Ts: timestamppb.New(cmd.When),
			Tz: cmd.Tz,
		},
	}
	resp, err := client.CreateEvent(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(template.New("eventTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventTpl))
	if err := t.Execute(os.Stdout, resp.Event); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EventsUpdateCmd struct {
	Name        *string    `name:"name"  help:"event name"`
	Description *string    `name:"description"  help:"event description"`
	When        *time.Time `name:"when"  help:"event start time"`
	Tz          *string    `name:"tz"  help:"event timezone"`
	RefID       string     `name:"ref-id" required:""`
}

func (cmd *EventsUpdateCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.UpdateEventRequest{
		RefId: cmd.RefID,
	}
	if cmd.Name != nil {
		req.Name = cmd.Name
	}
	if cmd.Description != nil {
		req.Description = cmd.Description
	}
	if (cmd.When != nil && cmd.Tz == nil) ||
		(cmd.When == nil && cmd.Tz != nil) {
		return fmt.Errorf("either both or neither of when and tz are required")
	}
	if cmd.When != nil {
		req.When = &icbt.TimestampTZ{
			Ts: timestamppb.New(*cmd.When),
			Tz: *cmd.Tz,
		}
	}
	if cmd.Name == nil && cmd.Description == nil && cmd.When == nil {
		return fmt.Errorf("at least one field must be included to update anything")
	}

	if _, err := client.UpdateEvent(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}

type EventsDeleteCmd struct {
	RefID string `name:"ref-id" required:""`
}

func (cmd *EventsDeleteCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.DeleteEventRequest{
		RefId: cmd.RefID,
	}
	if _, err := client.DeleteEvent(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}

type EventsGetDetailsCmd struct {
	RefID string `name:"ref-id" required:""`
}

func (cmd *EventsGetDetailsCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.GetEventDetailsRequest{
		RefId: cmd.RefID,
	}
	resp, err := client.GetEventDetails(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	outWriter := indent.NewWriterPipe(os.Stdout, 2, nil)

	fmt.Println("event:")
	t := template.Must(template.New("eventTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventTpl))
	if err := t.Execute(outWriter, resp.Event); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	fmt.Println("items:")
	if len(resp.Items) > 0 {
		t2 := template.Must(template.New("eventItemTpl").
			Funcs(sprig.FuncMap()).
			Parse(eventItemTpl))
		for _, item := range resp.Items {
			if err := t2.Execute(outWriter, item); err != nil {
				return fmt.Errorf("executing template: %w", err)
			}
		}
	}

	fmt.Println("earmarks:")
	if len(resp.Earmarks) > 0 {
		t2 := template.Must(template.New("earmarkTpl").
			Funcs(sprig.FuncMap()).
			Parse(earmarkTpl))
		for _, earmark := range resp.Earmarks {
			if err := t2.Execute(outWriter, earmark); err != nil {
				return fmt.Errorf("executing template: %w", err)
			}
		}
	}
	return nil
}

type EventsListItemsCmd struct {
	RefID string `name:"ref-id" required:""`
}

func (cmd *EventsListItemsCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.ListEventItemsRequest{
		RefId: cmd.RefID,
	}
	resp, err := client.ListEventItems(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	outWriter := indent.NewWriterPipe(os.Stdout, 2, nil)

	fmt.Println("items:")
	if len(resp.Items) > 0 {
		t2 := template.Must(template.New("eventItemTpl").
			Funcs(sprig.FuncMap()).
			Parse(eventItemTpl))
		for _, item := range resp.Items {
			if err := t2.Execute(outWriter, item); err != nil {
				return fmt.Errorf("executing template: %w", err)
			}
		}
	}
	return nil
}

type EventsListEarmarksCmd struct {
	RefID string `name:"ref-id" required:""`
}

func (cmd *EventsListEarmarksCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.ListEventEarmarksRequest{
		RefId: cmd.RefID,
	}
	resp, err := client.ListEventEarmarks(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	outWriter := indent.NewWriterPipe(os.Stdout, 2, nil)

	fmt.Println("earmarks:")
	if len(resp.Earmarks) > 0 {
		t2 := template.Must(template.New("earmarkTpl").
			Funcs(sprig.FuncMap()).
			Parse(earmarkTpl))
		for _, earmark := range resp.Earmarks {
			if err := t2.Execute(outWriter, earmark); err != nil {
				return fmt.Errorf("executing template: %w", err)
			}
		}
	}
	return nil
}
