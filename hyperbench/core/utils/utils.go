package utils

/**
 *  Copyright (C) 2021 HyperBench.
 *  SPDX-License-Identifier: Apache-2.0
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 * @brief Provide statistical and record-related tool functions
 * @file utils.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	fcom "github.com/hyperbench/hyperbench-common/common"

	"strconv"
)

// RemoteStatistic2CSV append data's fields to base,
func RemoteStatistic2CSV(base []string, data *fcom.RemoteStatistic) []string {
	if base == nil {
		base = make([]string, 0, 30)
	}
	base = append(base,
		"statistic",
		i2s(data.Start),
		i2s(data.End),
		i2s(data.BlockNum),
		i2s(data.Tps),
		i2s(data.CTps),
		i2s(data.Bps),
		i2s(data.SentTx),
		i2s(data.MissedTx),
		i2s(data.Tps))
	return base
}

// AggData2CSV append data's fields to base,
func AggData2CSV(base []string, t fcom.DataType, data fcom.AggData) []string {
	if base == nil {
		base = make([]string, 0, 30)
	}
	base = append(base,
		string(t),
		data.Label,
		i2s(data.Time),
		i2s(data.Duration),
		i2s(data.Num),
		i2s(data.Statuses[fcom.Failure]),
		i2s(data.Statuses[fcom.Success]),
		i2s(data.Statuses[fcom.Confirm]),
		i2s(data.Statuses[fcom.Unknown]))
	base = Latency2CSV(base, data.Send)
	base = Latency2CSV(base, data.Confirm)
	base = Latency2CSV(base, data.Write)
	return base
}

// Latency2CSV append latency's fields to base
func Latency2CSV(base []string, latency fcom.Latency) []string {
	if base == nil {
		base = make([]string, 0, 7)
	}
	base = append(base,
		i2s(latency.Avg),
		i2s(latency.P0),
		i2s(latency.P50),
		i2s(latency.P90),
		i2s(latency.P95),
		i2s(latency.P99),
		i2s(latency.P100))
	return base
}

func i2s(i interface{}) (s string) {
	switch v := i.(type) {
	case int64:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	}
	return ""
}

// DivideAndCeil a/b rounded up.
func DivideAndCeil(a, b int) int {
	return (a + b - 1) / b
}
