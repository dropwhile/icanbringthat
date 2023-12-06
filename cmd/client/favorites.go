package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/dropwhile/icbt/rpc"
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
	req := &rpc.ListFavoriteEventsRequest{
		Archived: &cmd.Archived,
	}
	resp, err := client.ListFavoriteEvents(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(template.New("eventTpl").
		Funcs(funcMap).
		Parse(eventTpl))
	for _, event := range resp.Events {
		if err := t.Execute(os.Stdout, event); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	return nil
}

type FavoritesAddCmd struct {
	EventRefID string `name:"event-ref-id" required:""`
}

func (cmd *FavoritesAddCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &rpc.CreateFavoriteRequest{
		EventRefId: cmd.EventRefID,
	}
	resp, err := client.AddFavorite(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(
		template.New("favoriteTpl").
			Funcs(funcMap).
			Parse(strings.TrimLeft(favoriteTpl, "\n")),
	)
	if err := t.Execute(os.Stdout, resp.Favorite); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

type FavoritesRemoveCmd struct {
	EventRefID string `name:"event-ref-id" required:""`
}

func (cmd *FavoritesRemoveCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &rpc.RemoveFavoriteRequest{
		EventRefId: cmd.EventRefID,
	}
	if _, err := client.RemoveFavorite(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}
