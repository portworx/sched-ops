package k8s

import (
	aut_v1alpaha1 "github.com/libopenstorage/autopilot-api/pkg/apis/autopilot/v1alpha1"
	"github.com/portworx/sched-ops/k8s/autopilot"
)

// AutopilotRule CRD - BEGIN

// CreateAutopilotRule creates the AutopilotRule object
func (k *k8sOps) CreateAutopilotRule(rule *aut_v1alpaha1.AutopilotRule) (*aut_v1alpaha1.AutopilotRule, error) {
	return autopilot.Instance().CreateAutopilotRule(rule)
}

// GetAutopilotRule gets the AutopilotRule for the provided name
func (k *k8sOps) GetAutopilotRule(name string) (*aut_v1alpaha1.AutopilotRule, error) {
	return autopilot.Instance().GetAutopilotRule(name)
}

// UpdateAutopilotRule updates the AutopilotRule
func (k *k8sOps) UpdateAutopilotRule(rule *aut_v1alpaha1.AutopilotRule) (*aut_v1alpaha1.AutopilotRule, error) {
	return autopilot.Instance().UpdateAutopilotRule(rule)
}

// DeleteAutopilotRule deletes the AutopilotRule of the given name
func (k *k8sOps) DeleteAutopilotRule(name string) error {
	return autopilot.Instance().DeleteAutopilotRule(name)
}

// ListAutopilotRules lists AutopilotRules
func (k *k8sOps) ListAutopilotRules() (*aut_v1alpaha1.AutopilotRuleList, error) {
	return autopilot.Instance().ListAutopilotRules()
}

// AutopilotRule CRD - END
