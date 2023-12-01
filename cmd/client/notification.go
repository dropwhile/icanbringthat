package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/dropwhile/icbt/rpc"
)

type NotificationListCmd struct{}

const notifTpl = `- ref_id: {{.RefId}}
  message: {{.Message}}
  created: {{.Created.AsTime.Format "2006-01-02T15:04:05Z07:00" }}
`

func (cmd *NotificationListCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &rpc.ListNotificationsRequest{}
	resp, err := client.ListNotifications(meta.ctx, req)
	if err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	t := template.Must(template.New("notifTpl").Parse(notifTpl))
	for _, notif := range resp.Notifications {
		if err := t.Execute(os.Stdout, notif); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}
	}
	return nil
}

type NotificationDeleteCmd struct {
	RefID string `nane:"ref-id" required:""`
}

func (cmd *NotificationDeleteCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &rpc.DeleteNotificationRequest{
		RefId: cmd.RefID,
	}
	if _, err := client.DeleteNotification(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}

	return nil
}

type NotificationDeleteAllCmd struct{}

func (cmd *NotificationDeleteAllCmd) Run(meta *RunArgs) error {
	client := meta.client
	req := &rpc.DeleteAllNotificationsRequest{}
	if _, err := client.DeleteAllNotifications(meta.ctx, req); err != nil {
		return fmt.Errorf("client request: %w", err)
	}
	return nil
}
