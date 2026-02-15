package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"pledge-backend/config"
	"pledge-backend/log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// InitRedis 初始化Redis
func InitRedis() *redis.Pool {
	log.Logger.Info("Init Redis")
	redisConf := config.Config.Redis
	// 建立连接池  RedisConn 全局变量，用于存储初始化后的连接池指针
	RedisConn = &redis.Pool{
		MaxIdle:     10,                // 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态。
		MaxActive:   0,                 // 最大的激活连接数，表示同时最多有N个连接   0 表示无穷大
		Wait:        true,              // 如果连接数不足则阻塞等待
		IdleTimeout: 180 * time.Second, // 连接空闲超过此时间将被自动回收
		// 定义建立物理连接的方法
		Dial: func() (redis.Conn, error) {
			// 连接服务器
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", redisConf.Address, redisConf.Port))
			if err != nil {
				return nil, err
			}
			// 验证密码
			_, err = c.Do("auth", redisConf.Password)
			if err != nil {
				panic("redis auth err " + err.Error())
			}
			// 选择db （Redis默认有16个db，通常使用0）
			_, err = c.Do("select", redisConf.Db)
			if err != nil {
				panic("redis select db err " + err.Error())
			}
			return c, nil
		},
	}
	err := RedisConn.Get().Err() // 获取一个连接测试一下，确保配置没写错
	if err != nil {
		panic("redis init err " + err.Error())
	}
	return RedisConn
}

/* ================== String（字符串）类型操作 ==================
   String 是 Redis 最基础的类型，常用于缓存对象、Token、验证码等。
   =========================================================== */
// RedisSet 设置 key、value、并支持设置过期秒数
// 学习点：interface{} 可以接收任何类型，但在存入Redis前必须序列化（如 JSON）
func RedisSet(key string, data interface{}, aliveSeconds int) error {
	conn := RedisConn.Get() // 从池中取一个连接
	defer func() {
		_ = conn.Close() // 函数结束时必须归还连接到池中
	}()

	// 将 Go 的结构体或 Map 转为 JSON 字符串
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// EX 参数代表过期时间（单位：秒）
	if aliveSeconds > 0 {
		_, err = conn.Do("set", key, value, "EX", aliveSeconds)
	} else {
		_, err = conn.Do("set", key, value)
	}
	if err != nil {
		return err
	}
	return nil
}

// RedisSetString  设置key、value、time
func RedisSetString(key string, data string, aliveSeconds int) error {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	var err error
	if aliveSeconds > 0 {
		_, err = redis.String(conn.Do("set", key, data, "EX", aliveSeconds))
	} else {
		_, err = redis.String(conn.Do("set", key, data))
	}
	if err != nil {
		return err
	}
	return nil
}

// RedisGet 获取Key 对应的原始字节数据
func RedisGet(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	// redis.Bytes 会将 Redis 返回的内容转为 Go 的 []byte
	reply, err := redis.Bytes(conn.Do("get", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// RedisGetString 获取Key对应的字符串
func RedisGetString(key string) (string, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	// redis.String 会将结果直接转为 string 类型
	reply, err := redis.String(conn.Do("get", key))
	if err != nil {
		return "", err
	}
	return reply, nil
}

// RedisSetInt64  专门存储 64位整数
func RedisSetInt64(key string, data int64, time int) error {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = redis.Int64(conn.Do("set", key, value))
	if err != nil {
		return err
	}
	if time != 0 {
		// expire 命令单独给 key 设置过期时间
		_, err = redis.Int64(conn.Do("expire", key, time))
		if err != nil {
			return err
		}
	}
	return nil
}

// RedisGetInt64 获取整数值
func RedisGetInt64(key string) (int64, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := redis.Int64(conn.Do("get", key))
	if err != nil {
		return -1, err
	}
	return reply, nil
}

// RedisDelete 删除Key
func RedisDelete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	return redis.Bool(conn.Do("del", key))
}

// RedisFlushDB 清空当前DB
func RedisFlushDB() error {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("flushdb")
	if err != nil {
		return err
	}
	return nil
}

/* ================== Hash（哈希）类型操作 ==================
   Hash 像 Go 里面的 Map，适合存储对象信息（比如用户信息：姓名、年龄、余额）。
   优点：可以只更新或读取对象中的某一个字段，不用频繁序列化整个大 JSON。
   ========================================================= */

// RedisGetHashOne 获取Heah其中一个值
func RedisGetHashOne(key, name string) (interface{}, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := conn.Do("hgetall", key, name)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// RedisSetHash 批量设置 Hash 字段
func RedisSetHash(key string, data map[string]string, time interface{}) error {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	// 循环发送命令到缓冲区，不立刻产生网络请求
	for k, v := range data {
		err := conn.Send("hset", key, k, v)
		if err != nil {
			return err
		}
	}
	// 将缓冲区命令一次性发送给 Redis 提升效率
	err := conn.Flush()
	if err != nil {
		return err
	}
	// 设置过期时间
	if time != nil {
		_, err = conn.Do("expire", key, time.(int))
		if err != nil {
			return err
		}
	}
	return nil
}

// RedisGetHash 获取 Hash 类型的所有键值对
func RedisGetHash(key string) (map[string]string, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	// redis.StringMap 是一个非常方便的辅助函数，直接把结果转为 map[string]string
	reply, err := redis.StringMap(conn.Do("hgetall", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// RedisDelHash 删除Hash
func RedisDelHash(key string) (bool, error) {

	return true, nil
}

// RedisExistsHash 检查Key是否存在
func RedisExistsHash(key string) bool {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	exists, err := redis.Bool(conn.Do("hexists", key))
	if err != nil {
		return false
	}
	return exists
}

// RedisExists 检查Key是否存在
func RedisExists(key string) bool {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	exists, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return false
	}
	return exists
}

// RedisGetTTL 获取Key剩余时间
func RedisGetTTL(key string) int64 {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := redis.Int64(conn.Do("ttl", key))
	if err != nil {
		return 0
	}
	return reply
}

/* ================== Set（集合）类型操作 ==================
   Set 是无序的、且不重复的字符串集合。常用于抽奖、共同好友、唯一名单等。
   ======================================================== */
// RedisSAdd 向集合中添加一个元素
func RedisSAdd(k, v string) int64 {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := conn.Do("SAdd", k, v)
	if err != nil {
		return -1
	}
	return reply.(int64)
}

// RedisSmembers 获取集合中的所有成员
func RedisSmembers(k string) ([]string, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	reply, err := redis.Strings(conn.Do("smembers", k))
	if err != nil {
		return []string{}, errors.New("读取set错误")
	}
	return reply, err
}

type RedisEncryptionTask struct {
	RecordOrderFlowId int32  `json:"recordOrderFlow"` //密码转账表ID
	Encryption        string `json:"encryption"`      //密码串
	EndTime           int64  `json:"endTime"`         //失效截止时间
}

/* ================== List（列表）类型操作 ==================
   List 是有序集合，可以作为队列（先进先出）或栈（先进后出）。
   ======================================================== */
// RedisListRpush 在列表右侧（队尾）插入数据
func RedisListRpush(listName string, encryption string) error {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("rpush", listName, encryption)
	return err
}

// RedisListLRange 取获取列表中的一段范围（0, -1 表示获取全部）
func RedisListLRange(listName string) ([]string, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	res, err := redis.Strings(conn.Do("lrange", listName, 0, -1))
	return res, err
}

// RedisListLRem 删除列表中指定元素
// 参数 1 表示只删除从左到右遇到的第 1 个匹配项
func RedisListLRem(listName string, encryption string) error {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("lrem", listName, 1, encryption)
	return err
}

// RedisListLength 列表长度
func RedisListLength(listName string) (interface{}, error) {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	len, err := conn.Do("llen", listName)
	return len, err
}

// RedisDelList list 删除整个列表
func RedisDelList(setName string) error {
	conn := RedisConn.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("del", setName)
	return err
}

// 数据选择：
// 简单的配置、小的 JSON 对象用 String。
// 需要频繁更新属性的大对象用 Hash。
// 任务队列、聊天消息记录用 List。
// 去重名单、点赞名单用 Set。
