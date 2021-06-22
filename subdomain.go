package main

import (
	"bytes"
	"context"
	"github.com/projectdiscovery/subfinder/v2/pkg/passive"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"io"
	"io/ioutil"
	"log"
)

func find(domain string) ([]byte, error) {
	resolve.DefaultResolvers = append(resolve.DefaultResolvers, "114.114.114.114", "114.114.115.115")
	config := runner.ConfigFile{
		Resolvers:  resolve.DefaultResolvers,
		Sources:    passive.DefaultSources,
		AllSources: passive.DefaultAllSources,
		Recursive:  passive.DefaultRecursiveSources,
	}
	runnerInstance, err := runner.NewRunner(&runner.Options{
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
		YAMLConfig:         config,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	buf := bytes.Buffer{}
	err = runnerInstance.EnumerateSingleDomain(context.Background(), domain, []io.Writer{&buf})
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(&buf)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return data, nil
}
