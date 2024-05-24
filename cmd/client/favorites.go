// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/Masterminds/sprig/v3"

	"github.com/dropwhile/icanbringthat/internal/util"
	"github.com/dropwhile/icanbringthat/rpc/icbt"
)

const favoriteTpl = `
{{- /* whitespace fix */ -}}
- event_ref_id: {{.EventRefId}}
  created: {{.Created.AsTime.Format "2006-01-02T15:04:05Z07:00"}}
`

type FavoritesListCmd struct {
	Archived bool `name:"archived" help:"show archived events"`
}

func (cmd *FavoritesListCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.FavoriteListEventsRequest{
		Archived: &cmd.Archived,
	}
	resp, err := client.FavoriteListEvents(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(template.New("eventTpl").
		Funcs(sprig.FuncMap()).
		Parse(eventTpl))
	for _, event := range resp.Events {
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
	req := &icbt.FavoriteCreateRequest{
		EventRefId: cmd.EventRefID,
	}
	resp, err := client.FavoriteAdd(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(
		template.New("favoriteTpl").
			Funcs(sprig.FuncMap()).
			Parse(strings.TrimLeft(favoriteTpl, "\n")),
	)
	if err := t.Execute(os.Stdout, resp.Favorite); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type FavoritesRemoveCmd struct {
	EventRefID string `name:"event-ref-id" arg:"" required:""`
}

func (cmd *FavoritesRemoveCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.FavoriteRemoveRequest{
		EventRefId: cmd.EventRefID,
	}
	if _, err := client.FavoriteRemove(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}
