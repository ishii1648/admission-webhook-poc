package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	// appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	// "github.com/imdario/mergo"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

// SidecarInjector annotates Pods
type SidecarInjector struct {
	Name          string
	client        client.Client
	decoder       *admission.Decoder
	SidecarConfig *Config
}

type Config struct {
	Spec corev1.PodSpec `json:"spec"`
}

// SidecarInjector adds an annotation to every incoming pods.
func (s *SidecarInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	log.Info("starting handle")

	pod := &corev1.Pod{}

	err := s.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

	shoudInjectSidecar := shoudInject(pod)

	if shoudInjectSidecar {
		log.Info("Injecting sidecar...")

		pod.Spec.Volumes = append(pod.Spec.Volumes, s.SidecarConfig.Spec.Volumes...)
		pod.Spec.Containers = append(pod.Spec.Containers, s.SidecarConfig.Spec.Containers...)
		pod.Annotations["webserver-sidecar-added"] = "true"

		log.Info("Sidecar ", s.Name, " injected.")
	} else {
		log.Info("Inject not needed.")
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// InjectClient injects the client.
func (s *SidecarInjector) InjectClient(c client.Client) error {
	s.client = c
	return nil
}

// SidecarInjector implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (s *SidecarInjector) InjectDecoder(d *admission.Decoder) error {
	s.decoder = d
	return nil
}

// SidecarInjector implements inject.Client.
// A client will be automatically injected.

// InjectDecoder injects the decoder.
func shoudInject(pod *corev1.Pod) bool {
	shouldInjectSidecar, err := strconv.ParseBool(pod.Annotations["webserver-injection"])

	if err != nil {
		shouldInjectSidecar = false
	}

	if shouldInjectSidecar {
		alreadyUpdated, err := strconv.ParseBool(pod.Annotations["webserver-sidecar-added"])

		if err == nil && alreadyUpdated {
			shouldInjectSidecar = false
		}
	}

	log.Info("Should Inject: ", shouldInjectSidecar)

	return shouldInjectSidecar
}
