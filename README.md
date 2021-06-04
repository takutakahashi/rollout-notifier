# rollout-notifier
Notify a completion of rollout for applications on kubernetes

## Installation

1. Install base components
Use kubectl:
```
$ kubectl apply -f https://github.com/takutakahashi/rollout-notifier/releases/download/v1.1.0/release.yaml
```

2. Create configuration secret

```
$ cat config.yaml
channelID: foobar # channel ID for Slack
token: xoxb-foo-bar # Slack app token
$
$ kubectl create secret generic rollout-notifier-secret --from-file=config.yaml
```

Latest release version: see below

https://github.com/takutakahashi/rollout-notifier/releases

## Usage

Once finished the setup, all of Deployment rollouts in all namespaces are notified to Slack automatically.  

## Features

1. Notify with Comment

When you add an annotiation `rollout-notifier.io/comment`, you can add a comment with Notification post.

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
  namespace: bar
  annotations:
    rollout-notifier.io/comment: "hello with rollout-notifier!"
...snip...
```
