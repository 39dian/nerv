package cmd

import (
	"fmt"
	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"github.com/ChaosXu/nerv/lib/env"
	"github.com/ChaosXu/nerv/cmd/cli/lib"
	"github.com/ChaosXu/nerv/lib/credential"
)

func init() {

	var credential = &cobra.Command{
		Use:    "credential [command] [flags]",
		Short:    "Manage credential",
		Long:    "Manage credential",
		RunE: credentialFn,
	}
	NervCmd.AddCommand(credential)

	//list
	var list = &cobra.Command{
		Use:    "list",
		Short:    "List all credentials",
		Long:    "List all credentails",
		RunE: listCredentials,
	}
	list.Flags().StringVarP(&flag_config, "config", "c", "../config/config.json", "The path of config.json. Default is ../config/config.json ")
	credential.AddCommand(list)

	//get
	var get = &cobra.Command{
		Use:    "get",
		Short:    "Get a credential",
		Long:    "Get all credential",
		RunE: getCredential,
	}
	get.Flags().UintVarP(&flag_id, "id", "i", 0, "Credential id")
	get.Flags().StringVarP(&flag_config, "config", "c", "../config/config.json", "The path of config.json. Default is ../config/config.json ")
	credential.AddCommand(get)


	//create
	var create = &cobra.Command{
		Use:    "create",
		Short:    "Create a credential",
		Long:    "Create a credential",
		RunE: createCredential,
	}
	create.Flags().StringVarP(&flag_data, "data", "d", "", "JSON of credential")
	create.Flags().StringVarP(&flag_config, "config", "c", "../config/config.json", "The path of config.json. Default is ../config/config.json ")
	credential.AddCommand(create)


	//delete
	var delete = &cobra.Command{
		Use:    "delete",
		Short:    "Delete a credential",
		Long:    "Delete a credential",
		RunE: deleteCredential,
	}
	delete.Flags().UintVarP(&flag_id, "id", "i", 0, "Credential id")
	delete.Flags().StringVarP(&flag_config, "config", "c", "../config/config.json", "The path of config.json. Default is ../config/config.json ")
	credential.AddCommand(delete)
}

func credentialFn(cmd *cobra.Command, args []string) error {
	return nil
}

func listCredentials(cmd *cobra.Command, args []string) error {
	//init
	env.InitByConfig(flag_config)
	gdb := lib.InitDB()
	defer gdb.Close()

	cres := []credential.Credential{}
	if err := gdb.Find(&cres).Error; err != nil {
		return err
	}

	fmt.Println("ID\tType\tName\tPassword\tCreateAt")
	for _, cre := range cres {
		fmt.Printf("%d\t%s\t%s\t%s\t%s\n", cre.ID, cre.Type, cre.Name, cre.Password, cre.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func getCredential(cmd *cobra.Command, args []string) error {
	env.InitByConfig(flag_config)
	gdb := lib.InitDB()
	defer gdb.Close()

	id := flag_id
	data := credential.Credential{}
	if err := gdb.First(&data, id).Error; err != nil {
		return err
	}
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(buf))
	return nil
}

func createCredential(cmd *cobra.Command, args []string) error {
	if flag_data == "" {
		return errors.New("--data -d is null")
	}
	fmt.Println(flag_data)

	//init
	env.InitByConfig(flag_config)
	db := lib.InitDB()
	defer db.Close()

	data := &credential.Credential{}

	err := json.Unmarshal([]byte(flag_data), data)
	if err != nil {
		return err
	}

	if err := db.Create(data).Error; err != nil {
		return err
	}

	fmt.Printf("Create credetial success. id=%d\n", data.ID)
	return nil
}

func deleteCredential(cmd *cobra.Command, args []string) error {
	if flag_id == 0 {
		return errors.New("--id -i is null")
	}
	//init
	env.InitByConfig(flag_config)
	gdb := lib.InitDB()
	defer gdb.Close()

	id := flag_id
	data := credential.Credential{}
	if err := gdb.First(&data, id).Error; err != nil {
		return err
	}

	if err := gdb.Unscoped().Delete(data).Error; err != nil {
		return err
	}

	fmt.Printf("Delete credential success. id=%d\n", id)

	return nil
}
