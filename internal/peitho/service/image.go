// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/marmotedu/errors"

	"github.com/tianrandailove/peitho/pkg/docker"
	"github.com/tianrandailove/peitho/pkg/log"
)

type DockerAuthentication struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	Serveraddress string `json:"serveraddress"`
}

type ImageSrv interface {
	Build(ctx context.Context, dockerfile string, tags []string, content io.Reader) (io.ReadCloser, error)
	Create(ctx context.Context, fromImage string) (io.ReadCloser, error)
	Inspect(ctx context.Context, imageID string) (interface{}, error)
	AddTag(ctx context.Context, image, newTag string) error
	Push(ctx context.Context, imageTag string) (io.ReadCloser, error)
}

type imageService struct {
	docker docker.DockerService
}

var _ ImageSrv = (*imageService)(nil)

func newImage(srv *service) *imageService {
	return &imageService{
		docker: srv.docker,
	}
}

func (i imageService) Build(
	ctx context.Context,
	dockerfile string,
	tags []string,
	content io.Reader,
) (io.ReadCloser, error) {
	imageOptions := types.ImageBuildOptions{
		Dockerfile: dockerfile,
		Tags:       tags,
		Remove:     true,
		Version:    "1",
	}

	if content == nil {
		log.Errorf("content is nil")

		return nil, errors.New("content is nil")
	}

	resp, err := i.docker.ImageBuild(ctx, content, imageOptions)
	if err != nil {
		log.Errorf("build image failed: %v", err)

		return nil, err
	}

	reader := bufio.NewReader(resp.Body)

	for {
		line, _, readErr := reader.ReadLine()
		if readErr != nil {
			break
		}
		log.Debugf(string(line))
	}

	log.Infof("build image %s success", tags[0])

	// async to push image
	go func() {
		oldTag := tags[0]
		newTag := fmt.Sprintf("%s/%s/%s", i.docker.GetServerAddress(), i.docker.GetProjectName(), tags[0])

		log.Debugf("oldTag:%s", oldTag)
		log.Debugf("newTag:%s", newTag)

		if tagErr := i.docker.ImageTag(ctx, oldTag, newTag); tagErr != nil {
			log.Errorf("add new tag failed: %v", tagErr)

			return
		}

		pushOpt := types.ImagePushOptions{}
		auth, authErr := i.docker.RegistryAuth()
		if err != nil {
			log.Errorf("get registryAuth failed: %v", authErr)
		} else {
			pushOpt.RegistryAuth = auth
		}

		log.Debugf("RegistryAuth: %s", pushOpt.RegistryAuth)

		readerCloser, pushErr := i.docker.ImagePush(ctx, newTag, pushOpt)
		if pushErr != nil {
			log.Errorf("push failed: %v", pushErr)

			return
		}

		reader := bufio.NewReader(readerCloser)
		for {
			line, _, readErr := reader.ReadLine()
			if readErr != nil {
				break
			}
			log.Debugf(string(line))
		}

		log.Infof("%s push success", newTag)
	}()

	return resp.Body, err
}

func (i imageService) Create(ctx context.Context, fromImage string) (io.ReadCloser, error) {
	resp, err := i.docker.ImagePull(ctx, fromImage, types.ImagePullOptions{})
	if err != nil {
		log.Errorf("pull image failed: %v", err)

		return nil, err
	}

	return resp, nil
}

func (i imageService) Inspect(ctx context.Context, imageID string) (interface{}, error) {
	imageInspect, _, err := i.docker.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		log.Errorf("inspect image failed %v", err)

		return nil, err
	}

	return imageInspect, nil
}

func (i imageService) AddTag(ctx context.Context, imageTag, newTag string) error {
	if err := i.docker.ImageTag(ctx, imageTag, newTag); err != nil {
		log.Errorf("add tag failed: %v", err)

		return err
	}

	return nil
}

func (i imageService) Push(ctx context.Context, imageTag string) (io.ReadCloser, error) {
	auth, err := i.docker.RegistryAuth()
	if err != nil {
		log.Errorf("get RegistryAuth failed: %v", err)
	}

	opts := types.ImagePushOptions{
		RegistryAuth: auth,
	}

	log.Debugf("RegistryAuth: %s", opts.RegistryAuth)

	response, err := i.docker.ImagePush(context.Background(), imageTag, opts)
	if err != nil {
		log.Errorf("pus image failed: %v", err)

		return nil, err
	}

	return response, nil
}
