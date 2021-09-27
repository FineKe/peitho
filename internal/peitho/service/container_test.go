// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"

	"github.com/tianrandailove/peitho/pkg/docker"
	"github.com/tianrandailove/peitho/pkg/k8s"
)

func TestContainerCreate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	k8sSrv := k8s.NewMockK8sService(ctrl)

	containerSrv := newContainer(&service{
		docker: dockerSrv,
		k8s:    k8sSrv,
	})

	ctx := context.Background()

	con := Container{Image: "hyperledger/fabric-ccenv-amd64:1.4.8"}

	config := &container.Config{
		Image: "hyperledger/fabric-ccenv-amd64:1.4.8",
	}

	hostConfig := &container.HostConfig{
		LogConfig: container.LogConfig{},
	}

	dockerSrv.EXPECT().ContainerCreate(ctx, config, hostConfig, nil, nil, "").Return(struct {
		ID       string   `json:"Id"`
		Warnings []string `json:"Warnings"`
	}{
		ID: "container",
	}, nil)

	t.Run("create universal container", func(t *testing.T) {
		result, err := containerSrv.Create(ctx, "", con)

		if err != nil {
			t.Errorf("create universal container failed:%v", err)
		}

		want := "container"
		if result.Id != want {
			t.Errorf("want containerID: %s got: %s", want, result.Id)
		}
	})

	ctx = context.Background()
	podName := "dev.peer0.org1"
	dockerSrv.EXPECT().GetServerAddress().Return("172.198.101.18:8099")
	dockerSrv.EXPECT().GetProjectName().Return("chaincode")

	dockerSrv.EXPECT().
		ImageInspectWithRaw(ctx, gomock.Eq("hyperledger/fabric-ccenv-amd64:1.4.8")).
		Return(types.ImageInspect{ID: "hyperledger/fabric-ccenv-amd64:1.4.8"}, nil, nil)
	dockerSrv.EXPECT().
		ImageInspectWithRaw(ctx, gomock.Eq("172.198.101.18:8099/chaincode/hyperledger/fabric-ccenv-amd64:1.4.8")).
		Return(types.ImageInspect{ID: "172.198.101.18:8099/chaincode/hyperledger/fabric-ccenv-amd64:1.4.8"}, nil, nil)

	k8sSrv.EXPECT().
		CreateChaincodeDeployment(ctx, "dev-peer0-org1", "172.198.101.18:8099/chaincode/hyperledger/fabric-ccenv-amd64:1.4.8", con.Env, con.Cmd).
		Return(nil)

	t.Run("create chaincode deployment", func(t *testing.T) {
		result, err := containerSrv.Create(ctx, podName, con)

		if err != nil {
			t.Errorf("create chaincode deployment failed:%v", err)
		}
		want := "dev-peer0-org1"
		if result.Id != want {
			t.Errorf("want pod name: %s, got: %s", want, result.Id)
		}
	})
}

func Test_newContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := &service{
		docker: docker.NewMockDockerService(ctrl),
		k8s:    k8s.NewMockK8sService(ctrl),
	}

	type args struct {
		srv *service
	}
	tests := []struct {
		name string
		args args
		want *containerService
	}{
		{
			name: "new container",
			args: args{srv: srv},
			want: newContainer(srv),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newContainer(tt.args.srv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containerService_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	k8sSrv := k8s.NewMockK8sService(ctrl)

	ctx := context.Background()
	dockerSrv.EXPECT().
		CopyToContainer(ctx, "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2", "chaincode", nil, gomock.Any()).
		Return(nil)

	type fields struct {
		docker docker.DockerService
		k8s    k8s.K8sService
	}
	type args struct {
		ctx         context.Context
		containerID string
		path        string
		content     io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "upload chiancode tar to container",
			fields: fields{
				docker: dockerSrv,
				k8s:    k8sSrv,
			},
			args: args{
				ctx:         ctx,
				containerID: "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2",
				path:        "chaincode",
				content:     nil,
			},
			wantErr: false,
		},

		{
			name: "upload nil tls tar to chaincode",
			fields: fields{
				docker: dockerSrv,
				k8s:    k8sSrv,
			},
			args: args{
				ctx:         ctx,
				containerID: "dev-peer0-org1-mycc-v-1-0",
				path:        "",
				content:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &containerService{
				docker: tt.fields.docker,
				k8s:    tt.fields.k8s,
			}
			if err := cs.Upload(tt.args.ctx, tt.args.containerID, tt.args.path, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("containerService.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	k8sSrv := k8s.NewMockK8sService(ctrl)

	ctx := context.Background()

	dockerSrv.EXPECT().
		CopyFromContainer(ctx, "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2", "/chaincode")

	type fields struct {
		docker docker.DockerService
		k8s    k8s.K8sService
	}
	type args struct {
		ctx         context.Context
		containerID string
		path        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		{
			name: "fetch chiancode bin",
			fields: fields{
				docker: dockerSrv,
				k8s:    k8sSrv,
			},
			args: args{
				ctx:         ctx,
				containerID: "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2",
				path:        "/chaincode",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &containerService{
				docker: tt.fields.docker,
				k8s:    tt.fields.k8s,
			}
			got, err := cs.Fetch(tt.args.ctx, tt.args.containerID, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("containerService.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("containerService.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containerService_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	k8sSrv := k8s.NewMockK8sService(ctrl)
	ctx := context.Background()

	dockerSrv.EXPECT().
		ContainerStart(ctx, "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2", types.ContainerStartOptions{}).
		Return(nil).
		AnyTimes()
	k8sSrv.EXPECT().QueryDeploymentStatus(ctx, "dev-peer0-org1-mycc").Return(true, nil)
	type fields struct {
		docker docker.DockerService
		k8s    k8s.K8sService
	}
	type args struct {
		ctx         context.Context
		containerID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "start build chaincode container",
			fields: fields{
				docker: dockerSrv,
				k8s:    k8sSrv,
			},
			args: args{
				ctx:         ctx,
				containerID: "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2",
			},
			wantErr: false,
		},
		{
			name: "block check deployment status",
			fields: fields{
				docker: dockerSrv,
				k8s:    k8sSrv,
			},
			args: args{
				ctx:         ctx,
				containerID: "dev-peer0-org1-mycc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &containerService{
				docker: tt.fields.docker,
				k8s:    tt.fields.k8s,
			}
			if err := cs.Start(tt.args.ctx, tt.args.containerID); (err != nil) != tt.wantErr {
				t.Errorf("containerService.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	ctx := context.Background()
	dockerSrv.EXPECT().
		ContainerStop(ctx, "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2", gomock.Any()).
		Return(nil)

	type fields struct {
		docker docker.DockerService
		k8s    k8s.K8sService
	}
	type args struct {
		ctx         context.Context
		containerID string
		timeout     time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "stop universal container",
			fields: fields{
				docker: dockerSrv,
				k8s:    nil,
			},
			args: args{
				ctx:         ctx,
				containerID: "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2",
				timeout:     0,
			},
			wantErr: false,
		},
		{
			name: "stop chiancode deployment",
			fields: fields{
				docker: dockerSrv,
				k8s:    nil,
			},
			args: args{
				ctx:         ctx,
				containerID: "dev-peer0-org1-mycc",
				timeout:     0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &containerService{
				docker: tt.fields.docker,
				k8s:    tt.fields.k8s,
			}
			if err := cs.Stop(tt.args.ctx, tt.args.containerID, tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("containerService.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_Kill(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	ctx := context.Background()
	dockerSrv.EXPECT().ContainerKill(ctx, "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2", "kill")
	type fields struct {
		docker docker.DockerService
		k8s    k8s.K8sService
	}
	type args struct {
		ctx         context.Context
		containerID string
		signal      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "stop universal container",
			fields: fields{
				docker: dockerSrv,
				k8s:    nil,
			},
			args: args{
				ctx:         ctx,
				containerID: "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2",
				signal:      "kill",
			},
			wantErr: false,
		},
		{
			name: "stop chiancode deployment",
			fields: fields{
				docker: dockerSrv,
				k8s:    nil,
			},
			args: args{
				ctx:         ctx,
				containerID: "dev-peer0-org1-mycc",
				signal:      "kill",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &containerService{
				docker: tt.fields.docker,
				k8s:    tt.fields.k8s,
			}
			if err := cs.Kill(tt.args.ctx, tt.args.containerID, tt.args.signal); (err != nil) != tt.wantErr {
				t.Errorf("containerService.Kill() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	k8sSrv := k8s.NewMockK8sService(ctrl)
	ctx := context.Background()

	dockerSrv.EXPECT().
		ContainerRemove(ctx, "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2", gomock.Any()).
		Return(nil)
	k8sSrv.EXPECT().DeleteChaincodeDeployment(ctx, "dev-peer0-org1-mycc").Return(nil)
	k8sSrv.EXPECT().DeleteConfigMapDeployment(ctx, "dev-peer0-org1-mycc").Return(nil)

	type fields struct {
		docker docker.DockerService
		k8s    k8s.K8sService
	}
	type args struct {
		ctx         context.Context
		containerID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "remove universal container",
			fields: fields{
				docker: dockerSrv,
				k8s:    k8sSrv,
			},
			args: args{
				ctx:         ctx,
				containerID: "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2",
			},
			wantErr: false,
		},
		{
			name: "remove chiancode deployment",
			fields: fields{
				docker: dockerSrv,
				k8s:    k8sSrv,
			},
			args: args{
				ctx:         ctx,
				containerID: "dev-peer0-org1-mycc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &containerService{
				docker: tt.fields.docker,
				k8s:    tt.fields.k8s,
			}
			if err := cs.Remove(tt.args.ctx, tt.args.containerID); (err != nil) != tt.wantErr {
				t.Errorf("containerService.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_Wait(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	ctx := context.Background()
	ch := make(chan container.ContainerWaitOKBody, 1)
	ch <- container.ContainerWaitOKBody{}
	dockerSrv.EXPECT().
		ContainerWait(ctx, "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2", gomock.Any()).
		Return(ch, nil)
	type fields struct {
		docker docker.DockerService
		k8s    k8s.K8sService
	}
	type args struct {
		ctx         context.Context
		containerID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "wait universal container",
			fields: fields{
				docker: dockerSrv,
				k8s:    nil,
			},
			args: args{
				ctx:         ctx,
				containerID: "ceced51c5497820e2e7fe4a12eb413c2cf8c8675553b5d84c56d9bb161e078f2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &containerService{
				docker: tt.fields.docker,
				k8s:    tt.fields.k8s,
			}
			if err := cs.Wait(tt.args.ctx, tt.args.containerID); (err != nil) != tt.wantErr {
				t.Errorf("containerService.Wait() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
