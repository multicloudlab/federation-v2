/*
Copyright 2018 The Kubernetes Authors.

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

package adapters

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/google/go-github/github"
	"github.com/kubernetes-sigs/federation-v2/pkg/controller/util"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
)

type GitHubAdapter struct {
	workQueue workqueue.Interface
	name      string

	repo        string
	sourceOwner string
	password    string
	sourceRef   string
	authorName  string
	authorEmail string
	message     string
}

type GithubItem struct {
	cluster string
	file    string
	content string
}

func NewGitHubAdapter() (*GitHubAdapter, error) {
	g := &GitHubAdapter{
		name: "GithubAdaptor",

		repo:        "??", // TODO: currently, repo is cluster name.
		sourceOwner: "??", // e.g. xunpan
		password:    "??", // password or your account
		sourceRef:   "refs/heads/master",
		authorName:  "federationv2", // TODO: polish this information
		authorEmail: "federationv2@k8s.io",
		message:     "this is pushed by Kubernetes federation v2 configuration generation",

		workQueue: workqueue.New(),
	}

	// TODO: use common stopChan if possible
	stopChan := make(chan struct{})
	g.Run(stopChan)

	return g, nil
}

func (g *GitHubAdapter) PushToHandler(ops []util.FederatedOperation) error {
	for _, op := range ops {
		y, err := yaml.Marshal(op.Obj)
		if err != nil {
			glog.Errorf(fmt.Sprintf("Github: Error in marshal object cluster <%s> type <%s> key <%s> error: %s", op.ClusterName, op.Type, op.Key, err.Error()))
		} else {
			glog.Infof(fmt.Sprintf("Github: cluster <%s> type <%s> key <%s>: %s", op.ClusterName, op.Type, op.Key, y))
		}

		switch op.Type {
		case util.OperationTypeAdd:
		case util.OperationTypeUpdate:
		case util.OperationTypeDelete:
		}

		filename := filepath.Join(op.Key, "target.yaml")
		glog.Infof(fmt.Sprintf("Github: create file <%s> for cluster <%s>", filename, op.ClusterName))

		item := &GithubItem{
			cluster: op.ClusterName,
			file:    filename,
			content: string(y),
		}

		g.workQueue.Add(item)
	}

	return nil
}

func (g *GitHubAdapter) Run(stopChan <-chan struct{}) {
	go wait.Until(g.worker, time.Second, stopChan)
}

func (g *GitHubAdapter) worker() {
	ctx := context.Background()
	auth := github.BasicAuthTransport{
		Username: g.sourceOwner,
		Password: g.password,
	}
	client := github.NewClient(auth.Client())

	glog.Infof(fmt.Sprintf("Github: a GitHubAdapter.worker is running"))

	for {
		obj, quit := g.workQueue.Get()
		if quit {
			return
		}
		item := obj.(*GithubItem)
		glog.Infof(fmt.Sprintf("Github: get item cluster <%s> file <%s> with content for handling", item.cluster, item.file))

		sourceRepo := github.String(item.cluster)

		ref, _, err := client.Git.GetRef(ctx, g.sourceOwner, *sourceRepo, g.sourceRef)
		if err != nil {
			glog.Errorf(fmt.Sprintf("Github: client.Git.GetRef <%s>\n", err.Error()))
			// TODO: error handling for all related part
		}

		entries := []github.TreeEntry{}
		entries = append(entries, github.TreeEntry{Path: github.String(item.file), Type: github.String("blob"), Content: github.String(item.content), Mode: github.String("100644")})
		tree, _, err := client.Git.CreateTree(ctx, g.sourceOwner, *sourceRepo, *ref.Object.SHA, entries)
		if err != nil {
			glog.Errorf(fmt.Sprintf("Github: client.Git.CreateTree <%s>\n", err.Error()))
		}

		// Get the parent commit to attach the commit to.
		parent, _, err := client.Repositories.GetCommit(ctx, g.sourceOwner, *sourceRepo, *ref.Object.SHA)
		if err != nil {
			glog.Errorf(fmt.Sprintf("Github: client.Repositories.GetCommit <%s>\n", err.Error()))
		}
		parent.Commit.SHA = parent.SHA

		// Create the commit using the tree.
		date := time.Now()
		author := &github.CommitAuthor{Date: &date, Name: github.String(g.authorName), Email: github.String(g.authorEmail)}
		commit := &github.Commit{Author: author, Message: github.String(g.message), Tree: tree, Parents: []github.Commit{*parent.Commit}}
		newCommit, _, err := client.Git.CreateCommit(ctx, g.sourceOwner, *sourceRepo, commit)
		if err != nil {
			glog.Errorf(fmt.Sprintf("Github: client.Git.CreateCommit <%s>\n", err.Error()))
		}

		// Attach the commit to the master branch.
		ref.Object.SHA = newCommit.SHA
		_, _, err = client.Git.UpdateRef(ctx, g.sourceOwner, *sourceRepo, ref, false)
		if err != nil {
			glog.Errorf(fmt.Sprintf("Github: client.Git.UpdateRef <%s>\n", err.Error()))
		}

		glog.Infof(fmt.Sprintf("Github: a github push request should be done for file <%s>", item.file))

		g.workQueue.Done(item)
	}
}
