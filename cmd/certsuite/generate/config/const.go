package config

// Menu names
//
//nolint:unused
const (
	// Main menu
	create = "Create"
	show   = "Show"
	save   = "Save"
	quit   = "Exit"
	// Create
	certSuiteResources = "resources"
	exceptions         = "Exceptions"
	collector          = "Collector"
	settings           = "Settings"
	previousMenu       = "\U0001F878"
	// CertSuite resources
	namespaces          = "Namespaces"
	pods                = "Pods"
	operators           = "Operators"
	crdFilters          = "CRD filters"
	managedDeployments  = "Managed Deployments"
	managedStatefulSets = "Managed StatefulSets"
	// Exceptions
	kernelTaints            = "Kernel taints"
	helmCharts              = "Helm charts"
	protocolNames           = "Protocol names"
	services                = "Services"
	nonScalableDeployments  = "Non-scalable Deployments"
	nonScalableStatefulSets = "Non-scalable StatefulSets"
	// Collector
	appEndPoint = "Application end point"
	executedBy  = "Certification executor"
	partnerName = "Partner name"
	appPassword = "Application password"
	// Settings
	probeDaemonSet = "Probe DaemonSet namespace"
)

// Menu help
const (
	// Main menu
	createConfigHelp = "Create a configuration for the Certification Suite"
	showConfigHelp   = "Show the current configuration in YAML format"
	saveConfigHelp   = `Save the current configuration to a YAML file (default "certsuite_config.yaml")`
	exitHelp         = "Exit the tool (changes not saved will be lost)"
	backHelp         = "Move to previous menu"
	// Create
	certSuiteResourcesHelp = `Configure the workload resources of the CertSuite to be verified.
Only the resources that the CertSuite uses are required to be configured. The rest can be left empty.
Usually a basic configuration includes "Namespaces" and "Pods" at least.`
	exceptionsdHelp = `Allow adding exceptions to skip several checks for different resources.
The exceptions must be justified in order to pass the CertSuite. Feedback
regarding the exceptions configured can be provided in an HTML page after loading
the claim.json file with the results.`
	collectordHelp = `Parameters required to send the CertSuite Certification Suite results to a data collector.`
	settingsHelp   = `Configure various settings for the CertSuite Certification Suite.`
	// CertSuite resources
	namespacesHelp = `The namespaces in which the workload under test will be deployed.`

	podLabelsHelp = `The labels that each Pod of the workload under test must have to be verified
by the Certification Suite.
If a new label is used for this purpose make sure it is added to the workload's Pods,
ideally in the pod's definition as the on-the-fly labels are lost if the Pod gets
rescheduled.
For Pods own by a Deployment, the same label as the one defined in the
"spec.selector.matchLabels" section of the Deployment can be used.`

	operatorLabelsHelp = `The labels that each operator's CSV of the workload under test must have to be verified
by the Certification Suite.
If a new label is used for this purpose make sure it is added to the workload operator's CSVs.`

	crdFiltersHelp = `The CRD name suffix used to filter the workload's CRDs among all the CRDs present in the cluster.
It must also be specified if the resources own by the CRD are scalable or not in order to avoid
some lifecycle test cases.`
	managedDeploymentsHelp = `The Deployments managed by a Custom Resource whose scaling is controlled using
the "scale" subresource of the CR.
The CRD defining that CR should be included in the CRD filters with the scalable
property set to true. If so, the test case "lifecycle-deployment-scaling" will be
skipped, otherwise it will fail.`
	managedStatefulSetsHelp = `The StatefulSets managed by a Custom Resource whose scaling is controlled using
the "scale" subresource of the CR.
The CRD defining that CR should be included in the CRD filters with the scalable
property set to true. If so, the test case "lifecycle-statefulset-scaling" will be
skipped, otherwise it will fail.`
	// Exceptions
	kernelTaintsHelp = `The list of kernel modules loaded by the workload that make the Linux kernel mark itself
as "tainted" but that should skip verification.
Test cases affected: platform-alteration-tainted-node-kernel.`
	helmChartsHelp = `The list of Helm charts that the workload uses whose certification status will not be verified.
If no exception is configured, the certification status for all Helm charts will be checked
in the OpenShift Helms Charts repository (see https://charts.openshift.io/).
Test cases affected: affiliated-certification-helmchart-is-certified`
	protocolNamesHelp = `The list of allowed protocol names to be used for container port names.
The name field of a container port must be of the form <protocol>[-<suffix>] where <protocol> must
be allowed by default or added to this list. The optional <suffix> can be chosen by the application.
Protocol names allowed by default: "grpc", "grpc-web", "http", "http2", "tcp", "udp".
Test cases affected: manageability-container-port-name-format.`
	servicesHelp = `The list of Services that will skip verification.
Services included in this list will be filtered out at the autodiscovery stage
and will not be subject to checks in any test case.
Tests cases affected: networking-dual-stack-service, access-control-service-type`
	nonScalableDeploymentsHelp = `The list of Deployments that do not support scale in/out operations.
Deployments included in this list will skip any scaling operation check.
Test cases affected: lifecycle-deployment-scaling`
	nonScalableStatefulSetsHelp = `The list of StatefulSets that do not support scale in/out operations.
StatefulSets included in this list will skip any scaling operation check.
Test cases affected: lifecycle-statefulset-scaling`
	// Collector (TODO)
	// Settings
	probeDaemonSetHelp = `Set the namespace where the probe DaemonSet will be deployed.
The namespace will be created in case it does not exist. If not set, the default namespace
is "certsuite".
This DaemonSet, called "certsuite-probe" is deployed and used internally by the Certification Suite
to issue some shell commands that are needed in certain test cases. Some of these test cases might
fail or be skipped in case it is not deployed correctly.`
)

// Prompts, syxtax, examples
const (
	// CertSuite resources
	namespacePrompt  = "Enter a comma-separated list of the namespaces in which the CertSuite is deploying its workload."
	namespaceSyntax  = "ns1[, ns2]..."
	namespaceExample = "workload, sample-workload"
	podsPrompt       = "Enter a comma-separated list of labels to identify the CertSuite's Pods under test."
	podsSyntax       = "pod-label-1[, pod-label-2]..."
	podsExample      = "redhat-best-practices-for-k8s.com/generic: target"
	operatorsPrompt  = "Enter a comma-separated list of labels to identify the CertSuite's operators under test."
	operatorsSyntax  = "operator-label-1[, operator-label-2]..."
	operatorsExample = "redhat-best-practices-for-k8s.com/operator1: target"
	crdFiltersPrompt = "Enter a comma-separated list of the CRD's name suffixes that the CertSuite contains. Also, specify if the\n" +
		"resources managed by those CRDs are scalable."
	crdFiltersSyntax           = "crd-name-suffix/{true|false}[,crd-name-suffix/{true|false}]..."
	crdFiltersExample          = "group1.test.com/true"
	managedDeploymentsPrompt   = "Enter a comma-separated list of Deployments that are managed by a Custom Resource."
	managedDeploymentsSyntax   = "managed-deployment1[, managed-deployment2]..."
	managedDeploymentsExample  = "group1-deployment"
	managedStatefulSetsPrompt  = "Enter a comma-separated list of StatefulSets that are managed by a Custom Resource."
	managedStatefulSetsSyntax  = "managed-statefulset1[, managed-statefulset2]..."
	managedStatefulSetsExample = "group1-statefulset"
	// Exceptions
	kernelTaintsPrompt             = "Enter a comma-separated list of kernel taints (modules)"
	kernelTaintsSyntax             = "mod1[,mod2]..."
	kernelTaintsExample            = "vboxsf, vboxguest"
	helmChartsPrompt               = "Enter a comma-separated list of Helm charts that will skip verification."
	helmChartsSyntax               = "chart1[,chart2]..."
	helmChartsExample              = "coredns"
	protocolNamesPrompt            = "Enter a comma-separated list of protocol names"
	protocolNamesSyntax            = "proto1[,proto2]..."
	protocolNamesExample           = "http3, sctp"
	servicesPrompt                 = "Enter a comma-separated list of Service names"
	servicesSyntax                 = "svc1[,svc2]..."
	servicesExample                = "hazelcast-platform-controller-manager-service, hazelcast-platform-webhook-service"
	nonScalableDeploymentsPrompt   = "Enter a comma-separated list of Deployments that do not support scaling operations."
	nonScalableDeploymentsSyntax   = "deployment1-name/deployment1-namespace[,deployment2-name/deployment2-namespace]..."
	nonScalableDeploymentsExample  = "deployment-test/certsuite-test"
	nonScalableStatefulSetsPrompt  = "Enter a comma-separated list of StatefulSets that do not support scaling operations."
	nonScalableStatefulSetsSyxtax  = "statefulset1-name/statefulset1-namespace[,statefulset2-name/statefulset2-namespace]..."
	nonScalableStatefulSetsExample = "statefulset-test-test/certsuite-test"
	// Collector (TODO)
	// Settings
	probeDaemonSetPrompt  = "Enter the namespace in which de probe daemonset will be deployed."
	probeDaemonSetSyntax  = "ds-namespace"
	probeDaemonSetExample = "certsuite-probe"
)

// Internal constants
const (
	defaultConfigFileName        = "certsuite_config.yml"
	defaultConfigFilePermissions = 0o644
)
