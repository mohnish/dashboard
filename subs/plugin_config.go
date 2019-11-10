package subs

type PluginConfig struct {
  Id string `json:"id"`
  Name string `json:"name"`
  Url string `json:"url"`
  Method string `json:"method"`
  Headers map[string]string `json:"headers"`
  Interval string `json:"interval"`
}
