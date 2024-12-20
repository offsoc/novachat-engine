
### Novachat_engine
Novachat-engine is a open source mtproto server for Telegram compatible Android, IOS, tdesktop(Mac/Windows/Linux/Web).

## Introduce [Telegram](https://telegram.org/)

#### Novachat_engine server
    Novachat_engine server base layer 122.

##### Features
- contacts
- sms
- channel/megagroup(unsupport broadcast)
- voice call
- sticketsets

#### ENV
    Golang version: GO 1.17.11

    go mod tidy -go=1.16 && go mod tidy -go=1.17

### Support Components
    mongodb 4.4.19
    mysql 5.7
    kafka (you can find in $ROOT_DIR/lib)
    redis (must redisbloom found in $ROOT_DIR/lib)
    etcd

#### Installing
- [Docker](https://github.com/novachat-engine-team/novachat-engine/blob/main/install/install_docker.md)
- [Ubuntu](https://github.com/novachat-engine-team/novachat-engine/blob/main/install/install_ubuntu.md)
- [CentOS7](https://github.com/novachat-engine-team/novachat-engine/blob/main/install/install_centos-7.md)


    MySQL default: root/123456

- Golang 1.17.11

      // install go
      wget https://go.dev/dl/go1.17.11.linux-amd64.tar.gz
      tar xvf go1.17.11.linux-amd64.tar.gz
      sudo mv go /opt/
      export PATH=$PATH:/opt/go/bin

#### Init DB

###### Docker Deploy

  - MySQL [enterprise.sql](https://github.com/novachat-engine-team/novachat-engine/blob/main/scripts/enterprise.sql)

        docker cp enterprise.sql mysql:/
  
        docker exec -it mysql bash
  
        mysql -uroot -h127.0.0.1 -p < enterprise.sql
  
  - MongoDB
  
        docker exec -it mongodb bash
        // start mongo client
        mongo

        // init mongo rs
        rs.initiate()

###### System Deploy

  - MySQL [enterprise.sql](https://github.com/novachat-engine-team/novachat-engine/blob/main/scripts/enterprise.sql)
  
        mysql -uroot -h127.0.0.1 -p < enterprise.sql
     
  - MongoDB
     
        // start mongo client
        mongo

        // init mongo rs
        rs.initiate()

  - Redis [redis.md](https://github.com/novachat-engine-team/novachat-engine/blob/main/doc/redis/redis.md)

#### Building and Running
    git clone https://github.com/novachat-engine-team/novachat-engine

    1. cd novachat-engine/gen
    2. mkdir out
    3. ./config.sh
    4. ./build.sh all
        Found error in building `INSTALL_fix`(https://github.com/novachat-engine-team/novachat-engine/blob/main/INSTALL_fix)

    5. cp start.sh restart.sh status.sh stop.sh out
    6. sudo bash -c "cat >> /etc/hosts" << EOF
    127.0.0.1 mysql-hosts 
    127.0.0.1 kafka-hosts 
    127.0.0.1 redis-hosts 
    127.0.0.1 etcd-hosts 
    127.0.0.1 mongodb-hosts
    EOF

    7. cd out
    8. ./start.sh all
       start serve
    9. ./status.sh
       check serve is running, `ls log` all serve logs
       eg:
          check auth
          test       20683       1  0 10:43 pts/0    00:00:00 ./auth auth.yaml
          check biz_server
          test       20684       1  0 10:43 pts/0    00:00:00 ./biz_server biz_server.yaml
          check gate
          test       20685       1  0 10:43 pts/0    00:00:00 ./gate gate.yaml
          check chat
          test       21010       1  0 10:44 pts/0    00:00:00 ./chat chat.yaml
          check session
          test       20687       1  0 10:43 pts/0    00:00:00 ./session session.yaml
          check msg
          test       20688       1  0 10:43 pts/0    00:00:00 ./msg msg.yaml
          check sync
          test       20689       1  0 10:43 pts/0    00:00:00 ./sync sync.yaml
          check sfs
          test       20690       1  0 10:43 pts/0    00:00:00 ./sfs sfs.yaml
          check relay
          test       20691       1  0 10:43 pts/0    00:00:00 ./relay relay.yaml

    10. connecting xxxx:12347
      xxxx is your IP

tips: building error fix in [INSTALL_fix](https://github.com/novachat-engine-team/novachat-engine/blob/main/INSTALL_fix)

#### Notes
Coming soon

- web
- videocall
- audiocall group call
- video group call(need LAYER 131)
- bots


#### Author 
PM [Author](https://t.me/bigM1223)

Groups [Telegram](https://t.me/novachat_telegram)
