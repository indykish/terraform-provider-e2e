package dbaas_postgress_test

import (
	"context"
	"errors"
	"testing"

	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/config"
	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_postgress"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockConfig creates a mock config with a mock goe2e client
func mockConfigWithPostgreSQLService(service goe2e.PostgreSQLService) *config.Config {
	cfg := &config.Config{
		DefaultProjectID: "test-project",
		DefaultRegion:    "Mumbai",
	}
	client := &goe2e.Client{
		PostgreSQL: service,
	}

	// Use the public testing method to set the client
	cfg.SetGoe2eClientForTesting(client)

	return cfg
}

// Helper to create a test PostgreSQL cluster
func createTestPostgreSQLCluster(id int, name, status string) *goe2e.PostgreSQLCluster {
	return &goe2e.PostgreSQLCluster{
		ID:                  id,
		Name:                name,
		Status:              status,
		IsEncryptionEnabled: false,
		MasterNode: goe2e.DBNode{
			PublicIPAddress:  "1.2.3.4",
			PrivateIPAddress: "10.0.0.1",
			Port:             "5432",
			Disk:             "100GB",
			Status:           status,
			Database: goe2e.DBCreds{
				PGDetail: goe2e.PGDetail{
					ID: 0,
				},
			},
			Plan: goe2e.Plan{
				Software: goe2e.Software{},
			},
		},
	}
}

// ============================================================================
// Tests for resourceCreatePostgresDB
// ============================================================================

func TestResourceCreatePostgresDB_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVersion:   "14",
		tfconstants.AttrPlan:      "E2E-2C-4GB",
		tfconstants.AttrDBaaSName: "test-postgres",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":     "testuser",
				"password": "testpass",
				"name":     "testdb",
			},
		},
	})

	// Mock API calls
	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwarePostgreSQL, "14", "").Return(1, nil)
	mockService.On("GetTemplateID", ctx, "E2E-2C-4GB", "1", "").Return(100, nil)
	mockService.On("CreateCluster", ctx, mock.AnythingOfType("*goe2e.PostgreSQLClusterCreateRequest")).Return(
		createTestPostgreSQLCluster(123, "test-postgres", goe2econstants.DBaaSStatusRunning),
		&goe2e.Response{},
		nil,
	)

	diags := resource.CreateContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "123", d.Id())
	mockService.AssertExpectations(t)
}

func TestResourceCreatePostgresDB_GetSoftwareIDError(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVersion:   "14",
		tfconstants.AttrPlan:      "E2E-2C-4GB",
		tfconstants.AttrDBaaSName: "test-postgres",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":     "testuser",
				"password": "testpass",
				"name":     "testdb",
			},
		},
	})

	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwarePostgreSQL, "14", "").Return(0, errors.New("software not found"))

	diags := resource.CreateContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceCreatePostgresDB_GetTemplateIDError(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVersion:   "14",
		tfconstants.AttrPlan:      "E2E-2C-4GB",
		tfconstants.AttrDBaaSName: "test-postgres",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":     "testuser",
				"password": "testpass",
				"name":     "testdb",
			},
		},
	})

	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwarePostgreSQL, "14", "").Return(1, nil)
	mockService.On("GetTemplateID", ctx, "E2E-2C-4GB", "1", "").Return(0, errors.New("template not found"))

	diags := resource.CreateContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestResourceCreatePostgresDB_CreateClusterError(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrVersion:   "14",
		tfconstants.AttrPlan:      "E2E-2C-4GB",
		tfconstants.AttrDBaaSName: "test-postgres",
		tfconstants.AttrDatabase: []interface{}{
			map[string]interface{}{
				"user":     "testuser",
				"password": "testpass",
				"name":     "testdb",
			},
		},
	})

	mockService.On("GetSoftwareID", ctx, goe2econstants.DBaaSSoftwarePostgreSQL, "14", "").Return(1, nil)
	mockService.On("GetTemplateID", ctx, "E2E-2C-4GB", "1", "").Return(100, nil)
	mockService.On("CreateCluster", ctx, mock.AnythingOfType("*goe2e.PostgreSQLClusterCreateRequest")).Return(
		nil,
		&goe2e.Response{},
		errors.New("create failed"),
	)

	diags := resource.CreateContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

// ============================================================================
// Tests for resourceReadPostgresDB
// ============================================================================

func TestResourceReadPostgresDB_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	cluster := createTestPostgreSQLCluster(123, "test-postgres", goe2econstants.DBaaSStatusRunning)
	cluster.MasterNode.Database.PGDetail.ID = 5
	cluster.MasterNode.Database.ID = 10
	cluster.MasterNode.Database.Username = "testuser"
	cluster.MasterNode.Database.Database = "testdb"
	cluster.MasterNode.Plan.Name = "E2E-2C-4GB"
	cluster.MasterNode.Plan.Software.Version = "14"
	cluster.MasterNode.Status = goe2econstants.DBaaSStatusRunning

	mockService.On("GetCluster", ctx, "123").Return(cluster, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "123", d.Id())
	// Verify mock was called
	mockService.AssertExpectations(t)
}

func TestResourceReadPostgresDB_NotFound(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	mockService.On("GetCluster", ctx, "123").Return(nil, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "", d.Id()) // ID should be cleared
	mockService.AssertExpectations(t)
}

func TestResourceReadPostgresDB_Error(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	mockService.On("GetCluster", ctx, "123").Return(nil, &goe2e.Response{}, errors.New("API error"))

	diags := resource.ReadContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

// ============================================================================
// Tests for resourceDeletePostgresDB
// ============================================================================

func TestResourceDeletePostgresDB_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	mockService.On("DeleteCluster", ctx, "123").Return(&goe2e.Response{}, nil)

	diags := resource.DeleteContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "", d.Id()) // ID should be cleared
	mockService.AssertExpectations(t)
}

func TestResourceDeletePostgresDB_AlreadyDeleted(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	// Simulate delete fails because cluster doesn't exist
	mockService.On("DeleteCluster", ctx, "123").Return(nil, errors.New("not found"))
	// Then ClusterExists confirms it's gone
	mockService.On("ClusterExists", ctx, "123").Return(false, &goe2e.Response{}, nil)

	diags := resource.DeleteContext(ctx, d, cfg)

	// Should succeed (idempotent) - no error expected
	assert.False(t, diags.HasError())
	assert.Equal(t, "", d.Id()) // ID should be cleared
	mockService.AssertExpectations(t)
}

func TestResourceDeletePostgresDB_Error(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.ResourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	d.SetId("123")

	mockService.On("DeleteCluster", ctx, "123").Return(nil, errors.New("delete failed"))
	// ClusterExists returns true (cluster still exists) so error should be returned
	mockService.On("ClusterExists", ctx, "123").Return(true, &goe2e.Response{}, nil)

	diags := resource.DeleteContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}
