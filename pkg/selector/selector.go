package selector

import (
	"fmt"

	"github.com/itmecho/kvm/pkg/github"
)

func Select(msg string, releases []github.Release) github.Release {
	for i, r := range releases {
		fmt.Printf("[%d] %s\n", i, r.Name)
	}

	fmt.Print(msg)

	// TODO find a nice way of doing this
	var choice int
	fmt.Scan(&choice)

	// TODO check that the choice was valid

	return releases[choice]
}
