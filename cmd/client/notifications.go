package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/dropwhile/icbt/rpc/icbt"
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
	req := &icbt.ListNotificationsRequest{}
	resp, err := client.ListNotifications(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(template.New("notifTpl").
		Funcs(funcMap).
		Parse(notifTpl))
	for _, notif := range resp.Notifications {
		if err := t.Execute(os.Stdout, notif); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	return nil
}

type NotificationsDeleteCmd struct {
	RefID string `name:"ref-id" required:""`
}

func (cmd *NotificationsDeleteCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.DeleteNotificationRequest{
		RefId: cmd.RefID,
	}
	if _, err := client.DeleteNotification(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	return nil
}

type NotificationsDeleteAllCmd struct{}

func (cmd *NotificationsDeleteAllCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &icbt.DeleteAllNotificationsRequest{}
	if _, err := client.DeleteAllNotifications(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}
