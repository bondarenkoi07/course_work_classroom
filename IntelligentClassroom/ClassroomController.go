package IntelligentClassroom

import "github.com/pkg/errors"

//struct ClassroomController contains variables,
//which defines rules of classroom's devices controller work.
// True value means that device is active, False means that device is turned off
type ClassroomController struct {
	AirConditioner bool `json:"conditioner"`
	Projector      bool `json:"projector"`
	Screen         bool `json:"screen"`
	Lock           bool `json:"lock"`
}

//Constructor NewClassroomController initialize object with start condition :
//Conditioner, projector and screen are turned off and classroom is locked
func NewClassroomController() *ClassroomController {
	return &ClassroomController{AirConditioner: false, Projector: false, Screen: false, Lock: true}
}

//Opening of classroom means that lock is inactive, and it means that all other devices should be inactive
func (c *ClassroomController) OpenLock() {
	(*c).Lock = false
}

//Before locking, all other devices should be inactive, so in other ways, CloseLock() will throw error
func (c *ClassroomController) CloseLock() error {
	var err error

	if (*c).AirConditioner == false && (*c).Screen == false && (*c).Projector == false {
		err = nil
		(*c).Lock = true
	} else {
		errorString := "Before closing door, you should also turn off: "
		if (*c).AirConditioner == true {
			errorString += " Air Conditioner"
		}
		if (*c).Projector == true {
			errorString += " projector"
		}
		if (*c).Screen == true {
			errorString += " screen"
		}
		err = errors.New(errorString)

		(*c).Lock = false
	}

	return err
}

//All devices besides lock could be activated only if lock is also activated
func (c *ClassroomController) ActivateConditioner() error {
	var err error
	if (*c).AirConditioner == true {
		if (*c).Lock == true {
			err = errors.New("Bad condition, system returns to start condition")
			(*c).AirConditioner = false
			(*c).Screen = false
			(*c).Projector = false
		} else {
			err = errors.New("Air Conditioner is already activated")
		}
	} else if (*c).Lock == true {
		err = errors.New("Classroom is locked")
	} else {
		err = nil
		(*c).AirConditioner = true
	}
	return err
}

func (c *ClassroomController) DeactivateConditioner() {
	(*c).AirConditioner = false
}

func (c *ClassroomController) ActivateScreen() error {
	var err error
	if (*c).Screen == true {
		if (*c).Lock == true {
			err = errors.New("Bad condition, system returns to start condition")
			(*c).AirConditioner = false
			(*c).Screen = false
			(*c).Projector = false
		} else {
			err = errors.New("Screen is already activated")
		}
	} else if (*c).Lock == true {
		err = errors.New("Classroom is locked")
	} else {
		err = nil
		(*c).Screen = true
	}
	return err
}

func (c *ClassroomController) DeactivateScreen() {
	(*c).Screen = false
}

func (c *ClassroomController) ActivateProjector() error {
	var err error
	if (*c).Projector == true {
		if (*c).Lock == true {
			err = errors.New("Bad condition, system returns to start condition")
			(*c).AirConditioner = false
			(*c).Screen = false
			(*c).Projector = false
		} else {
			err = errors.New("Projector is already activated")
		}
	} else if (*c).Lock == true {
		err = errors.New("Classroom is locked")
	} else {
		err = nil
		(*c).Projector = true
	}
	return err
}

func (c *ClassroomController) DeactivateProjector() {
	(*c).Projector = false
}
