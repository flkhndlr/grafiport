package restore

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	gapi "github.com/grafana/grafana-api-golang-client"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
)

func DataSources(username, password, url, directory string) error {
	var (
		filesInDir    []os.DirEntry
		rawDatasource []byte
	)
	folderName := "dataSources"
	userInfo := url2.UserPassword(username, password)
	config := gapi.Config{BasicAuth: userInfo}
	client, err := gapi.New(url, config)
	if err != nil {
		log.Error("Failed to create a client%s\n", err)
		return err
	}

	path := filepath.Join(directory, folderName)

	filesInDir, err = os.ReadDir(path)
	if err != nil {
		log.Error("Failed to read folder%s\n", err)
		return err
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawDatasource, err = os.ReadFile(filepath.Join(path, file.Name())); err != nil {
				log.Error(err)
				continue
			}

			var newDatasource gapi.DataSource
			if err = json.Unmarshal(rawDatasource, &newDatasource); err != nil {
				log.Error(err)
				continue
			}
			status, err := client.DataSourceByUID(newDatasource.UID)
			if err != nil {
				log.Error("Failed Status Check if Datasource already exists")
				continue
			}
			if status != nil {
				err = client.UpdateDataSource(&newDatasource)
				if err != nil {
					log.Error("Error updating Datasource", err)
				}

			} else {
				_, err = client.NewDataSource(&newDatasource)
				if err != nil {
					log.Error("Error creating Datasource", err)
				}
			}
		}
	}
	return nil
}
