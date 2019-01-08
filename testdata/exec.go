package main

import (
	"os/exec"
)

func executor() int {
	cmd := exec.Command("ls", "la");
	err := cmd.Run();
	if err != nil {
		return 1;
	}
	return 0;
}
