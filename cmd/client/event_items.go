package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/dropwhile/icbt/rpc/icbt"
)

const eventItemTpl = `
{{- /* whitespace fix */ -}}
- event_item_ref_id: {{.RefId}}
  description: {{.Description}}
  created: {{.Created.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

type EventItemsAddCmd struct {
	EventRefId  string `name:"event-ref-id" required:"" help:"event ref-id"`
	Description string `name:"description" required:"" help:"event item description"`
}

func (cmd *EventItemsAddCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.AddEventItemRequest{
		EventRefId:  cmd.EventRefId,
		Description: cmd.Description,
	}
	resp, err := client.AddEventItem(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(template.New("eventItemTpl").
		Funcs(funcMap).
		Parse(eventItemTpl))
	if err := t.Execute(os.Stdout, resp.EventItem); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EventItemsUpdateCmd struct {
	RefId       string `name:"ref-id" required:"" help:"event-item ref-id"`
	Description string `name:"description" required:"" help:"event item description"`
}

func (cmd *EventItemsUpdateCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.UpdateEventItemRequest{
		RefId:       cmd.RefId,
		Description: cmd.Description,
	}

	resp, err := client.UpdateEventItem(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(template.New("eventItemTpl").
		Funcs(funcMap).
		Parse(eventItemTpl))
	if err := t.Execute(os.Stdout, resp.EventItem); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EventItemsRemoveCmd struct {
	RefId string `name:"ref-id" required:"" help:"event-item ref-id"`
}

func (cmd *EventItemsRemoveCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.RemoveEventItemRequest{
		RefId: cmd.RefId,
	}
	if _, err := client.RemoveEventItem(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}
