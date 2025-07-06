package main

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/offen/docker-volume-backup/internal/labels"
)

func TestLoadConfigsFromLabels(t *testing.T) {
	origDocker := dockerClientFactory
	origScan := scanVolumeLabelsFn
	origBasic := parseBasicLabelsFn
	origAdvanced := parseAdvancedLabelsFn
	defer func() {
		dockerClientFactory = origDocker
		scanVolumeLabelsFn = origScan
		parseBasicLabelsFn = origBasic
		parseAdvancedLabelsFn = origAdvanced
	}()

	dockerClientFactory = func() (*client.Client, error) { return nil, nil }
	scanVolumeLabelsFn = func(_ interface {
		VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
	}) (map[string]map[string]string, error) {
		return map[string]map[string]string{
			"alpha": {"schedule": "@hourly", "target": "/data1"},
			"beta":  {"schedule": "@daily", "target": "/data2"},
		}, nil
	}
	parseBasicLabelsFn = func(ls map[string]string) (labels.Config, error) {
		return labels.Config{
			BackupCronExpression: ls["schedule"],
			BackupArchive:        ls["target"],
		}, nil
	}
	parseAdvancedLabelsFn = func(map[string]string, *labels.Config) error { return nil }

	cfgs, err := sourceConfiguration(configStrategyLabels)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfgs) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(cfgs))
	}

	results := map[string]*Config{}
	for _, c := range cfgs {
		results[c.source] = c
	}

	if results["alpha"].BackupCronExpression != "@hourly" || results["alpha"].BackupArchive != "/data1" {
		t.Errorf("unexpected config for alpha: %+v", results["alpha"])
	}
	if results["beta"].BackupCronExpression != "@daily" || results["beta"].BackupArchive != "/data2" {
		t.Errorf("unexpected config for beta: %+v", results["beta"])
	}
}

func TestLoadConfigsFromLabelsEmpty(t *testing.T) {
	origDocker := dockerClientFactory
	origScan := scanVolumeLabelsFn
	defer func() {
		dockerClientFactory = origDocker
		scanVolumeLabelsFn = origScan
	}()

	dockerClientFactory = func() (*client.Client, error) { return nil, nil }
	scanVolumeLabelsFn = func(_ interface {
		VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
	}) (map[string]map[string]string, error) {
		return map[string]map[string]string{}, nil
	}

	cfgs, err := sourceConfiguration(configStrategyLabels)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfgs) != 0 {
		t.Fatalf("expected 0 configs, got %d", len(cfgs))
	}
}
