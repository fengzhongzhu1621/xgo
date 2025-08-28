package handlers

import (
	"context"

	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/action"
	"github.com/fengzhongzhu1621/xgo/collections/messagebus/dew/commands/query"
)

// OrgHandler is a handler for organization commands.
type OrgHandler struct{}

// NewOrgHandler creates a new organization handler.
func NewOrgHandler() *OrgHandler {
	return &OrgHandler{}
}

func (h *OrgHandler) UpdateOrg(_ context.Context, command *action.UpdateOrgAction) error {
	return nil
}

func (h *OrgHandler) GetOrgDetails(_ context.Context, command *query.GetOrgDetailsQuery) error {
	command.Result = "Success"
	return nil
}
