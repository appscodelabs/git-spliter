package main

import (
	"fmt"

	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/appscode/go/types"
	shell "github.com/codeskyblue/go-sh"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func NewCmdInit() *cobra.Command {
	var (
		repoDir    string
		pkgDir     string
		branchName = "master"
		ghUser     string
		ghToken    = os.Getenv("GITHUB_AUTH_TOKEN")
		// ghLicense  = "apache-2.0" // https://help.github.com/articles/licensing-a-repository/#searching-github-by-license-type
		ghTopics   []string
	)
	cmd := &cobra.Command{
		Use:               "init",
		Short:             "Init",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			packages, err := ioutil.ReadDir(filepath.Join(repoDir, pkgDir))
			if err != nil {
				return err
			}
			for _, pkg := range packages {
				if !pkg.IsDir() {
					continue
				}
				fmt.Printf("GITHUB_AUTH_TOKEN: %s\n\n", ghToken)

				scratchDir, err := ioutil.TempDir(os.TempDir(), "git-spliter")
				if err != nil {
					return err
				}
				fmt.Printf("Using scratch dir: %s\n\n", scratchDir)

				sh := shell.NewSession()
				sh.ShowCMD = true

				err = sh.Command("cp", "--recursive", repoDir, scratchDir).Run()
				if err != nil {
					return err
				}
				sh.SetDir(filepath.Join(scratchDir, filepath.Base(repoDir)))
				fmt.Printf("Chnaged dir: %s\n\n", filepath.Join(scratchDir, filepath.Base(repoDir)))

				err = sh.Command("git", "filter-branch", "--prune-empty", "--subdirectory-filter", filepath.Join(pkgDir, pkg.Name()), branchName).Run()
				if err != nil {
					return err
				}

				ctx := context.Background()
				ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
				tc := oauth2.NewClient(ctx, ts)
				client := github.NewClient(tc)

				desc := pkg.Name()
				desc = strings.Replace(desc, "-", " ", -1)
				desc = strings.Replace(desc, "_", " ", -1)
				desc = strings.Title(desc)

				r := &github.Repository{
					Name:        types.StringP(pkg.Name()),
					Description: types.StringP(desc),
					AutoInit:    types.FalseP(),
					Private:     types.FalseP(),
					// HasIssues         *bool   `json:"has_issues,omitempty"`
					HasWiki:     types.FalseP(),
					HasPages:    types.FalseP(),
					HasProjects: types.FalseP(),
					// HasDownloads      *bool   `json:"has_downloads,omitempty"`
					// LicenseTemplate: types.StringP(ghLicense),
				}
				repo, resp, _ := client.Repositories.Get(ctx, ghUser, *r.Name)
				if resp.StatusCode != 200 {
					repo, _, err = client.Repositories.Create(ctx, ghUser, r)
					if err != nil {
						return err
					}
					fmt.Printf("Successfully created new repo: %v\n", repo.GetSSHURL())
					client.Repositories.ReplaceAllTopics(ctx, ghUser, *r.Name, ghTopics)
				} else {
					fmt.Printf("Repo %s exists\n", repo.GetSSHURL())
				}

				err = sh.Command("git", "remote", "set-url", "origin", repo.GetSSHURL()).Run()
				if err != nil {
					return err
				}
				//err = sh.Command("git", "pull", "--rebase", "origin", branchName).Run()
				//if err != nil {
				//	return err
				//}
				err = sh.Command("git", "push", "-u", "origin", branchName).Run()
				if err != nil {
					return err
				}

				// git remote set-url origin https://github.com/USERNAME/NEW-REPOSITORY-NAME.git
				// git push -u origin BRANCH-NAME
			}
			return err
		},
	}

	cmd.Flags().StringVar(&repoDir, "repo-dir", repoDir, "Directory where git repo is located")
	cmd.Flags().StringVar(&pkgDir, "pkg-dir", pkgDir, "Directory of git repo whose direct sub folders will be split into new repos")
	cmd.Flags().StringVar(&branchName, "branch-name", branchName, "Branch name from where new repo is created")

	cmd.Flags().StringVar(&ghUser, "github-username", ghUser, "Name of repo to create in authenticated user's GitHub account.")
	cmd.Flags().StringVar(&ghToken, "github-token", ghToken, "Name of repo to create in authenticated user's GitHub account.")
	cmd.Flags().StringSliceVar(&ghTopics, "github-topics", ghTopics, "Github repo topics")

	return cmd
}
