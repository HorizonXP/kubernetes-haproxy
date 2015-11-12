/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// A small script that converts the given open ssl public/private keys to
// a secret that it writes to stdout as json. Most common use case is to
// create a secret from self signed certificates used to authenticate with
// a devserver. Usage: go run make_secret.go -crt ca.crt -key priv.key > secret.json
package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/latest"
	"k8s.io/kubernetes/pkg/runtime"

	// This installs the legacy v1 API
	_ "k8s.io/kubernetes/pkg/api/install"
)

type Certificate struct {
	Path        string `json:"path"`
	Destination string `json:"destination"`
}

type Service struct {
	Name        string `json:"name"`
	Certificate Certificate `json:"certificate"`
}

type Stats struct {
	Certificate Certificate `json:"certificate"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type Config struct {
	VolName     string `json:"volName"`
	Certificate Certificate `json:"certificate"`
	Stats       Stats `json:"stats"`
	Services    []Service `json:"services"`
}

// TODO:
// Add a -o flag that writes to the specified destination file.
// Teach the script to create crt and key if -crt and -key aren't specified.
var (
	config = flag.String("config", "secrets.json", "path to secrets configuration.")
)

func read(file string) []byte {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Cannot read file %v, %v", file, err)
	}
	return b
}

func main() {
	flag.Parse()
	if *config == "" {
		log.Fatalf("Need to specify -config")
	}
	jsonBlob := read(*config)
	var cfg Config
	var err = json.Unmarshal(jsonBlob, &cfg)
	if err != nil {
		log.Fatalf("Unable to unmarshal json blob: %v", string(jsonBlob))
	}
	data := map[string][]byte{}

	if cfg.Stats.Username != "" && cfg.Stats.Password != "" {
		data["stats-username"] = []byte(cfg.Stats.Username)
		data["stats-password"] = []byte(cfg.Stats.Password)
	}
	if cfg.Stats.Certificate.Path != "" && cfg.Stats.Certificate.Destination != "" {
		data[cfg.Stats.Certificate.Destination] = read(cfg.Stats.Certificate.Path)
	}
	if cfg.Certificate.Path != "" && cfg.Certificate.Destination != "" {
		data[cfg.Certificate.Destination] = read(cfg.Certificate.Path)
	}
	for _, svc := range cfg.Services {
		if svc.Certificate.Path != "" && svc.Certificate.Destination != "" {
			data[svc.Certificate.Destination] = read(svc.Certificate.Path)
		}
	}

	secret := &api.Secret{
		ObjectMeta: api.ObjectMeta{
			Name: cfg.VolName,
		},
		Data: data, 
	}
	fmt.Printf(runtime.EncodeOrDie(latest.GroupOrDie("").Codec, secret))
}
