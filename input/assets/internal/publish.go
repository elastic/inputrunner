// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package internal

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	stateless "github.com/elastic/beats/v7/filebeat/input/v2/input-stateless"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/elastic-agent-libs/mapstr"
)

type AssetOption func(beat.Event) beat.Event

// Publish emits a `beat.Event` to the specified publisher, with the provided
// parameters
func Publish(publisher stateless.Publisher, opts ...AssetOption) {
	event := beat.Event{Fields: mapstr.M{}, Meta: mapstr.M{}}
	for _, o := range opts {
		event = o(event)
	}
	publisher.Publish(event)
}

func WithAssetCloudProvider(value string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["cloud.provider"] = value
		return e
	}
}

func WithAssetRegion(value string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["cloud.region"] = value
		return e
	}
}

func WithAssetAccountID(value string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["cloud.account.id"] = value
		return e
	}
}

func WithAssetTypeAndID(t, id string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["asset.type"] = t
		e.Fields["asset.id"] = id
		e.Fields["asset.ean"] = fmt.Sprintf("%s:%s", t, id)
		return e
	}
}

func WithAssetKind(value string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["asset.kind"] = value
		return e
	}
}

func WithAssetParents(value []string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["asset.parents"] = value
		return e
	}
}

func WithAssetChildren(value []string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["asset.children"] = value
		return e
	}
}

func WithAssetMetadata(value mapstr.M) AssetOption {
	return func(e beat.Event) beat.Event {
		flattenedValue := value.Flatten()
		for k, v := range flattenedValue {
			e.Fields["asset.metadata."+k] = v
		}
		return e
	}
}

func WithNodeData(name string, startTime *metav1.Time) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["kubernetes.node.name"] = name
		e.Fields["kubernetes.node.start_time"] = startTime
		return e
	}
}

func WithPodData(name, uid, namespace string, startTime *metav1.Time) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["kubernetes.pod.name"] = name
		e.Fields["kubernetes.pod.uid"] = uid
		e.Fields["kubernetes.pod.start_time"] = startTime
		e.Fields["kubernetes.namespace"] = namespace
		return e
	}
}

func WithContainerData(name, uid, namespace, state string, startTime *metav1.Time) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["kubernetes.container.name"] = name
		e.Fields["kubernetes.container.uid"] = uid
		e.Fields["kubernetes.container.start_time"] = startTime
		e.Fields["kubernetes.container.state"] = state
		e.Fields["kubernetes.namespace"] = namespace
		return e
	}
}

func WithHostData(hostname, architecture string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["host.hostname"] = hostname
		e.Fields["host.architecture"] = architecture
		return e
	}
}

func WithHostOsData(osBuild, osFamily, osKernel, osName, osPlatform, osType, osVersion string) AssetOption {
	return func(e beat.Event) beat.Event {
		e.Fields["host.os.build"] = osBuild
		e.Fields["host.os.family"] = osFamily
		e.Fields["host.os.kernel"] = osKernel
		e.Fields["host.os.name"] = osName
		e.Fields["host.os.platform"] = osPlatform
		e.Fields["host.os.type"] = osType
		e.Fields["host.os.version"] = osVersion
		return e
	}
}
func ToMapstr(input map[string]string) mapstr.M {
	out := mapstr.M{}
	for k, v := range input {
		out[k] = v
	}
	return out
}
