package data

import (
	"github.com/zbitech/controller/pkg/model"
	"math/rand"
)

var (
	letters       = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	boolean       = []interface{}{true, false}
	volumeTypes   = []interface{}{model.EphemeralDataVolume, model.PersistentDataVolume}
	sourceTypes   = []interface{}{model.NewDataSource, model.VolumeDataSource, model.SnapshotDataSource}
	versions      = []interface{}{"v1"}
	networkTypes  = []interface{}{model.NetworkTypeTest}
	instanceTypes = []interface{}{model.InstanceTypeZCASH, model.InstanceTypeLWD}
	//	roleTypes         = []interface{}{model.RoleUser, model.RoleOwner}
	//	subscriptionTypes = []interface{}{model.SubscriptionTeamMember, model.SubscriptionBronzeLevel}
	//	snapshotTypes     = []interface{}{model.SnapshotDataSource, model.SnapshotScheduleType}
	scheduleTypes  = []interface{}{model.HourlySnapshotScheduled, model.DailySnapshotSchedule, model.WeeklySnapshotSchedule, model.MonthlySnapshotSchedule}
	resourceStatus = []interface{}{"active", "deleted", "failed", "pending"}
)

func randomString(n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func randomValue(array []interface{}) interface{} {
	return array[rand.Intn(len(array))]
}

func getProperty(props map[string]interface{}, key string, _default interface{}) interface{} {

	if props != nil {
		val, ok := props[key]
		if ok {
			return val
		}
	}

	return _default
}
