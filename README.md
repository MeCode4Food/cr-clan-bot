# YAML config

- no ids allowed in the config
- no tokens allowed in the config
- values of sensitive/not allowed keys in the config replaced with ""
- sensitive config goes in secrets/\<environment\>.yaml
  - secrets/\<environment\>.yaml env gets merged during runtime
