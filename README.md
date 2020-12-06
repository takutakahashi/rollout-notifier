# rollout-notifier
Notify a completion of rollout for applications on kubernetes

## Usage

```
/daemon --namespace [namespace] --type [type] --config [config path]
```

|  parameter  |  means  |
| ---- | ---- |
|  namespace  |  The namespace with the deployment you want to notify  |
|  type  |  Notification type. slack and noop are available now |
|  config |  Configutation file path. config example are available in example/ |
