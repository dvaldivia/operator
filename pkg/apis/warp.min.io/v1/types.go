/*
 * Copyright (C) 2020, MinIO, Inc.
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3,
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License, version 3,
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 *
 */

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=warp,singular=warp

// Warp is a specification for a MinIO resource
type Warp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Scheduler WarpScheduler `json:"scheduler,omitempty"`
	Spec      WarpSpec      `json:"spec"`
	// Status provides details of the state of the Warp
	// +optional
	Status WarpStatus `json:"status"`
}

// WarpScheduler (`scheduler`) - Object describing Kubernetes Scheduler to use for deploying the Warp.
type WarpScheduler struct {
	// *Optional* +
	//
	// Specify the name of the https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/[Kubernetes scheduler] to be used to schedule Warp pods
	Name string `json:"name"`
}

// WarpSpec is the spec for a Warp resource
type WarpSpec struct {
	// Image defines the Warp Docker image.
	// +optional
	Image string `json:"image,omitempty"`
	// ImagePullSecret defines the secret to be used for pull image from a private Docker image.
	// +optional
	ImagePullSecret corev1.LocalObjectReference `json:"imagePullSecret,omitempty"`

	// HostTarget is the location of MinIO that will be benchmarked.
	HostTarget IOConfig `json:"hostTarget,omitempty"`

	// WarpClient is to warp clients and run benchmarks there, if this is specified, the underlying
	// statefulset for the warp resource won't be created for the run.
	// +optional
	WarpClient string `json:"warpClient,omitempty"`

	// Mixed is the configuration for a `mixed` benchmark. Only one benchmark configuration can be specified.
	// +optional
	Mixed *MixedConfiguration `json:"mixed,omitempty"`
	// Get is the configuration for a `get` benchmark. Only one benchmark configuration can be specified.
	// +optional
	Get *GetConfiguration `json:"get,omitempty"`
	// Put is the configuration for a `put` benchmark. Only one benchmark configuration can be specified.
	// +optional
	Put *PutConfiguration `json:"put,omitempty"`
	// Delete is the configuration for a `delete` benchmark. Only one benchmark configuration can be specified.
	// +optional
	Delete *DeleteConfiguration `json:"delete,omitempty"`
	// List is the configuration for a `list` benchmark. Only one benchmark configuration can be specified.
	// +optional
	List *ListConfiguration `json:"list,omitempty"`
	// Stat is the configuration for a `stat` benchmark. Only one benchmark configuration can be specified.
	// +optional
	Stat *StatConfiguration `json:"stat,omitempty"`
	// Select is the configuration for a `select` benchmark. Only one benchmark configuration can be specified.
	// +optional
	Select *SelectConfiguration `json:"select,omitempty"`
	// Versioned is the configuration for a `versioned` benchmark. Only one benchmark configuration can be specified.
	// +optional
	Versioned *VersionedConfiguration `json:"versioned,omitempty"`

	// Pod Management Policy for pod created by StatefulSet
	// +optional
	PodManagementPolicy appsv1.PodManagementPolicyType `json:"podManagementPolicy,omitempty"`
	// If provided, use these environment variables for Warp resource
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
	// Security Context allows user to set entries like runAsUser, privilege escalation etc.
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`
	// ServiceAccountName is the name of the ServiceAccount to use to run pods of all Warp
	// Pods created as a part of this Warp Resource.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// PriorityClassName indicates the Pod priority and hence importance of a Pod relative to other Pods.
	// This is applied to MinIO pods only.
	// Refer Kubernetes documentation for details https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/#priorityclass
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// Image pull policy. One of Always, Never, IfNotPresent.
	// This is applied to MinIO pods only.
	// Refer Kubernetes documentation for details https://kubernetes.io/docs/concepts/containers/images#updating-images
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
}

type IOConfig struct {

	// Insecure is disable TLS certificate verification
	// +optional
	Insecure string `json:"insecure,omitempty"`

	// Host is Host. Multiple hosts can be specified as a comma separated list. (default: "127.0.0.1:9000") [$WARP_HOST]
	// +optional
	Host string `json:"host,omitempty"`

	// CredsSecret references a secret with the access/secret keys for accessing the hosts
	// +optional
	CredsSecret *corev1.LocalObjectReference `json:"credsSecret,omitempty"`

	// Tls is Use TLS (HTTPS) for transport [$WARP_TLS]
	// +optional
	Tls string `json:"tls,omitempty"`

	// Region is Specify a custom Region [$WARP_REGION]
	// +optional
	Region string `json:"region,omitempty"`

	// Encrypt is Encrypt/decrypt objects (using server-side encryption with random keys)
	// +optional
	Encrypt string `json:"encrypt,omitempty"`

	// Bucket is Bucket to use for benchmark data. ALL DATA WILL BE DELETED IN BUCKET! (default: "warp-benchmark-bucket")
	// +optional
	Bucket string `json:"bucket,omitempty"`

	// Host-select is Host selection algorithm. Can be "weighed" or "roundrobin" (default: "weighed")
	// +optional
	HostSelect string `json:"host-select,omitempty"`

	// Concurrent is Run this many Concurrent operations (default: 12)
	// +optional
	Concurrent string `json:"concurrent,omitempty"`

	// Noprefix is Do not use separate prefix for each thread
	// +optional
	Noprefix string `json:"noprefix,omitempty"`

	// DisableMultipart is disable multipart uploads
	// +optional
	DisableMultipart string `json:"disable-multipart,omitempty"`

	// Md5 is Add MD5 sum to uploads
	// +optional
	Md5 string `json:"md5,omitempty"`

	// storage-class is Specify custom storage class, for instance 'STANDARD' or 'REDUCED_REDUNDANCY'.
	// +optional
	StorageClass string `json:"storage-class,omitempty"`
}

type GenConfig struct {
	// obj.generator is Use specific data generator (default: "random")
	// +optional
	ObjGenerator string `json:"obj.generator,omitempty"`

	// obj.randsize is Randomize size of Objects so they will be up to the specified size
	// +optional
	ObjRandsize string `json:"obj.randsize,omitempty"`
}

type BenchConfig struct {
	// Benchdata is Output benchmark+profile data to this file. By default unique filename is generated.
	// +optional
	Benchdata string `json:"benchdata,omitempty"`

	// Serverprof is Run MinIO server profiling during benchmark; possible values are 'cpu', 'mem', 'block', 'mutex' and 'trace'.
	// +optional
	Serverprof string `json:"serverprof,omitempty"`

	// Duration is Duration to run the benchmark. Use 's' and 'm' to specify seconds and minutes. (default: 5m0s)
	// +optional
	Duration string `json:"duration,omitempty"`

	// Autoterm is Auto terminate when benchmark is considered stable.
	// +optional
	Autoterm string `json:"autoterm,omitempty"`

	// autoterm.dur is Minimum Duration where output must have been stable to allow automatic termination. (default: 10s)
	// +optional
	AutotermDur string `json:"autoterm.dur,omitempty"`

	// autoterm.pct is The percentage the last 6/25 time blocks must be within current speed to auto terminate. (default: 7.5)
	// +optional
	AutotermPct string `json:"autoterm.pct,omitempty"`

	// Noclear is Do not clear Bucket before or after running benchmarks. Use when running multiple clients.
	// +optional
	Noclear string `json:"noclear,omitempty"`

	// Syncstart is Specify a benchmark start time. Time format is 'hh:mm' where hours are specified in 24h format, server TZ.
	// +optional
	Syncstart string `json:"syncstart,omitempty"`
}

type AnalyzeConfig struct {
	// analyze.dur is Split analysis into durations of this length. Can be '1s', '5s', '1m', etc.
	// +optional
	AnalyzeDur string `json:"analyze.dur,omitempty"`

	// analyze.out is Output aggregated data as to file
	// +optional
	AnalyzeOut string `json:"analyze.out,omitempty"`

	// analyze.op is Only output for this op. Can be GET/PUT/DELETE, etc.
	// +optional
	AnalyzeOp string `json:"analyze.op,omitempty"`

	// analyze.host is Only output for this Host.
	// +optional
	AnalyzeHost string `json:"analyze.host,omitempty"`

	// analyze.skip is Additional Duration to skip when analyzing data. (default: 0s)
	// +optional
	AnalyzeSkip string `json:"analyze.skip,omitempty"`

	// AnalyzeV is Display additional analysis data.
	// +optional
	AnalyzeV string `json:"analyze.v,omitempty"`
}

// MixedConfiguration exposes the configurations for a mixed benchmark
type MixedConfiguration struct {
	GenConfig
	BenchConfig
	AnalyzeConfig

	// Objects is Number of Objects to upload. (default: 2500)
	// +optional
	Objects int `json:"objects,omitempty"`

	// ObjSize is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	ObjSize string `json:"obj_size,omitempty"`

	// get-distrib is The amount of GET operations. (default: 45)
	// +optional
	GetDistrib string `json:"get-distrib,omitempty"`

	// stat-distrib is The amount of STAT operations. (default: 30)
	// +optional
	StatDistrib string `json:"stat-distrib,omitempty"`

	// put-distrib is The amount of PUT operations. (default: 15)
	// +optional
	PutDistrib string `json:"put-distrib,omitempty"`

	// delete-distrib is The amount of DELETE operations. Must be at least the same as PUT. (default: 10)
	// +optional
	DeleteDistrib string `json:"delete-distrib,omitempty"`
}

// GetConfiguration exposes the configurations for a mixed benchmark
type GetConfiguration struct {
	GenConfig
	BenchConfig
	AnalyzeConfig

	// Objects is Number of Objects to upload. (default: 2500)
	// +optional
	Objects int `json:"objects,omitempty"`

	// ObjSize is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	ObjSize string `json:"obj_size,omitempty"`
}

// PutConfiguration exposes the configurations for a mixed benchmark
type PutConfiguration struct {
	GenConfig
	BenchConfig
	AnalyzeConfig

	// ObjSize is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	ObjSize string `json:"obj_size,omitempty"`
}

// DeleteConfiguration exposes the configurations for a delete benchmark
type DeleteConfiguration struct {
	GenConfig
	BenchConfig
	AnalyzeConfig

	// Objects is Number of Objects to upload. (default: 2500)
	// +optional
	Objects int `json:"objects,omitempty"`

	// ObjSize is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	ObjSize string `json:"obj_size,omitempty"`

	// Batch is Number of DELETE operations per batch.
	// +optional
	Batch int `json:"batch,omitempty"`
}

// ListConfiguration exposes the configurations for a delete benchmark
type ListConfiguration struct {
	GenConfig
	BenchConfig
	AnalyzeConfig

	// Objects is Number of Objects to upload. (default: 2500)
	// +optional
	Objects int `json:"objects,omitempty"`

	// ObjSize is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	ObjSize string `json:"obj_size,omitempty"`
}

// StatConfiguration exposes the configurations for a stat benchmark
type StatConfiguration struct {
	GenConfig
	BenchConfig
	AnalyzeConfig

	// Objects is Number of Objects to upload. (default: 2500)
	// +optional
	Objects int `json:"objects,omitempty"`

	// ObjSize is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	ObjSize string `json:"obj_size,omitempty"`
}

// SelectConfiguration exposes the configurations for a select benchmark
type SelectConfiguration struct {
	GenConfig
	BenchConfig
	AnalyzeConfig

	// Objects is Number of Objects to upload. (default: 2500)
	// +optional
	Objects int `json:"objects,omitempty"`

	// ObjSize is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	ObjSize string `json:"obj_size,omitempty"`
	// Query is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	Query string `json:"query,omitempty"`
}

// VersionedConfiguration exposes the configurations for a versioned benchmark
type VersionedConfiguration struct {
	GenConfig
	BenchConfig
	AnalyzeConfig

	// Objects is Number of Objects to upload. (default: 2500)
	// +optional
	Objects int `json:"objects,omitempty"`

	// ObjSize is Size of each generated object. Can be a number or 10KiB/MiB/GiB. All sizes are base 2 binary. (default: "10MiB")
	// +optional
	ObjSize string `json:"obj_size,omitempty"`

	// get-distrib is The amount of GET operations. (default: 45)
	// +optional
	GetDistrib string `json:"get-distrib,omitempty"`

	// stat-distrib is The amount of STAT operations. (default: 30)
	// +optional
	StatDistrib string `json:"stat-distrib,omitempty"`

	// put-distrib is The amount of PUT operations. (default: 15)
	// +optional
	PutDistrib string `json:"put-distrib,omitempty"`

	// delete-distrib is The amount of DELETE operations. Must be at least the same as PUT. (default: 10)
	// +optional
	DeleteDistrib string `json:"delete-distrib,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WarpList is a list of Warp resources
type WarpList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Warp `json:"items"`
}

// WarpStatus is the status for a Warp resource
type WarpStatus struct {
	CurrentState string `json:"currentState"`
}
