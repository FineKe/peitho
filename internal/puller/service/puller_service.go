// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/tianrandailove/peitho/pkg/docker"
	"github.com/tianrandailove/peitho/pkg/log"
)

type PullerSrv interface {
	PullImage(image string, pullAddress string) error
}

type pullerService struct {
	docker docker.DockerService
}

var _ PullerSrv = (*pullerService)(nil)

func newPuller(srv *service) *pullerService {
	return &pullerService{
		docker: srv.docker,
	}
}

func (p *pullerService) PullImage(image string, pullAddress string) error {
	ctx := context.Background()
	// check exists
	_, _, err := p.docker.ImageInspectWithRaw(ctx, image)
	if err == nil {
		log.Infof("%s already exists the host", image)

		return nil
	}
	log.Infof("%s not exists the host", image)

	fileName := fmt.Sprintf("%s.tar", image)
	// if the file exists, then delete it
	_, err = os.Stat(fileName)
	if err == nil {
		log.Debugf("%s already exists, then delete it", fileName)
		err = os.Remove(fileName)
		if err != nil {
			log.Errorf("delete %s failed, cause by: %v", fileName, err)

			return err
		}
	}
	// download image.tar from pullAddress
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", pullAddress, image), nil)
	if err != nil {
		log.Infof("create http request failed: %v", err)

		return err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Errorf("download %s.tar failed, cause by:%v", image, err)

		return err
	}
	defer resp.Body.Close()

	//// create file
	//file, err := os.Create(fileName)
	//if err != nil {
	//	log.Errorf("create %s failed, cause by:%v", fileName, err)
	//
	//	return err
	//}
	//defer file.Close()
	//fileSize, err := io.Copy(file, resp.Body)
	//if err != nil {
	//	log.Errorf("copy file failed: %v", err)
	//
	//	return err
	//}

	// log.Infof("write file completed, file size: %d byte", fileSize)
	// load image.tar
	loadResp, err := p.docker.ImageLoad(ctx, resp.Body, false)
	if err != nil {
		log.Errorf("load %s failed: %v", fileName, err)

		return err
	}
	defer loadResp.Body.Close()
	reader := bufio.NewReader(loadResp.Body)
	for {
		line, _, err1 := reader.ReadLine()
		if err1 != nil {
			break
		}
		log.Infof("%s", line)
	}

	return err
}
