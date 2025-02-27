package shutdown

import (
	"context"

	"k8s.io/klog/v2"

	"github.com/openshift/origin/pkg/monitor/monitorapi"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Consumer interface {
	Consume(event *corev1.Event)
	Done()
}

func StartMonitoringGracefulShutdownEvents(stop context.Context, recorder monitorapi.RecorderWriter, client kubernetes.Interface) {
	consumer := newConsumer(recorder)
	for namespace := range namespaceToServer {
		startGatheringByNamespace(stop, client, namespace, consumer)
	}
}

func startGatheringByNamespace(stop context.Context, client kubernetes.Interface, namespace string, consumer Consumer) {
	lw := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "events", namespace, fields.Everything())
	store := &cache.FakeCustomStore{
		// ReplaceFunc called when we do our initial list on starting the reflector.
		// With no resync period, it should not get called again.
		ReplaceFunc: func(items []interface{}, rv string) error {
			for _, obj := range items {
				event, ok := obj.(*corev1.Event)
				if !ok {
					continue
				}
				consumer.Consume(event)
			}
			return nil
		},
		AddFunc: func(obj interface{}) error {
			event, ok := obj.(*corev1.Event)
			if !ok {
				return nil
			}
			consumer.Consume(event)
			return nil
		},
		UpdateFunc: func(obj interface{}) error {
			event, ok := obj.(*corev1.Event)
			if !ok {
				return nil
			}
			consumer.Consume(event)
			return nil
		},
	}
	reflector := cache.NewReflector(lw, &corev1.Event{}, store, 0)
	go func() {
		klog.V(4).Infof("GracefulShutdownEvent: watching events namespace=%s", namespace)
		reflector.Run(stop.Done())
		klog.V(4).Infof("GracefulShutdownEvent: event watch ended namespace=%s", namespace)
		consumer.Done()
	}()
}
