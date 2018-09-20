package storage

import (
	"time"

	st "github.com/Azure/azure-sdk-for-go/storage"
)

type StrikesInstance struct {
	PackageID         string // PartitionKey
	Name              string // RowKey
	ResourceGroup     string
	PackageName       string
	PackageVersion    string
	PackageParameters string
}

func (s *StrikesInstance) ConvertEntity(table *st.Table) *st.Entity {
	m := make(map[string]interface{})
	m["ResourceGroup"] = s.ResourceGroup
	m["PackageName"] = s.PackageName
	m["PackageVersion"] = s.PackageVersion
	m["PackageParameters"] = s.PackageParameters

	entity := &st.Entity{
		PartitionKey: s.PackageID,
		RowKey:       s.Name,
		TimeStamp:    time.Now(),
		Properties:   m,
		Table:        table,
	}
	return entity
}
