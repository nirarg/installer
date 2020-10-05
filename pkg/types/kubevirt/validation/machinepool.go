package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/kubevirt"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *kubevirt.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.CPU <= 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("cpu"), p.CPU, "CPU must be positive"))
	}

	if _, err := resource.ParseQuantity(p.StorageSize); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("storage"), p.StorageSize, "Storage size must be of Quantity type format"))
	}

	if _, err := resource.ParseQuantity(p.Memory); err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.Memory, "Memory must be of Quantity type format"))
	}

	return allErrs
}
