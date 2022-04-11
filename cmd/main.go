package main

import (
	"os"
	"time"

	"github.com/akumor/keepassnotifier/internal/entries"
	"github.com/akumor/keepassnotifier/internal/notifiers"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
)

type RootOptions struct {
	ConfigPath          string
	KeepassDatabasePath string
	CredentialsPath     string
}

func main() {
	// parse flags
	rootOpts := RootOptions{}
	rootCmd := &cobra.Command{
		Use:   "keepassnotifier",
		Short: "Notify based on keepass database information",
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			glog.Info("starting keepassnotifier")
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&rootOpts.ConfigPath, "config", "c", "", "path to config file for keepassnotifier")
	rootCmd.PersistentFlags().StringVarP(&rootOpts.KeepassDatabasePath, "database", "d", "", "path to keepass database file")
	rootCmd.PersistentFlags().StringVarP(&rootOpts.CredentialsPath, "credentials", "C", "", "path to json file for credentials")

	entriesCmd := &cobra.Command{
		Use:   "entries [arguments]",
		Short: "Notify based on attributes of keepass database entries",
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := NewConfig(rootOpts.ConfigPath)
			if err != nil {
				// TODO wrap the error and return it without logging
				glog.Errorf("failed to parse config: %v", err)
				return err
			}

			creds := &Credentials{}
			if rootOpts.CredentialsPath != "" {
				err = creds.LoadCredentialsFromFile(rootOpts.CredentialsPath)
				if err != nil {
					// TODO wrap the error and return it without logging
					glog.Errorf("failed to parse credentials: %v", err)
					return err
				}
			}
			err = creds.LoadCredentialsFromEnv()
			if err != nil {
				// TODO wrap the error and return it without logging
				glog.Errorf("failed to retrieve credentials from env: %v", err)
				return err
			}

			glog.Info("running entries command")
			now := time.Now()
			file, _ := os.Open(rootOpts.KeepassDatabasePath)

			db := gokeepasslib.NewDatabase()
			db.Credentials = gokeepasslib.NewPasswordCredentials(creds.Keepass.DatabasePassword)
			err = gokeepasslib.NewDecoder(file).Decode(db)
			if err != nil {
				// TODO wrap the error and return it without logging
				glog.Errorf("failed to decode specified keepass database file: %v", err)
				return err
			}

			err = db.UnlockProtectedEntries()
			if err != nil {
				// TODO wrap the error and return it without logging
				glog.Errorf("failed to unlock protected entries: %v", err)
				return err
			}

			nearlyExpiredEntries := make([]gokeepasslib.Entry, 0)
			for index := range db.Content.Root.Groups {
				newEntries, err := entries.GetNearlyExpiredEntries(db.Content.Root.Groups[index], now, cfg.Entries.Proximity)
				if err != nil {
					// TODO wrap the error and return it without logging
					glog.Errorf("failed to get nearly expired entries: %v", err)
					return err
				}
				nearlyExpiredEntries = append(nearlyExpiredEntries, newEntries...)
			}

			// log the nearly expired entries
			for index := range nearlyExpiredEntries {
				glog.Info(nearlyExpiredEntries[index].GetTitle())
			}
			glog.Infof("%v\n", nearlyExpiredEntries)

			// send email notifications via sendgrid if notifier is configured
			if cfg.Notifiers.SendGrid.FromEmailAddress != "" && cfg.Notifiers.SendGrid.ToEmailAddress != "" && creds.SendGrid.ApiKey != "" {
				sg, err := notifiers.NewSendGrid(creds.SendGrid.ApiKey, cfg.Notifiers.SendGrid.FromEmailAddress, cfg.Notifiers.SendGrid.ToEmailAddress)
				if err != nil {
					// TODO wrap the error and return it without logging
					glog.Errorf("failed to send notification via sendgrid: %v", err)
					return err
				}
				// build message to send via sendgrid
				message := "The following keepass database entries are close to expiring:<br>"
				for index := range nearlyExpiredEntries {
					message = message + nearlyExpiredEntries[index].GetTitle() + "\t" + nearlyExpiredEntries[index].Times.ExpiryTime.Time.String() + "<br>"
				}
				sg.Notify(message)
			}
			return nil
		},
	}

	rootCmd.AddCommand(entriesCmd)

	if err := rootCmd.Execute(); err != nil {
		// TODO wrap the error and log it before exiting
		os.Exit(1)
	}

}
