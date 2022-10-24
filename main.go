package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	policyv1 "github.com/cerbos/cerbos/api/genpb/cerbos/policy/v1"
	"github.com/cerbos/cerbos/client"
	yaml "github.com/goccy/go-yaml"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	username      = "cerbos"
	password      = "randomHash"
	clientAddress = "dns:///localhost:3593"
)

const (
	policyTypeResource    string = "resource"
	policyTypePrincipal   string = "principal"
	policyTypeDerivedRole string = "derivedRole"
)

func getKey(p *policyv1.Policy) string {
	var key string

	switch t := p.GetPolicyType().(type) {
	case *policyv1.Policy_ResourcePolicy:
		version := t.ResourcePolicy.Version
		if version == "" {
			version = "default"
		}
		key = policyTypeResource + "." + t.ResourcePolicy.Resource + ".v" + version
		if s := t.ResourcePolicy.Scope; s != "" {
			key += "/" + s
		}
	case *policyv1.Policy_PrincipalPolicy:
		version := t.PrincipalPolicy.Version
		if version == "" {
			version = "default"
		}
		key = policyTypePrincipal + "." + t.PrincipalPolicy.Principal + ".v" + version
		if s := t.PrincipalPolicy.Scope; s != "" {
			key += "/" + s
		}
	case *policyv1.Policy_DerivedRoles:
		key = policyTypeDerivedRole + "." + t.DerivedRoles.Name
	}

	return key
}

type ListPolicyIDResponse struct {
	IDs []string `json:"policyIds"`
}

func listPolicies(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//c, err := client.NewAdminClientWithCredentials("unix:/tmp/sock/cerbos.grpc", username, password, client.WithPlaintext())
	c, err := client.NewAdminClientWithCredentials(clientAddress, username, password, client.WithPlaintext())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	policies, err := c.ListPolicies(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(ListPolicyIDResponse{policies})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

type GetPoliciesResponse struct {
	Policies []string `json:"policies"`
}

func getPolicy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ids, ok := req.URL.Query()["id"]
	if !ok {
		return
	}

	c, err := client.NewAdminClientWithCredentials(clientAddress, username, password, client.WithPlaintext())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	policies, err := c.GetPolicy(context.Background(), ids...)
	if err != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var yamlPolicies []string
	for _, p := range policies {
		jsonBytes, err := protojson.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		yamlBytes, err := yaml.JSONToYAML(jsonBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		yamlPolicies = append(yamlPolicies, string(yamlBytes))
	}

	resp, err := json.Marshal(GetPoliciesResponse{yamlPolicies})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

type CreatePolicyPayload struct {
	PolicyKind string `json:"policyKind"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Scope      string `json:"scope"`
}

type PolicyKeyResponse struct {
	ID string `json:"id"`
}

func createPolicy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params CreatePolicyPayload
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch params.PolicyKind {
	case policyTypeResource, policyTypePrincipal, policyTypeDerivedRole:
	default:
		http.Error(w, "`policyKind` must be one of: "+strings.Join([]string{policyTypeResource, policyTypePrincipal, policyTypeDerivedRole}, ", "), http.StatusBadRequest)
		return
	}

	missingParams := []string{}
	if params.PolicyKind == "" {
		missingParams = append(missingParams, "policyType")
	}

	if params.Name == "" {
		missingParams = append(missingParams, "name")
	}

	if len(missingParams) > 0 {
		http.Error(w, "missing params: "+strings.Join(missingParams, ", "), http.StatusBadRequest)
		return
	}

	if params.Version == "" {
		params.Version = "default"
	}

	ps := client.NewPolicySet()
	p := &policyv1.Policy{
		ApiVersion: "api.cerbos.dev/v1",
	}

	switch params.PolicyKind {
	case policyTypeResource:
		p.PolicyType = &policyv1.Policy_ResourcePolicy{
			ResourcePolicy: &policyv1.ResourcePolicy{
				Resource: params.Name,
				Version:  params.Version,
				Scope:    params.Scope,
			},
		}
		//p := client.NewResourcePolicy(params.Name, params.Version)
		//ps = ps.AddResourcePolicies(p)
	case policyTypePrincipal:
		p.PolicyType = &policyv1.Policy_PrincipalPolicy{
			PrincipalPolicy: &policyv1.PrincipalPolicy{
				Principal: params.Name,
				Version:   params.Version,
				Scope:     params.Scope,
			},
		}
		//p := client.NewPrincipalPolicy(params.Name, params.Version)
		//ps = ps.AddPrincipalPolicies(p)
	case policyTypeDerivedRole:
		p.PolicyType = &policyv1.Policy_DerivedRoles{
			DerivedRoles: &policyv1.DerivedRoles{
				Name: params.Name,
			},
		}
		//dr := client.NewDerivedRoles(params.Name)
		//ps = ps.AddDerivedRoles(dr)
	default:
		http.Error(w, "`policyType` must be one of: "+strings.Join([]string{policyTypeResource, policyTypePrincipal, policyTypeDerivedRole}, ", "), http.StatusBadRequest)
		return
	}

	ps = ps.AddPolicies(p)

	c, err := client.NewAdminClientWithCredentials(clientAddress, username, password, client.WithPlaintext())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.AddOrUpdatePolicy(context.Background(), ps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(PolicyKeyResponse{getKey(p)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

type UpdatePolicyPayload struct {
	ID     string `json:"id"`
	Policy string `json:"policy"`
}

func loadAndValidatePolicy(params UpdatePolicyPayload) (*client.PolicySet, error) {
	ps := client.NewPolicySet()

	ps.AddPolicyFromReader(strings.NewReader(params.Policy))

	if err := ps.Validate(); err != nil {
		return ps, err
	}

	return ps, nil
}

func validatePolicy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params UpdatePolicyPayload
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := loadAndValidatePolicy(params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func updatePolicy(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params UpdatePolicyPayload
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ps, err := loadAndValidatePolicy(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := client.NewAdminClientWithCredentials(clientAddress, username, password, client.WithPlaintext())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.AddOrUpdatePolicy(context.Background(), ps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p := ps.GetPolicies()[0]

	resp, err := json.Marshal(PolicyKeyResponse{getKey(p)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)

}

func getAuditLog(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	c, err := client.NewAdminClientWithCredentials(clientAddress, username, password, client.WithPlaintext())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gen, err := c.AuditLogs(context.Background(), client.AuditLogOptions{
		Type: client.AccessLogs,
		Tail: 100,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for al := range gen {
		l, err := al.AccessLog()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, l)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./client/dist")))

	http.HandleFunc("/policies", listPolicies)
	http.HandleFunc("/policy", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			getPolicy(w, req)
		case "POST":
			createPolicy(w, req)
		case "PATCH":
			updatePolicy(w, req)
		default:
			http.Error(w, fmt.Sprintf("Method not supported: %s", req.Method), http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/validate", validatePolicy)
	http.HandleFunc("/auditlog", getAuditLog)

	http.ListenAndServe(":8090", nil)
}
