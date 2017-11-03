/*
Copyright 2017 HomeAway, Inc

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

package cobblerclient

import (
	"fmt"
	"reflect"
)

// Repo is a created repo.
type Repo struct {
	// These are internal fields and cannot be modified.
	Ctime         float64 `mapstructure:"ctime"          cobbler:"noupdate"` // TODO: convert to time
	Depth         int     `mapstructure:"depth"          cobbler:"noupdate"`
	ID            string  `mapstructure:"uid"            cobbler:"noupdate"`
	Mtime         float64 `mapstructure:"mtime"          cobbler:"noupdate"` // TODO: convert to time
	TreeBuildTime string  `mapstructure:tree_build_time" cobbler:"noupdate"`

	AptComponents   []string `mapstructure:"apt_components"`
	AptDists        []string `mapstructure:"apt_dists"`
	Arch            string   `mapstructure:"arch"`
	Breed           string   `mapstructure:"breed"`
	Comment         string   `mapstructure:"comment"`
	CreateRepoFlags string   `mapstructure:"createrepo_flags"`
	Environment     string   `mapstructure:"environment"`
	KeepUpdated     bool     `mapstructure:"keep_updated"`
	Mirror          string   `mapstructure:"mirror"`
	MirrorLocally   bool     `mapstructure:"mirror_locally"`
	Name            string   `mapstructure:"name"`
	Owners          []string `mapstructure:"owners"`
	Proxy           string   `mapstructure:"proxy" cobbler:"newfield"`
	RpmList         []string `mapstructure:"rpm_list"`
	//YumOpts                map[string]interface{} `mapstructure:"yumopts"`
}

// GetRepos returns all repos in Cobbler.
func (c *Client) GetRepos() ([]*Repo, error) {
	var repos []*Repo

	result, err := c.Call("get_repos", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	for _, r := range result.([]interface{}) {
		var repo Repo
		decodedResult, err := decodeCobblerItem(r, &repo)
		if err != nil {
			return nil, err
		}

		repos = append(repos, decodedResult.(*Repo))
	}

	return repos, nil
}

// GetRepo returns a single repo obtained by its name.
func (c *Client) GetRepo(name string) (*Repo, error) {
	var repo Repo

	result, err := c.Call("get_repo", name, c.Token)
	if result == "~" {
		return nil, fmt.Errorf("Repo %s not found.", name)
	}

	if err != nil {
		return nil, err
	}

	decodeResult, err := decodeCobblerItem(result, &repo)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*Repo), nil
}

// CreateRepo creates a repo.
func (c *Client) CreateRepo(repo Repo) (*Repo, error) {
	// Make sure a repo with the same name does not already exist
	if _, err := c.GetRepo(repo.Name); err == nil {
		return nil, fmt.Errorf("A Repo with the name %s already exists.", repo.Name)
	}

	result, err := c.Call("new_repo", c.Token)
	if err != nil {
		return nil, err
	}
	newId := result.(string)

	item := reflect.ValueOf(&repo).Elem()
	if err := c.updateCobblerFields("repo", item, newId); err != nil {
		return nil, err
	}

	if _, err := c.Call("save_repo", newId, c.Token); err != nil {
		return nil, err
	}

	return c.GetRepo(repo.Name)
}

// UpdateRepo updates a single repo.
func (c *Client) UpdateRepo(repo *Repo) error {
	item := reflect.ValueOf(repo).Elem()
	id, err := c.GetItemHandle("repo", repo.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("repo", item, id); err != nil {
		return err
	}

	if _, err := c.Call("save_repo", id, c.Token); err != nil {
		return err
	}

	return nil
}

// DeleteRepo deletes a single repo by its name.
func (c *Client) DeleteRepo(name string) error {
	_, err := c.Call("remove_repo", name, c.Token)
	return err
}
