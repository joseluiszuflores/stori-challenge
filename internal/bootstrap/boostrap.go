package bootstrap

import (
	"context"
	"github.com/golang/glog"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"
	"github.com/joseluiszuflores/stori-challenge/internal/config"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/email"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/file"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/storage/dynamo/client"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/storage/dynamo/conn"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/storage/dynamo/migration"
	transaction2 "github.com/joseluiszuflores/stori-challenge/internal/platform/storage/dynamo/transaction"
	"github.com/joseluiszuflores/stori-challenge/internal/transaction"
)

// Run call all configurations and connections to DB.
func Run(path, userId string) error {
	fileService := file.NewService(path)
	transactions, err := fileService.GetDataFile()
	if err != nil {
		glog.Error(err)

		return err
	}
	config.Init()
	return Setup(transactions, userId)
}

func Setup(transactions mooc.Transactions, userId string) error {
	key, secret := config.Config.AWSAccessKey, config.Config.AWSSecretAcessKey
	region, url := config.Config.AWSRegion, config.Config.AWSUrlDynamoDev

	// new configuration of AWS data.
	cnf, err := conn.NewAWSConfig(key, secret, region, url, config.Config.DevEnv)
	if err != nil {
		glog.Errorf("error loading the configuration of lambda: [%s]", err)

		return err
	}

	// Call to new instance of DynamoDB.
	//nolint:varnamelen
	db := conn.NewDynamoDBClient(cnf)
	// Instance of migrator to does the migration.
	m := migration.NewMigrator(db)
	// does migration of client table and created transaction table.
	if err := m.Do(); err != nil {
		glog.Errorf("error migrating: [%s]", err)

		return err
	}
	smtp := email.NewSMTPService(config.Config.SMTPHost, config.Config.SMTPPort, config.Config.SMTPUsername, config.Config.SMTPPassword, config.Config.SMTPUsername, config.Config.SMTPTemplatePath)

	clientRep := client.NewRepository(db)
	transactionRep := transaction2.NewRepository(db)

	service, err := transaction.NewService(userId, transactions, smtp, clientRep, transactionRep)
	if err != nil {
		return err
	}
	if err := service.SummaryTransaction(context.TODO()); err != nil {
		return err
	}
	return nil
}
