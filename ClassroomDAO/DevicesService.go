package ClassroomDAO

import (
	Classroom "awesomeProject1/IntelligentClassroom"
	"github.com/pkg/errors"
)

type DevicesDAO struct {
	repo *Classroom.ClassroomController
}

func (d *DevicesDAO) Init() {
	(*d).repo = Classroom.NewClassroomController()
}

func (d *DevicesDAO) GetRepo() *Classroom.ClassroomController {
	return (*d).repo
}

func (d *DevicesDAO) ActivateDevice(deviceName string) error {
	var err error

	switch deviceName {
	case "conditioner":
		err = (*d).repo.ActivateConditioner()
		break
	case "projector":
		err = (*d).repo.ActivateProjector()
		break
	case "screen":
		err = (*d).repo.ActivateScreen()
		break
	case "lock":
		err = (*d).repo.CloseLock()
		break
	default:
		err = errors.New("Unknown device")
	}
	return err
}

func (d *DevicesDAO) DeactivateDevice(deviceName string) error {
	var err error

	switch deviceName {
	case "conditioner":
		(*d).repo.DeactivateConditioner()
		break
	case "projector":
		(*d).repo.DeactivateProjector()
		break
	case "screen":
		(*d).repo.DeactivateScreen()
		break
	case "lock":
		(*d).repo.OpenLock()
		break
	default:
		err = errors.New("Unknown device")
	}
	return err
}
