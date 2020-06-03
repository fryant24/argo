package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wfv1 "git.code.oa.com/henrylwang/argo/pkg/apis/workflow/v1alpha1"
	"git.code.oa.com/henrylwang/argo/test"
)

// TestStepsFailedRetries ensures a steps template will recognize exhausted retries
func TestStepsFailedRetries(t *testing.T) {
	wf := test.LoadTestWorkflow("testdata/steps-failed-retries.yaml")
	woc := newWoc(*wf)
	woc.operate()
	assert.Equal(t, string(wfv1.NodeFailed), string(woc.wf.Status.Phase))
}
