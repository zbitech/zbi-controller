package manager

import (
	"context"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type AppResourceManager struct {
}

func NewAppResourceManager() interfaces.AppResourceManagerIF {
	return &AppResourceManager{}
}

// CreateVolumeResource gen
func (app *AppResourceManager) CreateVolumeResource(ctx context.Context, project, instance string, volumes ...model.VolumeSpec) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "app.CreateVolumeResource")
	defer func() { logger.LogServiceTime(log) }()

	fileTemplate := helper.Config.GetAppTemplate()

	var objects = make([]unstructured.Unstructured, 0, len(volumes))
	for _, volume := range volumes {
		properties := make(map[string]interface{})
		data, err := fileTemplate.ExecuteTemplate("VOLUME", volume)
		if err != nil {
			//			logger.Errorf(ctx, "volume template failed - %s", err)
			return nil, err
		}

		properties["volumeName"] = volume.VolumeName
		if volume.DataSourceType == model.VolumeDataSource {
			properties["source"] = "Volume"
		} else if volume.DataSourceType == model.SnapshotDataSource {
			properties["source"] = "Snapshot"
		}
		properties["size"] = volume.Size

		obj, err := helper.CreateYAMLObject(data)
		if err != nil {
			log.Errorf("volume template failed - %s", err)
			return nil, err
		}
		objects = append(objects, *obj)
	}

	return objects, nil
}

func (app *AppResourceManager) CreateSnapshotResource(ctx context.Context, project, instance string, req *model.SnapshotRequest) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "app.CreateSnapshotResource")
	defer func() { logger.LogServiceTime(log) }()

	fileTemplate := helper.Config.GetAppTemplate()

	var specArr []string
	var err error

	snapshotSpec := model.SnapshotSpec{
		SnapshotName:  req.VolumeName + "-" + utils.GenerateRandomString(5, true),
		Namespace:     req.Namespace,
		VolumeName:    req.VolumeName,
		SnapshotClass: req.SnapshotClass,
		Labels:        req.Labels,
	}

	specArr, err = fileTemplate.ExecuteTemplates([]string{"SNAPSHOT"}, snapshotSpec)

	if err != nil {
		log.Errorf("backup templates for version %s failed - %s", req.Labels["version"], err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	//properties := make(map[string]interface{})
	//properties["volumeName"] = req.VolumeName

	return helper.CreateYAMLObjects(specArr /*, project, instance, properties*/)
}

func (app *AppResourceManager) CreateSnapshotScheduleResource(ctx context.Context, project, instance string, req *model.SnapshotScheduleRequest) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "app.CreateSnapshotScheduleResource")
	defer func() { logger.LogServiceTime(log) }()

	fileTemplate := helper.Config.GetAppTemplate()

	var specArr []string
	var err error

	snapshotSpec := model.SnapshotScheduleSpec{
		ScheduleName:     req.VolumeName + "-" + string(req.Schedule),
		Namespace:        req.Namespace,
		SnapshotClass:    req.SnapshotClass,
		BackupExpiration: req.BackupExpiration,
		MaxBackupCount:   req.MaxBackupCount,
		ScheduleType:     req.Schedule,
		Schedule:         helper.CreateSnapshotSchedule(req.Schedule),
		Labels:           req.Labels,
	}

	specArr, err = fileTemplate.ExecuteTemplates([]string{"SCHEDULE_SNAPSHOT"}, snapshotSpec)

	if err != nil {
		log.Errorf("backup templates for version %s failed - %s", req.Labels["version"], err)
		//		return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	return helper.CreateYAMLObjects(specArr /*, project, instance*/)
}
