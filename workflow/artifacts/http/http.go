package http

import (
	"git.code.oa.com/henrylwang/argo/errors"
	wfv1 "git.code.oa.com/henrylwang/argo/pkg/apis/workflow/v1alpha1"
	"git.code.oa.com/henrylwang/argo/workflow/common"
)

// HTTPArtifactDriver is the artifact driver for a HTTP URL
type HTTPArtifactDriver struct{}

// Load download artifacts from an HTTP URL
func (h *HTTPArtifactDriver) Load(inputArtifact *wfv1.Artifact, path string) error {
	// Download the file to a local file path
	return common.RunCommand("curl", "-sS", "-L", "-o", path, inputArtifact.HTTP.URL)
}

func (h *HTTPArtifactDriver) Save(path string, outputArtifact *wfv1.Artifact) error {
	return errors.Errorf(errors.CodeBadRequest, "HTTP output artifacts unsupported")
}
