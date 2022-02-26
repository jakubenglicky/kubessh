package root

import (
	"context"
	"fmt"
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/spf13/cobra"
	"jakubenglicky/kubessh/utils/k8s"
)

var RootCmd = &cobra.Command{
	Use:   "kubessh",
	Short: "KubeSSH",
	Run: func(c *cobra.Command, args []string) {
		clientset, defaultNamespace, _ := k8s.KubernetesClient()

		var namespace = defaultNamespace
		var deployment = ""

		if len(args) == 1 {
			deployment = args[0]	
		}

		if len(args) == 2 {
			namespace = args[0]
			deployment = args[1]
		}
		
		podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			panic(err)
		}

		var podName = ""
		for _, pod := range podList.Items {
			if strings.HasPrefix(pod.GetName(), deployment) {
				podName = pod.GetName()
				break
			}
		}
	
		fmt.Println(podName)
		fmt.Println(deployment)
	
	},
}