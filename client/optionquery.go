/**
* @Author: xhzhang
* @Date: 2019/9/5 10:14
 */
package client

type queryOption struct {
	Ids           []int32 // id为int32类型的
	AgentIDs      []string
	ServiceIDs    []string
	Names         []string // 所有类型的name
	OrgNames      []string
	ProNames      []string
	EnvNames      []string
	GroupNames    []string
	ReleaseNames  []string
	MoudleNames   []string
	AgentIsOnLine bool
}

type QueryOption interface {
	apply(*queryOption)
}

type OptionFunc struct {
	f func(*queryOption)
}

func (fdo *OptionFunc) apply(do *queryOption) {
	fdo.f(do)
}

func newFuncQueryOption(f func(*queryOption)) *OptionFunc {
	return &OptionFunc{f: f}
}

//默认参数
func defaultQueryOption() queryOption {
	return queryOption{}
}

func WithIds(ids []int32) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.Ids = ids })
}

func WithNames(names []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.Names = names })
}

func WithOrgs(names []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.OrgNames = names })
}

func WithPros(names []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.ProNames = names })
}

func WithEnvs(names []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.EnvNames = names })
}

func WithGroups(names []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.GroupNames = names })
}

func WithReleases(names []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.ReleaseNames = names })
}

func WithMoudles(names []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.MoudleNames = names })
}


func WithAgents(ids []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.AgentIDs = ids })
}

func WithServices(ids []string) QueryOption {
	return newFuncQueryOption(func(o *queryOption) { o.ServiceIDs = ids })
}

func WithAgentStatus(status bool) QueryOption{
	return newFuncQueryOption(func(o *queryOption) { o.AgentIsOnLine = status })
}
