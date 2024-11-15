package kernel

type (
	vexation interface {
		vapour() inspection
		cause() string
		extent() string
	}

	vex struct {
		data            interface{}
		vap             inspection
		causeOfVexation string
		ofExtent        string
	}
)

func (v *vex) Data() interface{} {
	return v.data
}

func (v *vex) vapour() inspection {
	return v.vap
}

func (v *vex) cause() string {
	return v.causeOfVexation
}

func (v *vex) extent() string {
	return v.ofExtent
}
