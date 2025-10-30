package rpc_org_theprotocols_session

import (
	"fmt"

	"github.com/hereus-pbc/network-core/pkg/interfaces"
	"github.com/hereus-pbc/network-core/pkg/types"
)

func GetUserId(session interfaces.Session) (*types.UserId, error) {
	userId := session.GetUser().ToUserId()
	if userId == nil {
		return &types.UserId{}, fmt.Errorf("user ID not found")
	}
	return userId, nil
}
