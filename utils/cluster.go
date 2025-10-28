package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"net/http"
	"strings"
	"time"
)

func uniqueKey(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}

func CloneNodes(nodes []models.ClusterNodes) []models.ClusterNodes {
	out := make([]models.ClusterNodes, len(nodes))
	copy(out, nodes)
	return out
}

func EnsureUniqueNodes(nodes []models.ClusterNodes) []models.ClusterNodes {
	seen := make(map[string]struct{}, len(nodes))
	out := make([]models.ClusterNodes, 0, len(nodes))
	for _, n := range nodes {
		key := uniqueKey(n.Host, n.Port)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, n)
	}
	return out
}

func RemoveNode(nodes []models.ClusterNodes, host, port string) []models.ClusterNodes {
	key := uniqueKey(host, port)
	out := make([]models.ClusterNodes, 0, len(nodes))
	for _, n := range nodes {
		if uniqueKey(n.Host, n.Port) == key {
			continue
		}
		out = append(out, n)
	}
	return out
}

func BroadcastCluster(nodes []models.ClusterNodes) error {
	if len(nodes) == 0 {
		return nil
	}

	payload, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 4 * time.Second}
	var errorsList []string

	for _, n := range nodes {
		addr := BuildHTTPAddr(n.Host, n.Port)
		req, err := http.NewRequest(http.MethodPost, addr+"/cluster/fetch_nodes", bytes.NewReader(payload))
		if err != nil {
			errorsList = append(errorsList, err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			errorsList = append(errorsList, fmt.Sprintf("%s -> %v", addr, err))
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 400 {
			errorsList = append(errorsList, fmt.Sprintf("%s -> status %d", addr, resp.StatusCode))
		}
	}

	if len(errorsList) > 0 {
		return fmt.Errorf("broadcast errors: %s", strings.Join(errorsList, "; "))
	}
	return nil
}
