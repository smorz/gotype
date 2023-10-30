package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jimsmart/schema"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stoewer/go-strcase"
)

func main() {

	connector, err := mssql.NewConnector("sqlserver://m.razavian:1abomesvok@db1.corp.mabnadp.com?database=Beta_Rahavard365")
	if err != nil {
		log.Fatal(err)

	}
	connector.SessionInitSQL = "SET ANSI_NULLS ON"

	// Pass connector to sql.OpenDB to get a sql.DB object
	db := sql.OpenDB(connector)
	defer db.Close()

	schemaName := "Academy"
	os.Mkdir("type/", os.ModeDir|os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("type/" + strcase.SnakeCase(schemaName) + ".go")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString("package main\n")

	m, err := schema.Tables(db)
	if err != nil {
		log.Fatal(err)
	}
	for st, columns := range m {
		if st[0] != schemaName {
			continue
		}
		typeName := strcase.LowerCamelCase(st[1])
		file.WriteString("type " + typeName + " struct {\n")
		for _, c := range columns {
			name := strcase.UpperCamelCase(c.Name())
			if name == "Id" {
				file.WriteString("\tID\tint64 `gorm:\"column:Id\" gorm:\"primary_key,autoincrement\"`\n")
				continue
			}
			if strings.HasPrefix(name, "Record") {
				continue
			}
			fieldType := c.ScanType().String()
			if isNullable, _ := c.Nullable(); isNullable {
				fieldType = "*" + fieldType
			}
			file.WriteString("\t" + name + "\t" + fieldType + "\n")

		}
		file.WriteString(fmt.Sprintf(`}
func (%s) TableName() string{
	return "%s.%s"
}

`, typeName, schemaName, st[1]))

	}
}
