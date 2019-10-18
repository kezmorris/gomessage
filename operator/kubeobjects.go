package main

import (
	"os"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPodObject returns an apiv1 conference Pod object with name conferenceName
func GetPodObject(conferenceName string) *apiv1.Pod {
	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      conferenceName,
			Namespace: "default",
			Labels: map[string]string{
				"app": "call_pod",
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:            "conference",
					Image:           os.Getenv("CONF_IMAGE"),
					ImagePullPolicy: apiv1.PullIfNotPresent,
				},
			},
		},
	}
}
