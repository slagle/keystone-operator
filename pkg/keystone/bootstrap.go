package keystone

import (
	comv1 "github.com/openstack-k8s-operators/keystone-operator/pkg/apis/keystone/v1"
	util "github.com/openstack-k8s-operators/lib-common/pkg/util"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type bootstrapOptions struct {
	AdminPassword string
	APIEndpoint   string
	ServiceName   string
}

// BootstrapJob func
func BootstrapJob(cr *comv1.KeystoneAPI, configMapName string, APIEndpoint string) *batchv1.Job {

	// NOTE: as a convention the configmap is name the same as the service
	opts := bootstrapOptions{cr.Spec.AdminPassword, APIEndpoint, configMapName}
	runAsUser := int64(0)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName + "-bootstrap",
			Namespace: cr.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:      "OnFailure",
					ServiceAccountName: "keystone",
					Containers: []corev1.Container{
						{
							Name:    configMapName + "-bootstrap",
							Image:   cr.Spec.ContainerImage,
							Command: []string{"/bin/bash", "-c", util.ExecuteTemplateFile("bootstrap.sh", &opts)},
							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runAsUser,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "KOLLA_CONFIG_STRATEGY",
									Value: "COPY_ALWAYS",
								},
							},
							VolumeMounts: getVolumeMounts(),
						},
					},
				},
			},
		},
	}
	job.Spec.Template.Spec.Volumes = getVolumes(configMapName)
	return job
}
