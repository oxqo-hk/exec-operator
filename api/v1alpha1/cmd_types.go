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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CmdSpec defines the desired state of Cmd
type CmdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Command  string            `json:"command"`
	Selector map[string]string `json:"selector,omitempty"`
	IPs      []string          `json:"ips,omitempty"`
	Names    []string          `json:"keys,omitempty"`
}

// CmdStatus defines the observed state of Cmd
type CmdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Results map[string]CmdResult `json:"results,omitempty"`
	Done    bool                 `json:"done,omitempty"`
}

type CmdResult struct {
	Timestamp string `json:"timestamp,omitempty"`
	Stdout    string `json:"stdout,omitempty"`
	Stderr    string `json:"stderr,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cmd is the Schema for the cmds API
type Cmd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CmdSpec   `json:"spec,omitempty"`
	Status CmdStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CmdList contains a list of Cmd
type CmdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cmd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cmd{}, &CmdList{})
}
