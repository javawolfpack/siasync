package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  // "os/signal"
  // "syscall"
  // "time"
  "strings"

  "gitlab.com/NebulousLabs/Sia/build"
  sia "gitlab.com/NebulousLabs/Sia/node/api/client"
  // "github.com/takama/daemon"
  "gopkg.in/yaml.v2"
)

var (
  archive           bool
  debug             bool
  password          string
  prefix            string
  include           string
  includeExtensions []string
  exclude           string
  excludeExtensions []string
  siaDir            string
  daemon_cmd        string
  dataPieces        uint64
  parityPieces      uint64
)

// //    dependencies that are NOT required by the service, but might be used
// var dependencies = []string{"dummy.service"}
//
// const (
//
//     // name of the service
//     name        = "SiaSync"
//     description = "Sia Synchronization Service"
//
// )


type YamlSync struct {
    Version    string
    Sync []Sync `yaml:"sync"`
}

type Sync struct {
    Name string `yaml:"name"`
    Path string `yaml:"path"`
    Archive bool `yaml:"archive"`
    Prefix string `yaml:"siaDir"`
    DataPieces int `yaml:"dataPieces"`
    ParityPieces int `yaml:"parityPieces"`
    IncludeExtensions []string `yaml:"includeExtensions"`
    ExcludeExtensions []string `yaml:"excludeExtensions"`
}



// // Service has embedded daemon
// type Service struct {
//     daemon.Daemon
// }

var stdlog, errlog *log.Logger

func Usage() {
  fmt.Printf(`usage: siasync <flags>
  for example: ./siasync -password abcd123

`)
  flag.PrintDefaults()
}

// findApiPassword looks for the API password via a flag, env variable, or the default apipassword file
func findApiPassword() string {
  // password from cli -password flag
  if password != "" {
    return password
  } else {
    // password from environment variable
    envPassword := os.Getenv("SIA_API_PASSWORD")
    if envPassword != "" {
      return envPassword
    } else {
      // password from apipassword file
      APIPasswordFile, err := ioutil.ReadFile(build.APIPasswordFile(build.DefaultSiaDir()))
      if err != nil {
        fmt.Println("Could not read API password file:", err)
      }
      return strings.TrimSpace(string(APIPasswordFile))
    }
  }

}

func testConnection(sc *sia.Client) {
  // Get siad Version
  version, err := sc.DaemonVersionGet()
  if err != nil {
    panic(err)
  }
  log.Println("Connected to Sia ", version.Version)

  // Check Allowance
  rg, err := sc.RenterGet()
  if err != nil {
    log.Fatal("Could not get renter info:", err)
  }
  if rg.Settings.Allowance.Funds.IsZero() {
    log.Fatal("Cannot upload: No allowance available")
  }

  // Check Contracts
  rc, err := sc.RenterDisabledContractsGet()
  if err != nil {
    log.Fatal("Could not get renter contracts", err)
  }
  if len(rc.ActiveContracts) == 0 {
    log.Fatal("No active contracts")
  }
  var GoodForUpload = 0
  for _, c := range rc.ActiveContracts {
    if c.GoodForUpload {
      GoodForUpload += 1
    }
  }
  log.Println(GoodForUpload, " contracts are ready for upload")

}

func init() {
    f, err := os.OpenFile("text.log",
      os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
      log.Println(err)
    }
    stdlog = log.New(f, "stdout: ", log.Ldate|log.Ltime)
    errlog = log.New(f, "stderr: ", log.Ldate|log.Ltime)
}

func main() {
  flag.Usage = Usage
  address := flag.String("address", "127.0.0.1:9980", "Sia's API address")
  flag.StringVar(&password, "password", "", "Sia's API password")
  agent := flag.String("agent", "Sia-Agent", "Sia agent")
  // flag.BoolVar(&archive, "archive", false, "Files will not be removed from Sia, even if they are deleted locally")
  flag.BoolVar(&debug, "debug", false, "Enable debug mode. Warning: generates a lot of output.")
  // flag.StringVar(&prefix, "subfolder", "siasync", "Folder on Sia to sync files too")
  // flag.StringVar(&include, "include", "", "Comma separated list of file extensions to copy, all other files will be ignored.")
  // flag.StringVar(&daemon_cmd, "daemon", "none", "use one of the daemon commands: install | remove | start | stop | status")
  // flag.StringVar(&exclude, "exclude", "", "Comma separated list of file extensions to skip, all other files will be copied.")
  // flag.Uint64Var(&dataPieces, "data-pieces", 10, "Number of data pieces in erasure code")
  // flag.Uint64Var(&parityPieces, "parity-pieces", 30, "Number of parity pieces in erasure code")
  //
  flag.Parse()
  //
  sc := sia.New(*address)
  sc.Password = findApiPassword()
  sc.UserAgent = *agent
  // directory := os.Args[len(os.Args)-1]
  //
  // // Verify that we can talk to Sia and have valid contracts.
  testConnection(sc)
  //
  // includeExtensions = strings.Split(include, ",")
  // excludeExtensions = strings.Split(exclude, ",")
  //
  // sf, err := NewSiafolder(directory, sc)
  // if err != nil {
  //   log.Fatal(err)
  // }
  // defer sf.Close()
  //
  // log.Println("watching for changes to ", directory)
  //
  // done := make(chan os.Signal)
  // signal.Notify(done, os.Interrupt)
  // <-done
  // fmt.Println("caught quit signal, exiting...")

  y := YamlSync{}
  yamlFile, err := ioutil.ReadFile("test.yml")
  if err != nil {
    log.Fatalf("error: %v", err)
  }
  yaml.Unmarshal(yamlFile, &y)
  fmt.Printf("%+v\n", y.Sync)


}
