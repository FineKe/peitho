// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/marmotedu/errors"

	"github.com/tianrandailove/peitho/pkg/k8s"

	"github.com/tianrandailove/peitho/pkg/docker"
)

func Test_newImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	k8sSrv := k8s.NewMockK8sService(ctrl)
	srv := &service{
		docker: dockerSrv,
		k8s:    k8sSrv,
	}
	type args struct {
		srv *service
	}
	tests := []struct {
		name string
		args args
		want *imageService
	}{
		{
			name: "create image service",
			args: args{srv: srv},
			want: newImage(srv),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newImage(tt.args.srv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageService_Build(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	ctx := context.Background()
	dockerSrv.EXPECT().
		ImageBuild(ctx, nil, gomock.Any()).
		Return(types.ImageBuildResponse{}, errors.New("content is nil")).
		AnyTimes()
	type fields struct {
		docker docker.DockerService
	}
	type args struct {
		ctx        context.Context
		dockerfile string
		tags       []string
		content    io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		{
			name:   "build chiancode image",
			fields: fields{dockerSrv},
			args: args{
				ctx:        ctx,
				dockerfile: "",
				tags:       []string{"mycc:latest"},
				content:    nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := imageService{
				docker: tt.fields.docker,
			}
			got, err := i.Build(tt.args.ctx, tt.args.dockerfile, tt.args.tags, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageService.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageService.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	ctx := context.Background()
	dockerSrv.EXPECT().ImagePull(ctx, "hyperledger/fabric-ccenv:latest", gomock.Any()).Return(nil, nil)
	type fields struct {
		docker docker.DockerService
	}
	type args struct {
		ctx       context.Context
		fromImage string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		{
			name:   "pull image",
			fields: fields{dockerSrv},
			args: args{
				ctx:       ctx,
				fromImage: "hyperledger/fabric-ccenv:latest",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := imageService{
				docker: tt.fields.docker,
			}
			got, err := i.Create(tt.args.ctx, tt.args.fromImage)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageService_Inspect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	ctx := context.Background()
	dockerSrv.EXPECT().
		ImageInspectWithRaw(ctx, "hyperledger/fabric-ccenv:latest").
		Return(types.ImageInspect{ID: "1"}, nil, errors.New("no such image"))
	type fields struct {
		docker docker.DockerService
	}
	type args struct {
		ctx     context.Context
		imageID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "inspect image",
			fields: fields{dockerSrv},
			args: args{
				ctx:     ctx,
				imageID: "hyperledger/fabric-ccenv:latest",
			},
			want:    types.ImageInspect{ID: "1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := imageService{
				docker: tt.fields.docker,
			}
			_, err := i.Inspect(tt.args.ctx, tt.args.imageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageService.Inspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_imageService_AddTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	ctx := context.Background()
	dockerSrv.EXPECT().ImageTag(ctx, "dev-peer0-org1-mycc", "172.199.189.12:8099/chaincode/dev-peer0-org1-mycc")

	type fields struct {
		docker docker.DockerService
	}
	type args struct {
		ctx      context.Context
		imageTag string
		newTag   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "add tag",
			fields: fields{dockerSrv},
			args: args{
				ctx:      ctx,
				imageTag: "dev-peer0-org1-mycc",
				newTag:   "172.199.189.12:8099/chaincode/dev-peer0-org1-mycc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := imageService{
				docker: tt.fields.docker,
			}
			if err := i.AddTag(tt.args.ctx, tt.args.imageTag, tt.args.newTag); (err != nil) != tt.wantErr {
				t.Errorf("imageService.AddTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_imageService_Push(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerSrv := docker.NewMockDockerService(ctrl)
	ctx := context.Background()
	dockerSrv.EXPECT().RegistryAuth().Return("base64 auth", nil)
	dockerSrv.EXPECT().ImagePush(ctx, "172.198.161.22:8099/chaincode/mycc:latest", gomock.Any()).Return(nil, nil)

	type fields struct {
		docker docker.DockerService
	}
	type args struct {
		ctx      context.Context
		imageTag string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		{
			name:   "push image",
			fields: fields{dockerSrv},
			args: args{
				ctx:      ctx,
				imageTag: "172.198.161.22:8099/chaincode/mycc:latest",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := imageService{
				docker: tt.fields.docker,
			}
			got, err := i.Push(tt.args.ctx, tt.args.imageTag)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageService.Push() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageService.Push() = %v, want %v", got, tt.want)
			}
		})
	}
}
