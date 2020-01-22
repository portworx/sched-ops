package k8s

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/portworx/sched-ops/k8s/prometheus"
)

func (k *k8sOps) ListPrometheuses(namespace string) (*monitoringv1.PrometheusList, error) {
	return prometheus.Instance().ListPrometheuses(namespace)
}

func (k *k8sOps) GetPrometheus(name string, namespace string) (*monitoringv1.Prometheus, error) {
	return prometheus.Instance().GetPrometheus(name, namespace)
}

func (k *k8sOps) CreatePrometheus(p *monitoringv1.Prometheus) (*monitoringv1.Prometheus, error) {
	return prometheus.Instance().CreatePrometheus(p)
}

func (k *k8sOps) UpdatePrometheus(p *monitoringv1.Prometheus) (*monitoringv1.Prometheus, error) {
	return prometheus.Instance().UpdatePrometheus(p)
}

func (k *k8sOps) DeletePrometheus(name, namespace string) error {
	return prometheus.Instance().DeletePrometheus(name, namespace)
}

func (k *k8sOps) ListServiceMonitors(namespace string) (*monitoringv1.ServiceMonitorList, error) {
	return prometheus.Instance().ListServiceMonitors(namespace)
}

func (k *k8sOps) GetServiceMonitor(name string, namespace string) (*monitoringv1.ServiceMonitor, error) {
	return prometheus.Instance().GetServiceMonitor(name, namespace)
}

func (k *k8sOps) CreateServiceMonitor(serviceMonitor *monitoringv1.ServiceMonitor) (*monitoringv1.ServiceMonitor, error) {
	return prometheus.Instance().CreateServiceMonitor(serviceMonitor)
}

func (k *k8sOps) UpdateServiceMonitor(serviceMonitor *monitoringv1.ServiceMonitor) (*monitoringv1.ServiceMonitor, error) {
	return prometheus.Instance().UpdateServiceMonitor(serviceMonitor)
}

func (k *k8sOps) DeleteServiceMonitor(name, namespace string) error {
	return prometheus.Instance().DeleteServiceMonitor(name, namespace)
}

func (k *k8sOps) ListPrometheusRules(namespace string) (*monitoringv1.PrometheusRuleList, error) {
	return prometheus.Instance().ListPrometheusRules(namespace)
}

func (k *k8sOps) GetPrometheusRule(name string, namespace string) (*monitoringv1.PrometheusRule, error) {
	return prometheus.Instance().GetPrometheusRule(name, namespace)
}

func (k *k8sOps) CreatePrometheusRule(rule *monitoringv1.PrometheusRule) (*monitoringv1.PrometheusRule, error) {
	return prometheus.Instance().CreatePrometheusRule(rule)
}

func (k *k8sOps) UpdatePrometheusRule(rule *monitoringv1.PrometheusRule) (*monitoringv1.PrometheusRule, error) {
	return prometheus.Instance().UpdatePrometheusRule(rule)
}

func (k *k8sOps) DeletePrometheusRule(name, namespace string) error {
	return prometheus.Instance().DeletePrometheusRule(name, namespace)
}

func (k *k8sOps) ListAlertManagers(namespace string) (*monitoringv1.AlertmanagerList, error) {
	return prometheus.Instance().ListAlertManagers(namespace)
}

func (k *k8sOps) GetAlertManager(name string, namespace string) (*monitoringv1.Alertmanager, error) {
	return prometheus.Instance().GetAlertManager(name, namespace)
}

func (k *k8sOps) CreateAlertManager(alertmanager *monitoringv1.Alertmanager) (*monitoringv1.Alertmanager, error) {
	return prometheus.Instance().CreateAlertManager(alertmanager)
}

func (k *k8sOps) UpdateAlertManager(alertmanager *monitoringv1.Alertmanager) (*monitoringv1.Alertmanager, error) {
	return prometheus.Instance().UpdateAlertManager(alertmanager)
}

func (k *k8sOps) DeleteAlertManager(name, namespace string) error {
	return prometheus.Instance().DeleteAlertManager(name, namespace)
}
