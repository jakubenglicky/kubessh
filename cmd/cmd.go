package cmd

import (
	"context"
	"fmt"
	"io"
	"jakubenglicky/kubessh/utils/config"
	"jakubenglicky/kubessh/utils/k8s"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

var FlagNamespace string
var FlagDeployment string

var RootCmd = &cobra.Command{
	Use:   "kubessh",
	Short: "KubeSSH",
	Args:  cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		clientset, restConfig, _, _ := k8s.KubernetesClient()

		podList, err := clientset.CoreV1().Pods(ResolveNamespace()).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			panic(err)
		}

		var podName string
		for _, pod := range podList.Items {
			if strings.HasPrefix(pod.GetName(), FlagDeployment) {
				podName = pod.GetName()
				break
			}
		}

		fmt.Println(ResolveNamespace())
		fmt.Println(podName)

		var r io.Reader
		var w io.Writer
		ExecCmdExample(clientset, restConfig, podName, "ls", r, w, w)

	},
}

func Execute() {
	RootCmd.Flags().StringVarP(
		&FlagNamespace,
		"namespace",
		"n",
		"",
		"Namespace",
	)

	RootCmd.Flags().StringVarP(
		&FlagDeployment,
		"deployment",
		"d",
		"",
		"Deployment",
	)

	RootCmd.Execute()
}

func ResolveNamespace() string {
	_, _, defaultNamespace, _ := k8s.KubernetesClient()
	var config = config.GetConfig()

	var namespace = defaultNamespace

	if config.DefaultNamespace != "" {
		namespace = config.DefaultNamespace
	}

	if FlagNamespace != "" {
		namespace = FlagNamespace
	}

	return namespace
}

// ExecCmd exec command on specific pod and wait the command's output.
func ExecCmdExample(client kubernetes.Interface, config *restclient.Config, podName string,
	command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := []string{
		"sh",
		"-c",
		command,
	}
	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(ResolveNamespace()).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	if stdin == nil {
		option.Stdin = false
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return err
	}

	return nil
}
