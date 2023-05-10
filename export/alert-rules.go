package export

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
)

func AlertRules(username, password, url, directory string) error {
	var (
		err error
	)
	folderName := "alertRules"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}
	path := filepath.Join(directory, folderName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0760)
		if err != nil {
			log.Fatal("Error creating directory", err)
		}
	}
	alertRules, err := client.AlertRules()
	if err != nil {
		log.Error("Failed to get AlertRules", err)
		return err
	}
	for _, alertRule := range alertRules {
		jsonAlertRule, err := json.Marshal(alertRule)
		if err != nil {
			log.Error("Error unmarshalling json File", err)
		}
		err = os.WriteFile(filepath.Join(path, slug.Make(alertRule.Title))+".json", jsonAlertRule, os.FileMode(0666))
		if err != nil {
			log.Error("Couldn't write AlertRule to disk", err)
		}
	}
	return nil
}
