/*
Copyright 2022 oxqo.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/oxqo-hk/exec-operator/api/v1alpha1"
	execv1alpha1 "github.com/oxqo-hk/exec-operator/api/v1alpha1"
)

// CmdReconciler reconciles a Cmd object
type CmdReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	RESTClient rest.Interface
	RESTConfig *rest.Config
}

//+kubebuilder:rbac:groups=exec.github.com,resources=cmds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=exec.github.com,resources=cmds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=exec.github.com,resources=cmds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cmd object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *CmdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log.SetLogger(zap.New())
	logger := log.Log.WithValues("cmd_controller", req.NamespacedName)

	cmd := &execv1alpha1.Cmd{}
	err := r.Get(ctx, req.NamespacedName, cmd)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("failed to get cmd, must be deletion.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get cmd.")
		return ctrl.Result{}, err
	}
	//check if this cmd is already executed
	if cmd.Status.Done == true {
		return ctrl.Result{}, nil
	}

	pods := &corev1.PodList{}
	//List by selector
	if len(cmd.Spec.Selector) != 0 {
		selector := labels.SelectorFromSet(cmd.Spec.Selector)
		opts := &client.ListOptions{
			LabelSelector: selector,
			Namespace:     cmd.Namespace,
		}
		err = r.List(ctx, pods, opts)
		if err != nil {
			logger.Info("failed to get any pods with given selector")
		}
	}
	//Get Targets by IPs and Names
	all_pods := &corev1.PodList{}
	opts := &client.ListOptions{
		Namespace: cmd.Namespace,
	}
	err = r.List(ctx, all_pods, opts)
	if err != nil {
		logger.Error(err, fmt.Sprintf("failed to List pods"))
		return ctrl.Result{}, err
	}
	for _, pod := range all_pods.Items {
		for _, ip := range cmd.Spec.IPs {
			if pod.Status.PodIP == ip {
				AppendPodIfNotDup(pods, pod)
			}
		}
		for _, name := range cmd.Spec.Names {
			if pod.Name == name {
				AppendPodIfNotDup(pods, pod)
			}
		}
	}

	results := map[string]v1alpha1.CmdResult{}
	for _, pod := range pods.Items {
		//do exec from here
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		err = r.ExecuteCommand(&pod, []string{"sh", "-c", cmd.Spec.Command}, bufio.NewWriter(&stdout), bufio.NewWriter(&stderr))
		if err != nil {
			logger.Error(err, "error calling exec api")
		}
		logger.Info(string(stdout.String()) + ":from pod - " + pod.Name)
		results[pod.Namespace+"/"+pod.Name] = v1alpha1.CmdResult{
			Stdout:    stdout.String(),
			Stderr:    stderr.String(),
			Timestamp: time.Now().String(),
		}
	}
	cmd_copy := cmd.DeepCopy()
	cmd_copy.Status.Results = results
	cmd_copy.Status.Done = true
	err = r.Status().Update(ctx, cmd_copy)
	if err != nil {
		logger.Error(err, "failed to update cmd status: "+cmd.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func AppendPodIfNotDup(pods *corev1.PodList, el corev1.Pod) {
	for _, pod := range pods.Items {
		if pod.Name == el.Name && pod.Namespace == el.Namespace {
			return
		}
	}
	pods.Items = append(pods.Items, el)
	return
}

func (r *CmdReconciler) ExecuteCommand(pod *corev1.Pod, command []string, stdout io.Writer, stderr io.Writer) error {
	exec := r.RESTClient.
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			//which Container should run command??
			Container: pod.Spec.Containers[0].Name,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
		}, runtime.NewParameterCodec(r.Scheme))

	executor, err := remotecommand.NewSPDYExecutor(r.RESTConfig, "POST", exec.URL())
	if err != nil {
		return err
	}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdout: stdout,
		Stderr: stderr,
		Tty:    false,
	})
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CmdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&execv1alpha1.Cmd{}).
		Complete(r)
}
