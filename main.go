package main

import (
  "os"
  "fmt"
  "log"
  "time"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/redis.v3"
)

var mongoCollection = getMongoCollection()
var client = redis.NewClient(&redis.Options{Addr: initRedisUri(), DB: 12})

type Event struct {
  Category string `json:"category"`
  Action string `json:"action"`
  Label string `json:"label"`
  MetaData map[string]string `json:"meta_data"`
  CreatedAt time.Time
}

func initRedisUri() string{
  redis_host := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	redis_port := os.Getenv("REDIS_PORT_6379_TCP_PORT")
	if len(redis_host) == 0{
		redis_host = "localhost"
	}
	if len(redis_port) == 0{
		redis_port = "6379"
	}
	redis_info := fmt.Sprintf("%s:%s", redis_host, redis_port)
	return redis_info
}

//监听队列消息
func listenRedisChannel(channel string){
  pubsub, err := client.Subscribe(channel)
  if err != nil {
    panic(err)
  }
  log.Println("监听Redis队列 ", channel)
  for {
    msg, err := pubsub.ReceiveMessage()
    if err != nil{
      log.Println(err)
    }
    data := []byte(msg.Payload)
    go insertEvent(data);
  }
}

func getMongoCollection() * mgo.Collection{
  session, err := mgo.Dial(initMongoAddr())
  if err != nil{
    log.Fatal(err)
  }

  session.SetMode(mgo.Monotonic, true)

  return session.DB("hxf_server").C("events")
}

func initMongoAddr() string{
  mongo_host := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	mongo_port := os.Getenv("MONGO_PORT_27017_TCP_PORT")
	if len(mongo_host) == 0{
		mongo_host = "localhost"
	}
	if len(mongo_port) == 0{
		mongo_port = "27017"
	}
	mongo_info := fmt.Sprintf("%s:%s", mongo_host, mongo_port)
	return mongo_info
}

func insertEvent(data []byte) error{
  var event Event

  err := json.Unmarshal(data, &event)
  event.CreatedAt = time.Now()
  if err != nil{
    return err
  }

  err = mongoCollection.Insert(&event)

  if err != nil {
    return err
  }

  log.Printf("插入事件 %+v\n", event)
  return nil
}

func main() {
  //监听事件队列
  listenRedisChannel("hxf.push.events")
}
