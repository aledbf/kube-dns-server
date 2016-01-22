package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"text/template"

	"github.com/golang/glog"
	k8sexec "k8s.io/kubernetes/pkg/util/exec"
)

const (
	nsdTmpl = `
`
)

type record struct {
	name string
	ip   net.IP
}

type nsd struct {
	ns []record
	a  []record
}

func (k *nsd) WriteCfg(svcs []vip) error {
	w, err := os.Create("/etc/nsd/nsd.conf.")
	if err != nil {
		return err
	}
	defer w.Close()

	t, err := template.New("nsd").Parse(nsdTmpl)
	if err != nil {
		return err
	}

	conf := make(map[string]interface{})
	conf["ns"] = k.iface
	conf["a"] = k.ip

	b, _ := json.Marshal(conf)
	glog.Infof("%v", string(b))

	return t.Execute(w, conf)
}

func (k *nsd) Start() {
	cmd := exec.Command("/usr/sbin/nsd",
		"-d",
		"-P", "/nsd.pid")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		glog.Errorf("nsd error: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		glog.Fatalf("nsd error: %v", err)
	}
}

func (k *nsd) Reload() error {
	glog.Info("reloading nsd server")
	_, err := k8sexec.New().Command("killall", "-1", "nsd").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error reloading nsd: %v", err)
	}

	return nil
}
