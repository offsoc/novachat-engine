#!/bin/bash

CUR_DIR=$(cd "$(dirname "$0")";pwd)
ROOT_DIR=$CUR_DIR/../../../../

PLAT=`uname`
MAC="Darwin"
if [[ $PLAT =~ $MAC ]]; then
  PATH=$ROOT_DIR/tools/protoc/mac/:$PATH
else
  PATH=$ROOT_DIR/tools/protoc/linux/:$PATH
fi

SRC_DIR=$CUR_DIR

PB_DST_DIR=.

if [ ! -d "$PB_DST_DIR" ]; then
    mkdir $PB_DST_DIR
fi

GOGOPROTO_PATH=$ROOT_DIR:$ROOT_DIR/github.com/gogo/protobuf/:$ROOT_DIR/github.com/gogo/protobuf/protobuf/:$ROOT_DIR/schema

protoc -I=$SRC_DIR --proto_path=$GOGOPROTO_PATH:./ \
    --gogo_out=plugins=grpc,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mpkg/rpc/rpc/rpc.proto=novachat_engine/pkg/rpc/rpc,Mschema.tl.core_types.proto=novachat_engine/mtproto,Mschema.tl.sync_service.proto=novachat_engine/mtproto,Mschema.tl.sync.proto=novachat_engine/mtproto:$PB_DST_DIR \
    $SRC_DIR/*.proto

gofmt -w $PB_DST_DIR/*.go
