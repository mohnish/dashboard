package component

import (
  "io/ioutil"
  "encoding/json"
  "time"
)

type Component struct {
  Id string `json:"id"`
  App string `json:"app"`
  Category string `json:"category"`
  Environment string `json:"environment"`
  Branch string `json:"branch"`
  Committer string `json:"committer"`
  Commit string `json:"commit"`
  Status string `json:"status"`
  UpdatedAt string `json:"updatedAt"`
}

const (
  StorageType string = "json"
  Dir string = "storage"
)

func (c *Component) Save() error {
  c.UpdatedAt = time.Now().Format(time.UnixDate)
  return ioutil.WriteFile(c.FileName(), []byte(c.String()), 0600)
}

func (c *Component) Load() {
  return
}

func (c *Component) String() string {
  str, _ := json.Marshal(c)

  return string(str)
}

func (c *Component) FileName() string {
  return Dir + "/" + c.Id + "." + StorageType
}
