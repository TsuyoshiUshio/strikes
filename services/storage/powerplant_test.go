package storage

import (
	"reflect"
	"testing"
	"time"

	st "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {

	ExpectedInstanceID := "SomeGUID"
	ExpectedPackageID := "foo"
	ExpectedResourceGroup := "bar"
	ExpectedName := "baz"
	ExpectedPackageName := "qux"
	ExpectedPackageVersion := "1.0.0"
	ExpectedPackageParameters := "--set foo=bar --set bar=baz"
	ExpectedTableName := "quux"

	instance := &StrikesInstance{
		InstanceID:        ExpectedInstanceID,
		PackageID:         ExpectedPackageID,
		Name:              ExpectedName,
		ResourceGroup:     ExpectedResourceGroup,
		PackageName:       ExpectedPackageName,
		PackageVersion:    ExpectedPackageVersion,
		PackageParameters: ExpectedPackageParameters,
	}
	entity := instance.ConvertEntity(&st.Table{Name: ExpectedTableName})
	assert.Equal(t, ExpectedTableName, entity.Table.Name, "Table instance hasn't set correctly")
	assert.Equal(t, ExpectedInstanceID, entity.PartitionKey, "Wrong PartitionKey")
	assert.Equal(t, ExpectedName, entity.RowKey, "Wrong RowKey")
	assert.Equal(t, ExpectedResourceGroup, entity.Properties["ResourceGroup"], "Wrong ResourceGroup")
	assert.Equal(t, ExpectedPackageID, entity.Properties["PackageID"], "Wrong PackageID")
	assert.Equal(t, ExpectedPackageName, entity.Properties["PackageName"], "Wrong PackageName")
	assert.Equal(t, ExpectedPackageVersion, entity.Properties["PackageVersion"], "Wrong Package Version")
	assert.Equal(t, ExpectedPackageParameters, entity.Properties["PackageParameters"], "Wrong Package parameters")
}

func TestList(t *testing.T) {
	ExpectedTimeStamp := time.Now()

	ExpectedInstanceId01 := "foo01"
	ExpectedName01 := "bar01"
	ExpectedResourceGroup01 := "baz01"
	ExpectedPackageID01 := "qux01"
	ExpectedPackageName01 := "quux01"
	ExpectedVersion01 := "1.0.0"
	ExpectedPackageParameters01 := "{\"foo\":\"bar\"}"

	ExpectedInstanceId02 := "foo02"
	ExpectedName02 := "bar02"
	ExpectedResourceGroup02 := "baz02"
	ExpectedPackageID02 := "qux02"
	ExpectedPackageName02 := "quux02"
	ExpectedVersion02 := "1.0.1"
	ExpectedPackageParameters02 := "{\"bar\":\"baz\"}"

	m1 := make(map[string]interface{})
	m1["ResourceGroup"] = ExpectedResourceGroup01
	m1["PackageID"] = ExpectedPackageID01
	m1["PackageName"] = ExpectedPackageName01
	m1["PackageVersion"] = ExpectedVersion01
	m1["PackageParameters"] = ExpectedPackageParameters01

	m2 := make(map[string]interface{})
	m2["ResourceGroup"] = ExpectedResourceGroup02
	m2["PackageID"] = ExpectedPackageID02
	m2["PackageName"] = ExpectedPackageName02
	m2["PackageVersion"] = ExpectedVersion02
	m2["PackageParameters"] = ExpectedPackageParameters02

	entities := []*st.Entity{
		&st.Entity{
			PartitionKey: ExpectedInstanceId01,
			RowKey:       ExpectedName01,
			TimeStamp:    ExpectedTimeStamp,
			Properties:   m1,
		},
		&st.Entity{
			PartitionKey: ExpectedInstanceId02,
			RowKey:       ExpectedName02,
			TimeStamp:    ExpectedTimeStamp,
			Properties:   m2,
		},
	}

	fakeGetTableReference := func() (*st.Table, error) {
		return &st.Table{}, nil
	}
	fakeQueryEntities := func(t *st.Table, timeout uint, ml st.MetadataLevel, options *st.QueryOptions) (*st.EntityQueryResult, error) {
		result := st.EntityQueryResult{
			Entities: entities,
		}
		return &result, nil
	}
	monkey.Patch(getTableReference, fakeGetTableReference)
	var table *st.Table
	monkey.PatchInstanceMethod(reflect.TypeOf(table), "QueryEntities", fakeQueryEntities)
	defer monkey.UnpatchAll()
	results, err := List()
	assert.Nil(t, err)
	assert.Equal(t, ExpectedInstanceId01, (*results)[0].InstanceID)
	assert.Equal(t, ExpectedName01, (*results)[0].Name)
	assert.Equal(t, ExpectedResourceGroup01, (*results)[0].ResourceGroup)
	assert.Equal(t, ExpectedPackageID01, (*results)[0].PackageID)
	assert.Equal(t, ExpectedPackageName01, (*results)[0].PackageName)
	assert.Equal(t, ExpectedPackageParameters01, (*results)[0].PackageParameters)
	assert.Equal(t, ExpectedTimeStamp, (*results)[0].TimeStamp)

	assert.Equal(t, ExpectedInstanceId02, (*results)[1].InstanceID)
	assert.Equal(t, ExpectedName02, (*results)[1].Name)
	assert.Equal(t, ExpectedResourceGroup02, (*results)[1].ResourceGroup)
	assert.Equal(t, ExpectedPackageID02, (*results)[1].PackageID)
	assert.Equal(t, ExpectedPackageName02, (*results)[1].PackageName)
	assert.Equal(t, ExpectedPackageParameters02, (*results)[1].PackageParameters)
	assert.Equal(t, ExpectedTimeStamp, (*results)[1].TimeStamp)
}
