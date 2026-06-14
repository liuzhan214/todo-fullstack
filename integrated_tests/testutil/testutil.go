package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

var (
	apiAddr   = envOrDefault("TODO_API", "http://localhost:20000")
	redisAddr = envOrDefault("REDIS_ADDR", "localhost:6380")
	rdb       = redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx       = context.Background()
)

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ---------- HTTP 请求封装 ----------

func Post[T any](t *testing.T, path string, body any) *T {
	t.Helper()
	return doRequest[T](t, http.MethodPost, path, body)
}

func Put[T any](t *testing.T, path string, body any) *T {
	t.Helper()
	return doRequest[T](t, http.MethodPut, path, body)
}

func Get[T any](t *testing.T, path string) *T {
	t.Helper()
	return doRequest[T](t, http.MethodGet, path, nil)
}

func Delete[T any](t *testing.T, path string, body any) *T {
	t.Helper()
	return doRequest[T](t, http.MethodDelete, path, body)
}

func doRequest[T any](t *testing.T, method, path string, body any) *T {
	t.Helper()

	url := apiAddr + path
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HTTP %d (期望 200), 请求: %s %s",
			resp.StatusCode, method, url)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return &result
}

// ---------- Redis 验证封装 ----------

// HGet 读取 Redis hash 指定 field 的值。
// 如果 key 或 field 不存在，t 会 Fail。
func HGet(t *testing.T, key, field string) string {
	t.Helper()
	val, err := rdb.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		t.Fatalf("Redis [HGET %s %s] → key 或 field 不存在", key, field)
	}
	if err != nil {
		t.Fatalf("Redis [HGET %s %s] 错误: %v", key, field, err)
	}
	return val
}

// HGetAll 读取 Redis hash 的所有 field-value 对
func HGetAll(t *testing.T, key string) map[string]string {
	t.Helper()
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		t.Fatalf("Redis [HGETALL %s] 错误: %v", key, err)
	}
	return result
}

// HDel 删除 Redis hash 中的指定 field（测试清理用）
func HDel(t *testing.T, key string, fields ...string) {
	t.Helper()
	if err := rdb.HDel(ctx, key, fields...).Err(); err != nil {
		t.Fatalf("Redis [HDEL %s %v] 错误: %v", key, fields, err)
	}
}

// IncrVal 读取 Redis 计数器当前值
func IncrVal(t *testing.T, key string) int64 {
	t.Helper()
	val, err := rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0
	}
	if err != nil {
		t.Fatalf("Redis [GET %s] 错误: %v", key, err)
	}
	return val
}

// DelKey 删除 Redis key（测试清理用）
func DelKey(t *testing.T, keys ...string) {
	t.Helper()
	if err := rdb.Del(ctx, keys...).Err(); err != nil {
		t.Fatalf("Redis [DEL %v] 错误: %v", keys, err)
	}
}

// ---------- 断言 ----------

func AssertEqual[T comparable](t *testing.T, want, got T) {
	t.Helper()
	if want != got {
		t.Fatalf("期望 %v, 实际 %v", want, got)
	}
}

func AssertNotEmpty(t *testing.T, label, value string) {
	t.Helper()
	if value == "" {
		t.Fatalf("%s 不应为空", label)
	}
}

// PrintJSON 打印结构化 JSON（辅助调试用）
func PrintJSON(t *testing.T, label string, v any) {
	t.Helper()
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("--- %s ---\n%s\n", label, string(b))
}
