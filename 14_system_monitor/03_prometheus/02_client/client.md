<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [客户端](#%E5%AE%A2%E6%88%B7%E7%AB%AF)
  - [调用](#%E8%B0%83%E7%94%A8)
  - [api 接口](#api-%E6%8E%A5%E5%8F%A3)
    - [/api/v1/query](#apiv1query)
    - [/api/v1/query_range](#apiv1query_range)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 客户端

## 调用



## api 接口

### /api/v1/query
```shell
/*
1. GET /api/v1/query 单点查询数据格式 https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries
{
  "resultType": "matrix" | "vector" | "scalar" | "string",
  "result": <value>
}

<value> 详细格式: https://prometheus.io/docs/prometheus/latest/querying/api/#expression-query-result-formats

vector 类型

[
  {
    "metric": { "<label_name>": "<label_value>", ... },
    "value": [ <unix_time>, "<sample_value>" ],
    "histogram": [ <unix_time>, <histogram> ]
  },
  ...
]


```



查询
```go
// github.com/prometheus/client_golang@v1.11.0/api/prometheus/v1/api.go
func (h *httpAPI) Query(ctx context.Context, query string, ts time.Time) (model.Value, Warnings, error) {
	u := h.client.URL(epQuery, nil)
	q := u.Query()

	q.Set("query", query)
	if !ts.IsZero() {
		q.Set("time", formatTime(ts))
	}

	_, body, warnings, err := h.client.DoGetFallback(ctx, u, q)
	if err != nil {
		return nil, warnings, err
	}

	// 类型转换 
	var qres queryResult
	return model.Value(qres.v), warnings, json.Unmarshal(body, &qres)
}
```

```go
func (h *apiClientImpl) DoGetFallback(ctx context.Context, u *url.URL, args url.Values) (*http.Response, []byte, Warnings, error) {
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(args.Encode()))
	if err != nil {
		return nil, nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, body, warnings, err := h.Do(ctx, req)
	if resp != nil && (resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotImplemented) {
		u.RawQuery = args.Encode()
		req, err = http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, nil, warnings, err
		}

	} else {
		if err != nil {
			return resp, body, warnings, err
		}
		return resp, body, warnings, nil
	}
	return h.Do(ctx, req)
}
```

返回结果
```go
type queryResult struct {
	Type   model.ValueType `json:"resultType"`
	Result interface{}     `json:"result"`

	// The decoded value.
	v model.Value
}
```




具体四种类型
```go
type ValueType int

const (
	ValNone ValueType = iota
	ValScalar
	ValVector
	ValMatrix
	ValString
)

// MarshalJSON implements json.Marshaler.
func (et ValueType) MarshalJSON() ([]byte, error) {
	return json.Marshal(et.String())
}

func (et *ValueType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "<ValNone>":
		*et = ValNone
	case "scalar":
		*et = ValScalar
	case "vector":
		*et = ValVector
	case "matrix":
		*et = ValMatrix
	case "string":
		*et = ValString
	default:
		return fmt.Errorf("unknown value type %q", s)
	}
	return nil
}

```



Vector: 具有相同的时间戳
```go
// Vector is basically only an alias for Samples, but the
// contract is that in a Vector, all Samples have the same timestamp.
type Vector []*Sample

// MarshalJSON implements json.Marshaler.
func (s SamplePair) MarshalJSON() ([]byte, error) {
	t, err := json.Marshal(s.Timestamp)
	if err != nil {
		return nil, err
	}
	v, err := json.Marshal(s.Value)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SamplePair) UnmarshalJSON(b []byte) error {
	v := [...]json.Unmarshaler{&s.Timestamp, &s.Value}
	return json.Unmarshal(b, &v)
}

// Equal returns true if this SamplePair and o have equal Values and equal
// Timestamps. The semantics of Value equality is defined by SampleValue.Equal.
func (s *SamplePair) Equal(o *SamplePair) bool {
	return s == o || (s.Value.Equal(o.Value) && s.Timestamp.Equal(o.Timestamp))
}

func (s SamplePair) String() string {
	return fmt.Sprintf("%s @[%s]", s.Value, s.Timestamp)
}

// Sample is a sample pair associated with a metric.
type Sample struct {
	Metric    Metric      `json:"metric"`
	Value     SampleValue `json:"value"`
	Timestamp Time        `json:"timestamp"`
}

```


### /api/v1/query_range
```shell
2. GET /api/v1/query_range 范围查询数据格式 https://prometheus.io/docs/prometheus/latest/querying/api/#range-queries

{
  "resultType": "matrix",
  "result": <value>
}

<value> 详细格式: https://prometheus.io/docs/prometheus/latest/querying/api/#range-vectors
[
  {
    "metric": { "<label_name>": "<label_value>", ... },
    "values": [ [ <unix_time>, "<sample_value>" ], ... ],
    "histograms": [ [ <unix_time>, <histogram> ], ... ]
  },
  ...
]


*/
```





类型
```go
type Value interface {
	Type() ValueType
	String() string
}

func (Matrix) Type() ValueType  { return ValMatrix }
func (Vector) Type() ValueType  { return ValVector }
func (*Scalar) Type() ValueType { return ValScalar }
func (*String) Type() ValueType { return ValString }

```



反序列化实现
```go
type Scalar struct {
	Value     SampleValue `json:"value"`
	Timestamp Time        `json:"timestamp"`
}


func (s Scalar) String() string {
	return fmt.Sprintf("scalar: %v @[%v]", s.Value, s.Timestamp)
}

// MarshalJSON implements json.Marshaler.
func (s Scalar) MarshalJSON() ([]byte, error) {
	v := strconv.FormatFloat(float64(s.Value), 'f', -1, 64)
	return json.Marshal([...]interface{}{s.Timestamp, string(v)})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *Scalar) UnmarshalJSON(b []byte) error {
	var f string
	v := [...]interface{}{&s.Timestamp, &f}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	value, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return fmt.Errorf("error parsing sample value: %s", err)
	}
	s.Value = SampleValue(value)
	return nil
}
```



