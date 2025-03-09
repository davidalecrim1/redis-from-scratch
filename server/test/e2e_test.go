package e2e_test

import (
	"context"
	"log"
	"os"
	"testing"

	"redis-from-scratch/server"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var redisClient *redis.Client

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	s := startRedisServer()
	defer s.Close()

	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     "localhost:8103",
			Password: "", // no password
			DB:       0,  // default DB
		})
	defer redisClient.Close()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func startRedisServer() *server.Server {
	ctx := context.Background()
	s := server.NewServer(
		server.Config{
			ListenAddr: ":" + "8103",
		})

	go func() {
		err := s.Start(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return s
}

func TestRedisOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("Set and Get existing key", func(t *testing.T) {
		t.Parallel()

		// Test SET operation
		err := redisClient.Set(ctx, "test_key", "test_value", 0).Err()
		assert.NoError(t, err)

		// Test GET operation for existing key
		val, err := redisClient.Get(ctx, "test_key").Result()
		assert.NoError(t, err)
		assert.Equal(t, "test_value", val)
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		t.Parallel()

		// Test GET operation for non-existent key
		val, err := redisClient.Get(ctx, "non_existent_key").Result()
		assert.ErrorIs(t, err, redis.Nil)
		assert.Empty(t, val)
	})
}
