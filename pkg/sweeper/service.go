package sweeper

import (
	"context"
	"github.com/tianrandailove/peitho/pkg/k8s"
	"github.com/tianrandailove/peitho/pkg/log"
	"github.com/tianrandailove/peitho/pkg/options"
	"time"
)

type SweepService interface {
	Start()
	Stop()
}

type PeithoSweeper struct {
	enable   bool
	interval int
	k8s      k8s.K8sService
	ch       chan struct{}
}

func NewPeithoSweeper(k8s k8s.K8sService, option *options.SweeperOption) (*PeithoSweeper, error) {
	return &PeithoSweeper{
		enable:   option.Enable,
		interval: option.Interval,
		k8s:      k8s,
		ch:       make(chan struct{}),
	}, nil
}

func (ps *PeithoSweeper) Start() {
	if !ps.enable {
		log.Info("needen't to start sweeper")

		return
	}
	log.Info("starting sweeper")

	ticker := time.NewTicker(time.Duration(ps.interval) * time.Second)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			ds, err := ps.k8s.ListDeploymentByPrefix(ctx, k8s.ChaincodePrefix)
			if err != nil {
				log.Errorf("list chaincode deployment failed: %v", err)

				continue
			}
			for _, deployment := range ds {
				// chaincode unavailable
				// to delete it
				if deployment.Status.UnavailableReplicas > 0 {
					log.Infof("to delete %s", deployment.Name)
					err := ps.k8s.DeleteChaincodeDeployment(ctx, deployment.Name)
					if err != nil {
						log.Errorf("failed to delete %s, cause by: %v", deployment.Name, err)
					} else {
						log.Infof("delete %s success", deployment.Name)
					}
					// delete configmap if exists
					// ignore err
					_ = ps.k8s.DeleteConfigMapDeployment(ctx, deployment.Name)
				}

				time.Sleep(time.Second)
			}
		case <-ps.ch:
			break
		}
	}
}

func (ps *PeithoSweeper) Stop() {
	ps.ch <- struct{}{}
}
