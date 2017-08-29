///go:generate const_list -type=SenarioType ./wepay
///go:generate mapconst -type=SenarioType ./wepay

//go:generate stringer -type=SenarioType ./wepay/consts-ts.go
package cement
