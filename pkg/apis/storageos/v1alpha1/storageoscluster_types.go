package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterPhase is the phase of the storageos cluster at a given point in time.
type ClusterPhase string

// Constants for operator defaults values and different phases.
const (
	ClusterPhaseInitial ClusterPhase = ""
	ClusterPhaseRunning              = "Running"

	DefaultNamespace = "storageos"

	DefaultServiceName         = "storageos"
	DefaultServiceType         = "ClusterIP"
	DefaultServiceExternalPort = 5705
	DefaultServiceInternalPort = 5705

	DefaultIngressHostname = "storageos.local"

	DefaultNodeContainerImage                   = "storageos/node:1.0.0"
	DefaultCSIDriverRegistrarContainerImage     = "quay.io/k8scsi/driver-registrar:v0.4.1"
	DefaultCSIExternalProvisionerContainerImage = "quay.io/k8scsi/csi-provisioner:v0.4.0"
	DefaultCSIExternalAttacherContainerImage    = "quay.io/k8scsi/csi-attacher:v0.4.0"
	DefaultInitContainerImage                   = "storageos/init:0.1"

	DefaultCSIEndpoint                 = "unix://var/lib/kubelet/plugins/storageos/csi.sock"
	DefaultCSIRegistrarSocketDir       = "/var/lib/kubelet/device-plugins/"
	DefaultCSIKubeletDir               = "/var/lib/kubelet"
	DefaultCSIPluginDir                = "/var/lib/kubelet/plugins/storageos/"
	DefaultCSIDeviceDir                = "/dev"
	DefaultCSIRegistrationDir          = "/var/lib/kubelet/plugins"
	DefaultCSIKubeletRegistrationPath  = "/var/lib/kubelet/plugins/storageos/csi.sock"
	DefaultCSIDriverRegistrationMode   = "node-register"
	DefaultCSIDriverRequiresAttachment = "true"
)

// StorageOSClusterSpec defines the desired state of StorageOSCluster
type StorageOSClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file

	// Join is the join token used for service discovery.
	Join string `json:"join"`

	// CSI defines the configurations for CSI.
	CSI StorageOSClusterCSI `json:"csi"`

	// ResourceNS is the kubernetes Namespace where storageos resources are
	// provisioned.
	ResourceNS string `json:"namespace"`

	// Service is the Service configuration for the cluster nodes.
	Service StorageOSClusterService `json:"service"`

	// SecretRefName is the name of the secret object that contains all the
	// sensitive cluster configurations.
	SecretRefName string `json:"secretRefName"`

	// SecretRefNamespace is the namespace of the secret reference.
	SecretRefNamespace string `json:"secretRefNamespace"`

	// SharedDir is the shared directory to be used when the kubelet is running
	// in a container.
	// Typically: "/var/lib/kubelet/plugins/kubernetes.io~storageos".
	// If not set, defaults will be used.
	SharedDir string `json:"sharedDir"`

	// Ingress defines the ingress configurations used in the cluster.
	Ingress StorageOSClusterIngress `json:"ingress"`

	// Images defines the various container images used in the cluster.
	Images ContainerImages `json:"images"`

	// KVBackend defines the key-value store backend used in the cluster.
	KVBackend StorageOSClusterKVBackend `json:"kvBackend"`

	// Pause is to pause the operator for the cluster.
	Pause bool `json:"pause"`

	// Debug is to set debug mode of the cluster.
	Debug bool `json:"debug"`

	// NodeSelectorTerms is to set the placement of storageos pods using
	// node affinity requiredDuringSchedulingIgnoredDuringExecution.
	NodeSelectorTerms []corev1.NodeSelectorTerm `json:"nodeSelectorTerms"`
	// Resources is to set the resource requirements of the storageos containers.
	Resources corev1.ResourceRequirements `json:"resources"`
}

// StorageOSClusterStatus defines the observed state of StorageOSCluster
type StorageOSClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file

	Phase            ClusterPhase          `json:"phase"`
	NodeHealthStatus map[string]NodeHealth `json:"nodeHealthStatus,omitempty"`
	Nodes            []string              `json:"nodes"`
	Ready            string                `json:"ready"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StorageOSCluster is the Schema for the storageosclusters API
// +k8s:openapi-gen=true
type StorageOSCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StorageOSClusterSpec   `json:"spec,omitempty"`
	Status StorageOSClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StorageOSClusterList contains a list of StorageOSCluster
type StorageOSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageOSCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StorageOSCluster{}, &StorageOSClusterList{})
}

// GetResourceNS returns the namespace where all the resources should be provisioned.
func (s StorageOSClusterSpec) GetResourceNS() string {
	if s.ResourceNS != "" {
		return s.ResourceNS
	}
	return DefaultNamespace
}

// GetNodeContainerImage returns node container image.
func (s StorageOSClusterSpec) GetNodeContainerImage() string {
	if s.Images.NodeContainer != "" {
		return s.Images.NodeContainer
	}
	return DefaultNodeContainerImage
}

// GetInitContainerImage returns init container image.
func (s StorageOSClusterSpec) GetInitContainerImage() string {
	if s.Images.InitContainer != "" {
		return s.Images.InitContainer
	}
	return DefaultInitContainerImage
}

// GetCSIDriverRegistrarImage returns CSI driver registrar container image.
func (s StorageOSClusterSpec) GetCSIDriverRegistrarImage() string {
	if s.Images.CSIDriverRegistrarContainer != "" {
		return s.Images.CSIDriverRegistrarContainer
	}
	return DefaultCSIDriverRegistrarContainerImage
}

// GetCSIExternalProvisionerImage returns CSI external provisioner container image.
func (s StorageOSClusterSpec) GetCSIExternalProvisionerImage() string {
	if s.Images.CSIExternalProvisionerContainer != "" {
		return s.Images.CSIExternalProvisionerContainer
	}
	return DefaultCSIExternalProvisionerContainerImage
}

// GetCSIExternalAttacherImage returns CSI external attacher container image.
func (s StorageOSClusterSpec) GetCSIExternalAttacherImage() string {
	if s.Images.CSIExternalAttacherContainer != "" {
		return s.Images.CSIExternalAttacherContainer
	}
	return DefaultCSIExternalAttacherContainerImage
}

// GetServiceName returns the service name.
func (s StorageOSClusterSpec) GetServiceName() string {
	if s.Service.Name != "" {
		return s.Service.Name
	}
	return DefaultServiceName
}

// GetServiceType returns the service type.
func (s StorageOSClusterSpec) GetServiceType() string {
	if s.Service.Type != "" {
		return s.Service.Type
	}
	return DefaultServiceType
}

// GetServiceExternalPort returns the service external port.
func (s StorageOSClusterSpec) GetServiceExternalPort() int {
	if s.Service.ExternalPort != 0 {
		return s.Service.ExternalPort
	}
	return DefaultServiceExternalPort
}

// GetServiceInternalPort returns the service internal port.
func (s StorageOSClusterSpec) GetServiceInternalPort() int {
	if s.Service.InternalPort != 0 {
		return s.Service.InternalPort
	}
	return DefaultServiceInternalPort
}

// GetIngressHostname returns the ingress host name.
func (s StorageOSClusterSpec) GetIngressHostname() string {
	if s.Ingress.Hostname != "" {
		return s.Ingress.Hostname
	}
	return DefaultIngressHostname
}

// GetCSIEndpoint returns the CSI unix socket endpoint path.
func (s StorageOSClusterSpec) GetCSIEndpoint() string {
	if s.CSI.Endpoint != "" {
		return s.CSI.Endpoint
	}
	return DefaultCSIEndpoint
}

// GetCSIRegistrarSocketDir returns the CSI registrar socket dir.
func (s StorageOSClusterSpec) GetCSIRegistrarSocketDir() string {
	if s.CSI.RegistrarSocketDir != "" {
		return s.CSI.RegistrarSocketDir
	}
	return DefaultCSIRegistrarSocketDir
}

// GetCSIKubeletDir returns the Kubelet dir.
func (s StorageOSClusterSpec) GetCSIKubeletDir() string {
	if s.CSI.KubeletDir != "" {
		return s.CSI.KubeletDir
	}
	return DefaultCSIKubeletDir
}

// GetCSIPluginDir returns the CSI plugin dir.
func (s StorageOSClusterSpec) GetCSIPluginDir() string {
	if s.CSI.PluginDir != "" {
		return s.CSI.PluginDir
	}
	return DefaultCSIPluginDir
}

// GetCSIDeviceDir returns the CSI device dir.
func (s StorageOSClusterSpec) GetCSIDeviceDir() string {
	if s.CSI.DeviceDir != "" {
		return s.CSI.DeviceDir
	}
	return DefaultCSIDeviceDir
}

// GetCSIRegistrationDir returns the CSI registration dir.
func (s StorageOSClusterSpec) GetCSIRegistrationDir() string {
	if s.CSI.RegistrationDir != "" {
		return s.CSI.RegistrationDir
	}
	return DefaultCSIRegistrationDir
}

// GetCSIKubeletRegistrationPath returns the CSI Kubelet Registration Path.
func (s StorageOSClusterSpec) GetCSIKubeletRegistrationPath() string {
	if s.CSI.KubeletRegistrationPath != "" {
		return s.CSI.KubeletRegistrationPath
	}
	return DefaultCSIKubeletRegistrationPath
}

// GetCSIDriverRegistrationMode returns the CSI Driver Registration Mode.
func (s StorageOSClusterSpec) GetCSIDriverRegistrationMode() string {
	if s.CSI.DriverRegistrationMode != "" {
		return s.CSI.DriverRegistrationMode
	}
	return DefaultCSIDriverRegistrationMode
}

// GetCSIDriverRequiresAttachment returns the CSI Driver Requires Attachment
func (s StorageOSClusterSpec) GetCSIDriverRequiresAttachment() string {
	if s.CSI.DriverRequiresAttachment != "" {
		return s.CSI.DriverRequiresAttachment
	}
	return DefaultCSIDriverRequiresAttachment
}

// ContainerImages contains image names of all the containers used by the operator.
type ContainerImages struct {
	NodeContainer                   string `json:"nodeContainer"`
	InitContainer                   string `json:"initContainer"`
	CSIDriverRegistrarContainer     string `json:"csiDriverRegistrarContainer"`
	CSIExternalProvisionerContainer string `json:"csiExternalProvisionerContainer"`
	CSIExternalAttacherContainer    string `json:"csiExternalAttacherContainer"`
}

// StorageOSClusterCSI contains CSI configurations.
type StorageOSClusterCSI struct {
	Enable                       bool   `json:"enable"`
	Endpoint                     string `json:"endpoint"`
	EnableProvisionCreds         bool   `json:"enableProvisionCreds"`
	EnableControllerPublishCreds bool   `json:"enableControllerPublishCreds"`
	EnableNodePublishCreds       bool   `json:"enableNodePublishCreds"`
	RegistrarSocketDir           string `json:"registrarSocketDir"`
	KubeletDir                   string `json:"kubeletDir"`
	PluginDir                    string `json:"pluginDir"`
	DeviceDir                    string `json:"deviceDir"`
	RegistrationDir              string `json:"registrationDir"`
	KubeletRegistrationPath      string `json:"kubeletRegistrationPath"`
	DriverRegistrationMode       string `json:"driverRegisterationMode"`
	DriverRequiresAttachment     string `json:"driverRequiresAttachment"`
}

// StorageOSClusterService contains Service configurations.
type StorageOSClusterService struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	ExternalPort int               `json:"externalPort"`
	InternalPort int               `json:"internalPort"`
	Annotations  map[string]string `json:"annotations"`
}

// StorageOSClusterIngress contains Ingress configurations.
type StorageOSClusterIngress struct {
	Enable      bool              `json:"enable"`
	Hostname    string            `json:"hostname"`
	TLS         bool              `json:"tls"`
	Annotations map[string]string `json:"annotations"`
}

// NodeHealth contains health status of a node.
type NodeHealth struct {
	DirectfsInitiator string `json:"directfsInitiator"`
	Director          string `json:"director"`
	KV                string `json:"kv"`
	KVWrite           string `json:"kvWrite"`
	Nats              string `json:"nats"`
	Presentation      string `json:"presentation"`
	Rdb               string `json:"rdb"`
}

// StorageOSClusterKVBackend stores key-value store backend configurations.
type StorageOSClusterKVBackend struct {
	Address string `json:"address"`
	Backend string `json:"backend"`
}
