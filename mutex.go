package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Mutex struct {
	client *redis.Client
	name   string
	value  string
	expire time.Duration
}

func NewMutex(client *redis.Client, name string, value string, expire time.Duration) *Mutex {
	return &Mutex{
		client: client,
		name:   name,
		value:  value,
		expire: expire,
	}
}

func (m *Mutex) Lock() bool {
	res, err := m.client.SetNX(context.Background(), m.name, m.value, m.expire).Result()
	//setnex 只有在键不存在的时候才进行设置。如果键名已经存在不会设置键，而是直接返回一个bool值。所以可以保证只有一个客户端可以获得锁
	//SexNX 参数1：一个上下文对象用于调用redis的上下文（当前状态），参数2 ：键的名字，参数3:键的值，参数4:过期时间
	return err == nil && res
}

func (m *Mutex) Unlock() bool {
	//即使使用setnx保证了锁获取的唯一性，
	//但不能保证只有一个客户端在释放锁。可能出现锁过期客户端还继续释放锁的情况，
	//又或者客户端运行阻塞没有正确释放，其他线程会尝试获得锁，释放时就出现了多个客户端释放锁的情况
	// 判断锁是否存在
	exists, err := m.client.Exists(context.Background(), m.name).Result()
	if err != nil {
		return false
	} //操作是否成功
	if exists == 0 { //判断键是否存在
		return false
	}

	// 获取锁的值并比较
	currentValue, err := m.client.Get(context.Background(), m.name).Result()
	if err != nil {
		return false
	}
	if currentValue != m.value {
		return false
	}

	// 删除锁
	_, err = m.client.Del(context.Background(), m.name).Result()
	if err != nil {
		return false
	}

	return true
} //使用del命令来释放锁虽然 DEL 命令本身是原子性的，但是在释放锁时，如果使用 DEL 命令来释放锁，就可能会出现以下情况：
//某个客户端刚刚检查到锁存在，但是在执行 DEL 命令之前，锁已经被另一个客户端释放了，这时候这个客户端就会错误地释放了另一个客户端所持有的锁。
//多个客户端同时尝试释放同一个锁，如果它们使用的是相同的值，那么它们有可能会释放其他客户端所持有的锁，导致锁的不确定性

//使用lua脚本释放锁  暂时不会
//func (m *Mutex) Unlock() bool {
//	script := `
//       if redis.call("get", KEYS[1]) == ARGV[1] then
//           return redis.call("del", KEYS[1])
//       else
//           return 0
//       end
//   `
//	res, err := m.client.Eval(context.Background(), script, []string{m.name}, m.value).Result()
//	return err == nil && res.(int64) == 1
//}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	m := NewMutex(client, "my_mutex", "my_value", 10*time.Second)

	if m.Lock() {
		fmt.Println("Got lock!")
		time.Sleep(1 * time.Second) //释放锁需要时间完成，如果不给个反应时间其他客户端会认为锁已经释放，有可能导致锁失效出现死锁现象
		m.Unlock()
		fmt.Println("Released lock!")
	} else {
		fmt.Println("Failed to get lock!")
	}
}
