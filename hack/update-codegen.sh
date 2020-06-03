#!/bin/bash
set -eux -o pipefail

go get k8s.io/code-generator/cmd/go-to-protobuf@v0.16.7-beta.0

bash ${GOPATH}/pkg/mod/k8s.io/code-generator@v0.16.7-beta.0/generate-groups.sh \
  "deepcopy,client,informer,lister" \
  git.code.oa.com/henrylwang/argo/pkg/client git.code.oa.com/henrylwang/argo/pkg/apis \
  workflow:v1alpha1 \
  --go-header-file ./hack/custom-boilerplate.go.txt
