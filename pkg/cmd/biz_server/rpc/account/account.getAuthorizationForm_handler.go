/*
 * Copyright (c) 2021-present,  NovaChat-Engine.
 *  All rights reserved.
 *
 *
 * @Author: Coder (coderxw@gmail.com)
 * @Time : 2021/04/13 16:46
 * @File : account.getAuthorizationForm_handler.go
 * @Desc :
 *
 */

package rpc

import (
	"context"
	"fmt"
	"novachat_engine/mtproto"
	"novachat_engine/pkg/log"
	"novachat_engine/pkg/rpc/metadata"
)

//  account.getAuthorizationForm#b86ba8e1 bot_id:int scope:string public_key:string = account.AuthorizationForm;
//
func (s *AccountServiceImpl) AccountGetAuthorizationForm(ctx context.Context, request *mtproto.TLAccountGetAuthorizationForm) (*mtproto.Account_AuthorizationForm, error) {
	md := metadata.RpcMetaDataFromContext(ctx)
	log.Infof("AccountGetAuthorizationForm %v, request: %v", metadata.RpcMetaDataDebug(md), request)

	// Impl AccountGetAuthorizationForm logic

	return nil, fmt.Errorf("%s", "Not impl AccountGetAuthorizationForm")
}
