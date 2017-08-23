package main

import "log"
import "os"
import "strconv"
import "io/ioutil"
import "path/filepath"
import "gopkg.in/yaml.v2"

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: " + os.Args[0] + " [config/storm.yml]")
	}
	filename, _ := filepath.Abs(os.Args[1])
	log.Println("Reading storm.yml file at: " + filename)
	yamlStr, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println("Saving a backup copy in: " + filename + "-orig")
	_ = ioutil.WriteFile(filename+"-orig", yamlStr, 0644)

	log.Println("Parsing YAML...")
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(yamlStr), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	{
		delete(m, "storm.zookeeper.servers")

		zookeeperServers := []string{}

		serverNumber := 1
		stop := false
		for serverNumber < 15 && !stop {
			hostEnvVar := "ZK_SERVER_" + strconv.Itoa(serverNumber) + "_SERVICE_HOST"
			log.Println("Checking for Zookeeper host in environment variable: " + hostEnvVar)
			hostEnv := os.Getenv(hostEnvVar)
			if len(hostEnv) == 0 {
				log.Println("Not Found. Stopping here")
				stop = true
			} else {
				zookeeperServers = append(zookeeperServers, hostEnv)
				log.Printf("Found server: %v\n", hostEnv)
			}
			serverNumber += 1
		}

		if len(zookeeperServers) > 0 {
			m["storm.zookeeper.servers"] = zookeeperServers
		} else {
			log.Println("WARNING - No Zookeeper servers found. Keeping as the default localhost")
		}
	}

	{
		log.Println("Checking for zookeeper root in env APACHE_STORM_ZK_ROOT")
		zkRoot := os.Getenv("APACHE_STORM_ZK_ROOT")
		if len(zkRoot) == 0 {
			log.Println("No zookeeper root found, keeping the default")
		} else {
			log.Println("Using zookeeper root " + zkRoot)
			m["storm.zookeeper.root"] = zkRoot
		}
	}

	{
		log.Println("Checking for nimbus server in env APACHE_STORM_NIMBUS_SERVICE_PORT")
		thriftPort := os.Getenv("APACHE_STORM_NIMBUS_SERVICE_PORT")
		if len(thriftPort) == 0 {
			log.Println("No nimbus thrift port found, keeping the default")
		} else {
			log.Println("Using nimbus thrift port " + thriftPort)
			port, _ := strconv.Atoi(thriftPort)
			m["nimbus.thrift.port"] = port
		}
	}

	{
		log.Println("Checking for nimbus server in env APACHE_STORM_NIMBUS_SERVICE_HOST")
		nimbusServer := os.Getenv("APACHE_STORM_NIMBUS_SERVICE_HOST")
		if len(nimbusServer) == 0 {
			log.Println("No nimbus servers found, keeping the default")
		} else {
			log.Println("Using nimbus server " + nimbusServer)
			
			if os.Getenv("STORM_CMD") == "nimbus" {
			        m["storm.local.hostname"] = nimbusServer
			} else {
			        m["nimbus.seeds"] = []string{nimbusServer}
			}
		}
	}
	
	if os.Getenv("STORM_CMD") == "ui" {
        	m["ui.port"] = 8080
	}
	
	if os.Getenv("STORM_CMD") == "drpc" {
        	m["drpc.port"] = 3772
	}

	{
		log.Println("Trying to save modified config")
		d, err := yaml.Marshal(&m)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = ioutil.WriteFile(filename, d, 0644)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}
