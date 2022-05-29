package server

import (
	"fmt"
	"net/url"

	"github.com/devict/job-board/pkg/config"
	"github.com/devict/job-board/pkg/data"
)

func signatureForItem(item data.DataModel, secret string) string {
	return item.Hash(secret)
}

func signedRoute(itemType string, action string, item data.DataModel, c config.Config) string {
	return fmt.Sprintf(
		"%s/%s/%s/%s?token=%s",
		c.URL,
		itemType,
		item.ItemID(),
		action,
		url.QueryEscape(signatureForItem(item, c.AppSecret)),
	)
}
