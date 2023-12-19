package main

import (
	"strings"
	"testing"

	"github.com/dustin/go-humanize"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestParseHumanizedQuntity(t *testing.T) {
	memMb := humanize.IBytes(uint64(3600) * 1024 * 1024)
	memMb = strings.ReplaceAll(memMb, " ", "")
	memMb = memMb[:len(memMb)-2]

	q, lerr := resource.ParseQuantity(memMb)
	if lerr != nil {
		t.Fatal(lerr)
	}

	t.Logf("%s", q.String())
}

func TestParseHumanizedCpuQuntity(t *testing.T) {
	coresStr := humanize.SI(5300, "m")

	cores := resource.MustParse(coresStr)
	t.Logf("milliCores = %v (%v)\n", cores.MilliValue(), cores.Format)
}
