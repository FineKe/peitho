// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package docker

//go:generate mockgen -self_package=github.com/tianrandailove/peitho/pkg/docker -destination mock_service.go -package docker github.com/tianrandailove/peitho/pkg/docker DockerService
import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/tianrandailove/peitho/pkg/log"
	"github.com/tianrandailove/peitho/pkg/options"
)

// Docker.
type Docker struct {
	DockerClient *client.Client
	Registry     *Registry
}

type DockerService interface {
	RegistryAuth() (string, error)
	GetServerAddress() string
	GetProjectName() string

	ContainerAttach(
		ctx context.Context,
		container string,
		options types.ContainerAttachOptions,
	) (types.HijackedResponse, error)
	ContainerCreate(
		ctx context.Context,
		config *containertypes.Config,
		hostConfig *containertypes.HostConfig,
		networkingConfig *networktypes.NetworkingConfig,
		platform *specs.Platform,
		containerName string,
	) (containertypes.ContainerCreateCreatedBody, error)
	ContainerKill(ctx context.Context, container, signal string) error
	ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error
	ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error
	ContainerStop(ctx context.Context, container string, timeout *time.Duration) error
	ContainerWait(
		ctx context.Context,
		container string,
		condition containertypes.WaitCondition,
	) (<-chan containertypes.ContainerWaitOKBody, <-chan error)
	CopyFromContainer(ctx context.Context, container, srcPath string) (io.ReadCloser, types.ContainerPathStat, error)
	CopyToContainer(
		ctx context.Context,
		container, path string,
		content io.Reader,
		options types.CopyToContainerOptions,
	) error

	ImageBuild(
		ctx context.Context,
		context io.Reader,
		options types.ImageBuildOptions,
	) (types.ImageBuildResponse, error)
	ImageInspectWithRaw(ctx context.Context, image string) (types.ImageInspect, []byte, error)
	ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error)
	ImageTag(ctx context.Context, image, ref string) error
}

// NewDocker new docker client from opt.
func newDocker(opt *options.DockerOption) (*Docker, error) {
	docker, err := client.NewClientWithOpts(client.WithHost(opt.Endpoint), client.WithVersion("1.38"))
	if err != nil {
		log.Errorf("new docker client failed: %v", err)

		return nil, err
	}

	go func() {
		for {
			_, err := docker.Ping(context.Background())
			if err != nil {
				log.Errorf("ping docker server failed: %v", err)
			}
			log.Info("docker server alive")

			time.Sleep(5 * time.Second)
		}
	}()

	return &Docker{
		DockerClient: docker,
		Registry: &Registry{
			Username:      opt.Registry.Username,
			Password:      opt.Registry.Password,
			Email:         opt.Registry.Email,
			Serveraddress: opt.Registry.Serveraddress,
			Project:       opt.Registry.Project,
		},
	}, nil
}

// NewDockerService new docker service instance.
func NewDockerService(opt *options.DockerOption) (*Docker, error) {
	return newDocker(opt)
}

type Registry struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	Serveraddress string `json:"serveraddress"`
	Project       string `json:"project"`
}

// RegistryAuth registry Authentication encoded by base64.
func (d *Docker) RegistryAuth() (string, error) {
	s := struct {
		Username      string `json:"username"`
		Password      string `json:"password"`
		Email         string `json:"email"`
		Serveraddress string `json:"serveraddress"`
	}{
		Username:      d.Registry.Username,
		Password:      d.Registry.Password,
		Email:         "",
		Serveraddress: d.Registry.Serveraddress,
	}

	jsonBytes, err := json.Marshal(s)
	if err != nil {
		log.Errorf("json marshal failed: %v", err)

		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

func (d *Docker) GetServerAddress() string {
	return d.Registry.Serveraddress
}

func (d *Docker) GetProjectName() string {
	return d.Registry.Project
}

func (d *Docker) ContainerAttach(
	ctx context.Context,
	container string,
	options types.ContainerAttachOptions,
) (types.HijackedResponse, error) {
	return d.DockerClient.ContainerAttach(ctx, container, options)
}

func (d *Docker) ContainerCreate(
	ctx context.Context,
	config *containertypes.Config,
	hostConfig *containertypes.HostConfig,
	networkingConfig *networktypes.NetworkingConfig,
	platform *specs.Platform,
	containerName string,
) (containertypes.ContainerCreateCreatedBody, error) {
	return d.DockerClient.ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName)
}

func (d *Docker) ContainerKill(ctx context.Context, container, signal string) error {
	return d.DockerClient.ContainerKill(ctx, container, signal)
}

func (d *Docker) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	return d.DockerClient.ContainerRemove(ctx, container, options)
}

func (d *Docker) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	return d.DockerClient.ContainerStart(ctx, container, options)
}

func (d *Docker) ContainerStop(ctx context.Context, container string, timeout *time.Duration) error {
	return d.DockerClient.ContainerStop(ctx, container, timeout)
}

func (d *Docker) ContainerWait(
	ctx context.Context,
	container string,
	condition containertypes.WaitCondition,
) (<-chan containertypes.ContainerWaitOKBody, <-chan error) {
	return d.DockerClient.ContainerWait(ctx, container, condition)
}

func (d *Docker) CopyFromContainer(
	ctx context.Context,
	container, srcPath string,
) (io.ReadCloser, types.ContainerPathStat, error) {
	return d.DockerClient.CopyFromContainer(ctx, container, srcPath)
}

func (d *Docker) CopyToContainer(
	ctx context.Context,
	container, path string,
	content io.Reader,
	options types.CopyToContainerOptions,
) error {
	return d.DockerClient.CopyToContainer(ctx, container, path, content, options)
}

func (d *Docker) ImageBuild(
	ctx context.Context,
	context io.Reader,
	options types.ImageBuildOptions,
) (types.ImageBuildResponse, error) {
	return d.DockerClient.ImageBuild(ctx, context, options)
}

func (d *Docker) ImageInspectWithRaw(ctx context.Context, image string) (types.ImageInspect, []byte, error) {
	return d.DockerClient.ImageInspectWithRaw(ctx, image)
}

func (d *Docker) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return d.DockerClient.ImagePull(ctx, ref, options)
}

func (d *Docker) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	return d.DockerClient.ImagePush(ctx, ref, options)
}

func (d *Docker) ImageTag(ctx context.Context, image, ref string) error {
	return d.DockerClient.ImageTag(ctx, image, ref)
}
