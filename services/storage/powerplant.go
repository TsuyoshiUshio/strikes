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
	TimeStamp         time.Time
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

func convertFromEntity(entity *st.Entity) *StrikesInstance {
	instance := StrikesInstance{
		InstanceID:        entity.PartitionKey,
		Name:              entity.RowKey,
		ResourceGroup:     entity.Properties["ResourceGroup"].(string),
		PackageID:         entity.Properties["PackageID"].(string),
		PackageName:       entity.Properties["PackageName"].(string),
		PackageVersion:    entity.Properties["PackageVersion"].(string),
		PackageParameters: entity.Properties["PackageParameters"].(string),
		TimeStamp:         entity.TimeStamp,
	}
	return &instance
}

func List() (*[]StrikesInstance, error) {
	table, err := getTableReference()
	if err != nil {
		return nil, err
	}
	options := st.QueryOptions{}
	results, err := table.QueryEntities(30, st.NoMetadata, &options)
	instances := make([]StrikesInstance, 0)
	for _, entity := range results.Entities {
		instances = append(instances, *(convertFromEntity(entity)))
	}
	return &instances, nil
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
