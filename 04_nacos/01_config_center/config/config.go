package config

/*
	{
	    "age": 10,
	    "name": "Danny"
	}
*/
type ServerConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Age  int    `mapstructure:"age" json:"age"`
}
