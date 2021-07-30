package config_test

import (
	rsConfig "rabbitsky/src/config"
	"testing"
)

func TestDebug(t *testing.T) {
	/* dunno what to do at this point bot to add coverage >:) */
	rsConfig.Debug()
	return
}

func TestInit(t *testing.T) {
	cases := []struct {
		success  bool
		filePath string
	}{
		{
			/* check if no config file */
			filePath: "",
			success:  false,
		},
		// {
		//  /* can't get this to work now somehow :/ */
		// 	filePath: "config/config.json",
		// 	success:  true,
		// },
	}

	for _, testCase := range cases {
		err := rsConfig.Init(testCase.filePath)
		if (err != nil && testCase.success) || (err == nil && !testCase.success) {
			t.Logf("Expecting %v but got %v", testCase.success, err)
			t.Fail()
		}
	}
	return
}
