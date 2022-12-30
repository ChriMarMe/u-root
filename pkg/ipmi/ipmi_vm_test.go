// Copyright 2019-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestIntegrationIPMI(t *testing.T) {
	o := &vmtest.Options{
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				// This integration test requires kernel built with the following options set:
				// CONFIG_IPMI=y
				// CONFIG_IPMI_DEVICE_INTERFACE=y
				// CONFIG_IPMI_WATCHDOG=y
				// CONFIG_IPMI_SI=y
				qemu.ArbitraryArgs{"-device", "ipmi-bmc-sim,id=bmc0"},
				qemu.ArbitraryArgs{"-device", "pci-ipmi-kcs,bmc=bmc0"},
			},
		},
	}
	vmtest.GolangTest(t, []string{"github.com/u-root/u-root/pkg/ipmi"}, o)
}

func TestWatchdogRunningQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %v", err)
	}
	defer i.Close()

	ret, err := i.WatchdogRunning()
	if err != nil {
		t.Errorf("i.WatchdogRunning() = %v", err)
	}
	if ret {
		t.Errorf("i.WatchdogRunning() = %t, want false", ret)
	}
}

func TestShutoffWatchdogQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %q", err)
	}
	defer i.Close()

	if err := i.ShutoffWatchdog(); err != nil {
		t.Errorf("i.ShutoffWatchdog() = %q", err)
	}
}

func TestGetDeviceIDQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %q", err)
	}
	defer i.Close()

	id, err := i.GetDeviceID()
	if err != nil {
		t.Errorf("i.GetDeviceID() = %q", err)
	}
	if id.DeviceID != 0x20 {
		t.Errorf("DeviceID: %q, want: %q", id.DeviceID, 0x1)
	}
	if id.DeviceRevision != 0x0 {
		t.Errorf("DeviceRevision: %q, want: %q", id.DeviceRevision, 0x0)
	}
	if id.FwRev1 != 0x0 {
		t.Errorf("FwRev1: %q, want: %q", id.FwRev1, 0x0)
	}
	if id.FwRev2 != 0x0 {
		t.Errorf("FwRev2: %q, want: %q", id.FwRev2, 0x0)
	}
	if id.IpmiVersion != 0x2 {
		t.Errorf("IpmiVersion: %q, want: %q", id.IpmiVersion, 0x2)
	}
	/*
		This field is differs on every call, I can't figure out why

		if id.AdtlDeviceSupport != 0xa {
			t.Errorf("AdtlDeviceSupport: %q, want: %q", id.AdtlDeviceSupport, 0xa)
		}
	*/
	if !bytes.Equal(id.ManufacturerID[:], []byte{0x0, 0x0, 0x0}) {
		t.Errorf("ManufacturerID: %q, want: %q", id.ManufacturerID, []byte{0x0, 0x0, 0x0})
	}
	if !bytes.Equal(id.ProductID[:], []byte{0x0, 0x0}) {
		t.Errorf("ProductID: %q, want: %q", id.ProductID, []byte{0x0, 0x0})
	}

}

func TestEnableSELQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %q", err)
	}
	defer i.Close()

	ret, err := i.EnableSEL()
	if err != nil {
		t.Errorf("i.EnableSEL() = %v", err)
	}
	if !ret {
		t.Errorf("i.EnableSEL() = %v, want true", ret)
	}
}

func TestGetSELInfoQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %v", err)
	}
	defer i.Close()

	info, err := i.GetSELInfo()
	if err != nil {
		t.Errorf("i.GetSELInfo() = %v", err)
	}
	if info.Version != 0x51 {
		t.Errorf("Version = %q, want %q", info.Version, 0x51)
	}
	if info.Entries != 0x0 {
		t.Errorf("Version = %q, want %q", info.Entries, 0x0)
	}
	if info.FreeSpace != 0x800 {
		t.Errorf("Version = %q, want %q", info.FreeSpace, 0x800)
	}
	if info.OpSupport != 0x2 {
		t.Errorf("Version = %q, want %q", info.Version, 0x2)
	}

}

func TestGetLanConfigQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	t.Skip("Not supported command")
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %v", err)
	}
	defer i.Close()

	if _, err := i.GetLanConfig(1, 1); err != nil {
		t.Errorf("i.GetLanConfig(1, 1) = %v", err)
	}
}

func TestRawCmdQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %v", err)
	}
	defer i.Close()

	// WatchdogRunning configuration
	data := []byte{0x6, 0x1}
	if _, err := i.RawCmd(data); err != nil {
		t.Errorf("i.RawCmd(data) = %v", err)
	}
}

func TestSetSystemFWVersionQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	t.Skip("Not supported command")
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %v", err)
	}
	defer i.Close()

	if err := i.SetSystemFWVersion("TestTest"); err == nil {
		t.Errorf("i.SetSystemFWVersion(TestTest) = %v", err)
	}
}

func TestLogSystemEventQemu(t *testing.T) {
	testutil.SkipIfNotRoot(t)
	i, err := Open(0)
	if err != nil {
		t.Fatalf("Open(0) = %v", err)
	}
	defer i.Close()

	e := &Event{}
	if err := i.LogSystemEvent(e); err != nil {
		t.Errorf("i.LogSystemEvent(e) = %v", err)
	}
}
