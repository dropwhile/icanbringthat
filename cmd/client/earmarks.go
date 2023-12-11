package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/Masterminds/sprig/v3"

	"github.com/dropwhile/icbt/rpc/icbt"
)

const earmarkTpl = `
{{- /* whitespace fix */ -}}
- ref_id: {{.RefId}}
  event_item_ref_id: {{.EventItemRefId}}
  note: {{.Note}}
  owner: {{.Owner}}
  created: {{.Created.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

const earmarkDetailTpl = `
{{- /* whitespace fix */ -}}
- ref_id: {{.Earmark.RefId}}
  event_item_ref_id: {{.Earmark.EventItemRefId}}
  event_ref_id: {{.EventRefId}}
  note: {{.Earmark.Note}}
  owner: {{.Earmark.Owner}}
  created: {{.Earmark.Created.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

type EarmarksCreateCmd struct {
	EventItemRefID string `name:"event-item-ref-id" required:"" help:"event item ref-id"`
	Note           string `name:"note" required:"" help:"earmark note"`
}

func (cmd *EarmarksCreateCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.CreateEarmarkRequest{
		EventItemRefId: cmd.EventItemRefID,
		Note:           cmd.Note,
	}
	resp, err := client.CreateEarmark(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(template.New("earmarkTpl").
		Funcs(sprig.FuncMap()).
		Parse(earmarkTpl))
	if err := t.Execute(os.Stdout, resp.Earmark); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EarmarksGetDetailsCmd struct {
	RefID string `name:"ref-id" required:""`
}

func (cmd *EarmarksGetDetailsCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.GetEarmarkDetailsRequest{
		RefId: cmd.RefID,
	}
	resp, err := client.GetEarmarkDetails(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t2 := template.Must(template.New("earmarkDetailTpl").
		Funcs(sprig.FuncMap()).
		Parse(earmarkDetailTpl))
	if err := t2.Execute(os.Stdout,
		map[string]interface{}{
			"Earmark":    resp.Earmark,
			"EventRefId": resp.EventRefId,
		}); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type EarmarksRemoveCmd struct {
	RefID string `name:"ref-id" required:"" help:"earmark ref-id"`
}

func (cmd *EarmarksRemoveCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.RemoveEarmarkRequest{
		RefId: cmd.RefID,
	}
	if _, err := client.RemoveEarmark(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}

type EarmarksListCmd struct {
	Archived bool `name:"archived" help:"show archived events"`
}

func (cmd *EarmarksListCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.ListEarmarksRequest{
		Archived: &cmd.Archived,
	}
	resp, err := client.ListEarmarks(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t2 := template.Must(template.New("earmarkTpl").
		Funcs(sprig.FuncMap()).
		Parse(earmarkTpl))
	for _, earmark := range resp.Earmarks {
		if err := t2.Execute(os.Stdout, earmark); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	return nil
}
