package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var REDIS *redis.Client

type ToDo struct {
	gorm.Model
	TaskId          string `json:"taskId"`
	TaskDescription string `json:"taskDescription"`
}

func getDatabaseUri() string {
	var dbUser = os.Getenv("POSTGRES_USER")
	var dbPassword = os.Getenv("POSTGRES_PASSWORD")
	var db = os.Getenv("POSTGRES_DB")
	var dbHost = os.Getenv("POSTGRES_HOST")
	var dbPort = os.Getenv("POSTGRES_PORT")
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		dbHost, dbPort, dbUser, db, dbPassword)
}

func connectPostgresDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(getDatabaseUri()), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Error: Cannot connect to postgres db")
	}
	if migrate_err := DB.AutoMigrate(&ToDo{}); migrate_err != nil {
		fmt.Println(migrate_err.Error())
		panic("Error: Unable to AutoMigrate postgres db")
	}
	fmt.Println("Postgres DB init automigration completed")
}

func getCacheUri() string {
	var cacheHost = os.Getenv("REDIS_HOST")
	var cachePort = os.Getenv("REDIS_PORT")
	return fmt.Sprintf("%s:%s", cacheHost, cachePort)
}

func connectRedisCache() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     getCacheUri(),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if _, redis_err := redisClient.Ping().Result(); redis_err != nil {
		fmt.Println(redis_err.Error())
		panic("Error: Unable to connect to Redis")
	}
	REDIS = redisClient
	fmt.Println("Redis cache init was completed")
}

func SetInCache(c *redis.Client, key string, value interface{}) bool {
	marshalledValue, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to set element in cache")
		return false
	}
	c.Set(key, marshalledValue, 0)
	return true
}

func GetFromCache(c *redis.Client, key string) interface{} {
	value, err := c.Get(key).Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return value
}

func DeleteFromCache(c *redis.Client, key string) {
	c.Del(key)
}

func InitialMigration() {
	connectPostgresDB()
	connectRedisCache()
}
