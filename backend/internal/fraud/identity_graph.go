package fraud

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// IdentityRelationship represents a relationship between two user identities
type IdentityRelationship struct {
	ID                int64
	UserID1           int64
	UserID2           int64
	RelationshipType  string // 'shared_ip', 'shared_device', 'shared_payout', 'shared_fingerprint'
	RelationshipScore int    // 0-100
	Strength          string // 'weak', 'medium', 'strong'
	Evidence          map[string]interface{}
	LastDetected      time.Time
}

// FraudCluster represents a detected fraud cluster
type FraudCluster struct {
	ID          int64
	ClusterName string
	ClusterType string // 'fake_merchants', 'bot_creators', 'payout_fraud', 'coupon_abuse'
	MemberCount int
	RiskScore   int    // 0-100
	Status      string // 'detected', 'investigating', 'confirmed', 'resolved'
	CreatedAt   time.Time
	ResolvedAt  *time.Time
}

// FraudGraphStorage defines storage interface for identity graph
type FraudGraphStorage interface {
	SaveRelationship(ctx context.Context, rel *IdentityRelationship) error
	GetRelationships(ctx context.Context, userID int64) ([]*IdentityRelationship, error)
	GetRelationship(ctx context.Context, userID1, userID2 int64) (*IdentityRelationship, error)
	DeleteRelationship(ctx context.Context, userID1, userID2 int64) error

	SaveCluster(ctx context.Context, cluster *FraudCluster) error
	GetCluster(ctx context.Context, clusterID int64) (*FraudCluster, error)
	GetClustersByType(ctx context.Context, clusterType string) ([]*FraudCluster, error)
	GetActiveClusters(ctx context.Context) ([]*FraudCluster, error)
	UpdateCluster(ctx context.Context, cluster *FraudCluster) error

	GetClusterMembers(ctx context.Context, clusterID int64) ([]int64, error)
	AddClusterMember(ctx context.Context, clusterID, userID int64, confidence int) error
	RemoveClusterMember(ctx context.Context, clusterID, userID int64) error

	GetGraphMetrics(ctx context.Context, date time.Time) (map[string]interface{}, error)
}

// IdentityGraph manages the identity graph for fraud detection
type IdentityGraph struct {
	storage FraudGraphStorage
	config  *IdentityGraphConfig
}

// IdentityGraphConfig holds identity graph configuration
type IdentityGraphConfig struct {
	RelationshipScoreThreshold int
	ClusterMinSize             int
	SharedIPWeight             int
	SharedDeviceWeight         int
	SharedPayoutWeight         int
	SharedTaxIDWeight          int
	SharedFingerprintWeight    int
}

// NewIdentityGraph creates a new identity graph
func NewIdentityGraph(storage FraudGraphStorage, config *IdentityGraphConfig) *IdentityGraph {
	if config == nil {
		config = &IdentityGraphConfig{
			RelationshipScoreThreshold: 60,
			ClusterMinSize:             5,
			SharedIPWeight:             30,
			SharedDeviceWeight:         25,
			SharedPayoutWeight:         40,
			SharedTaxIDWeight:          45,
			SharedFingerprintWeight:    35,
		}
	}

	return &IdentityGraph{
		storage: storage,
		config:  config,
	}
}

// DetectSharedIP creates a relationship for shared IP
func (g *IdentityGraph) DetectSharedIP(ctx context.Context, userID1, userID2 int64, ipAddress string, timeframe time.Duration) error {
	if userID1 == 0 || userID2 == 0 {
		return errors.New("user IDs are required")
	}

	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}

	score := g.config.SharedIPWeight
	rel := &IdentityRelationship{
		UserID1:           userID1,
		UserID2:           userID2,
		RelationshipType:  "shared_ip",
		RelationshipScore: score,
		Strength:          g.getStrength(score),
		Evidence: map[string]interface{}{
			"ip_address": ipAddress,
			"timeframe":  timeframe.String(),
		},
		LastDetected: time.Now(),
	}

	return g.storage.SaveRelationship(ctx, rel)
}

// DetectSharedDevice creates a relationship for shared device fingerprint
func (g *IdentityGraph) DetectSharedDevice(ctx context.Context, userID1, userID2 int64, fingerprint string) error {
	if userID1 == 0 || userID2 == 0 {
		return errors.New("user IDs are required")
	}

	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}

	score := g.config.SharedDeviceWeight
	rel := &IdentityRelationship{
		UserID1:           userID1,
		UserID2:           userID2,
		RelationshipType:  "shared_device",
		RelationshipScore: score,
		Strength:          g.getStrength(score),
		Evidence: map[string]interface{}{
			"device_fingerprint": fingerprint,
		},
		LastDetected: time.Now(),
	}

	return g.storage.SaveRelationship(ctx, rel)
}

// DetectSharedPayout creates a relationship for shared payout account
func (g *IdentityGraph) DetectSharedPayout(ctx context.Context, userID1, userID2 int64, payoutID string) error {
	if userID1 == 0 || userID2 == 0 {
		return errors.New("user IDs are required")
	}

	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}

	score := g.config.SharedPayoutWeight
	rel := &IdentityRelationship{
		UserID1:           userID1,
		UserID2:           userID2,
		RelationshipType:  "shared_payout",
		RelationshipScore: score,
		Strength:          g.getStrength(score),
		Evidence: map[string]interface{}{
			"payout_id": payoutID,
		},
		LastDetected: time.Now(),
	}

	return g.storage.SaveRelationship(ctx, rel)
}

// GetRelationships gets all relationships for a user
func (g *IdentityGraph) GetRelationships(ctx context.Context, userID int64) ([]*IdentityRelationship, error) {
	return g.storage.GetRelationships(ctx, userID)
}

// DetectClusters detects fraud clusters in the graph
func (g *IdentityGraph) DetectClusters(ctx context.Context, clusterType string) ([]*FraudCluster, error) {
	// This is a simplified implementation - in production, use proper graph clustering algorithms
	// like Louvain community detection or label propagation

	relationships, err := g.storage.GetRelationships(ctx, 0) // Get all relationships
	if err != nil {
		return nil, fmt.Errorf("failed to get relationships: %w", err)
	}

	// Group users by strong relationships
	userGroups := make(map[int64][]int64)
	for _, rel := range relationships {
		if rel.RelationshipScore >= g.config.RelationshipScoreThreshold {
			userGroups[rel.UserID1] = append(userGroups[rel.UserID1], rel.UserID2)
			userGroups[rel.UserID2] = append(userGroups[rel.UserID2], rel.UserID1)
		}
	}

	// Create clusters from groups
	var clusters []*FraudCluster
	for userID, connectedUsers := range userGroups {
		if len(connectedUsers) >= g.config.ClusterMinSize {
			cluster := &FraudCluster{
				ClusterName: fmt.Sprintf("Cluster_%d", userID),
				ClusterType: clusterType,
				MemberCount: len(connectedUsers),
				Status:      "detected",
				CreatedAt:   time.Now(),
			}

			if err := g.storage.SaveCluster(ctx, cluster); err != nil {
				return nil, fmt.Errorf("failed to save cluster: %w", err)
			}

			// Add members
			for _, member := range connectedUsers {
				_ = g.storage.AddClusterMember(ctx, cluster.ID, member, 80) // Default 80% confidence
			}

			clusters = append(clusters, cluster)
		}
	}

	return clusters, nil
}

// GetCluster gets a cluster by ID
func (g *IdentityGraph) GetCluster(ctx context.Context, clusterID int64) (*FraudCluster, error) {
	return g.storage.GetCluster(ctx, clusterID)
}

// GetActiveClusters gets all active fraud clusters
func (g *IdentityGraph) GetActiveClusters(ctx context.Context) ([]*FraudCluster, error) {
	return g.storage.GetActiveClusters(ctx)
}

// MarkClusterResolved marks a cluster as resolved
func (g *IdentityGraph) MarkClusterResolved(ctx context.Context, clusterID int64) error {
	cluster, err := g.storage.GetCluster(ctx, clusterID)
	if err != nil {
		return fmt.Errorf("cluster not found: %w", err)
	}

	if cluster == nil {
		return errors.New("cluster not found")
	}

	now := time.Now()
	cluster.Status = "resolved"
	cluster.ResolvedAt = &now

	return g.storage.UpdateCluster(ctx, cluster)
}

// GetGraphMetrics gets metrics for the identity graph
func (g *IdentityGraph) GetGraphMetrics(ctx context.Context, date time.Time) (map[string]interface{}, error) {
	return g.storage.GetGraphMetrics(ctx, date)
}

// getStrength determines relationship strength from score
func (g *IdentityGraph) getStrength(score int) string {
	switch {
	case score >= 80:
		return "strong"
	case score >= 50:
		return "medium"
	default:
		return "weak"
	}
}
