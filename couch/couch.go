package couch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/byuoitav/common/db/couch"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
)

// ConfigFile represents a generic config file for shipwright
type ConfigFile struct {
	Path string          `json:"path"`
	File json.RawMessage `json:"contents"`
}

type configFileQueryResponse struct {
	Docs []struct {
		Rev string `json:"_rev,omitempty"`
		*ConfigFile
	} `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}

// UpdateConfigFiles updates the config files on disk from couchdb, using database db
func UpdateConfigFiles(ctx context.Context, db string) *nerr.E {
	errMsg := "unable to update config files"

	if len(db) == 0 {
		return nerr.Createf("", "%s: must pass a valid db name", errMsg)
	}

	addr := os.Getenv("COUCH_ADDRESS")
	if len(addr) == 0 {
		return nerr.Createf("", "%s: COUCH_ADDRESS is not set", errMsg)
	}

	url := fmt.Sprintf("%s/%s/_find", strings.Trim(addr, "/"), db)
	log.L.Infof("Updating config files from %s", url)

	query := []byte(`{
	"selector": {
		"_id": {
			"$regex": ""
		}
	}
	}`)

	// build request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(query))
	if err != nil {
		return nerr.Translate(err).Add(errMsg)
	}

	req = req.WithContext(ctx)
	req.Header.Add("content-type", "application/json")

	// add auth
	uname := os.Getenv("COUCH_USERNAME")
	pass := os.Getenv("COUCH_PASSWORD")
	if len(uname) > 0 && len(pass) > 0 {
		log.L.Infof("%v-%v", uname, pass)
		req.SetBasicAuth(uname, pass)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nerr.Translate(err).Add(errMsg)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nerr.Translate(err).Add(errMsg)
	}

	if resp.StatusCode/100 != 2 {
		ce := couch.CouchError{}
		err = json.Unmarshal(b, &ce)
		if err != nil {
			return nerr.Translate(err).Addf("%s: received a %d response from %s. body: %s", errMsg, resp.StatusCode, url, b)
		}

		return nerr.Createf(ce.Error, "%s: %+v", errMsg, ce)
	}

	// read the docs from the response
	docs := configFileQueryResponse{}
	err = json.Unmarshal(b, &docs)
	if err != nil {
		return nerr.Translate(err).Add(errMsg)
	}

	configs := []ConfigFile{}
	for i := range docs.Docs {
		if docs.Docs[i].ConfigFile != nil {
			configs = append(configs, *docs.Docs[i].ConfigFile)
		}
	}

	return WriteFilesToDisk(configs)
}

// WriteFilesToDisk writes each of the config files to disk
func WriteFilesToDisk(configs []ConfigFile) *nerr.E {
	for _, config := range configs {
		log.L.Infof("Writing new config file to %s", config.Path)

		// create dirs in case they don't exist
		dir := filepath.Dir(config.Path)
		err := os.MkdirAll(dir, 0775)
		if err != nil {
			return nerr.Translate(err).Addf("unable to write %s to disk", config.Path)
		}

		f, err := os.Create(config.Path)
		if err != nil {
			return nerr.Translate(err).Addf("unable to write %s to disk", config.Path)
		}

		_, err = f.Write(config.File)
		if err != nil {
			return nerr.Translate(err).Addf("unable to write %s to disk", config.Path)
		}
	}

	return nil
}
