package dbaas_postgress_test

import (
	"context"
	"errors"
	"testing"

	tfconstants "github.com/e2eterraformprovider/terraform-provider-e2e/e2e/constants"
	"github.com/e2eterraformprovider/terraform-provider-e2e/e2e/dbaas_postgress"
	"github.com/e2eterraformprovider/terraform-provider-e2e/goe2e"
	goe2econstants "github.com/e2eterraformprovider/terraform-provider-e2e/goe2e/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// createTestPostgreSQLClusterForDataSource creates a test cluster
func createTestPostgreSQLClusterForDataSource(id int, name, status string) *goe2e.PostgreSQLCluster {
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
// Tests for dataSourceReadPostgreSQL
// ============================================================================

func TestDataSourceReadPostgreSQL_Success(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	cluster := createTestPostgreSQLClusterForDataSource(123, "test-postgres", goe2econstants.DBaaSStatusRunning)
	cluster.MasterNode.Database.PGDetail.ID = 5
	cluster.MasterNode.Database.ID = 10
	cluster.MasterNode.Database.Database = "testdb"
	cluster.MasterNode.Database.Username = "testuser"
	cluster.MasterNode.Plan.Name = "E2E-2C-4GB"
	cluster.MasterNode.Plan.Software.Version = "14"
	cluster.MasterNode.Status = goe2econstants.DBaaSStatusRunning

	mockService.On("GetCluster", ctx, "123").Return(cluster, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, "123", d.Id())
	assert.Equal(t, 10, d.Get(tfconstants.AttrDatabaseID))
	assert.Equal(t, "testdb", d.Get(tfconstants.AttrDatabaseName))
	assert.Equal(t, "testuser", d.Get(tfconstants.AttrDatabaseUser))
	assert.Equal(t, goe2econstants.DBaaSStatusRunning, d.Get(tfconstants.AttrStatus))
	mockService.AssertExpectations(t)
}

func TestDataSourceReadPostgreSQL_NotFound(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	mockService.On("GetCluster", ctx, "123").Return(nil, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestDataSourceReadPostgreSQL_Error(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	mockService.On("GetCluster", ctx, "123").Return(nil, &goe2e.Response{}, errors.New("API error"))

	diags := resource.ReadContext(ctx, d, cfg)

	assert.True(t, diags.HasError())
	mockService.AssertExpectations(t)
}

func TestDataSourceReadPostgreSQL_NoParameterGroup(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	cluster := createTestPostgreSQLClusterForDataSource(123, "test-postgres", goe2econstants.DBaaSStatusRunning)
	cluster.MasterNode.Database.PGDetail.ID = 0 // No parameter group
	cluster.MasterNode.Database.ID = 10
	cluster.MasterNode.Database.Database = "testdb"
	cluster.MasterNode.Database.Username = "testuser"
	cluster.MasterNode.Plan.Name = "E2E-2C-4GB"
	cluster.MasterNode.Plan.Software.Version = "14"
	cluster.MasterNode.Status = goe2econstants.DBaaSStatusRunning

	mockService.On("GetCluster", ctx, "123").Return(cluster, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, tfconstants.DBaaSDefaultParameterGroupID, d.Get(tfconstants.AttrParameterGroupID))
	mockService.AssertExpectations(t)
}

func TestDataSourceReadPostgreSQL_NoPublicIP(t *testing.T) {
	ctx := context.Background()
	mockService := &dbaas_postgress.MockDBaaSPostgreSQLService{}
	cfg := mockConfigWithPostgreSQLService(mockService)

	resource := dbaas_postgress.DataSourcePostgresDBaaS()
	d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		tfconstants.AttrID: "123",
	})

	cluster := createTestPostgreSQLClusterForDataSource(123, "test-postgres", goe2econstants.DBaaSStatusRunning)
	cluster.MasterNode.PublicIPAddress = "" // No public IP
	cluster.MasterNode.Database.ID = 10
	cluster.MasterNode.Database.Database = "testdb"
	cluster.MasterNode.Database.Username = "testuser"
	cluster.MasterNode.Plan.Name = "E2E-2C-4GB"
	cluster.MasterNode.Plan.Software.Version = "14"
	cluster.MasterNode.Status = goe2econstants.DBaaSStatusRunning

	mockService.On("GetCluster", ctx, "123").Return(cluster, &goe2e.Response{}, nil)

	diags := resource.ReadContext(ctx, d, cfg)

	assert.False(t, diags.HasError())
	assert.Equal(t, false, d.Get(tfconstants.AttrIsPublicIPAttached))
	mockService.AssertExpectations(t)
}
