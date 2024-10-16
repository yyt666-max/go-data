package permit

type IPermitInitialize interface {
	// Grants 返回权限定义
	// map[domain]map[access][]target
	Grants() map[string]map[string][]string
}
