package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"auth-get/pkg/auth"
)

var (
	rootCmd = &cobra.Command{
		Use: 	"auth-get",
		Short: 	"A tool for getting auth file.",
		Long: 	"A tool for getting host authentication file. \nWork by compass machine auth api." ,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := auth.NewClient(MasterIP, Username, Password, Sshport)
			nodesMap, err := c.GetNodes()
			if err != nil {
				logrus.Errorf("get nodes error: %s\n", err)
				return err
			}
			logrus.Debugf("nodesMap: %v\n", nodesMap)
			switch {
			case Output == "yaml":
				if err = c.OutputYaml(nodesMap); err != nil {
					logrus.Errorf("output error: %s\n", err)
					return err
				}
			case Output == "json":
				if err = c.OutputJson(nodesMap); err != nil {
					logrus.Errorf("output error: %s\n", err)
					return err
				}
			case Output == "inventory":
				if err = c.OutputInventory(nodesMap); err != nil {
					logrus.Errorf("output error: %s\n", err)
					return err
				}
			}
			return nil
		},
	}
	MasterIP, Username, Password string
	Sshport string
	Output string
)


func init()  {
	//set loglevel
	logrus.SetLevel(logrus.DebugLevel)
	
	// cobra flag
	rootCmd.PersistentFlags().StringVarP(&MasterIP, "masterIp", "m", "",
		"Master VIP or control cluster's master IP ")
	rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "admin",
		"Compass username")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "Pwd123456",
		"Master VIP or control cluster's master IP ")
	rootCmd.PersistentFlags().StringVarP(&Sshport, "port", "", "22",
		"Set SSH port in output file ")
	rootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "yaml",
		"output auth file by yaml, json or ansible inventory (yaml, json, inventory) ")


}
func Execute()  {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}