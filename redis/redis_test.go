package redis

import (
	"context"
	"testing"
	"time"
	"conan/log"
)

func TestExpTime(t *testing.T) {
	for i := 0; i < 10; i++ {
		num := getRandomTime()
		log.Debug("getRandomTime is %d", num)
		time.Sleep(time.Duration(num + 1))
	}
}

func TestRedis(t *testing.T) {
	newRedisCli(nil)
	err := Set(context.Background(), "test", "1111", 0)

	if err != nil {
		log.Error("Redis Set Fail err is %s", err.Error())
		return
	}

}

func TestRedisCmd(t *testing.T) {
	newRedisCli(nil)
	//
	//err := HSet(context.TODO(), "appkey", "test", 11111)
	_, err := HSets(context.TODO(), "appkey", "test1", 111, "test2", 222, "test3", 333)
	if err != nil {
		log.Error("Redis Set Fail Err is %s", err.Error())
		return
	}

	//var res string
	//err := HGet(context.TODO(), "appkey", "test", &res)
	//
	//if err != nil {
	//	log.Error("Redis Get Fail Err is %s", err.Error())
	//}

	//log.Info("Redis Get res is %s", res)
	//var str string
	//err := Get(context.TODO(), "appkey", &str)
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//err := Exist(context.TODO(), "asda")
	//fmt.Println(err)

	//affRow, err := Del(context.TODO(), "test")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(affRow)
	//
	//keys, err := HKEYS(context.TODO(), "zxczx")
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(keys)
	//val, err := HValues(context.TODO(), "qwe")
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(val)

	//a1 := model.PlatformInfo{
	//	SignUrl: "http:11.11.11.11",
	//}
	//b1, _ := a1.Marshal()
	//
	//HSet(context.TODO(), "asd", "a1", b1)
	//
	//a2 := model.PlatformInfo{
	//	SignUrl: "http:11.232",
	//}
	//b2, _ := a2.Marshal()
	//HSet(context.TODO(), "asd", "a2", b2)
	//a3 := model.PlatformInfo{
	//	SignUrl: "http:11.qweqw",
	//}
	//b3, _ := a3.Marshal()
	//HSet(context.TODO(), "asd", "a3", b3)
	//
	//res, err := HValues(context.TODO(), "asd")
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//tmp := model.PlatformInfo{}
	//
	//tmp.Unmarshal([]byte(res[1]))
	//
	//fmt.Println(tmp.SignUrl)

}
