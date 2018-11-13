package godo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// DiskConfig is used when creating a disk
type DiskConfig struct {
	MachineID   int    `json:"machineId"`
	DiskName    string `json:"diskName"`
	Description string `json:"description"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
	SSDSize     int    `json:"ssdSize"`
	IOPS        int    `json:"iops"`
	DiskID      int    `json:"diskId"`
	Detach      bool   `json:"detach"`
	Permanently string `json:"permanently"`
}

// DiskDeleteConfig is used when deleting a disk
type DiskDeleteConfig struct {
	DiskID      int    `json:"diskId"`
	Detach      bool   `json:"detach"`
	Permanently string `json:"permanently"`
}

// DiskInfo contains all information related to a disk
type DiskInfo struct {
	ReferenceID         string        `json:"referenceId"`
	DiskPath            string        `json:"diskPath"`
	Images              []interface{} `json:"images"`
	GUID                int           `json:"guid"`
	ID                  int           `json:"id"`
	AccountID           int           `json:"accountId"`
	SizeUsed            int           `json:"sizeUsed"`
	Descr               string        `json:"descr"`
	Gid                 int           `json:"gid"`
	Role                string        `json:"role"`
	Params              string        `json:"params"`
	Type                string        `json:"type"`
	Status              string        `json:"status"`
	RealityDeviceNumber int           `json:"realityDeviceNumber"`
	Passwd              string        `json:"passwd"`
	Iotune              struct {
		TotalIopsSec int `json:"total_iops_sec"`
	} `json:"iotune"`
	Name    string        `json:"name"`
	SizeMax int           `json:"sizeMax"`
	Meta    []interface{} `json:"_meta"`
	ACL     struct {
	} `json:"acl"`
	Iqn           string `json:"iqn"`
	BootPartition int    `json:"bootPartition"`
	Login         string `json:"login"`
	Order         int    `json:"order"`
	Ckey          string `json:"_ckey"`
}

// DiskList is a list of disks
// Returned when using the List method
type DiskList []struct {
	Username    interface{} `json:"username"`
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Size        int         `json:"size"`
	Type        string      `json:"type"`
	ID          int         `json:"id"`
	AccountID   int         `json:"accountId"`
}

// DiskService is an interface for interfacing with the Disk
// endpoints of the OVC API
// See: https://ch-lug-dc01-001.gig.tech/g8vdc/#/ApiDocs
type DiskService interface {
	List(int) (*DiskList, error)
	Get(string) (*DiskInfo, error)
	Create(*DiskConfig) error
	Delete(*DiskDeleteConfig) error
}

// DiskServiceOp handles communication with the machine related methods of the
// OVC API
type DiskServiceOp struct {
	client *OvcClient
}

var _ DiskService = &DiskServiceOp{}

// List all disks
func (s *DiskServiceOp) List(accountID int) (*DiskList, error) {
	diskMap := make(map[string]interface{})
	diskMap["accountId"] = accountID
	diskJSON, err := json.Marshal(diskMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/disks/list", bytes.NewBuffer(diskJSON))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	var disks = new(DiskList)
	err = json.Unmarshal(body, &disks)
	if err != nil {
		return nil, err
	}
	return disks, nil
}

// Create a new Disk
func (s *DiskServiceOp) Create(diskConfig *DiskConfig) error {
	diskConfigJSON, err := json.Marshal(*diskConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/machines/addDisk", bytes.NewBuffer(diskConfigJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Delete an existing Disk
func (s *DiskServiceOp) Delete(diskConfig *DiskDeleteConfig) error {
	diskConfigJSON, err := json.Marshal(*diskConfig)
	if err != nil {
		return err
	}
	log.Println(string(diskConfigJSON))
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/disks/delete", bytes.NewBuffer(diskConfigJSON))
	if err != nil {
		return err
	}
	_, err = s.client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// Get individual Disk
func (s *DiskServiceOp) Get(diskID string) (*DiskInfo, error) {
	diskIDMap := make(map[string]interface{})
	diskIDInt, err := strconv.Atoi(diskID)
	if err != nil {
		return nil, err
	}
	diskIDMap["diskId"] = diskIDInt
	diskIDJson, err := json.Marshal(diskIDMap)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.client.ServerURL+"/cloudapi/disks/get", bytes.NewBuffer(diskIDJson))
	if err != nil {
		return nil, err
	}
	body, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	var diskInfo = new(DiskInfo)
	err = json.Unmarshal(body, &diskInfo)
	if err != nil {
		return nil, err
	}
	return diskInfo, nil

}
