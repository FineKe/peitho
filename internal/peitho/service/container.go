// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/marmotedu/errors"

	"github.com/tianrandailove/peitho/internal/peitho/util"
	"github.com/tianrandailove/peitho/pkg/docker"
	"github.com/tianrandailove/peitho/pkg/k8s"
	"github.com/tianrandailove/peitho/pkg/log"
)

type Container struct {
	Env          []string   `json:"Env,omitempty"   yaml:"Env,omitempty"   toml:"Env,omitempty"`
	Cmd          []string   `json:"Cmd"             yaml:"Cmd"             toml:"Cmd"`
	Image        string     `json:"Image,omitempty" yaml:"Image,omitempty" toml:"Image,omitempty"`
	Entrypoint   string     `json:"entrypoint"`
	AttachStdout bool       `json:"AttachStdout"`
	AttachStderr bool       `json:"AttachStderr"`
	HostCfg      HostConfig `json:"HostConfig"`
}

type HostConfig struct {
	NetworkMode string `json:"NetworkMode"`
	Memory      int    `json:"Memory"`
}

type ContainerResult struct {
	Id       string   `json:"Id"`
	Warnings []string `json:"Warnings"`
}

// ContainerSrv defines functions used to handle user request.
type ContainerSrv interface {
	Create(ctx context.Context, containerID string, container Container) (*ContainerResult, error)
	Upload(ctx context.Context, containerID string, path string, content io.Reader) error
	Fetch(ctx context.Context, containerID string, path string) (io.ReadCloser, error)
	Start(ctx context.Context, containerID string) error
	Stop(ctx context.Context, containerID string, timeout time.Duration) error
	Kill(ctx context.Context, containerID string, signal string) error
	Remove(ctx context.Context, containerID string) error
	Wait(ctx context.Context, containerID string) error
}

type containerService struct {
	docker docker.DockerService
	k8s    k8s.K8sService
}

var _ ContainerSrv = (*containerService)(nil)

func newContainer(srv *service) *containerService {
	return &containerService{
		docker: srv.docker,
		k8s:    srv.k8s,
	}
}

// Create create universal container or k8s deployment
func (cs *containerService) Create(ctx context.Context, containerID string, c Container) (*ContainerResult, error) {
	// if containterID == "" , it occurs in building chaincode binary package phase
	if containerID == "" {
		config := &container.Config{
			AttachStdout: c.AttachStdout,
			AttachStderr: c.AttachStderr,
			Env:          c.Env,
			Cmd:          c.Cmd,
			Image:        c.Image,
		}
		hostConfig := &container.HostConfig{
			LogConfig: container.LogConfig{},
		}

		response, err := cs.docker.ContainerCreate(ctx, config, hostConfig, nil, nil, containerID)
		if err != nil {
			log.Errorf("create container failed: %v", err)

			return nil, err
		}

		return &ContainerResult{Id: response.ID, Warnings: response.Warnings}, err
	}

	// deployment image tag
	imageTag := fmt.Sprintf("%s/%s/%s", cs.docker.GetServerAddress(), cs.docker.GetProjectName(), c.Image)

	// ensure registry has the image
	// try pull
	registryAuth, err := cs.docker.RegistryAuth()
	output, err := cs.docker.ImagePull(ctx, imageTag, types.ImagePullOptions{
		RegistryAuth: registryAuth,
	})
	if err != nil {
		return nil, ErrNoSuchImage
	}
	defer output.Close()

	log.Debugf("image %s exists", imageTag)

	// in create chaincode containter phase
	// use k8sapi to create deployment
	podName := util.GetDeploymentName(containerID)
	log.Infof("create chiancode deployment, podname: %s.", podName)

	// create chaincode deployment
	if err := cs.k8s.CreateChaincodeDeployment(ctx, podName, imageTag, c.Env, c.Cmd); err != nil {
		return nil, err
	}

	return &ContainerResult{Id: podName, Warnings: nil}, nil
}

// Upload upload archive, like contract source code
func (cs *containerService) Upload(ctx context.Context, containerID string, path string, content io.Reader) error {
	// it's not chaincode container id
	if util.IsContainerID(containerID) {
		opts := types.CopyToContainerOptions{
			AllowOverwriteDirWithFile: false,
			CopyUIDGID:                false,
		}

		// copy content to the container directly
		if err := cs.docker.CopyToContainer(ctx, containerID, path, content, opts); err != nil {
			log.Errorf("copy archive to container failed: %v", err)

			return err
		}

		return nil
	}

	// id is chaincode id
	// uncompress it
	// and make it to configmap
	if content == nil {
		log.Errorf("content is nil")

		return errors.New("content is nil")
	}

	gzipReader, err := gzip.NewReader(content)
	if err != nil {
		log.Errorf("uncompress archive failed: ", err)

		return err
	}

	tarReader := tar.NewReader(gzipReader)
	// key file
	files := make(map[string]string, 3)

	// traverse file
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			// end of tar archive
			break
		}
		log.Debugf("UnTarGzing file..." + header.Name)

		// check if it is diretory or file
		if header.Typeflag != tar.TypeDir {
			bf := bytes.NewBuffer(make([]byte, 0, 1024))

			_, err := io.Copy(bf, tarReader)
			if err != nil {
				log.Errorf("copy content failed: %v", err)

				return err
			}

			strs := strings.Split(header.Name, "/")
			fileName := strs[len(strs)-1]
			files[fileName] = bf.String()

			log.Debugf("file name:%s\n file content:%s\n", fileName, files[fileName])
		}
	}

	// create tls configmap
	name := util.GetDeploymentName(containerID)
	if err := cs.k8s.CreateConfigMap(ctx, name, files); err != nil {
		return err
	}

	if len(files) > 3 {
		ctx = context.WithValue(ctx, "version", "v2.0.0")
	}
	// update chaincode deployment
	if err := cs.k8s.UpdateDeployment(context.WithValue(ctx, "version", "v2.0.0"), name); err != nil {
		return err
	}

	return nil
}

// Fetch fetch contract bin
func (cs *containerService) Fetch(ctx context.Context, containerID string, path string) (io.ReadCloser, error) {
	reader, _, err := cs.docker.CopyFromContainer(ctx, containerID, path)
	if err != nil {
		log.Errorf("fetch archive failed: %v", err)

		return nil, err
	}

	return reader, nil
}

// Start start a universal container or waitting for deployment be ok
func (cs *containerService) Start(ctx context.Context, containerID string) error {
	if util.IsContainerID(containerID) {
		err := cs.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{
			CheckpointID:  "",
			CheckpointDir: "",
		})
		if err != nil {
			log.Errorf("start chaincode failed: %v", err)

			return err
		}

		return nil
	}

	log.Info("start check chaincode deployment status....")

	podName := util.GetDeploymentName(containerID)
	// check 100 time
	for i := 0; i < 100; i++ {
		ok, _ := cs.k8s.QueryDeploymentStatus(ctx, podName)
		if ok {
			log.Info("check chaincode deployment ok")

			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return errors.New("check chaincode deployment status timeout")
}

// Stop stop a universal container
func (cs *containerService) Stop(ctx context.Context, containerID string, timeout time.Duration) error {
	if util.IsContainerID(containerID) {
		err := cs.docker.ContainerStop(ctx, containerID, &timeout)
		if err != nil {
			log.Errorf("停止容器失败: %v,", err)

			return err
		}

		return nil
	}

	return nil
}

// Kill kill a universal container
func (cs *containerService) Kill(ctx context.Context, containerID string, signal string) error {
	if util.IsContainerID(containerID) {
		err := cs.docker.ContainerKill(ctx, containerID, signal)
		if err != nil {
			log.Errorf("kill: %v", err)

			return err
		}

		return nil
	}

	return nil
}

// Remove delete universal container and chaincode deployment
func (cs *containerService) Remove(ctx context.Context, containerID string) error {
	if util.IsContainerID(containerID) {
		opts := types.ContainerRemoveOptions{}

		err := cs.docker.ContainerRemove(ctx, containerID, opts)
		if err != nil {
			log.Errorf("remove contaienr failed: %v", err)

			return err
		}

		return nil
	}

	// delete configmap and deployment
	// ignore error
	name := util.GetDeploymentName(containerID)
	_ = cs.k8s.DeleteChaincodeDeployment(ctx, name)
	_ = cs.k8s.DeleteConfigMapDeployment(ctx, name)

	return nil
}

// Wait wait for universal container
func (cs *containerService) Wait(ctx context.Context, containerID string) error {
	if util.IsContainerID(containerID) {
		okc, errc := cs.docker.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
		select {
		case <-okc:
			return nil
		case e := <-errc:
			return e
		}
	}

	return nil
}
