package storage

import (
	"time"

	st "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/TsuyoshiUshio/strikes/config"
)

type StrikesInstance struct {
	InstanceID        string // PartitionKey
	Name              string // RowKey
	ResourceGroup     string
	PackageID         string
	PackageName       string
	PackageVersion    string
	PackageParameters string
}

func (s *StrikesInstance) ConvertEntity(table *st.Table) *st.Entity {
	m := make(map[string]interface{})
	m["ResourceGroup"] = s.ResourceGroup
	m["PackageID"] = s.PackageID
	m["PackageName"] = s.PackageName
	m["PackageVersion"] = s.PackageVersion
	m["PackageParameters"] = s.PackageParameters

	entity := &st.Entity{
		PartitionKey: s.InstanceID,
		RowKey:       s.Name,
		TimeStamp:    time.Now(),
		Properties:   m,
		Table:        table,
	}
	return entity
}

func InsertOrUpdate(instance *StrikesInstance) error {
	table, err := getTableReference()
	if err != nil {
		return err
	}
	entity := instance.ConvertEntity(table)
	err = entity.InsertOrReplace(&st.EntityOptions{})
	if err != nil {
		return err
	}
	return nil
}
func getTableReference() (*st.Table, error) {
	// read the powerplant config
	powerPlantConfig, err := getPowerPlantConfig()
	if err != nil {
		return nil, nil
	}
	client, err := st.NewBasicClient(powerPlantConfig.StorageAccountName, powerPlantConfig.StorageAccountKey)
	if err != nil {
		return nil, err
	}
	tableService := client.GetTableService()
	table := tableService.GetTableReference(DEFAULT_STORAGE_TABLE_NAME)
	return table, nil
}

func getPowerPlantConfig() (*config.PowerPlantConfig, error) {
	configContext, err := config.NewConfigContext()
	if err != nil {
		return nil, err
	}
	return configContext.GetPowerPlantConfig()
}
