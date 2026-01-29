package feature

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// FeatureFlag represents a feature flag
type FeatureFlag struct {
	Key         string    `json:"key"`
	Enabled     bool      `json:"enabled"`
	Description string    `json:"description"`
	Rollout     int       `json:"rollout"` // Percentage 0-100
	Rules       []Rule    `json:"rules"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Rule defines targeting rules for feature flags
type Rule struct {
	Attribute string   `json:"attribute"` // user_id, email, country, etc.
	Operator  string   `json:"operator"`  // equals, contains, in
	Values    []string `json:"values"`
}

// FeatureFlagManager manages feature flags
type FeatureFlagManager struct {
	flags       map[string]*FeatureFlag
	redisClient *redis.Client
	mu          sync.RWMutex
}

// NewFeatureFlagManager creates a new feature flag manager
func NewFeatureFlagManager(redisClient *redis.Client) *FeatureFlagManager {
	manager := &FeatureFlagManager{
		flags:       make(map[string]*FeatureFlag),
		redisClient: redisClient,
	}

	// Load flags from Redis
	if redisClient != nil {
		manager.loadFromRedis()
	}

	return manager
}

// RegisterFlag registers a new feature flag
func (m *FeatureFlagManager) RegisterFlag(flag *FeatureFlag) {
	m.mu.Lock()
	defer m.mu.Unlock()

	flag.UpdatedAt = time.Now()
	if flag.CreatedAt.IsZero() {
		flag.CreatedAt = time.Now()
	}

	m.flags[flag.Key] = flag

	// Persist to Redis
	if m.redisClient != nil {
		m.saveToRedis(flag)
	}
}

// IsEnabled checks if a feature is enabled for a given context
func (m *FeatureFlagManager) IsEnabled(key string, ctx map[string]string) bool {
	m.mu.RLock()
	flag, exists := m.flags[key]
	m.mu.RUnlock()

	if !exists {
		return false
	}

	if !flag.Enabled {
		return false
	}

	// Check rules
	if len(flag.Rules) > 0 {
		if !m.matchesRules(flag.Rules, ctx) {
			return false
		}
	}

	// Check rollout percentage
	if flag.Rollout < 100 {
		// Use consistent hashing for gradual rollout
		userID := ctx["user_id"]
		if userID == "" {
			userID = ctx["session_id"]
		}
		if userID == "" {
			return false
		}

		hash := hashString(flag.Key + userID)
		return hash%100 < uint32(flag.Rollout)
	}

	return true
}

// matchesRules checks if context matches all rules
func (m *FeatureFlagManager) matchesRules(rules []Rule, ctx map[string]string) bool {
	for _, rule := range rules {
		value, exists := ctx[rule.Attribute]
		if !exists {
			return false
		}

		switch rule.Operator {
		case "equals":
			if len(rule.Values) > 0 && value != rule.Values[0] {
				return false
			}
		case "contains":
			found := false
			for _, v := range rule.Values {
				if value == v {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case "in":
			found := false
			for _, v := range rule.Values {
				if value == v {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

// UpdateFlag updates an existing feature flag
func (m *FeatureFlagManager) UpdateFlag(key string, enabled bool, rollout int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	flag, exists := m.flags[key]
	if !exists {
		return fmt.Errorf("feature flag not found: %s", key)
	}

	flag.Enabled = enabled
	flag.Rollout = rollout
	flag.UpdatedAt = time.Now()

	// Persist to Redis
	if m.redisClient != nil {
		m.saveToRedis(flag)
	}

	return nil
}

// GetFlag retrieves a feature flag
func (m *FeatureFlagManager) GetFlag(key string) (*FeatureFlag, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	flag, exists := m.flags[key]
	if !exists {
		return nil, fmt.Errorf("feature flag not found: %s", key)
	}

	return flag, nil
}

// ListFlags returns all feature flags
func (m *FeatureFlagManager) ListFlags() []*FeatureFlag {
	m.mu.RLock()
	defer m.mu.RUnlock()

	flags := make([]*FeatureFlag, 0, len(m.flags))
	for _, flag := range m.flags {
		flags = append(flags, flag)
	}
	return flags
}

// saveToRedis saves a flag to Redis
func (m *FeatureFlagManager) saveToRedis(flag *FeatureFlag) error {
	ctx := context.Background()
	data, err := json.Marshal(flag)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("feature_flag:%s", flag.Key)
	return m.redisClient.Set(ctx, key, data, 0).Err()
}

// loadFromRedis loads all flags from Redis
func (m *FeatureFlagManager) loadFromRedis() error {
	ctx := context.Background()
	keys, err := m.redisClient.Keys(ctx, "feature_flag:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		data, err := m.redisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var flag FeatureFlag
		if err := json.Unmarshal([]byte(data), &flag); err != nil {
			continue
		}

		m.flags[flag.Key] = &flag
	}

	return nil
}

// hashString creates a hash for consistent rollout
func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// ABTest represents an A/B test configuration
type ABTest struct {
	Name       string            `json:"name"`
	Variations []ABTestVariation `json:"variations"`
	Traffic    map[string]int    `json:"traffic"` // variation -> percentage
}

// ABTestVariation represents a test variation
type ABTestVariation struct {
	Name  string                 `json:"name"`
	Value map[string]interface{} `json:"value"`
}

// ABTestManager manages A/B tests
type ABTestManager struct {
	tests map[string]*ABTest
	mu    sync.RWMutex
}

// NewABTestManager creates a new A/B test manager
func NewABTestManager() *ABTestManager {
	return &ABTestManager{
		tests: make(map[string]*ABTest),
	}
}

// RegisterTest registers a new A/B test
func (m *ABTestManager) RegisterTest(test *ABTest) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tests[test.Name] = test
}

// GetVariation returns the variation for a given user
func (m *ABTestManager) GetVariation(testName, userID string) *ABTestVariation {
	m.mu.RLock()
	test, exists := m.tests[testName]
	m.mu.RUnlock()

	if !exists {
		return nil
	}

	// Use consistent hashing
	hash := hashString(testName + userID)
	percentage := int(hash % 100)

	cumulative := 0
	for _, variation := range test.Variations {
		traffic := test.Traffic[variation.Name]
		cumulative += traffic
		if percentage < cumulative {
			return &variation
		}
	}

	// Default to first variation
	if len(test.Variations) > 0 {
		return &test.Variations[0]
	}
	return nil
}
