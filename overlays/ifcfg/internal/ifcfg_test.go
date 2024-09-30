package ifcfg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/app/wwctl/overlay/show"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_ifcfgOverlay(t *testing.T) {
	env := testenv.New(t)
	defer env.RemoveAll(t)
	env.ImportFile(t, "etc/warewulf/nodes.conf", "nodes.conf")
	env.ImportFile(t, "var/lib/warewulf/overlays/ifcfg/rootfs/etc/sysconfig/network-scripts/ifcfg.ww", "../rootfs/etc/sysconfig/network-scripts/ifcfg.ww")
	env.ImportFile(t, "var/lib/warewulf/overlays/ifcfg/rootfs/etc/sysconfig/network.ww", "../rootfs/etc/sysconfig/network.ww")

	tests := []struct {
		name string
		args []string
		log  string
	}{
		{
			name: "ifcfg:ifcfg.ww",
			args: []string{"--render", "node1", "ifcfg", "etc/sysconfig/network-scripts/ifcfg.ww"},
			log:  ifcfg,
		},
		{
			name: "ifcfg:network.ww",
			args: []string{"--render", "node1", "ifcfg", "etc/sysconfig/network.ww"},
			log:  ifcfg_network,
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

const ifcfg string = `backupFile: true
writeFile: true
Filename: ifcfg-default.conf
# This file is autogenerated by warewulf
TYPE=ethernet
DEVICE=wwnet0
NAME=default
BOOTPROTO=static
DEVTIMEOUT=10
IPADDR=192.168.3.21
NETMASK=255.255.255.0
GATEWAY=192.168.3.1
HWADDR=e6:92:39:49:7b:03
TYPE=ethernet
ONBOOT=true
IPV6INIT=yes
IPV6_AUTOCONF=yes
IPV6_DEFROUTE=yes
IPV6_FAILURE_FATAL=no
backupFile: true
writeFile: true
Filename: ifcfg-secondary.conf
# This file is autogenerated by warewulf
TYPE=ethernet
DEVICE=wwnet1
NAME=secondary
BOOTPROTO=static
DEVTIMEOUT=10
IPADDR=192.168.3.22
NETMASK=255.255.255.0
GATEWAY=192.168.3.1
HWADDR=9a:77:29:73:14:f1
TYPE=ethernet
ONBOOT=true
IPV6INIT=yes
IPV6_AUTOCONF=yes
IPV6_DEFROUTE=yes
IPV6_FAILURE_FATAL=no
DNS1=8.8.8.8
DNS2=8.8.4.4
`

const ifcfg_network string = `backupFile: true
writeFile: true
Filename: etc/sysconfig/network
NETWORKING=yes
HOSTNAME=node1
`