//revive:disable:package-comments,exported
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/debug"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

func main() {
	ctx := context.Background()

	// Main project
	createAndDeployStack(ctx, "trade", "dev-figaro", "localproject_main", auto.ConfigMap{
		"github:token": auto.ConfigValue{
			Value:  os.Getenv("GITHUB_TOKEN"),
			Secret: true,
		},
		"github:owner": auto.ConfigValue{
			Value:  os.Getenv("GITHUB_OWNER"),
			Secret: true,
		},
	})

	// Mirrored project
	createAndDeployStack(ctx, "gitlab-mirror", "dev-figaro-mirrored", "localproject_mirrored", auto.ConfigMap{
		"gitlab:token": auto.ConfigValue{
			Value:  os.Getenv("GITLAB_TOKEN"),
			Secret: true,
		},
	})

}

func createAndDeployStack(ctx context.Context, projectName, stackStr, workDir string, configMap auto.ConfigMap) {
	org := os.Getenv("PULUMI_ORG_NAME")
	pat := os.Getenv("PULUMI_ACCESS_TOKEN")
	stackName := auto.FullyQualifiedStackName(org, projectName, stackStr)
	workDirPath := filepath.Join(workDir)

	stack, err := auto.NewStackLocalSource(ctx, stackName, workDirPath)
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackName + " already exists")
	} else if err != nil {
		panic(err)
	}

	err = stack.Workspace().SetEnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_CONFIG_PASSPHRASE": "",
		"PULUMI_ACCESS_TOKEN":      pat,
	})
	if err != nil {
		panic(err)
	}

	err = stack.SetAllConfig(ctx, configMap)
	if err != nil {
		panic(err)
	}

	refr, err := stack.Refresh(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(refr.StdOut)

	prev, err := stack.Preview(ctx, optpreview.DebugLogging(debug.LoggingOptions{Debug: true}))
	if err != nil {
		panic(err)
	}
	log.Println(prev.StdOut)

	up, err := stack.Up(ctx, optup.DebugLogging(debug.LoggingOptions{Debug: true}))
	if err != nil {
		panic(err)
	}
	log.Println(up.StdOut)
}
