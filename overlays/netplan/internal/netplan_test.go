package netplan 

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/app/wwctl/overlay/show"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_wickedOverlay(t *testing.T) {
	env := testenv.New(t)
	defer env.RemoveAll(t)
	env.ImportFile(t, "etc/warewulf/nodes.conf", "nodes.conf")
	env.ImportFile(t, "var/lib/warewulf/overlays/netplan/rootfs/etc/netplan/01-netcfg.yaml.ww", "../rootfs/etc/netplan/01-netcfg.yaml.ww")

	tests := []struct {
		name string
		args []string
		log  string
	}{
		{
			name: "netplan",
			args: []string{"--render", "node1", "netplan", "etc/netplan/01-netcfg.yaml.ww"},
			log:  netplan,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := show.GetCommand()
			cmd.SetArgs(tt.args)
			stdout := bytes.NewBufferString("")
			stderr := bytes.NewBufferString("")
			logbuf := bytes.NewBufferString("")
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			wwlog.SetLogWriter(logbuf)
			err := cmd.Execute()
			assert.NoError(t, err)
			assert.Empty(t, stdout.String())
			assert.Empty(t, stderr.String())
			assert.Equal(t, tt.log, logbuf.String())
		})
	}
}

const netplan string = `backupFile: true
writeFile: true
Filename: default
# This file is autogenerated by warewulf
network:
  version: 2
  renderer: networkd
  ethernets:
         wwnet0:
         addresses:
         - 192.168.3.21/255.255.255.0
         
backupFile: true
writeFile: true
Filename: secondary
# This file is autogenerated by warewulf
network:
  version: 2
  renderer: networkd
  ethernets:
         wwnet1:
         addresses:
         - 192.168.3.22/255.255.255.0
         mtu: 9000
`