package main

import "testing"

func TestStartNodeExporter(t *testing.T) {
	conf := &NodeExporterConfig{}
	StartNodeExporter(conf)
}
