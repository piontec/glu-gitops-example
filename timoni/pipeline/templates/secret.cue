package templates

import (
	corev1 "k8s.io/api/core/v1"
)

#Secret: corev1.#Secret & {
	#config:    #Config
	apiVersion: "v1"
	kind:       "Secret"
	metadata:   #config.metadata
	stringData: github_password: #config.pipeline.password
}
