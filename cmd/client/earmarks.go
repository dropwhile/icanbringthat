// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"html/template"
	"os"

	"connectrpc.com/connect"
	"github.com/Masterminds/sprig/v3"

	"github.com/dropwhile/icanbringthat/internal/util"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

const earmarkTpl = `
{{- /* whitespace fix */ -}}
- ref_id: {{.GetRefId}}
  event_item_ref_id: {{.GetEventItemRefId}}
  note: {{.GetNote}}
  owner: {{.GetOwner}}
  created: {{.GetCreated.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

const earmarkDetailTpl = `
{{- /* whitespace fix */ -}}
- ref_id: {{.GetEarmark.GetRefId}}
  event_item_ref_id: {{.GetEarmark.GetEventItemRefId}}
  event_ref_id: {{.GetEventRefId}}
  note: {{.GetEarmark.GetNote}}
  owner: {{.GetEarmark.GetOwner}}
  created: {{.GetEarmark.GetCreated.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

type EarmarksCreateCmd struct {
	EventItemRefID string `name:"event-item-ref-id" arg:"" required:"" help:"event item ref-id"`
	Note           string `name:"note" required:"" help:"earmark note"`
}

func (cmd *EarmarksCreateCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EarmarkCreateRequest_builder{
		EventItemRefId: cmd.EventItemRefID,
		Note:           cmd.Note,
	}.Build()
	resp, err := client.EarmarkCreate(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(template.New("earmarkTpl").
		Funcs(sprig.FuncMap()).
		Parse(earmarkTpl))
	if err := t.Execute(os.Stdout, resp.Msg.GetEarmark()); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EarmarksGetDetailsCmd struct {
	RefID string `name:"ref-id" arg:"" required:""`
}

func (cmd *EarmarksGetDetailsCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EarmarkGetDetailsRequest_builder{
		RefId: cmd.RefID,
	}.Build()
	resp, err := client.EarmarkGetDetails(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t2 := util.Must(template.New("earmarkDetailTpl").
		Funcs(sprig.FuncMap()).
		Parse(earmarkDetailTpl))
	if err := t2.Execute(os.Stdout,
		map[string]any{
			"Earmark":    resp.Msg.GetEarmark(),
			"EventRefId": resp.Msg.GetEventRefId(),
		}); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EarmarksRemoveCmd struct {
	RefID string `name:"ref-id" arg:"" required:"" help:"earmark ref-id"`
}

func (cmd *EarmarksRemoveCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EarmarkRemoveRequest_builder{
		RefId: cmd.RefID,
	}.Build()
	if _, err := client.EarmarkRemove(meta.ctx, connect.NewRequest(req)); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}

type EarmarksListCmd struct {
	Archived bool `name:"archived" help:"show archived events"`
}

func (cmd *EarmarksListCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.EarmarksListRequest_builder{
		Archived: &cmd.Archived,
	}.Build()
	resp, err := client.EarmarksList(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t2 := util.Must(template.New("earmarkTpl").
		Funcs(sprig.FuncMap()).
		Parse(earmarkTpl))
	for _, earmark := range resp.Msg.GetEarmarks() {
		if err := t2.Execute(os.Stdout, earmark); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	return nil
}
