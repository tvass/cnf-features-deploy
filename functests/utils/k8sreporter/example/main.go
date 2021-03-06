package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/openshift-kni/cnf-features-deploy/functests/utils/k8sreporter"
	"github.com/openshift-kni/cnf-features-deploy/functests/utils/namespaces"

	perfUtils "github.com/openshift-kni/performance-addon-operators/functests/utils"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	ptpUtils "github.com/openshift/ptp-operator/test/utils"
	sriovNamespaces "github.com/openshift/sriov-network-operator/test/util/namespaces"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	performancev1alpha1 "github.com/openshift-kni/performance-addon-operators/pkg/apis/performance/v1alpha1"
	ptpv1 "github.com/openshift/ptp-operator/pkg/apis/ptp/v1"
	sriovv1 "github.com/openshift/sriov-network-operator/pkg/apis/sriovnetwork/v1"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "the kubeconfig path")
	report := flag.String("report", "report.log", "the file name used for the report")

	flag.Parse()

	addToScheme := func(s *runtime.Scheme) {
		ptpv1.AddToScheme(s)
		mcfgv1.AddToScheme(s)
		performancev1alpha1.SchemeBuilder.AddToScheme(s)
		sriovv1.AddToScheme(s)

	}

	namespacesToDump := map[string]bool{
		"openshift-performance-addon":      true,
		"openshift-ptp":                    true,
		"openshift-sriov-network-operator": true,
		"cnf-features-testing":             true,
		perfUtils.NamespaceTesting:         true,
		namespaces.DpdkTest:                true,
		sriovNamespaces.Test:               true,
		ptpUtils.NamespaceTesting:          true,
	}

	crds := []k8sreporter.CRData{
		{Cr: &mcfgv1.MachineConfigPoolList{}},
		{Cr: &ptpv1.PtpConfigList{}},
		{Cr: &ptpv1.NodePtpDeviceList{}},
		{Cr: &ptpv1.PtpOperatorConfigList{}},
		{Cr: &performancev1alpha1.PerformanceProfileList{}},
		{Cr: &sriovv1.SriovNetworkNodePolicyList{}},
		{Cr: &sriovv1.SriovNetworkList{}},
		{Cr: &sriovv1.SriovNetworkNodePolicyList{}},
		{Cr: &sriovv1.SriovOperatorConfigList{}},
	}

	skipPods := func(pod *v1.Pod) bool {
		found := namespacesToDump[pod.Namespace]
		return !found
	}

	f, err := os.OpenFile(*report, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open the file: %v\n", err)
		return
	}
	defer f.Close()

	reporter, err := k8sreporter.New(*kubeconfig, addToScheme, skipPods, f, crds...)
	if err != nil {
		log.Fatalf("Failed to initialize the reporter %s", err)
	}
	reporter.Dump(10 * time.Minute)
}
