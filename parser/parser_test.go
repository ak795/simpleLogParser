package parser

import (
	"testing"
)

func TestGetTopURL(t *testing.T) {
	factory := [][]string{
		{
			"1399202",
			"/a/1",
			"get",
			"45",
			"200",
		},
		{
			"9393002",
			"/a/1",
			"get",
			"22",
			"200",
		},
		{
			"3939929",
			"/a/1/ret",
			"delete",
			"102",
			"500",
		},
		{
			"203939211",
			"/a/b/c/d",
			"get",
			"939",
			"302",
		},
		{
			"129293",
			"/a/1/b",
			"put",
			"23.4",
			"200",
		},
		{
			"2390390",
			"/a/2",
			"post",
			"32",
			"200",
		},
	}
	expectedTopHits := []TopHits{
		{
			Method: "get",
			Url: "/a/{id}",
			Freq: 2,
		},
		{
			Method: "put",
			Url: "/a/{id}/b",
			Freq: 1,
		},
		{
			Method: "post",
			Url: "/a/{id}",
			Freq: 1,
		},
		{
			Method: "delete",
			Url: "/a/{id}/ret",
			Freq: 1,
		},
		{
			Method: "get",
			Url: "/a/b/c/d",
			Freq: 1,
		},
	}

	output := GetTopURL(factory)
	if len(output) != len(expectedTopHits) {
		t.Errorf("Length was incorrect, got: %d, want: %d.",len(output), len(expectedTopHits))
	}
	if len(output) != 5 {
		t.Errorf("Could not find top 5 hits, got: %d, want: %d.",len(output), len(expectedTopHits))
	}

	cnt := 0
	for i := 0; i < len(output); i++ {
		for j := 0; j < len(expectedTopHits); j++ {
			if output[i] == expectedTopHits[j] {
				cnt = cnt + 1
				break
			}
		}
	}

	if cnt != len(output) {
		t.Errorf("Incorrect Top Url, got: %v, want %v.", cnt, len(output))
	}
}

func TestGetLatencyMetric(t *testing.T) {
	factory := [][]string{
		{
			"1399202",
			"/a/1",
			"get",
			"45",
			"200",
		},
		{
			"9393002",
			"/a/1",
			"get",
			"22",
			"200",
		},
		{
			"3939929",
			"/a/1/ret",
			"delete",
			"102",
			"500",
		},
		{
			"203939211",
			"/a/b/c/d",
			"get",
			"939",
			"302",
		},
		{
			"129293",
			"/a/1/b",
			"put",
			"23.4",
			"200",
		},
		{
			"2390390",
			"/a/2",
			"post",
			"32",
			"200",
		},
	}
	expectedLatencyMetric := []Latency{
		{
			Method: "get",
			Url: "/a/{id}",
			MinTime: 22,
			MaxTime: 45,
			AvgTime: 33.5,
		},
		{
			Method: "delete",
			Url: "/a/{id}/ret",
			MinTime: 102,
			MaxTime: 102,
			AvgTime: 102,
		},
		{
			Method: "get",
			Url: "/a/b/c/d",
			MinTime: 939,
			MaxTime: 939,
			AvgTime: 939,
		},
		{
			Method: "put",
			Url: "/a/{id}/b",
			MaxTime: 23.4,
			MinTime: 23.4,
			AvgTime: 23.4,
		},
		{
			Method: "post",
			Url: "/a/{id}",
			MinTime: 32,
			MaxTime: 32,
			AvgTime: 32,
		},
	}

	output := GetLatencyMetric(factory)
	if len(output) != len(expectedLatencyMetric) {
		t.Errorf("Length was incorrect, got: %d, want: %d.",len(output), len(expectedLatencyMetric))
	}

	cnt := 0
	for i := 0; i < len(output); i++ {
		for j := 0; j < len(expectedLatencyMetric); j++ {
			if output[i] == expectedLatencyMetric[j] {
				cnt = cnt + 1
				break
			}
		}
	}

	if cnt != len(expectedLatencyMetric) {
		t.Errorf("Incorrect Top Url, got: %v, want %v.", cnt, len(expectedLatencyMetric))
	}
}
