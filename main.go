package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/bitrise-step-generate-universal-apk/apkexporter"
	"github.com/bitrise-steplib/bitrise-step-generate-universal-apk/bundletool"
	"github.com/bitrise-steplib/bitrise-step-generate-universal-apk/filedownloader"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config is defining the input arguments required by the Step.
type Config struct {
	DeployDir        string `env:"BITRISE_DEPLOY_DIR"`
	AABPath          string `env:"aab_path,required"`
	KeystoreURL      string `env:"keystore_url"`
	KeystotePassword string `env:"keystore_password"`
	KeyAlias         string `env:"key_alias"`
	KeyPassword      string `env:"private_key_password"`
}

func main() {
	var config Config
	if err := stepconf.Parse(&config); err != nil {
		failf("Error: %s \n", err)
	}
	stepconf.Print(config)
	fmt.Println()

	bundletoolTool, err := bundletool.New("0.15.0", filedownloader.New(http.DefaultClient))
	log.Infof("bundletool path created at: %s", bundletoolTool.Path())
	if err != nil {
		failf("Failed to initialize bundletool: %s \n", err)
	}

	exporter := apkexporter.New(bundletoolTool, filedownloader.New(http.DefaultClient))
	keystoreCfg := parseKeystoreConfig(config)
	apkPath, err := exporter.ExportUniversalAPK(config.AABPath, config.DeployDir, keystoreCfg)
	if err != nil {
		failf("Failed to export apk, error: %s \n", err)
	}

	if err = tools.ExportEnvironmentWithEnvman("APK_PATH", apkPath); err != nil {
		failf("Failed to export APK_PATH, error: %s \n", err)
	}

	log.Donef("Success APK exported to: %s", apkPath)
	os.Exit(0)
}

func parseKeystoreConfig(config Config) *bundletool.KeystoreConfig {
	if config.KeystoreURL == "" ||
		config.KeystotePassword == "" ||
		config.KeyAlias == "" ||
		config.KeyPassword == "" {
		return nil
	}

	return &bundletool.KeystoreConfig{
		Path:               config.KeystoreURL,
		KeystorePassword:   config.KeystotePassword,
		SigningKeyAlias:    config.KeyAlias,
		SigningKeyPassword: config.KeyPassword}
}

func failf(s string, a ...interface{}) {
	log.Errorf(s, a...)
	os.Exit(1)
}
