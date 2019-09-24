/**
* @Author: xhzhang
* @Date: 2019/9/5 15:24
 */
package client
type addOption struct {
	OrgName string
	ProName string
	EnvName string
}

type AddOption interface {
	apply(*addOption)
}

type funcOptionGroup struct {
	f func(*addOption)
}

func (fdo *funcOptionGroup) apply(do *addOption) {
	fdo.f(do)
}

func newFuncAddOption(f func(*addOption)) *funcOptionGroup {
	return &funcOptionGroup{f: f}
}

func WithOrg(name string) AddOption {
	return newFuncAddOption(func(o *addOption) { o.OrgName = name })
}

func WithPro(name string) AddOption {
	return newFuncAddOption(func(o *addOption) { o.ProName = name })
}

func WithEnv(name string) AddOption {
	return newFuncAddOption(func(o *addOption) { o.EnvName = name })
}

//默认参数
func defaultAddOption() addOption {
	return addOption{}
}

