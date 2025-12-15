package dbaas_postgress

import (
	"context"

	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	"github.com/stretchr/testify/mock"
)

// MockDBaaSPostgreSQLService is a unified mock implementation for DBaaSPostgreSQL service
type MockDBaaSPostgreSQLService struct {
	mock.Mock
}

// ExpandPostgresVPCList expands a list of VPC IDs to their metadata
func (m *MockDBaaSPostgreSQLService) ExpandPostgresVPCList(ctx context.Context, vpcIDs []string) ([]goe2e.VPCMetadata, error) {
	args := m.Called(ctx, vpcIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]goe2e.VPCMetadata), args.Error(1)
}

// ClusterExists checks if a cluster exists
func (m *MockDBaaSPostgreSQLService) ClusterExists(ctx context.Context, id string) (bool, *goe2e.Response, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Get(1).(*goe2e.Response), args.Error(2)
}

// CreateCluster creates a new PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) CreateCluster(ctx context.Context, req *goe2e.PostgreSQLClusterCreateRequest) (*goe2e.PostgreSQLCluster, *goe2e.Response, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*goe2e.Response), args.Error(2)
	}
	return args.Get(0).(*goe2e.PostgreSQLCluster), args.Get(1).(*goe2e.Response), args.Error(2)
}

// GetCluster retrieves a PostgreSQL cluster by ID
func (m *MockDBaaSPostgreSQLService) GetCluster(ctx context.Context, id string) (*goe2e.PostgreSQLCluster, *goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*goe2e.Response), args.Error(2)
	}
	return args.Get(0).(*goe2e.PostgreSQLCluster), args.Get(1).(*goe2e.Response), args.Error(2)
}

// DeleteCluster deletes a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) DeleteCluster(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// StartCluster starts a stopped PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) StartCluster(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// StopCluster stops a running PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) StopCluster(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// RestartCluster restarts a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) RestartCluster(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// AttachVPC attaches a VPC to a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) AttachVPC(ctx context.Context, id string, req *goe2e.PostgreSQLVPCAttachRequest) (*goe2e.Response, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// DetachVPC detaches a VPC from a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) DetachVPC(ctx context.Context, id string, req *goe2e.PostgreSQLVPCAttachRequest) (*goe2e.Response, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// AttachParameterGroup attaches a parameter group to a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) AttachParameterGroup(ctx context.Context, id, pgID string) (*goe2e.Response, error) {
	args := m.Called(ctx, id, pgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// DetachParameterGroup detaches a parameter group from a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) DetachParameterGroup(ctx context.Context, id, pgID string) (*goe2e.Response, error) {
	args := m.Called(ctx, id, pgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// AttachPublicIP attaches a public IP to a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) AttachPublicIP(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// DetachPublicIP detaches a public IP from a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) DetachPublicIP(ctx context.Context, id string) (*goe2e.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// UpgradePlan upgrades a PostgreSQL cluster's plan
func (m *MockDBaaSPostgreSQLService) UpgradePlan(ctx context.Context, id string, req *goe2e.PostgreSQLPlanUpgradeRequest) (*goe2e.Response, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// ExpandDisk expands the disk size of a PostgreSQL cluster
func (m *MockDBaaSPostgreSQLService) ExpandDisk(ctx context.Context, id string, req *goe2e.DiskExpansionRequest) (*goe2e.Response, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*goe2e.Response), args.Error(1)
}

// GetSoftwareID retrieves the software ID for a given name, dbType, and version
func (m *MockDBaaSPostgreSQLService) GetSoftwareID(ctx context.Context, name, dbType, version string) (int, error) {
	args := m.Called(ctx, name, dbType, version)
	return args.Int(0), args.Error(1)
}

// GetTemplateID retrieves the template ID for a given plan, dbType, and version
func (m *MockDBaaSPostgreSQLService) GetTemplateID(ctx context.Context, planName, dbType, version string) (int, error) {
	args := m.Called(ctx, planName, dbType, version)
	return args.Int(0), args.Error(1)
}
