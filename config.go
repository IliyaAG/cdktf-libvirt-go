package main

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type VmConfig struct {
    Hostname  string `yaml:"hostname"`
    Memory    int    `yaml:"memory"`
    Vcpu      int    `yaml:"vcpu"`
    DiskSize  int    `yaml:"diskSize"`
    IpAddress string `yaml:"ipAddress"`
}

type MyStackConfig struct {
    Image   string              `yaml:"image"`
    SshKeys []string            `yaml:"sshKeys"`
    Count   int                 `yaml:"count"`
    Vms     map[string]VmConfig `yaml:"vms"`
}

func SaveConfig(path string, cfg *MyStackConfig) error {
    data, err := yaml.Marshal(cfg)
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

func LoadConfig(path string) (*MyStackConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg MyStackConfig
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    if cfg.Count > 0 {
        cfg.Vms = make(map[string]VmConfig)
        for i := 1; i <= cfg.Count; i++ {
            cfg.Vms[fmt.Sprintf("vm%d", i)] = VmConfig{
                Hostname:  fmt.Sprintf("vm%d", i),
                Memory:    2048,
                Vcpu:      2,
                DiskSize:  20,
                IpAddress: fmt.Sprintf("192.168.122.%d", 100+i),
            }
        }
        cfg.Count = 0
        if err := SaveConfig(path, &cfg); err != nil {
            return nil, err
        }
    }

    for key, vm := range cfg.Vms {
        if vm.Hostname == "" {
            vm.Hostname = key
        }
        if vm.Memory == 0 {
            vm.Memory = 2048
        }
        if vm.Vcpu == 0 {
            vm.Vcpu = 2
        }
        if vm.DiskSize == 0 {
            vm.DiskSize = 20
        }
        if vm.IpAddress == "" {
            vm.IpAddress = fmt.Sprintf("192.168.122.%d", 100+len(cfg.Vms))
        }
        cfg.Vms[key] = vm
    }

    return &cfg, nil
}
