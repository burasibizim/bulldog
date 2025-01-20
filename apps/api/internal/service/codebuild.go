package service

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	"github.com/aws/aws-sdk-go-v2/service/codebuild/types"
)

type CodebuildService struct {
	client *codebuild.Client
}

func NewCodebuildService(client *codebuild.Client) *CodebuildService {
	return &CodebuildService{client: client}
}

func (s *CodebuildService) CreateCodebuild(organization_id string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load configuration, %v", err)
	}

	client := codebuild.NewFromConfig(cfg)

	createInput := &codebuild.CreateProjectInput{
		Name: aws.String(organization_id),
		Environment: &types.ProjectEnvironment{
			Type:           types.EnvironmentTypeLinuxContainer,
			ComputeType:    types.ComputeTypeBuildGeneral1Small,
			Image:          aws.String("aws/codebuild/standard:6.0"),
			PrivilegedMode: aws.Bool(false),
			EnvironmentVariables: []types.EnvironmentVariable{
				{
					Name:  aws.String("ORGANIZATION_ID"),
					Value: aws.String(organization_id),
				},
			},
		},
		ServiceRole: aws.String(os.Getenv("CODEBUILD_SERVICE_ROLE")),
		Source: &types.ProjectSource{
			Type:     types.SourceTypeGithub,
			Location: aws.String(os.Getenv("CODEBUILD_SOURCE_LOCATION")),
		},
		Artifacts: &types.ProjectArtifacts{
			Type: types.ArtifactsTypeNoArtifacts,
		},
	}

	output, err := client.CreateProject(context.TODO(), createInput)
	if err != nil {
		fmt.Printf("Error details: %+v\n", err)
		return fmt.Errorf("failed to create project, %v", err)
	}

	fmt.Printf("Created project: %s\n", *output.Project.Name)

	return nil
}

func (s *CodebuildService) StartCodebuild(projectName string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load configuration, %v", err)
	}

	client := codebuild.NewFromConfig(cfg)

	startBuildInput := &codebuild.StartBuildInput{
		ProjectName: aws.String(projectName),
	}

	output, err := client.StartBuild(context.TODO(), startBuildInput)
	if err != nil {
		fmt.Printf("Error details: %+v\n", err)
		return fmt.Errorf("failed to start build, %v", err)
	}

	fmt.Printf("Started build: %s\n", *output.Build.Id)

	return nil
}
