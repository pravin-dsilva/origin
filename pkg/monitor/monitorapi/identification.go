package monitorapi

import (
	"fmt"
	"strings"
)

type BackendConnectionType string

const (
	NewConnectionType    BackendConnectionType = "new"
	ReusedConnectionType BackendConnectionType = "reused"
)

func LocateRouteForDisruptionCheck(ns, name, disruptionBackendName string, connectionType BackendConnectionType) string {
	return fmt.Sprintf("ns/%s route/%s disruption/%s connection/%s", ns, name, disruptionBackendName, connectionType)
}

func LocateDisruptionCheck(disruptionBackendName string, connectionType BackendConnectionType) string {
	return fmt.Sprintf("disruption/%s connection/%s", disruptionBackendName, connectionType)
}

func IsE2ETest(l Locator) bool {
	_, ret := E2ETestFromLocator(l)
	return ret
}

func E2ETestFromLocator(l Locator) (string, bool) {
	test, ok := l.Keys[LocatorE2ETestKey]
	return test, ok
}

func NodeLocator(testName string) string {
	return fmt.Sprintf("node/%v", testName)
}

func IsNode(locator string) bool {
	_, ret := NodeFromLocator(locator)
	return ret
}

func NodeFromLocator(locator string) (string, bool) {
	ret := NodeFrom(LocatorParts(locator))
	return ret, len(ret) > 0
}

func OperatorLocator(operatorName string) string {
	return fmt.Sprintf("clusteroperator/%v", operatorName)
}

func OperatorFromLocator(locator string) (string, bool) {
	ret := OperatorFrom(LocatorParts(locator))
	return ret, len(ret) > 0
}

func NamespaceFromLocator(locator string) string {
	locatorParts := LocatorParts(locator)
	if ns, ok := locatorParts["ns"]; ok {
		return ns
	}
	if ns, ok := locatorParts["namespace"]; ok {
		return ns
	}
	return ""
}

func AlertFromLocator(locator string) string {
	return AlertFrom(LocatorParts(locator))
}

func LocatorParts(locator string) map[string]string {
	parts := map[string]string{}

	tags := strings.Split(locator, " ")
	for _, tag := range tags {
		keyValue := strings.SplitN(tag, "/", 2)
		if len(keyValue) == 1 {
			parts[keyValue[0]] = ""
		} else {
			parts[keyValue[0]] = keyValue[1]
		}
	}

	return parts
}

func NamespaceFrom(locatorParts map[string]string) string {
	if ns, ok := locatorParts["ns"]; ok {
		return ns
	}
	if ns, ok := locatorParts["namespace"]; ok {
		return ns
	}
	return ""
}

func NodeFrom(locatorParts map[string]string) string {
	return locatorParts["node"]
}

func OperatorFrom(locatorParts map[string]string) string {
	return locatorParts["clusteroperator"]
}

func AlertFrom(locatorParts map[string]string) string {
	return locatorParts["alert"]
}

func DisruptionFrom(locatorParts map[string]string) string {
	return locatorParts["disruption"]
}

func DisruptionConnectionTypeFrom(locatorParts map[string]string) string {
	return locatorParts["connection"]
}

func DisruptionLoadBalancerTypeFrom(locatorParts map[string]string) string {
	return locatorParts["load-balancer"]
}

func DisruptionProtocolFrom(locatorParts map[string]string) string {
	return locatorParts["protocol"]
}

func DisruptionTargetAPIFrom(locatorParts map[string]string) string {
	return locatorParts["target"]
}

func IsEventForLocator(locator string) EventIntervalMatchesFunc {
	return func(eventInterval Interval) bool {
		if eventInterval.Locator == locator {
			return true
		}
		return false
	}
}

type NamespacedReference struct {
	Namespace string
	Name      string
	UID       string
}
