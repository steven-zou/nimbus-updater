package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	nimbusCommand = "/mts/git/bin/nimbus-ctl"
)

func main() {
	nimbusPod := flag.String("p", "", "The pod name of the testbed located on")
	nimbusUser := flag.String("u", "", "The nimbus user used to update")
	testbed := flag.String("t", "", "The name of testbed")
	nimbusLocation := flag.String("l", "sc", "The location of the nimbus")
	interval := flag.Int64("i", 5, "The interval of updating action (days)")
	days := flag.Int("d", 7, "The extend days")
	flag.Parse()

	if len(strings.TrimSpace(*nimbusPod)) == 0 ||
		len(strings.TrimSpace(*testbed)) == 0 ||
		len(strings.TrimSpace(*nimbusUser)) == 0 {
		usage()
		return
	}

	//Start the ticket
	intervalByDays := time.Duration(*interval) * 24 * time.Hour
	tk := time.NewTicker(intervalByDays)
	fmt.Println("Updater is started...")

	update(*nimbusPod, *nimbusUser, *nimbusLocation, *testbed, *days)
	for _ = range tk.C {
		update(*nimbusPod, *nimbusUser, *nimbusLocation, *testbed, *days)
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("nimbus-updater [options]")
	fmt.Println("-p The pod name of the testbed located on")
	fmt.Println("-u The nimbus user used to update")
	fmt.Println("-t The name of testbed")
	fmt.Println("-l The location of the nimbus")
	fmt.Println("-i The interval of updating action (s)")
	fmt.Println("-d The extend days")
}

func update(pod, user, nimbusLocation, testbed string, days int) {
	args := []string{
		fmt.Sprintf("%s%s", "--nimbusLocation=", nimbusLocation),
		"--lease",
		fmt.Sprintf("%d", days),
		"--testbed",
		"extend_lease",
		testbed,
	}

	cmd := exec.Command(nimbusCommand, args...)
	env := os.Environ()
	env = append(env, fmt.Sprintf("NIMBUS=%s", pod))
	env = append(env, fmt.Sprintf("USER=%s", user))
	cmd.Env = env

	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	cmd.Start()
	if out, err := ioutil.ReadAll(stdoutReader); err != nil {
		fmt.Printf("%s\n", out)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error:%s\n", err.Error())
	}

	fmt.Println("Updated")
}
