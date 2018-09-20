package storage

import (
	"testing"

	st "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {

	ExpectedPackageID := "foo"
	ExpectedResourceGroup := "bar"
	ExpectedName := "baz"
	ExpectedPackageName := "qux"
	ExpectedPackageVersion := "1.0.0"
	ExpectedPackageParameters := "--set foo=bar --set bar=baz"
	ExpectedTableName := "quux"

	instance := &StrikesInstance{
		PackageID:         ExpectedPackageID,
		Name:              ExpectedName,
		ResourceGroup:     ExpectedResourceGroup,
		PackageName:       ExpectedPackageName,
		PackageVersion:    ExpectedPackageVersion,
		PackageParameters: ExpectedPackageParameters,
	}
	entity := instance.ConvertEntity(&st.Table{Name: ExpectedTableName})
	assert.Equal(t, ExpectedTableName, entity.Table.Name, "Table instance hasn't set correctly")
	assert.Equal(t, ExpectedPackageID, entity.PartitionKey, "Wrong PartitionKey")
	assert.Equal(t, ExpectedName, entity.RowKey, "Wrong RowKey")
	assert.Equal(t, ExpectedResourceGroup, entity.Properties["ResourceGroup"], "Wrong ResourceGroup")
	assert.Equal(t, ExpectedPackageName, entity.Properties["PackageName"], "Wrong PackageName")
	assert.Equal(t, ExpectedPackageVersion, entity.Properties["PackageVersion"], "Wrong Package Version")
	assert.Equal(t, ExpectedPackageParameters, entity.Properties["PackageParameters"], "Wrong Package parameters")
}
