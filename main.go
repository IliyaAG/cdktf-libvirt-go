package main

import (
    "fmt"

    cdktf "github.com/hashicorp/terraform-cdk-go/cdktf"
    libvirt "github.com/dmacvicar/terraform-provider-libvirt/cdktf/provider/libvirt"
    constructs "github.com/hashicorp/constructs-go/constructs/v10"
)

type MyStack struct {
    cdktf.TerraformStack
}

func NewMyStack(scope constructs.Construct, id string, cfg *MyStackConfig) cdktf.TerraformStack {
    stack := cdktf.NewTerraformStack(scope, &id)

    libvirt.NewLibvirtProvider(stack, jsii.String("libvirt"), &libvirt.LibvirtProviderConfig{
        Uri: jsii.String("qemu:///system"),
    })

    for key, vm := range cfg.Vms {
        libvirt.NewVolume(stack, jsii.String(fmt.Sprintf("%s-disk", key)), &libvirt.VolumeConfig{
            Name:   jsii.String(fmt.Sprintf("%s-ubuntu-disk.qcow2", key)),
            Source: jsii.String(cfg.Image),
            Format: jsii.String("qcow2"),
        })

        cloudinit := libvirt.NewCloudInitDisk(stack, jsii.String(fmt.Sprintf("%s-cloudinit", key)), &libvirt.CloudInitDiskConfig{
            Name: jsii.String(fmt.Sprintf("%s-commoninit.iso", key)),
            UserData: cdktf.Fn_Templatefile(jsii.String("templates/cloud_init.cfg.tpl"), &map[string]interface{}{
                "sshKeys": cfg.SshKeys,
            }),
            NetworkConfig: cdktf.Fn_Templatefile(jsii.String("templates/network_config.yml.tpl"), &map[string]interface{}{
                "ip": vm.IpAddress,
            }),
        })

        libvirt.NewDomain(stack, jsii.String(fmt.Sprintf("%s-domain", key)), &libvirt.DomainConfig{
            Name:   jsii.String(vm.Hostname),
            Memory: jsii.Number(float64(vm.Memory)),
            Vcpu:   jsii.Number(float64(vm.Vcpu)),
            Cloudinit: cloudinit.Id(),

            NetworkInterface: []*libvirt.DomainNetworkInterface{{
                NetworkName:   jsii.String("default"),
                WaitForLease:  jsii.Bool(true),
                Hostname:      jsii.String(vm.Hostname),
            }},

            Disk: []*libvirt.DomainDisk{{
                VolumeId: jsii.String(fmt.Sprintf("%s-ubuntu-disk.qcow2", key)),
            }},
        })
    }

    return stack
}

func main() {
    cfg, err := LoadConfig("config.yaml")
    if err != nil {
        panic(err)
    }

    app := cdktf.NewApp(nil)
    NewMyStack(app, "cdktf-libvirt-go", cfg)
    app.Synth()
}
