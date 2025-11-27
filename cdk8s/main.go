package main

import (
	"example.com/cdk8s/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, ns string, appLabel string) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), &cdk8s.ChartProps{
		Namespace: jsii.String(ns),
	})

	labels := map[string]*string{
		"app": jsii.String(appLabel),
	}

	k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
		Spec: &k8s.DeploymentSpec{
			Replicas: jsii.Number(3),
			Selector: &k8s.LabelSelector{
				MatchLabels: &labels,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &labels,
				},
				Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{{
						Name:  jsii.String("app-container"),
						Image: jsii.String("nginx:1.19.10"),
						Ports: &[]*k8s.ContainerPort{{
							ContainerPort: jsii.Number(80),
						}},
					}},
				},
			},
		},
	})

	k8s.NewKubeService(chart, jsii.String("service"), &k8s.KubeServiceProps{
		Spec: &k8s.ServiceSpec{
			Ports: &[]*k8s.ServicePort{
				{
					Port:       jsii.Number(8080),
					Name:       jsii.String("http"),
					TargetPort: k8s.IntOrString_FromNumber(jsii.Number(80)),
				},
			},
			Type:     jsii.String("ClusterIP"),
			Selector: &labels,
		},
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "Getting Started", "default", "my-app")
	app.Synth()
}
