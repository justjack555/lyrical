package parser

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"log"
	"strings"
)

func getRedisClient() (*redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
}

func buildRedisKey(k string) string {
	var sb strings.Builder

	sb.WriteString("SongWords:")
	sb.WriteString(k)

	return sb.String()
}

func addToRedisZSet(redisDb *redis.Client, k string, score int64, v []byte) error {
	z := redis.Z {
		Score: float64(score),
		Member: v,
	}

	_, err := redisDb.ZAdd(k, z).Result()
	if err != nil {
		log.Println("addToRedisZSet(): ERR: ", err)
		return err
	}

	return nil
}

func storeSongWord(redisDb *redis.Client, k string, v *songWord) error {
	score := v.Song.GetId()

//	log.Println("storeSongWord(): DEBUG: Indices", v.Indices, "; Song: ", v.Song.GetName())

	blob, err := json.Marshal(*v)
	if err != nil {
		log.Println("storeSongWord(): ERR: Unable to marshal list of Song words: ", err)
		return err
	}

	err = addToRedisZSet(redisDb, k, score, blob)

//	log.Println("storeSongWord(): Key: ", k, "; Score: ", score, "; Raw JSON is: ", string(blob))

	return err

}

func storeWord(redisDb *redis.Client, k string, v []*songWord) int {
	redisKey := buildRedisKey(k)
	errCount := 0

	for _, song := range v {
		err := storeSongWord(redisDb, redisKey, song)
		if err != nil {
			log.Println("storeWord(): ERR: ", err)
			errCount++
		}
	}

	return errCount
}

func (m *Master) storeIndex(redisDb *redis.Client) int {
	errCount := 0

	for k, v := range m.invIndex {
		errCount += storeWord(redisDb, k, v)
	}

	return errCount
}

func (m *Master) StoreIndex() error {
	redisDb := getRedisClient()

	errCount := m.storeIndex(redisDb)
	if errCount > 0 {
		return errors.New("Greater than 0 errors occurred while storing Indices.")
	}

	log.Println("StoreIndex(): INFO: Zero errors while storing inverted index")

	return nil
}
