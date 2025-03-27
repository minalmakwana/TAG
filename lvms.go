package main

import (
	"fmt"
	"errors"
)

// LogicalVolume represents a logical volume in the system
type LogicalVolume struct {
	Name       string
	Size       int // Size in megabytes
	IsAllocated bool
}

// VolumeGroup represents a volume group that contains logical volumes
type VolumeGroup struct {
	Name             string
	TotalSize        int    // Total size in megabytes
	FreeSize         int    // Free size in megabytes
	LogicalVolumes   []LogicalVolume
}

// CreateVolumeGroup creates a new volume group
func CreateVolumeGroup(name string, totalSize int) VolumeGroup {
	return VolumeGroup{
		Name:             name,
		TotalSize:        totalSize,
		FreeSize:         totalSize,
		LogicalVolumes:   []LogicalVolume{},
	}
}

// CreateLogicalVolume creates a new logical volume if there is enough free space
func (vg *VolumeGroup) CreateLogicalVolume(name string, size int) error {
	if size > vg.FreeSize {
		return errors.New("not enough free space to create logical volume")
	}
	vg.LogicalVolumes = append(vg.LogicalVolumes, LogicalVolume{
		Name:       name,
		Size:       size,
		IsAllocated: true,
	})
	vg.FreeSize -= size
	return nil
}

// DeleteLogicalVolume deletes an existing logical volume and frees up space
func (vg *VolumeGroup) DeleteLogicalVolume(name string) error {
	for i, lv := range vg.LogicalVolumes {
		if lv.Name == name {
			// Remove the logical volume
			vg.LogicalVolumes = append(vg.LogicalVolumes[:i], vg.LogicalVolumes[i+1:]...)
			vg.FreeSize += lv.Size
			return nil
		}
	}
	return errors.New("logical volume not found")
}

// ListVolumes lists all logical volumes in the volume group
func (vg *VolumeGroup) ListVolumes() {
	fmt.Printf("Logical Volumes in %s:\n", vg.Name)
	for _, lv := range vg.LogicalVolumes {
		status := "Not Allocated"
		if lv.IsAllocated {
			status = "Allocated"
		}
		fmt.Printf("Name: %s, Size: %d MB, Status: %s\n", lv.Name, lv.Size, status)
	}
}

func main() {
	// Create a new Volume Group with 1000MB of total space
	vg := CreateVolumeGroup("VG1", 1000)

	// Try to create logical volumes
	if err := vg.CreateLogicalVolume("LV1", 200); err != nil {
		fmt.Println("Error creating LV1:", err)
	}

	if err := vg.CreateLogicalVolume("LV2", 300); err != nil {
		fmt.Println("Error creating LV2:", err)
	}

	// List volumes
	vg.ListVolumes()

	// Try to delete a logical volume
	if err := vg.DeleteLogicalVolume("LV1"); err != nil {
		fmt.Println("Error deleting LV1:", err)
	}

	// List volumes again after deletion
	vg.ListVolumes()
}

