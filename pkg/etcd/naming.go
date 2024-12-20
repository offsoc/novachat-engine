/*
 * Copyright (c) 2021-present,  NovaChat-Engine.
 *  All rights reserved.
 *
 * @Author: Coder (coderxw@gmail.com)
 * @Time : 2021/3/22 23:11
 * @File : naming.go
 */

package etcd

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"novachat_engine/pkg/config"
	"novachat_engine/pkg/log"
	"strings"
	"time"
)

// etcd naming

type EtcdNamingServer struct {
	cc         resolver.ClientConn
	cli        *clientv3.Client
	etcdConfig config.EtcdServerConfig
}

func NewEtcdNamingServer(config config.EtcdServerConfig) *EtcdNamingServer {
	return &EtcdNamingServer{
		cc:         nil,
		cli:        nil,
		etcdConfig: config,
	}
}

func (s *EtcdNamingServer) Register(name, addr string) error {
	var err error

	if s.cli == nil {
		if s.etcdConfig.MaxCallSendMsgSize <= 0 {
			s.etcdConfig.MaxCallSendMsgSize = grpcMaxSendMsgSize
		}
		if s.etcdConfig.MaxCallRecvMsgSize <= 0 {
			s.etcdConfig.MaxCallRecvMsgSize = grpcMaxCallMsgSize
		}

		if s.etcdConfig.DialTimeout <= 0 {
			s.etcdConfig.DialTimeout = 3
		}

		s.cli, err = clientv3.New(clientv3.Config{
			Endpoints:            strings.Split(s.etcdConfig.EtcdAddr, ";"),
			AutoSyncInterval:     0,
			DialTimeout:          time.Duration(s.etcdConfig.DialTimeout) * time.Second,
			DialKeepAliveTime:    time.Duration(s.etcdConfig.DialKeepAliveTime) * time.Second,
			DialKeepAliveTimeout: time.Duration(s.etcdConfig.DialKeepAliveTimeout) * time.Second,
			MaxCallSendMsgSize:   int(s.etcdConfig.MaxCallSendMsgSize),
			MaxCallRecvMsgSize:   int(s.etcdConfig.MaxCallRecvMsgSize),
			TLS:                  nil,
			Username:             s.etcdConfig.UserName,
			Password:             s.etcdConfig.Password,
			RejectOldCluster:     false,
			DialOptions: []grpc.DialOption{
				grpc.WithInsecure(),
				grpc.WithInitialWindowSize(grpcInitialWindowSize),
				grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(s.etcdConfig.MaxCallRecvMsgSize))),
				grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(s.etcdConfig.MaxCallSendMsgSize))),
				grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
				grpc.WithKeepaliveParams(keepalive.ClientParameters{
					Time:                grpcKeepAliveTime,
					Timeout:             grpcKeepAliveTimeout,
					PermitWithoutStream: true,
				}),
			},
			LogConfig:           nil,
			Context:             nil,
			PermitWithoutStream: true,
		})

		if err != nil {
			return err
		}
	}
	go func() {
		for {
			leaseKeepAliveResponse, err1 := s.withAlive(name, addr, s.etcdConfig.TTL)
			if err1 != nil {
				log.Fatalf("withAlive etcd error:%s", err1.Error())
				time.Sleep(3)
				continue
			}
			for _ = range leaseKeepAliveResponse {
				// do nothing
			}
			log.Errorf("withAlive leaseKeepAliveResponse close")
		}
	}()

	return nil
}

func (s *EtcdNamingServer) withAlive(name string, addr string, ttl int64) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	leaseResp, err := s.cli.Grant(context.Background(), ttl)
	if err != nil {
		return nil, err
	}

	_, err = s.cli.Put(context.Background(), MakeServerKey(s.etcdConfig.Schema, name, addr), addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return nil, err
	}

	var leaseKeepAliveResponse <-chan *clientv3.LeaseKeepAliveResponse
	leaseKeepAliveResponse, err = s.cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return nil, err
	}
	return leaseKeepAliveResponse, nil
}

// UnRegister remove service from etcd
func (s *EtcdNamingServer) UnRegister(serverName string, addr string) {
	if s.cli != nil {
		closeChan := make(chan struct{})
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		go func() {
			s.cli.Delete(ctx, MakeServerKey(s.etcdConfig.Schema, serverName, addr))
			closeChan <- struct{}{}
		}()

		select {
		case <-closeChan:
			break
		case <-time.NewTimer(time.Second * 2).C:
			cancel()
			break
		}

		s.cli.Close()
	}
}
