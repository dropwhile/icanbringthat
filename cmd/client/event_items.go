// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/Masterminds/sprig/v3"

	"github.com/dropwhile/icanbringthat/internal/util"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
)

const eventItemTpl = `
{{- /* whitespace fix */ -}}
- event_item_ref_id: {{.RefId}}
  description: {{.Description}}
  created: {{.Created.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

type EventItemsAddCmd struct {
	EventRefId  string `name:"event-ref-id" arg:"" required:"" help:"event ref-id"`
	Description string `name:"description" required:"" help:"event item description"`
}

func (cmd *EventItemsAddCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.EventAddItemRequest{
		EventRefId:  cmd.EventRefId,
		Description: cmd.Description,
	}
	resp, err := client.EventAddItem(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(template.New("eventItemTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventItemTpl))
	if err := t.Execute(os.Stdout, resp.EventItem); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EventItemsUpdateCmd struct {
	RefId       string `name:"ref-id" arg:"" required:"" help:"event-item ref-id"`
	Description string `name:"description" required:"" help:"event item description"`
}

func (cmd *EventItemsUpdateCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.EventUpdateItemRequest{
		RefId:       cmd.RefId,
		Description: cmd.Description,
	}

	resp, err := client.EventUpdateItem(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(template.New("eventItemTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventItemTpl))
	if err := t.Execute(os.Stdout, resp.EventItem); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EventItemsRemoveCmd struct {
	RefId string `name:"ref-id" arg:"" required:"" help:"event-item ref-id"`
}

func (cmd *EventItemsRemoveCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.EventRemoveItemRequest{
		RefId: cmd.RefId,
	}
	if _, err := client.EventRemoveItem(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}
