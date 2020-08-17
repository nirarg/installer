package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/kubevirt"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *kubevirt.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// TODO <nargaman> Is it enough to check file exists, or need to check its valid kubeconfig
	if err := validate.FileInfo(p.InfraClusterKubeConfig); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("InfraClusterKubeConfig"), p.InfraClusterKubeConfig, err.Error()))
	}

	if err := validate.DomainName(p.InfraClusterIngressDomain, true); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("InfraClusterIngressDomain"), p.InfraClusterIngressDomain, err.Error()))
	}

	// TODO <nargaman> Add InfraClusterNamespace validation -check that namespace exists in InfraCluster
	if p.InfraClusterNamespace == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("InfraClusterNamespace"), p.InfraClusterNamespace, "InfraClusterNamespace can't be empty"))
	}

	// TODO <nargaman> InfraClusterStorageClass - Is there a way to validate this value? can be empty

	if p.SourcePvcName == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("SourcePvcName"), p.SourcePvcName, "SourcePvcName can't be empty"))
	}

	return allErrs
}
