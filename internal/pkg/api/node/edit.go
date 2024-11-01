package apinode

import (
	"os"

	"github.com/pkg/errors"
	"github.com/warewulf/warewulf/internal/pkg/api/routes/wwapiv1"
	"github.com/warewulf/warewulf/internal/pkg/node"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
	"gopkg.in/yaml.v3"
)

/*
Returns filtered list of nodes
*/
func FilteredNodes(nodeList *wwapiv1.NodeList) *wwapiv1.NodeYaml {
	nodeDB, err := node.New()
	if err != nil {
		wwlog.Error("Could not open nodeDB: %s\n", err)
		os.Exit(1)
	}
	nodeMap, _ := nodeDB.FindAllNodes()
	nodeMap = node.FilterNodeListByName(nodeMap, nodeList.Output)
	buffer, _ := yaml.Marshal(nodeMap)
	retVal := wwapiv1.NodeYaml{
		NodeConfMapYaml: string(buffer),
		Hash:            nodeDB.StringHash(),
	}
	return &retVal
}

/*
Add nodes from yaml
*/
func NodeAddFromYaml(nodeList *wwapiv1.NodeYaml) (err error) {
	nodeDB, err := node.New()
	if err != nil {
		return errors.Wrap(err, "Could not open NodeDB: %s\n")
	}
	nodeMap := make(map[string]*node.NodeConf)
	err = yaml.Unmarshal([]byte(nodeList.NodeConfMapYaml), nodeMap)
	if err != nil {
		return errors.Wrap(err, "Could not unmarshal Yaml: %s\n")
	}
	for nodeName, node := range nodeMap {
		err = nodeDB.SetNode(nodeName, *node)
		if err != nil {
			return errors.Wrap(err, "couldn't set node")
		}
	}
	err = nodeDB.Persist()
	if err != nil {
		return errors.Wrap(err, "failed to persist nodedb")
	}
	return nil
}
