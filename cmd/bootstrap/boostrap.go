package bootstrap

import (
	"github.com/golang/glog"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/storage/dynamo/conn"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/storage/dynamo/migration"
)

// Run call all configurations and connections to DB.
func Run() error {
	// new configuration of AWS data.
	cnf, err := conn.NewAWSConfig()
	if err != nil {
		glog.Errorf("error loading the configuration of aws: [%s]", err)

		return err
	}

	// Call to new instance of DynamoDB.
	db := conn.NewDynamoDBClient(cnf)
	// Instance of migrator to does the migration.
	m := migration.NewMigrator(db)
	// does migration of client table and created transaction table.
	if err := m.Do(); err != nil {
		glog.Errorf("error migrating: [%s]", err)

		return err
	}

	return nil
}
