// Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package k8s

//go:generate mockgen -self_package=github.com/tianrandailove/peitho/pkg/k8s -destination mock_service.go -package k8s github.com/tianrandailove/peitho/pkg/k8s K8sService
import (
	"context"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/tianrandailove/peitho/pkg/log"
	"github.com/tianrandailove/peitho/pkg/options"
)

const (
	// Mutual TLS auth client key and cert paths in the chaincode container.
	TLSClientKeyPath      string = "/etc/hyperledger/fabric/client.key"
	TLSClientCertPath     string = "/etc/hyperledger/fabric/client.crt"
	TLSClientRootCertPath string = "/etc/hyperledger/fabric/peer.crt"

	TLSClientKeyFile      string = "/etc/hyperledger/fabric/client_pem.key"
	TLSClientCertFile     string = "/etc/hyperledger/fabric/client_pem.crt"
	TLSClientRootCertFile string = "/etc/hyperledger/fabric/peer.crt"
)

type K8sClient struct {
	k8sClientSet *kubernetes.Clientset
	namespace    string
	dns          []string
}

type K8sService interface {
	CreateChaincodeDeployment(ctx context.Context, name string, image string, env []string, cmd []string) error
	UpdateDeployment(ctx context.Context, name string) error
	CreateConfigMap(ctx context.Context, name string, data map[string]string) error
	DeleteChaincodeDeployment(ctx context.Context, name string) error
	DeleteConfigMapDeployment(ctx context.Context, name string) error
	QueryDeploymentStatus(ctx context.Context, name string) (bool, error)
}

// NewK8sClient new k8sclient from opt.
func newK8sClient(opt *options.K8sOption) (*K8sClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", opt.Kubeconfig)
	if err != nil {
		log.Errorf("load kubeconfig failed: %v", err)

		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("init kubernates client failed: %v", err.Error())

		return nil, err
	}

	log.Infof("init kubernates client suceess ...")

	return &K8sClient{
		k8sClientSet: client,
		namespace:    opt.Namespace,
		dns:          opt.DNS,
	}, nil
}

// NewK8sService.
func NewK8sService(opt *options.K8sOption) (K8sService, error) {
	return newK8sClient(opt)
}

func (k8s *K8sClient) CreateChaincodeDeployment(
	ctx context.Context,
	name string,
	image string,
	env []string,
	cmd []string,
) error {
	// replicas
	replicas := int32(0)

	// build environment variables
	envs := make([]v1.EnvVar, 0, len(env))
	if len(env) > 0 {
		for _, s := range env {
			array := strings.Split(s, "=")
			envs = append(envs, v1.EnvVar{
				Name:  array[0],
				Value: array[1],
			})

			if array[0] == "CORE_PEER_TLS_ENABLED" && array[1] == "false" {
				replicas = 1
			}
		}
	}

	// dns
	var hostAlias []v1.HostAlias
	if len(k8s.dns) > 0 {
		for _, s := range k8s.dns {
			strs := strings.Split(s, ":")
			hostAlias = append(hostAlias, v1.HostAlias{
				IP:        strs[0],
				Hostnames: []string{strs[1]},
			})
		}
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   name,
					Labels: map[string]string{"app": name},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            name,
							Image:           image,
							ImagePullPolicy: v1.PullIfNotPresent,
							Env:             envs,
							Command:         cmd,
						},
					},
					HostAliases: hostAlias,
				},
			},
		},
	}

	opts := metav1.CreateOptions{}
	_, err := k8s.k8sClientSet.AppsV1().Deployments(k8s.namespace).Create(ctx, deployment, opts)
	if err != nil {
		log.Errorf("create deployment: %v", err)

		return err
	}

	return nil
}

func (k8s *K8sClient) UpdateDeployment(ctx context.Context, name string) error {
	// query deployment
	deployment, err := k8s.k8sClientSet.AppsV1().Deployments(k8s.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("query chaincode deployment failed: %v", err)

		return err
	}

	// upscale deployment
	replicas := int32(1)
	deployment.Spec.Replicas = &replicas

	items := []v1.KeyToPath{
		{
			Key:  "client.key",
			Path: "client.key",
		},
		{
			Key:  "client.crt",
			Path: "client.crt",
		},
		{
			Key:  "peer.crt",
			Path: "peer.crt",
		},
	}

	volumeMounts := []v1.VolumeMount{
		{
			Name:      name + "-config",
			MountPath: TLSClientKeyPath,
			SubPath:   "client.key",
		},
		{
			Name:      name + "-config",
			MountPath: TLSClientCertPath,
			SubPath:   "client.crt",
		},
		{
			Name:      name + "-config",
			MountPath: TLSClientRootCertPath,
			SubPath:   "peer.crt",
		},
	}

	if ctx.Value("version") == "v2.0.0" {
		items = append(items, v1.KeyToPath{
			Key:  "client_pem.key",
			Path: "client_pem.key",
		})
		items = append(items, v1.KeyToPath{
			Key:  "client_pem.crt",
			Path: "client_pem.crt",
		})

		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      name + "-config",
			MountPath: TLSClientKeyFile,
			SubPath:   "client_pem.key",
		})
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      name + "-config",
			MountPath: TLSClientCertFile,
			SubPath:   "client_pem.crt",
		})
	}

	// append volumes
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, v1.Volume{
		Name: name + "-config",
		VolumeSource: v1.VolumeSource{
			ConfigMap: &v1.ConfigMapVolumeSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name + "-configmap",
				},
				Items: items,
			},
		},
	})

	// mount volume
	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
		volumeMounts...,
	)

	deployment.ResourceVersion = ""

	// update deployment
	_, err = k8s.k8sClientSet.AppsV1().
		Deployments(k8s.namespace).
		Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Errorf("update chaincode deployment failed: %v", err)

		return err
	}

	return nil
}

func (k8s *K8sClient) CreateConfigMap(ctx context.Context, name string, data map[string]string) error {
	tlsConfigMap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name + "-configmap",
		},
		Data: data,
	}

	opts := metav1.CreateOptions{}

	_, err := k8s.k8sClientSet.CoreV1().ConfigMaps(k8s.namespace).Create(ctx, tlsConfigMap, opts)
	if err != nil {
		log.Errorf("create tlsConfigMap failed: %v", err)

		return err
	}

	return nil
}

func (k8s *K8sClient) DeleteChaincodeDeployment(ctx context.Context, name string) error {
	err := k8s.k8sClientSet.AppsV1().
		Deployments(k8s.namespace).
		Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("delete chaincode deployment failed: %v", err)
	}

	return nil
}

func (k8s *K8sClient) DeleteConfigMapDeployment(ctx context.Context, name string) error {
	err := k8s.k8sClientSet.CoreV1().
		ConfigMaps(k8s.namespace).
		Delete(context.Background(), name+"-configmap", metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("delete tls configmap failed: %v", err)
	}

	return nil
}

func (k8s *K8sClient) QueryDeploymentStatus(ctx context.Context, name string) (bool, error) {
	deployment, err := k8s.k8sClientSet.AppsV1().Deployments(k8s.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("get deployment %s failed: %v", name, err)

		return false, err
	}

	if deployment.Status.AvailableReplicas > 0 {
		return true, nil
	}

	return false, err
}
