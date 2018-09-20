package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Azure/azure-sdk-for-go/storage"
	s "github.com/TsuyoshiUshio/strikes/services/storage"
)

type PowerPlantConfig struct {
	ResourceGroup      string `json:"resourceGroup"`
	StorageAccountName string `json:"storageAccountName"`
	StorageAccountKey  string `json:"storageAccountKey"`
}

type StrikesInstance struct {
	*storage.Entity
	Foo string
	Bar string
	Baz time.Time
}

func (s *StrikesInstance) Setup(table *storage.Table) {
	s.Entity = &storage.Entity{}
	//	s.Entity.OdataMetadata = "https://sastrikesg857c0jo.table.core.windows.net/$metadata#powerplantstatus" // we don't need this.
	s.Entity.PartitionKey = "1"
	s.Entity.RowKey = "0"
	s.Entity.TimeStamp = time.Now()
	m := make(map[string]interface{})

	m["Foo"] = s.Foo
	//	m["Foo@odata.type"] = "Edm.String" // you can omit this.  // we don't need to add this. this is added by the library.
	m["Bar"] = s.Bar
	//	m["Bar@odata.type"] = "Edm.String" // you can omit this as well.
	m["Baz"] = s.Baz
	//	m["Baz@odata.type"] = "Edm.DateTime" // you can't omit this.
	s.Entity.Properties = m
	s.Entity.Table = table
}

func main() {
	content, err := ioutil.ReadFile("powerplant")
	if err != nil {
		panic(err)
	}
	var powerPlantConfig PowerPlantConfig

	err = json.Unmarshal(content, &powerPlantConfig)
	if err != nil {
		panic(err)
	}
	client, err := storage.NewBasicClient(powerPlantConfig.StorageAccountName, powerPlantConfig.StorageAccountKey)
	if err != nil {
		panic(err)
	}
	tableService := client.GetTableService()
	table := tableService.GetTableReference(s.DEFAULT_STORAGE_TABLE_NAME)
	// tableBatch := table.NewBatch()

	// We can't use Embedded for this purpose.
	instance := &StrikesInstance{
		Foo: "baz",
		Bar: "foo",
		Baz: time.Now(),
	}
	instance.Setup(table)
	fmt.Println("inserting entity")
	// tableBatch.InsertOrReplaceEntityByForce(instance.Entity)
	// err = tableBatch.ExecuteBatch()
	// if err != nil {
	// 	panic(err)
	// }
	// entity := table.GetEntityReference("1", "0")
	// jsonString, _ := json.Marshal(entity)
	// fmt.Println(string(jsonString))
	instance.Entity.InsertOrReplace(&storage.EntityOptions{}) // TimeOut is second

	fmt.Println("done!")
}
