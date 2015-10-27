package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var bodyMsg = Event{
	Category:  "button",
	Action:    "click",
	Label:     "nav_1",
	MetaData:  map[string]string{"from_user": "1"},
	CreatedAt: time.Now(),
}

// Http服务测试
func Test_HttpService(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(insertMongoByHttp))
	defer ts.Close()
	b, _ := json.Marshal(bodyMsg)
	body := bytes.NewBuffer([]byte(b))
	resp, err := http.Post(ts.URL, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		t.Log("Http服务测试成功")
	} else {
		t.Error("Http服务测试失败")
	}
}

//MongoDB插入测试
func Test_MongoDbInsert(t *testing.T) {
	mongoCollection.DropCollection()
	var prevCount, afterCount int
	prevCount, _ = mongoCollection.Count()
	b, _ := json.Marshal(bodyMsg)
	insertEvent(b)
	afterCount, _ = mongoCollection.Count()
	log.Println("PrevCount", prevCount)
	log.Println("AfterCount", afterCount)
	if val := afterCount - prevCount; val == 1 {
		t.Log("MongoDB插入测试成功")
	} else {
		t.Log("MongoDB插入测试失败")
	}
}
