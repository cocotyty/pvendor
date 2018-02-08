package main

import (
	"github.com/naoina/toml"
	"flag"
	"io/ioutil"
	"fmt"
	"regexp"
	"strings"
	"os"
	"os/exec"
)

type configFile struct {
	Dep []Dep
}
type Dep struct {
	Git    string
	Branch string
	Path   string
}

var gitPath = []*regexp.Regexp{
	regexp.MustCompile(`git@(.*):(.*)/(.*)\.git`),
	regexp.MustCompile(`git@(.*):(.*)/(.*)`),
	regexp.MustCompile(`https://(.*)/(.*)/(.*)\.git`),
	regexp.MustCompile(`https://(.*)/(.*)/(.*)`),
}
var S = string(os.PathSeparator)
var conf = flag.String("conf", "."+S+"vendor.toml", "-conf='./vendor.toml'")

func main() {
	flag.Parse()
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file, err := ioutil.ReadFile(*conf)
	if err != nil {
		if *conf == "."+S+"vendor.toml" {
			ioutil.WriteFile("."+S+"vendor.toml", def, 0777)
			return
		}
		fmt.Println("配置文件读取错误", *conf)
		return
	}
	deps := &configFile{}
	err = toml.Unmarshal(file, deps)
	if err != nil {
		panic(err)
	}
	for _, v := range deps.Dep {
		if v.Path == "" {
			v.Path = parseGitPath(v.Git)
			if v.Path == "" {
				fmt.Println("错误的git路径", v.Git)
				return
			}
		}
		cmd("rm", "-rf", dir+S+"vendor"+S+v.Path)
		cmd("git", "clone", v.Git, dir+S+"vendor"+S+v.Path)
		if v.Branch != "" {
			wdCmd(dir+S+"vendor"+S+v.Path, "git", "checkout", v.Branch)
		}
		cmd("rm","-rf", dir+S+"vendor"+S+v.Path+S+".git")

	}
}
func wdCmd(workDir string, args ... string) {
	c := exec.Command(args[0], args[1:]...)
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Dir = workDir
	c.Run()
}
func cmd(args ... string) {
	c := exec.Command(args[0], args[1:]...)
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Run()
}
func parseGitPath(git string) (string) {
	for _, v := range gitPath {
		finds := v.FindStringSubmatch(git)
		if len(finds) == 4 {
			return strings.Join(finds[1:], string(os.PathSeparator))
		}
	}
	return ""
}

var def = []byte(`
# [[dep]]
# git=""
# branch=""
# path=""
`)
