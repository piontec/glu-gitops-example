// Code generated by timoni.
// Note that this file must have no imports and all values must be concrete.

@if(!debug)

package main

// Defaults
values: {
	image: {
		repository: "ghcr.io/get-glu/gitops-example/pipeline"
		digest:     ""
		tag:        "latest"
	}
	test: image: {
		repository: "cgr.dev/chainguard/curl"
		digest:     ""
		tag:        "latest"
	}
	pipeline: password: ""
}
