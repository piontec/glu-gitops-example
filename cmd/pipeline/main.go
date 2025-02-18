package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/get-glu/glu"
	"github.com/get-glu/glu/pkg/fs"
	"github.com/get-glu/glu/pkg/pipelines"
	"github.com/get-glu/glu/ui"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	appsv1 "k8s.io/api/apps/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"oras.land/oras-go/v2/registry"
	"sigs.k8s.io/yaml"
)

func main() {
	slog.Info("starting system")
	if err := run(context.Background()); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("system exiting", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	system := glu.NewSystem(ctx, glu.Name("gitops-example"), glu.WithUI(ui.FS()))
	if err := pipelines.NewBuilder(system, glu.Name("gitops-example-app"), func() *AppResource {
		return &AppResource{
			Image: "ghcr.io/get-glu/gitops-example/app",
		}
	}).
		NewPhase(pipelines.OCIPhase[*AppResource](glu.Name("oci"), "app")).
		PromotesTo(pipelines.GitPhase[*AppResource](
			glu.Name("staging", glu.Label("url", "http://0.0.0.0:30081")),
			"gitopsexample",
		), NewCronScheduler(WithInterval(10*time.Second))).
		// ), schedule.New(schedule.WithInterval(30*time.Second))).
		PromotesTo(pipelines.GitPhase[*AppResource](
			glu.Name("production", glu.Label("url", "http://0.0.0.0:30082")),
			"gitopsexample",
		)).
		Build(); err != nil {
		return err
	}

	return system.Run()
	/*
		builder.
			func(b builder.PipelineBuilder[*AppResource]) error {
				// fetch the configured OCI repositority source named "checkout"
				ociSource, err := builder.OCISource(b, "app")
				if err != nil {
					return err
				}

				// fetch the configured Git repository source named "checkout"
				gitSource, err := builder.GitSource(b, "gitopsexample")
				if err != nil {
					return err
				}

				// build a phase which sources from the OCI repository
				ociPhase, err := b.NewPhase(glu.Name("oci"), ociSource)
				if err != nil {
					return err
				}

				// build a phase for the staging environment which source from the git repository
				// configure it to promote from the OCI phase
				staging, err := b.NewPhase(glu.Name("staging", glu.Label("url", "http://0.0.0.0:30081")),
					gitSource, core.PromotesFrom(ociPhase))
				if err != nil {
					return err
				}

				// build a phase for the production environment which source from the git repository
				// configure it to promote from the staging git phase
				_, err = b.NewPhase(glu.Name("production", glu.Label("url", "http://0.0.0.0:30082")),
					gitSource, core.PromotesFrom(staging))
				if err != nil {
					return err
				}

				// return configured pipeline to the system
				return nil
			}
			// AddTrigger(
			// cron_schedule.New(
			// cron_schedule.WithInterval(10 * time.Second),
			//		schedule.MatchesLabel("env", "staging"),
			//		// alternatively, the phase instance can be target directly with:
			//		// glu.ScheduleMatchesPhase(gitStaging),
			// ),
			// ).
			return nil
	*/
}

// AppResource is a custom envelope for carrying our specific repository configuration
// from one source to the next in our pipeline.
type AppResource struct {
	Image       string
	ImageDigest string
}

// Digest is a core required function for implementing glu.Resource
// It should return a unique digest for the state of the resource.
// In this instance we happen to be reading a unique digest from the source
// and so we can lean into that.
// This will be used for comparisons in the phase to decided whether or not
// a change has occurred when deciding if to update the target source.
func (c *AppResource) Digest() (string, error) {
	return c.ImageDigest, nil
}

// ReadFromOCIDescriptor is an OCI specific resource requirement.
// Its purpose is to read the resources state from a target OCI metadata descriptor.
// Here we're reading out the images digest from the metadata.
func (r *AppResource) ReadFromOCIDescriptor(d v1.Descriptor) error {
	r.ImageDigest = d.Digest.String()
	return nil
}

// ReadFrom is a Git specific resource requirement.
// It specifies how to read the resource from a target Filesystem.
// The type should navigate and source the relevant state from the fileystem provided.
// The function is also provided with metadata for the calling phase.
// This allows the defining type to adjust behaviour based on the context of the phase.
// Here we are reading a yaml file from a directory identified by the name of the phase.
func (r *AppResource) ReadFrom(_ context.Context, phase glu.Descriptor, fs fs.Filesystem) error {
	deployment, err := readDeployment(fs, fmt.Sprintf("env/%s/deployment.yaml", phase.Metadata.Name))
	if err != nil {
		return err
	}

	if containers := deployment.Spec.Template.Spec.Containers; len(containers) > 0 {
		ref, err := registry.ParseReference(containers[0].Image)
		if err != nil {
			return err
		}

		digest, err := ref.Digest()
		if err != nil {
			return err
		}

		r.ImageDigest = digest.String()
	}

	return nil
}

// WriteTo is a Git specific resource requirement.
// It specifies how to write the resource to a target Filesystem.
// The type should navigate and encode the state of the resource to the target Filesystem.
// The function is also provided with metadata for the calling phase.
// This allows the defining type to adjust behaviour based on the context of the phase.
// Here we are writing to a yaml file in a directory identified by the name of the phase.
func (r *AppResource) WriteTo(ctx context.Context, phase glu.Descriptor, fs fs.Filesystem) error {
	path := fmt.Sprintf("env/%s/deployment.yaml", phase.Metadata.Name)
	deployment, err := readDeployment(fs, path)
	if err != nil {
		return err
	}

	if containers := deployment.Spec.Template.Spec.Containers; len(containers) > 0 {
		containers[0].Image = fmt.Sprintf("%s@%s", r.Image, r.ImageDigest)

		for i := range containers[0].Env {
			if containers[0].Env[i].Name == "APP_IMAGE_DIGEST" {
				containers[0].Env[i].Value = r.ImageDigest
			}
		}
	}

	fi, err := fs.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}

	defer fi.Close()

	data, err := yaml.Marshal(deployment)
	if err != nil {
		return err
	}

	_, err = io.Copy(fi, bytes.NewReader(data))
	return err
}

func readDeployment(fs fs.Filesystem, path string) (*appsv1.Deployment, error) {
	fi, err := fs.OpenFile(
		path,
		os.O_RDONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	defer fi.Close()

	deployment := &appsv1.Deployment{}
	dec := k8syaml.NewYAMLOrJSONDecoder(fi, 1000)
	if err := dec.Decode(&deployment); err != nil {
		return nil, err
	}

	return deployment, nil
}
