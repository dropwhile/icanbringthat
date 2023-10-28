POST https://{{host}}/webhooks/pm
{
  "RecordType": "SubscriptionChange",
  "MessageID": "00000000-0000-0000-0000-000000000000",
  "ServerID": 23,
  "MessageStream": "outbound",
  "ChangedAt": "2023-10-28T00:47:26Z",
  "Recipient": "{{recipient}}",
  "Origin": "Recipient",
  "SuppressSending": true,
  "SuppressionReason": "HardBounce",
  "Tag": "welcome-email"
}

HTTP 200
