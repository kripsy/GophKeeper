//nolint:testpackage
package usecase

import (
	"reflect"
	"testing"
	"time"

	"github.com/kripsy/GophKeeper/internal/models"
)

func Test_findDifferences(t *testing.T) {
	testTime := time.Now()
	tests := []struct {
		name             string
		local            models.UserMeta
		server           models.UserMeta
		wantNeedDownload models.MetaData
		wantNeedUpload   models.MetaData
	}{
		{
			name: "download & upload",
			local: models.UserMeta{Data: models.MetaData{
				"data1": models.DataInfo{DataID: "data1", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}}},
			server: models.UserMeta{Data: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}}},
			wantNeedDownload: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime},
			},
			wantNeedUpload: models.MetaData{
				"data1": models.DataInfo{DataID: "data1", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime},
			},
		},
		{
			name: "upload",
			local: models.UserMeta{Data: models.MetaData{
				"data1": models.DataInfo{DataID: "data1", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}}},
			server:           models.UserMeta{Data: models.MetaData{}},
			wantNeedDownload: models.MetaData{},
			wantNeedUpload: models.MetaData{
				"data1": models.DataInfo{DataID: "data1", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}},
		},
		{
			name:  "download",
			local: models.UserMeta{Data: models.MetaData{}},
			server: models.UserMeta{Data: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}}},
			wantNeedDownload: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}},
			wantNeedUpload: models.MetaData{},
		},
		{
			name: "upload with deleted data",
			local: models.UserMeta{Data: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}},
			},
			server:           models.UserMeta{DeletedData: models.Deleted{"data2": struct{}{}}},
			wantNeedDownload: models.MetaData{},
			wantNeedUpload: models.MetaData{
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}},
		},
		{
			name: "download with deleted",
			local: models.UserMeta{Data: models.MetaData{},
				DeletedData: models.Deleted{"data2": struct{}{}}},
			server: models.UserMeta{Data: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}}},
			wantNeedDownload: models.MetaData{
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}},
			wantNeedUpload: models.MetaData{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNeedDownload, gotNeedUpload := findDifferences(tt.local, tt.server)
			if !reflect.DeepEqual(gotNeedDownload, tt.wantNeedDownload) {
				t.Errorf("findDifferences() gotNeedDownload = %v, want %v", gotNeedDownload, tt.wantNeedDownload)
			}
			if !reflect.DeepEqual(gotNeedUpload, tt.wantNeedUpload) {
				t.Errorf("findDifferences() gotNeedUpload = %v, want %v", gotNeedUpload, tt.wantNeedUpload)
			}
		})
	}
}
