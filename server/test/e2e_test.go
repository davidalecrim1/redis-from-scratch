package e2e_test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
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
			Addr:     s.ListenAddr,
			Password: "", // no password
			DB:       0,  // default DB
		})
	defer redisClient.Close()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func startRedisServer() *server.Server {
	randomPort := rand.Intn(1000) + 6000 // random from 6000 to 6999

	ctx := context.Background()
	s := server.NewServer(
		server.Config{
			ListenAddr: fmt.Sprintf(":%d", randomPort),
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

	t.Run("should be able to SET and GET a command of a valid key", func(t *testing.T) {
		t.Parallel()

		// Test SET operation
		err := redisClient.Set(ctx, "test_key", "test_value", 0).Err()
		assert.NoError(t, err)

		// Test GET operation for existing key
		val, err := redisClient.Get(ctx, "test_key").Result()
		assert.NoError(t, err)
		assert.Equal(t, "test_value", val)
	})

	t.Run("should return nothing on GET command of non-existent key", func(t *testing.T) {
		t.Parallel()

		val, err := redisClient.Get(ctx, "non_existent_key").Result()
		assert.ErrorIs(t, err, redis.Nil)
		assert.Empty(t, val)
	})

	t.Run("Should return a REPL enconded message from ECHO command", func(t *testing.T) {
		t.Parallel()

		message := "hey"
		resp, err := redisClient.Echo(ctx, message).Result()
		assert.Nil(t, err)
		assert.Equal(t, resp, message)
	})

	t.Cleanup(func() {
		assert.NoError(t, redisClient.Close())
	})
}
