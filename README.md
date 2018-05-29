# ec2-disable-source-dest

Small helper to [disable the network interface source/destination
check](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-eni.html#change_source_dest_check)
from within an AWS EC2 instance.

## Why?

Imagine you're setting up a [Kubernetes](http://kubernetes.io/)
cluster, one that scales automatically. On [Amazon
EC2](https://aws.amazon.com/ec2/), you'd do that with an autoscaling
group.

Depending on your network configuration, you may need to disable
source/destination checking on your instance. You can do this manually
in the console or at launch time for individual instances, but you
can't do it in the launch configuration of an autoscaling group.

This helper program reconfigures your instances once they have booted,
so you don't have to worry about it.

## Usage

Simply run it inside your EC2 instance:

```shell
docker run  \
  --net=host \
  --rm \
  rubenv/ec2-disable-source-dest
```
Note that this image does include the SSL certificates required to
trust the AWS EC2 API hosts, but if you'd rather use your own, you can
[mount certificates you trust into the
container](https://docs.docker.com/storage/bind-mounts/#use-a-read-only-bind-mount)
and set the "SSL_CERT_FILE" environment variable to point to your
mounted file instead:

```shell
docker run  \
  --mount type=bind,source=/etc/ssl/certs/ca-certificates.crt,destination=/etc/ssl/certs/ca-certificates.crt,readonly \
  --env 'SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt' \
  --net=host \
  --rm \
  rubenv/ec2-disable-source-dest
```

You'll need to set up an appropriate IAM role/policy for these
instances that is capable of using [the `ec2:ModifyInstanceAttribute`
action](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_ModifyInstanceAttribute.html).

## Running on boot (Container Linux)

Add the following systemd unit file to your cloud-config or Ignition
configuration:
```
[Unit]
Description=AWS EC2 instance network configuration adjustment
After=network-online.target
Requires=network-online.target
ConditionPathExists=!/.disabled-src-dest-check

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/sudo /usr/bin/rkt --insecure-options=image run \
  --net=host \
  --dns=host \
  docker://rubenv/ec2-disable-source-dest
ExecStartPost=/usr/bin/touch /.disabled-src-dest-check
ExecStartPost=-/usr/bin/sudo /usr/bin/rkt gc --grace-period=0s

[Install]
WantedBy=multi-user.target
```
Consider pulling [the _rubenv/ec2-disable-source-dest_ image](https://hub.docker.com/r/rubenv/ec2-disable-source-dest/) ahead of
time when building your AMI:
```shell
#!/bin/sh

set -e -u

for image in \
  docker://rubenv/ec2-disable-source-dest;
do
  case "${image}" in
    docker://*) opts='--insecure-options=image';;
    *)          opts=;;
  esac
  rkt "${opts}" fetch "${image}"
done
```
Note that it is possible to go even further and use [_rkt prepare_](https://coreos.com/rkt/docs/latest/subcommands/prepare.html) when
building your AMI and [_rkt
run-prepared_](https://coreos.com/rkt/docs/latest/subcommands/run-prepared.html)
at boot time in lieu of _rkt run_, but that allows the systemd unit to
run no more than once&mdash;absent another intervening invocation of
_rkt prepare_. If it fails, it's then harder to run it again.
