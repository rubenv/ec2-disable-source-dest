# ec2-disable-source-dest

Small helper to disable the EC2 source/dest check from within an instance.

## Why?

Imagine you're setting up a [Kubernetes](http://kubernetes.io/) cluster, one
that scales automatically. On [Amazon EC2](https://aws.amazon.com/ec2/), you'd
do that with an autoscaling group.

Depending on your network configuration, you'll need to disable source/dest
checking on your instance. You can do this manually in the config panel. But
you can't do it in the launch configuration of the autoscaling group.

This little helper reconfigures your instances once they have booted, so you
don't have to worry about it.

## Usage

Simply run it inside your EC2 instance:

```
docker run --rm rubenv/ec2-disable-source-dest
```

You'll need to set up an appropriate IAM role/policy that is capable of using
`ec2:ModifyInstanceAttribute`.

## Running on boot (CoreOS)

Add the following unit file to your cloud-config:

```
[Unit]
Description=Disable EC2 source/dest check
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
Restart=on-failure
ExecStart=/usr/bin/docker run --rm rubenv/ec2-disable-source-dest
```
