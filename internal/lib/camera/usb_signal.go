package camera

import (
	"fmt"
	"io/ioutil"
	"os"
)

// usb eg. 1.2, 1.3
func EnsureUsbExists(usb string) error {
	valuePath := fmt.Sprintf("/sys/devices/platform/fe380000.usb/usb5/5-1/5-%s/bConfigurationValue", usb)
	_, err := os.Stat(valuePath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(valuePath, []byte("1"), 0644)
}

func CheckDevicePath(path string) error {
	_, err := os.Stat(path)
	return err
}
