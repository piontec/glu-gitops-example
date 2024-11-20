bundle: {
	_pass: string @timoni(runtime:string:CONFIGURATION_REPOSITORY_PASSWORD)

	apiVersion: "v1alpha1"
	name:       "flux-aio"
	instances: {
		"flux": {
			module: url: "oci://ghcr.io/stefanprodan/modules/flux-aio"
			namespace: "flux-system"
			values: {
				controllers: {
					helm: enabled:         true
					kustomize: enabled:    true
					notification: enabled: true
				}
				hostNetwork:     false
				securityProfile: "privileged"
			}
		}
		"staging": {
			module: url: "oci://ghcr.io/stefanprodan/modules/flux-git-sync"
			namespace: "flux-system"
			values: {
				git: {
					url:  "https://github.com/get-glu/gitops-example"
					ref:  "refs/heads/main"
					path: "./env/staging"
				}
				sync: {
					targetNamespace: "default"
					wait:            true
				}
			}
		}
		"production": {
			module: url: "oci://ghcr.io/stefanprodan/modules/flux-git-sync"
			namespace: "flux-system"
			values: {
				git: {
					url:  "https://github.com/get-glu/gitops-example"
					ref:  "refs/heads/main"
					path: "./env/production"
				}
				sync: {
					targetNamespace: "default"
					wait:            true
				}
			}
		}
		"pipeline": {
			module: url: "file://pipeline"
			namespace: "glu"
			values: {
				image: pullPolicy:  "Always"
				pipeline: password: _pass
			}
		}
	}
}
