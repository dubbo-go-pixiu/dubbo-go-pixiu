// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetry

import (
	"strconv"
	"strings"

	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/networking/util"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/serviceregistry/provider"
	"github.com/apache/dubbo-go-pixiu/pkg/config/host"
)

var (
	// StatName patterns
	serviceStatPattern         = "%SERVICE%"
	serviceFQDNStatPattern     = "%SERVICE_FQDN%"
	servicePortStatPattern     = "%SERVICE_PORT%"
	servicePortNameStatPattern = "%SERVICE_PORT_NAME%"
	subsetNameStatPattern      = "%SUBSET_NAME%"
)

// BuildStatPrefix builds a stat prefix based on the stat pattern.
func BuildStatPrefix(statPattern string, host string, subset string, port *model.Port, attributes *model.ServiceAttributes) string {
	prefix := strings.ReplaceAll(statPattern, serviceStatPattern, shortHostName(host, attributes))
	prefix = strings.ReplaceAll(prefix, serviceFQDNStatPattern, host)
	prefix = strings.ReplaceAll(prefix, subsetNameStatPattern, subset)
	prefix = strings.ReplaceAll(prefix, servicePortStatPattern, strconv.Itoa(port.Port))
	prefix = strings.ReplaceAll(prefix, servicePortNameStatPattern, port.Name)
	return prefix
}

// BuildInboundStatPrefix builds a stat prefix based on the stat pattern and filter chain telemetry data.
func BuildInboundStatPrefix(statPattern string, tm FilterChainMetadata, subset string, port uint32, portName string) string {
	prefix := strings.ReplaceAll(statPattern, serviceStatPattern, tm.ShortHostname())
	prefix = strings.ReplaceAll(prefix, serviceFQDNStatPattern, tm.InstanceHostname.String())
	prefix = strings.ReplaceAll(prefix, subsetNameStatPattern, subset)
	prefix = strings.ReplaceAll(prefix, servicePortStatPattern, strconv.Itoa(int(port)))
	prefix = strings.ReplaceAll(prefix, servicePortNameStatPattern, portName)
	return prefix
}

// shortHostName constructs the name from kubernetes hosts based on attributes (name and namespace).
// For other hosts like VMs, this method does not do any thing - just returns the passed in host as is.
func shortHostName(host string, attributes *model.ServiceAttributes) string {
	if attributes.ServiceRegistry == provider.Kubernetes {
		return attributes.Name + "." + attributes.Namespace
	}
	return host
}

// TraceOperation builds the string format: "%s:%d/*" for a given host and port
func TraceOperation(host string, port int) string {
	// Format : "%s:%d/*"
	return util.DomainName(host, port) + "/*"
}

// FilterChainMetadata defines additional metadata for telemetry use for a filter chain.
type FilterChainMetadata struct {
	// InstanceHostname defines the hostname of the service this filter chain is built for.
	// Note: This is best effort; this may be empty if generated by Sidecar config, and there may be multiple
	// Services that make up the filter chain.
	InstanceHostname host.Name
	// KubernetesServiceNamespace is the namespace the service is defined in, if it is for a Kubernetes Service.
	// Note: This is best effort; this may be empty if generated by Sidecar config, and there may be multiple
	// Services that make up the filter chain.
	KubernetesServiceNamespace string
	// KubernetesServiceName is the name of service, if it is for a Kubernetes Service.
	// Note: This is best effort; this may be empty if generated by Sidecar config, and there may be multiple
	// Services that make up the filter chain.
	KubernetesServiceName string
}

// ShortHostname constructs the name from kubernetes service name if available or just uses instance host name.
func (tm FilterChainMetadata) ShortHostname() string {
	if tm.KubernetesServiceName != "" {
		return tm.KubernetesServiceName + "." + tm.KubernetesServiceNamespace
	}
	return tm.InstanceHostname.String()
}
