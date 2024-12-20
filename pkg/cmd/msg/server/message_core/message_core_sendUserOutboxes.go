/*
 * Copyright (c) 2021-present,  NovaChat-Engine.
 *  All rights reserved.
 *
 * @Author: Coder (coderxw@gmail.com)
 * @Time :
 * @File :
 */

package message_core

import (
	"context"
	"fmt"
	"novachat_engine/mtproto"
	msgService "novachat_engine/pkg/cmd/msg/rpc_client"
	syncClient "novachat_engine/pkg/cmd/sync/rpc_client"
	"novachat_engine/pkg/log"
	"novachat_engine/service/constants"
)

func (m *MessageCoreService) SendOutboxUserMessages(authKeyId int64, userId int64, peerId int64, list []*msgService.SendMessageData, draft bool) (*mtproto.Updates, error) {
	log.Debugf("SendOutboxUserMessages userId:%d peerId:%d", userId, peerId)

	ptsList, globalMessageIdList, err := m.outboxesCore.SendMessages(userId, peerId, constants.PeerTypeUser, list, draft)
	if err != nil {
		log.Errorf("SendOutboxUserMessages userId:%d peerId:%d error:%s", userId, peerId, err.Error())
		return nil, fmt.Errorf("SendOutboxUserMessages userId:%d peerId:%d error", userId, peerId)
	}

	if userId != peerId {
		err = m.SendUserMessage(&msgService.SendMessages{
			FromUserId:          userId,
			PeerId:              peerId,
			PeerType:            constants.PeerTypeUser.ToInt32(),
			DataList:            list,
			GlobalMessageIdList: globalMessageIdList,
		})
		if err != nil {
			log.Errorf("SendOutboxUserMessages SendUserMessage userId:%d peerId:%d error:%s", userId, peerId, err.Error())
			return nil, err
		}
	}

	messageIdList := make([]int32, 0, len(list))
	updates := &mtproto.Updates{
		Updates: make([]*mtproto.Update, 0, len(ptsList)*2),
	}
	for idx := range list {
		messageIdList = append(messageIdList, list[idx].Message.Id)
		//  updateMessageID#4e90bfd6 id:int random_id:long = Update;
		updates.Updates = append(updates.Updates, mtproto.NewTLUpdateMessageID(&mtproto.Update{
			Id4E90BFD671: list[idx].Message.Id,
			RandomId:     list[idx].RandomId,
		}).To_Update())
		//  updateNewMessage#1f2b0afd message:Message pts:int pts_count:int = Update;
		updates.Updates = append(updates.Updates, mtproto.NewTLUpdateNewMessage(&mtproto.Update{
			Message1F2B0AFD71: list[idx].Message,
			Pts:               ptsList[idx],
			PtsCount:          1,
		}).To_Update())
	}

	updates = mtproto.NewTLUpdates(updates).To_Updates()
	_, err = syncClient.GetSyncClientById(userId).ReqSyncUpdate(context.TODO(), &syncClient.SyncUpdate{
		UserId:          userId,
		IgnoreAuthKeyId: authKeyId,
		Updates:         updates,
		PeerType:        constants.PeerTypeUser.ToInt32(),
	})
	if err != nil {
		log.Warnf("SendOutboxUserMessages userId:%d peerId:%d error:%s", userId, peerId, err.Error())
	}

	log.Infof("SendOutboxUserMessages userId:%d peerId:%d messageIdList:%v", userId, peerId, messageIdList)
	return updates, nil
}
