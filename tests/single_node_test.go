package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"kubevirt.io/controller-lifecycle-operator-sdk/api"
	sspv1beta1 "kubevirt.io/ssp-operator/api/v1beta1"
)

var _ = Describe("Single Node Topology", func() {
	BeforeEach(func() {
		strategy.SkipUnlessSingleReplicaTopologyMode()
		waitUntilDeployed()
	})

	AfterEach(func() {
		strategy.RevertToOriginalSspCr()
	})

	It("[test_id:7792] Number of Template Validator replicas is not bigger then one", func() {
		watch, err := StartWatch(sspListerWatcher)
		Expect(err).ToNot(HaveOccurred())
		defer watch.Stop()

		var newValidatorReplicas int32 = 3
		updateSsp(func(foundSsp *sspv1beta1.SSP) {
			foundSsp.Spec.TemplateValidator.Replicas = &newValidatorReplicas
		})

		// Watch changes until above change
		err = WatchChangesUntil(watch, func(updatedSsp *sspv1beta1.SSP) bool {
			return *updatedSsp.Spec.TemplateValidator.Replicas == newValidatorReplicas &&
				updatedSsp.Generation > updatedSsp.Status.ObservedGeneration
		}, shortTimeout)
		Expect(err).ToNot(HaveOccurred())

		// Watch changes until SSP operator updates ObservedGeneration
		err = WatchChangesUntil(watch, func(updatedSsp *sspv1beta1.SSP) bool {
			return *updatedSsp.Spec.TemplateValidator.Replicas == newValidatorReplicas &&
				updatedSsp.Generation == updatedSsp.Status.ObservedGeneration && updatedSsp.Status.Phase == api.PhaseDeployed
		}, shortTimeout)
		Expect(err).ToNot(HaveOccurred())

		deployment := getTemplateValidatorDeployment()
		Expect(int(deployment.Status.Replicas)).Should(Equal(1), "In Single Mode Topology the number of replicas is at most 1")
	})

	It("[test_id:7793] Number of Template Validator replicas can be set to 0", func() {
		watch, err := StartWatch(sspListerWatcher)
		Expect(err).ToNot(HaveOccurred())
		defer watch.Stop()

		var newValidatorReplicas int32 = 0
		updateSsp(func(foundSsp *sspv1beta1.SSP) {
			foundSsp.Spec.TemplateValidator.Replicas = &newValidatorReplicas
		})

		err = WatchChangesUntil(watch, func(updatedSsp *sspv1beta1.SSP) bool {
			return *updatedSsp.Spec.TemplateValidator.Replicas == newValidatorReplicas &&
				updatedSsp.Generation > updatedSsp.Status.ObservedGeneration
		}, shortTimeout)
		Expect(err).ToNot(HaveOccurred())

		// Watch changes until SSP operator updates ObservedGeneration
		err = WatchChangesUntil(watch, func(updatedSsp *sspv1beta1.SSP) bool {
			return *updatedSsp.Spec.TemplateValidator.Replicas == newValidatorReplicas &&
				updatedSsp.Generation == updatedSsp.Status.ObservedGeneration && updatedSsp.Status.Phase == api.PhaseDeployed
		}, shortTimeout)
		Expect(err).ToNot(HaveOccurred())
		deployment := getTemplateValidatorDeployment()
		Expect(int(deployment.Status.Replicas)).Should(Equal(0), "In Single Mode Topology the number of replicas is at most 1")
	})
})
