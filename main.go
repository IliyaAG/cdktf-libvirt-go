package main

import (
    "cdk.tf/go/stack/generated/dmacvicar/libvirt/domain"
    "cdk.tf/go/stack/generated/dmacvicar/libvirt/volume"
    "cdk.tf/go/stack/generated/dmacvicar/libvirt/cloudinitdisk"
    "github.com/hashicorp/terraform-cdk-go/cdktf"
    "github.com/aws/constructs-go/constructs/v10"
)

type MyStackConfig struct {
    cdktf.TerraformStack
}

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
    stack := cdktf.NewTerraformStack(scope, &id)

    provider := domain.NewLibvirtProvider(stack, cdktf.String("libvirt"), &domain.LibvirtProviderConfig{
        Uri: cdktf.String("qemu:///system"),
    })

    vol := volume.NewVolume(stack, cdktf.String("ubuntu_qcow2"), &volume.VolumeConfig{
        Name:   cdktf.String("vm1-ubuntu-disk.qcow2"),
        Source: cdktf.String("iso-images/jammy-server-cloudimg-amd64.img"),
        Format: cdktf.String("qcow2"),
        Size:   cdktf.Number(10737418240),
        Provider: provider,
    })

    ci := cloudinitdisk.NewCloudinitDisk(stack, cdktf.String("commoninit"), &cloudinitdisk.CloudinitDiskConfig{
        Name:          cdktf.String("vm1-commoninit.iso"),
        UserData:      cdktf.String(stringFromFile("config/cloud_init.cfg")),
        NetworkConfig: cdktf.String(stringFromFile("config/network_config.yml")),
        Provider:      provider,
    })

    domain.NewDomain(stack, cdktf.String("domain_ubuntu"), &domain.DomainConfig{
        Name:   cdktf.String("vm1"),
        Memory: cdktf.Number(2048),
        Vcpu:   cdktf.Number(2),
        Cloudinit: ci.Id(),
        NetworkInterface: []domain.DomainNetworkInterface{
            {
                NetworkName:   cdktf.String("default"),
                WaitForLease:  cdktf.Bool(true),
                Hostname:      cdktf.String("vm1"),
            },
        },
        Disk: []domain.DomainDisk{
            {
                VolumeId: vol.Id(),
            },
        },
        Graphics: []domain.DomainGraphics{
            {
                Type:       cdktf.String("spice"),
                ListenType: cdktf.String("address"),
                Autoport:   cdktf.Bool(true),
            },
        },
    })

    return stack
}

func stringFromFile(path string) string {
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }
    return string(data)
}

func main() {
    app := cdktf.NewApp(nil)

    NewMyStack(app, "libvirt-vms")

    app.Synth()
}
