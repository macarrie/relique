package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/macarrie/relique/internal/types/displayable"
)

type ClientDisplay struct {
	Name    string                    `json:"name"`
	Address string                    `json:"address"`
	Port    uint32                    `json:"port"`
	Modules []displayable.Displayable `json:"modules"`
}

func (c Client) Display() displayable.Struct {
	var modDisplay []displayable.Displayable
	for _, mod := range c.Modules {
		modDisplay = append(modDisplay, mod)
	}
	var d displayable.Struct = ClientDisplay{
		Name:    c.Name,
		Address: c.Address,
		Port:    c.Port,
		Modules: modDisplay,
	}

	return d
}

func (d ClientDisplay) Summary() string {
	var moduleNames []string
	for _, mod := range d.Modules {
		moduleNames = append(moduleNames, mod.Display().Summary())
	}
	return fmt.Sprintf("%v (%v) with modules %v", d.Name, d.Address, strings.Join(moduleNames, ", "))
}

func (d ClientDisplay) Details() string {
	moduleList := ""
	moduleDetailsList := ""
	for _, mod := range d.Modules {
		moduleList = fmt.Sprintf("%s\t- %s\n", moduleList, mod.Display().Summary())
		moduleDetailsList = fmt.Sprintf("%s- %s\n\n", moduleDetailsList, mod.Display().Details())
	}

	clientDetails := fmt.Sprintf(
		`CLIENT
------ 

Name: 		%s
Address: 	%s
Port: 		%d
Modules: 	
%s`,
		d.Name,
		d.Address,
		d.Port,
		moduleList,
	)

	moduleDetails := fmt.Sprintf(
		`MODULES DETAILS
---------------

%s`,
		moduleDetailsList,
	)

	return fmt.Sprintf("%s\n\n%s", clientDetails, moduleDetails)
}

func (d ClientDisplay) TableHeaders() []string {
	return []string{"Name", "Address", "Port", "Modules"}
}

func (d ClientDisplay) TableRow() []string {
	var moduleNames []string
	for _, mod := range d.Modules {
		moduleNames = append(moduleNames, mod.Display().Summary())
	}
	return []string{d.Name, d.Address, strconv.Itoa(int(d.Port)), strings.Join(moduleNames, ",")}
}
