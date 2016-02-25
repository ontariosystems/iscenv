/*
Copyright 2016 Ontario Systems

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

package cmd

//import (
//	"os"
//
//	"github.com/ontariosystems/iscenv/iscenv"
//	"github.com/ontariosystems/iscenv/internal/app"
//
//	docker "github.com/fsouza/go-dockerclient"
//	"github.com/spf13/cobra"
//)
//
//// TODO: This command must go.
//var pullCmd = &cobra.Command{
//	Use:   "pull",
//	Short: "Pull the latest ISC product versions",
//	Long:  "Pull the latest versions of the ISC product images.",
//	Run:   pull,
//}
//
//func init() {
//	rootCmd.AddCommand(pullCmd)
//
//}
//
//func pull(_ *cobra.Command, _ []string) {
//	imgopts := docker.PullImageOptions{Registry: iscenv.Registry, Repository: iscenv.Repository, OutputStream: os.Stdout}
//	authcfg := app.GetAuthConfig()
//	err := app.DockerClient.PullImage(imgopts, authcfg)
//	if err != nil {
//		app.Fatalf("Could not pull latest ISC product version images, error: %s\n", err)
//	}
//}
