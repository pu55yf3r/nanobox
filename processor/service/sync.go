package service

import (
	"fmt"
	"regexp"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

type (

	// service ...
	service struct {
		label string
		name  string
		image string
	}

	// processServiceSync ...
	processServiceSync struct {
		control    processor.ProcessControl
		fail       bool
		newBoxfile models.Boxfile
		oldBoxfile models.Boxfile
	}
)

//
func init() {
	processor.Register("service_sync", serviceSyncFunc)
}

//
func serviceSyncFunc(control processor.ProcessControl) (processor.Processor, error) {
	return &processServiceSync{control: control}, nil
}

//
func (serviceSync processServiceSync) Results() processor.ProcessControl {
	return serviceSync.control
}

//
func (serviceSync *processServiceSync) Process() error {

	if err := serviceSync.loadNewBoxfile(); err != nil {
		return err
	}

	if err := serviceSync.loadOldBoxfile(); err != nil {
		return err
	}

	if err := serviceSync.purgeDeltaServices(); err != nil {
		return err
	}

	if err := serviceSync.provisionDataServices(); err != nil {
		return err
	}

	if err := serviceSync.replaceOldBoxfile(); err != nil {
		return err
	}

	return nil
}

// loadNewBoxfile loads the new build boxfile from the database
func (serviceSync *processServiceSync) loadNewBoxfile() error {

	if err := data.Get(config.AppName()+"_meta", "build_boxfile", &serviceSync.newBoxfile); err != nil {
		return err
	}

	return nil
}

// loadOldBoxfile loads the old boxfile from the database
func (serviceSync *processServiceSync) loadOldBoxfile() error {

	// we don't care about the error here because this could be the first build
	data.Get(config.AppName()+"_meta", "old_build_boxfile", &serviceSync.oldBoxfile)

	return nil
}

// replaceOldBoxfile replaces the old boxfile in the database with the new boxfile
func (serviceSync *processServiceSync) replaceOldBoxfile() error {

	if err := data.Put(config.AppName()+"_meta", "old_build_boxfile", &serviceSync.newBoxfile); err != nil {
		return err
	}

	return nil
}

// purgeDeltaServices will purge the services that were removed from the boxfile
func (serviceSync *processServiceSync) purgeDeltaServices() error {

	// convert the data into boxfile library helpers
	oldBoxfile := boxfile.New(serviceSync.oldBoxfile.Data)
	newBoxfile := boxfile.New(serviceSync.newBoxfile.Data)

	// fetch the services
	uids, err := data.Keys(config.AppName())
	if err != nil {
		return err
	}

	for _, uid := range uids {

		// ignore platform services
		if isPlatformUID(uid) {
			continue
		}

		// fetch the nodes
		newNode := newBoxfile.Node(uid)
		oldNode := oldBoxfile.Node(uid)

		// skip if the nodes are the same
		if newNode.Equal(oldNode) {
			continue
		}

		if err := serviceSync.purgeService(uid); err != nil {
			return err
		}

	}

	return nil
}

// purgeService will purge a service from the nanobox
func (serviceSync *processServiceSync) purgeService(uid string) error {
	service := processor.ProcessControl{
		DevMode: serviceSync.control.DevMode,
		Verbose: serviceSync.control.Verbose,
		Meta: map[string]string{
			"name": uid,
		},
	}

	if err := processor.Run("service_remove", service); err != nil {
		fmt.Printf("service_remove (%s): %s\n", uid, err.Error())
		return err
	}

	return nil
}

// provisionServices will provision services that are defined in the boxfile
// but not running on nanobox
func (serviceSync *processServiceSync) provisionDataServices() error {

	// convert the data into boxfile library helpers
	newBoxfile := boxfile.New(serviceSync.newBoxfile.Data)

	// grab all of the data nodes
	dataServices := newBoxfile.Nodes("data")

	for _, uid := range dataServices {
		image := newBoxfile.Node(uid).StringValue("image")

		if image == "" {
			serviceType := regexp.MustCompile(`.+\.`).ReplaceAllString(uid, "")
			image = "nanobox/" + serviceType
		}

		config := processor.ProcessControl{
			DevMode:      serviceSync.control.DevMode,
			Verbose:      serviceSync.control.Verbose,
			DisplayLevel: serviceSync.control.DisplayLevel + 1,
			Meta: map[string]string{
				"name":  uid,
				"image": image,
			},
		}

		// setup the service
		if err := processor.Run("service_setup", config); err != nil {
			return err
		}

		// and configure it
		if err := processor.Run("service_configure", config); err != nil {
			return err
		}

	}

	return nil
}

// isPlatform will return true if the uid matches a platform service
func isPlatformUID(uid string) bool {
	return uid == PORTAL || uid == HOARDER || uid == MIST || uid == LOGVAC
}