// This file is part of MinIO Operator
// Copyright (c) 2024 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package sidecar

import (
	"crypto/sha256"
	"encoding/hex"
	miniov2 "github.com/minio/operator/pkg/apis/minio.min.io/v2"
	"github.com/minio/operator/pkg/configuration"
	"github.com/minio/operator/sidecar/pkg/validator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/minio/operator/pkg/common"
)

func configureSidecarServer(c *Controller) *http.Server {
	router := mux.NewRouter().SkipClean(true).UseEncodedPath()

	router.Methods(http.MethodPost).
		Path(common.SidecarAPIConfigEndpoint).
		HandlerFunc(c.CheckConfigHandler).
		Queries(restQueries("c")...)

	router.NotFoundHandler = http.NotFoundHandler()

	s := &http.Server{
		Addr:           "0.0.0.0:" + common.SidecarHTTPPort,
		Handler:        router,
		ReadTimeout:    time.Minute,
		WriteTimeout:   time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}

// CheckConfigHandler - POST /sidecar/v1/config?c={hash}
//
// Checks the configuration hash and regenerates the configuration
// if it does not match, it does the following:
//
// 1. check if the current local stored configuration matches the hash
// 2. if not try to re-generate the configuration using listers
// 3. if not try to re-generate the configuration using live api results
func (c *Controller) CheckConfigHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["c"]
	log.Println("Checking config hash: ", hash)
	// check if the current local stored configuration matches the hash
	_, _, tmpFileContents, err := validator.ReadTmpConfig(miniov2.CfgFile)
	if err != nil {
		log.Println("Error reading tmp config: ", err)
		http.Error(w, "Error reading tmp config", http.StatusInternalServerError)
		return
	}
	hasher := sha256.New()
	hasher.Write([]byte(tmpFileContents))
	cfgHash := hex.EncodeToString(hasher.Sum(nil))
	if cfgHash == hash {
		// no changes to local config
		log.Println("Configuration matches, case 1")
		w.WriteHeader(http.StatusOK)
		return
	}

	// if not try to re-generate the configuration using listers
	tenant, err := c.tenantInformer.Lister().Tenants(c.namespace).Get(c.tenantName)
	if err != nil {
		log.Println("Error getting tenant: ", err)
		http.Error(w, "Error getting tenant", http.StatusInternalServerError)
		return
	}
	if tenant.Spec.Configuration == nil {
		log.Println("Tenant configuration not found")
		http.Error(w, "Tenant configuration not found", http.StatusInternalServerError)
		return
	}
	tenantSecret, err := c.secretInformer.Lister().Secrets(c.namespace).Get(tenant.Spec.Configuration.Name)
	if err != nil {
		log.Println("Error getting tenant secret: ", err)
		http.Error(w, "Error getting tenant secret", http.StatusInternalServerError)
		return
	}
	cfg, rootUserFound, rootPwdFound := configuration.GetFullTenantConfig(tenant, tenantSecret)
	if !rootUserFound || !rootPwdFound {
		log.Println("Root user or password not found in tmp config")
		http.Error(w, "Root user or password not found in tmp config", http.StatusInternalServerError)
		return
	}
	hasherCached := sha256.New()
	hasherCached.Write([]byte(cfg))
	cachedCfgHash := hex.EncodeToString(hasherCached.Sum(nil))
	if cachedCfgHash == hash {
		log.Println("Configuration matches case 2")
		err = os.WriteFile(miniov2.CfgFile, []byte(cfg), 0o644)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error writing config file", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	// if not try to re-generate the configuration using live api results
	tenant, err = c.controllerClient.MinioV2().Tenants(c.namespace).Get(r.Context(), c.tenantName, metav1.GetOptions{})
	if err != nil {
		log.Println("Error getting tenant: ", err)
		http.Error(w, "Error getting tenant", http.StatusInternalServerError)
		return
	}
	if tenant.Spec.Configuration == nil {
		log.Println("Tenant configuration not found")
		http.Error(w, "Tenant configuration not found", http.StatusInternalServerError)
		return
	}
	tenantSecret, err = c.kubeClient.CoreV1().Secrets(c.namespace).Get(r.Context(), tenant.Spec.Configuration.Name, metav1.GetOptions{})
	if err != nil {
		log.Println("Error getting tenant secret: ", err)
		http.Error(w, "Error getting tenant secret", http.StatusInternalServerError)
		return
	}
	cfg, rootUserFound, rootPwdFound = configuration.GetFullTenantConfig(tenant, tenantSecret)
	if !rootUserFound || !rootPwdFound {
		log.Println("Root user or password not found in tmp config")
		http.Error(w, "Root user or password not found in tmp config", http.StatusInternalServerError)
		return
	}
	hasherLive := sha256.New()
	hasherLive.Write([]byte(cfg))
	liveCfgHash := hex.EncodeToString(hasherLive.Sum(nil))
	if liveCfgHash != hash {
		// if the configuration still does not match, the hash we were asked about, makes no sense
		log.Println("Configuration does not match, case 3")
		http.Error(w, "Configuration does not match", http.StatusInternalServerError)
		return
	}
	// store the newly generated config
	err = os.WriteFile(miniov2.CfgFile, []byte(cfg), 0o644)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error writing config file", http.StatusInternalServerError)
		return
	}
	log.Println("Configuration regenerated")
	w.WriteHeader(http.StatusOK)
}
