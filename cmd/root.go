/*
Copyright © 2023 jinrui.cui

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

import (
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/infra-kubernetes/download/pkg/exec"
	"github.com/infra-kubernetes/download/pkg/file"
	"github.com/infra-kubernetes/download/pkg/version"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"strings"
)

var downloadPath string
var cleanup bool
var downloadTxt string
var registryPassword, registryUserName, registryDomain string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "download",
	Version: version.Get().String(),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("download path is %s", downloadPath)
		if cleanup {
			file.CleanDir(downloadPath)
		}
		_ = file.MkDirs(downloadPath)
		extraPackages()
		downloadPackages()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if downloadTxt == "" {
			logger.Error("download txt file is empty")
			os.Exit(1)
		}
		if !file.IsExist(downloadPath) {
			file.MkDirs(downloadPath)
			logger.Warn("download path %s is not exist, create it", downloadPath)
		}
		if registryPassword == "" {
			logger.Error("registry password is empty")
			os.Exit(1)
		}
		if registryUserName == "" {
			logger.Error("registry username is empty")
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(func() {
		logger.Cfg(true, false)
	})

	rootCmd.PersistentFlags().StringVarP(&downloadPath, "download-path", "d", "/tmp/images", "download path")
	rootCmd.PersistentFlags().BoolVarP(&cleanup, "cleanup", "c", false, "cleanup download package")
	rootCmd.PersistentFlags().StringVarP(&downloadTxt, "download-txt", "f", "", "download images txt file")
	rootCmd.PersistentFlags().StringVarP(&registryUserName, "registry-username", "u", "cuisongliu", "registry username")
	rootCmd.PersistentFlags().StringVarP(&registryPassword, "registry-password", "p", "", "registry password")
	rootCmd.PersistentFlags().StringVarP(&registryDomain, "registry-domain", "r", "ghcr.io", "registry domain")
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.install.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
func extraPackages() {
	fname := fmt.Sprintf("files/sealos_4.3.8_linux_%s.tar.gz", runtime.GOARCH)
	if !file.IsExist(fname) {
		logger.Fatal("sealos package %s not exist", fname)
	}
	if err := exec.RunCommand(fmt.Sprintf("tar -xvf %s sealos && mv sealos files/ && sudo chmod a+x files/sealos", fname)); err != nil {
		logger.Fatal("unzip %s error %s", fname, err)
	}
}
func downloadPackages() {
	lines := getDownloadImages(downloadTxt)
	if lines == nil {
		return
	}
	if err := exec.ExecShellForAny()([]any{
		exec.Logger("sealos login " + registryDomain),
		exec.SecretShell(fmt.Sprintf("sudo files/sealos login -u %s -p %s %s", registryUserName, registryPassword, registryDomain)),
	}); err != nil {
		logger.Fatal("sealos login error %s", err)
		return
	}
	for _, line := range lines {
		logger.Info("pull image: %s", line)
		err := exec.RunCommand(fmt.Sprintf("sudo files/sealos pull --policy=always %s", line))
		if err != nil {
			logger.Painc("pull image %s error %s", line, err)
			continue
		}
		//save tar
		tarName := processImage2TarName(line)
		if err = exec.RunCommand(fmt.Sprintf("sudo files/sealos save -o  %s/%s.tar %s", downloadPath, tarName, line)); err != nil {
			logger.Painc("save image %s error: %s", line, err)
		}
	}
	_ = exec.RunCommand("sudo files/sealos rmi `sudo files/sealos images -aq `")
	logger.Info("download package [%d] success, path is %s", len(lines), downloadPath)
}

func processImage2TarName(src string) string {
	all := strings.Split(src, "/")
	if len(all) != 3 {
		return src
	}
	name := strings.Split(all[2], ":")
	if len(name) != 2 {
		return ""
	}
	return fmt.Sprintf("%s_%s.tar", name[0], name[1])
}

func getDownloadImages(f string) []string {
	var lines []string
	lines, err := file.ReadLines(f)
	if err != nil {
		logger.Error("get download urls error %s", err)
		return nil
	}
	var newLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			logger.Info("skip line %s", line)
			continue
		}
		newLines = append(newLines, line)
	}
	return newLines
}
