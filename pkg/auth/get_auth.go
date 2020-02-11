package auth

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
	yaml "gopkg.in/yaml.v2"
	"text/template"
)

func NewClient(masterIP, username, password, sshport string) *client {
	conf := config{
		masterIP:   masterIP,
		username: 	username,
		password: 	password,
		sshport:    sshport,
	}
	return &client{
		httpClient:		&http.Client{},
		config: 		conf,
	}
}

// get nodes list
func (c *client)GetNodes() (map[string][]string, error) {
	url := "http://"+c.config.masterIP+":6002"+APIEndpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.config.username, c.config.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("get api error %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()


	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("read error %s\n", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("get response status is not 200:[%s] %s\n", resp.StatusCode, data)
		return nil, err
	}

	nodesMap := make(map[string][]string)

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		c, _ := jsonparser.GetString(value, "spec", "cluster")
		node, _ := jsonparser.GetString(value, "spec", "address", "[0]", "address")
		nodesMap[c] = append(nodesMap[c], node)
	}, "items")
	if err != nil {
		err = fmt.Errorf("parse json error %s\n", err)
		return nil, err
	}

	return nodesMap, nil
}

// get node auth
func (c *client)getNodeAuth(node string) (map[string]string, error) {
	url := "http://"+c.config.masterIP+":6002"+APIEndpoint+"/"+node+"/auth"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.config.username, c.config.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("get api error %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("read error %s\n", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("get response status is not 200:[%s] %s\n", resp.StatusCode, data)
		return nil, err
	}
	user, err := jsonparser.GetString(data, "user")
	if err != nil {
		err = fmt.Errorf("parse json error %s\n", err)
		return nil, err
	}
	password, err := jsonparser.GetString(data, "password")
	if err != nil {
		err = fmt.Errorf("parse json error %s\n", err)
		return nil, err
	}
	key, err := jsonparser.GetString(data, "key")
	if err != nil {
		err = fmt.Errorf("parse json error %s\n", err)
		return nil, err
	}
	return map[string]string{
		"host": 		node,
		"user": 		user,
		"password": 	password,
		"key": 			key,
	}, nil

}

// output

func (c *client)output(nodesMap map[string][]string) (map[string]map[string][]string, error) {
	outMap := make(map[string]map[string][]string)
	for k, v := range nodesMap {
		authList := make([]string, 0)
		for _, no := range v {
			auth, err := c.getNodeAuth(no)
			if err != nil {
				err = fmt.Errorf("get node auth error %s\n", err)
				return nil, err
			}
			logrus.Infof("get auth: [%s] %v", no, auth)
			//TODO handle auth by private key
			authList = append(authList, fmt.Sprintf("%s %s %s %s", no, c.config.sshport, auth["user"], auth["password"]))
		}
		outMap[k] = map[string][]string{
			"auth": authList,
		}
	}
	return outMap, nil
}

func (c *client)OutputJson(nodesMap map[string][]string) error {
	outMap, err := c.output(nodesMap)
	if err != nil {
		err = fmt.Errorf("get output error: %s\n", err)
		return err
	}
	jsonOutput, err := json.Marshal(outMap)
	if err != nil {
		err = fmt.Errorf("json marshal error %s\n", err)
		return err
	}
	fmt.Printf("%s\n",jsonOutput)

	return nil
}

func (c *client)OutputYaml(nodesMap map[string][]string) error {
	outMap, err := c.output(nodesMap)
	if err != nil {
		err = fmt.Errorf("get output error: %s\n", err)
		return err
	}
	yamlOutput, err := yaml.Marshal(outMap)
	if err != nil {
		err = fmt.Errorf("yaml marshal error %s\n", err)
		return err
	}
	fmt.Printf("%s\n",yamlOutput)

	return nil
}

func (c *client)OutputInventory(nodesMap map[string][]string) error {
	outMap, err := c.outputInvMap(nodesMap)
	if err != nil {
		err = fmt.Errorf("get output error: %s\n", err)
		return err
	}

	tmplMap := map[string]interface{}{
		"sshport": c.config.sshport,
		"outMap": outMap,
	}

	tmpl, err := template.New("inventory").Parse(inventory)
	if err != nil {
		err = fmt.Errorf("create template failed: %s:\n", err)
		return err
	}
	if err = tmpl.Execute(os.Stdout, tmplMap); err != nil {
		err = fmt.Errorf("render template failed: %s\n", err)
		return err
	}
	return nil
}


func (c *client)outputInvMap(nodesMap map[string][]string) (map[string][]map[string]string, error) {
	outMap := make(map[string][]map[string]string)
	for k, v := range nodesMap {
		authList := make([]map[string]string, 0)
		for _, no := range v {
			auth, err := c.getNodeAuth(no)
			if err != nil {
				err = fmt.Errorf("get node auth error %s\n", err)
				return nil, err
			}
			logrus.Infof("get auth: %v", auth)
			//TODO handle auth by private key
			authList = append(authList, auth)
		}
		outMap[k] = authList
	}
	return outMap, nil
}