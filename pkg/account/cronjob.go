package account

import (
	"fmt"
	"gmail_backup/pkg/gmail"
	"log"
)

// Start will run and create the cronjobs
func (a *Account) Start() error {
	a.cronjob.Start()
	err := a.createCronjobs()
	if err != nil {
		return err
	}

	return nil
}

// Reset will recreate the cronjobs
func (a *Account) Reset() error {
	for _, entry := range a.cronjob.Entries() {
		a.cronjob.Remove(entry.ID)
	}
	err := a.createCronjobs()
	if err != nil {
		return err
	}
	return nil
}

// Close will close the cronjobs
func (a *Account) Close() {
	a.cronjob.Stop()
}

func (a *Account) createCronjobs() error {
	accounts, err := a.getAllAccounts()
	if err != nil {
		return err
	}

	for _, ac := range accounts {

		expression := ac.CronExpression
		account := ac

		// fmt.Printf("%s  ---   %s\n", account.CronExpression, account.Email)

		job := func() {
			g, err := gmail.New(a.config, a.db)
			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Println("Backing up account " + account.Email)
			g.Backup(account, a.storage)
		}

		_, err := a.cronjob.AddFunc(expression, job)

		if err != nil {
			return err
		}
	}

	// fmt.Printf("%+v\n", a.cronjob.Entries())

	return nil
}
