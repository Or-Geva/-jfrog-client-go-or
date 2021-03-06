package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/http/jfroghttpclient"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

type SystemService struct {
	client     *jfroghttpclient.JfrogHttpClient
	artDetails *auth.ServiceDetails
}

func NewSystemService(artDetails auth.ServiceDetails, client *jfroghttpclient.JfrogHttpClient) *SystemService {
	return &SystemService{artDetails: &artDetails, client: client}
}

func (ss *SystemService) GetArtifactoryDetails() auth.ServiceDetails {
	return *ss.artDetails
}

func (ss *SystemService) GetJfrogHttpClient() *jfroghttpclient.JfrogHttpClient {
	return ss.client
}

func (ss *SystemService) IsDryRun() bool {
	return false
}

func (ss *SystemService) GetVersion() (string, error) {
	httpDetails := (*ss.artDetails).CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet((*ss.artDetails).GetUrl()+"api/system/version", true, &httpDetails)
	if err != nil {
		return "", err
	}

	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}
	var version artifactoryVersion
	err = json.Unmarshal(body, &version)
	if err != nil {
		return "", errorutils.CheckError(err)
	}
	return strings.TrimSpace(version.Version), nil
}

func (ss *SystemService) GetServiceId() (string, error) {
	httpDetails := (*ss.artDetails).CreateHttpClientDetails()
	resp, body, _, err := ss.client.SendGet((*ss.artDetails).GetUrl()+"api/system/service_id", true, &httpDetails)
	if err != nil {
		return "", err
	}
	if err = errorutils.CheckResponseStatus(resp, http.StatusOK); err != nil {
		return "", errorutils.CheckError(errorutils.GenerateResponseError(resp.Status, clientutils.IndentJson(body)))
	}

	return string(body), nil
}

type artifactoryVersion struct {
	Version string `json:"version,omitempty"`
}
