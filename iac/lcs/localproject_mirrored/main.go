//revive:disable:package-comments,exported
package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gitlab/sdk/v6/go/gitlab"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		groupName := "mirror-e/github-softwaredevelop"
		groupID, err := gitlab.LookupGroup(ctx, &gitlab.LookupGroupArgs{
			FullPath: pulumi.StringRef(groupName),
		}, nil)
		if err != nil {
			panic(err)
		}

		projectName := "figaro"
		projectDescription := fmt.Sprintf("A GitLab project for mirroring %s GitHub repository.", projectName)
		project, err := gitlab.NewProject(ctx, "newProjectFigaro", &gitlab.ProjectArgs{
			AutoCancelPendingPipelines:       pulumi.String("enabled"),
			BuildsAccessLevel:                pulumi.String("private"),
			Description:                      pulumi.String(projectDescription),
			IssuesEnabled:                    pulumi.Bool(true),
			LfsEnabled:                       pulumi.Bool(true),
			MergeMethod:                      pulumi.String("merge"),
			MergeRequestsEnabled:             pulumi.Bool(true),
			Name:                             pulumi.String(projectName),
			NamespaceId:                      pulumi.Int(groupID.GroupId),
			OnlyAllowMergeIfPipelineSucceeds: pulumi.Bool(true),
			RemoveSourceBranchAfterMerge:     pulumi.Bool(true),
			SharedRunnersEnabled:             pulumi.Bool(true),
			Topics: pulumi.StringArray{
				pulumi.String("dagger"),
				pulumi.String("figaro"),
				pulumi.String("fxblue"),
				pulumi.String("github"),
				pulumi.String("gitlab"),
				pulumi.String("go"),
				pulumi.String("golang"),
				pulumi.String("javascript"),
				pulumi.String("mirror"),
				pulumi.String("pulumi"),
				pulumi.String("vscode"),
			},
			VisibilityLevel: pulumi.String("private"),
		}, pulumi.Protect(false))
		if err != nil {
			panic(err)
		}

		ctx.Export("projectName", project.Name)
		ctx.Export("projectWebUrl", project.WebUrl)

		return nil
	})
}
