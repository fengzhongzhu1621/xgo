package cgroup

import "testing"

func TestGetDockerMemoryUsageInfos(t *testing.T) {
	tests := []struct {
		name      string
		wantUsage float64
		wantTotal uint64
		wantRss   uint64
		wantErr   bool
	}{
		{
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUsage, gotTotal, gotRss, err := GetDockerMemoryUsageInfos()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDockerMemoryUsageInfos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUsage <= tt.wantUsage {
				t.Errorf(
					"GetDockerMemoryUsageInfos() gotUsage = %v, want %v",
					gotUsage,
					tt.wantUsage,
				)
			}
			if gotTotal <= tt.wantTotal {
				t.Errorf(
					"GetDockerMemoryUsageInfos() gotTotal = %v, want %v",
					gotTotal,
					tt.wantTotal,
				)
			}
			if gotRss <= tt.wantRss {
				t.Errorf("GetDockerMemoryUsageInfos() gotRss = %v, want %v", gotRss, tt.wantRss)
			}
		})
	}
}

func Test_getMachineMemoryTotal(t *testing.T) {
	tests := []struct {
		name      string
		wantTotal uint64
		wantErr   bool
	}{
		{
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTotal, err := getMachineMemoryTotal()
			if (err != nil) != tt.wantErr {
				t.Errorf("getMachineMemoryTotal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTotal <= tt.wantTotal {
				t.Errorf("getMachineMemoryTotal() = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}
