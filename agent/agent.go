package agent

type Agent interface {
	Run()
	Process() bool
}
