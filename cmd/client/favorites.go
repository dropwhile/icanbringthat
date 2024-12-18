// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/Masterminds/sprig/v3"

	"github.com/dropwhile/icanbringthat/internal/util"
	icbt "github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1"
)

const favoriteTpl = `
{{- /* whitespace fix */ -}}
- event_ref_id: {{.GetEventRefId}}
  created: {{.GetCreated.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

type FavoritesListCmd struct {
	Archived bool `name:"archived" help:"show archived events"`
}

func (cmd *FavoritesListCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.FavoriteListEventsRequest_builder{
		Archived: &cmd.Archived,
	}.Build()
	resp, err := client.FavoriteListEvents(meta.ctx, connect.NewRequest(req))
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

type FavoritesAddCmd struct {
	EventRefID string `name:"event-ref-id" arg:"" required:""`
}

func (cmd *FavoritesAddCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.FavoriteAddRequest_builder{
		EventRefId: cmd.EventRefID,
	}.Build()
	resp, err := client.FavoriteAdd(meta.ctx, connect.NewRequest(req))
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(
		template.New("favoriteTpl").
			Funcs(sprig.FuncMap()).
			Parse(strings.TrimLeft(favoriteTpl, "\n")),
	)
	if err := t.Execute(os.Stdout, resp.Msg.GetFavorite()); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type FavoritesRemoveCmd struct {
	EventRefID string `name:"event-ref-id" arg:"" required:""`
}

func (cmd *FavoritesRemoveCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := icbt.FavoriteRemoveRequest_builder{
		EventRefId: cmd.EventRefID,
	}.Build()
	if _, err := client.FavoriteRemove(meta.ctx, connect.NewRequest(req)); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}
