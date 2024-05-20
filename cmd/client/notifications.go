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

const notifTpl = `
{{- /* whitespace fix */ -}}
- ref_id: {{.RefId}}
  message: {{.Message}}
  created: {{.Created.AsTime.Format "2006-01-02T15:04:05Z07:00" }}
`

type NotificationsListCmd struct{}

func (cmd *NotificationsListCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.NotificationsListRequest{}
	resp, err := client.NotificationsList(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := util.Must(template.New("notifTpl").
		Funcs(sprig.FuncMap()).
		Parse(notifTpl))
	for _, notif := range resp.Notifications {
		if err := t.Execute(os.Stdout, notif); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	return nil
}

type NotificationsDeleteCmd struct {
	RefID string `name:"ref-id" arg:"" required:""`
}

func (cmd *NotificationsDeleteCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.NotificationDeleteRequest{
		RefId: cmd.RefID,
	}
	if _, err := client.NotificationDelete(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	return nil
}

type NotificationsDeleteAllCmd struct{}

func (cmd *NotificationsDeleteAllCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.Empty{}
	if _, err := client.NotificationsDeleteAll(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}
