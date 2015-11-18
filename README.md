# Brightbox Cloud driver for Docker Machine

![](/docs/img/logo.png)

Install this driver in your PATH and you can create docker hosts with
ease on [Brightbox Cloud](https://www.brightbox.com).

If you are new to Brightbox Cloud you can [sign up in 2
minutes](https://manage.brightbox.com/signup) and get your user
credentials. You'll get a £20 credit to get you started creating
Docker hosts.

## Installation

### From a Release

Official release versions of the driver include a binary for Linux,
MacOS and Windows. You can find them on the [GitHub releases
page](https://github.com/brightbox/docker-machine-driver-brightbox/releases).

Pick the binary you require, download it into a directory on your
PATH as a file called `docker-machine-driver-brightbox` and make it
executable.

Linux:

```
curl -sSL -o ~/bin/docker-machine-driver-brightbox \
https://github.com/brightbox/docker-machine-driver-brightbox/releases/download/v0.0.4/bin.docker-machine-driver-brightbox_linux-amd64 && \
chmod 755 ~/bin/docker-machine-driver-brightbox

```

Mac OSX

```
sudo curl -sSL -o /usr/local/bin/docker-machine-driver-brightbox https://github.com/brightbox/docker-machine-driver-brightbox/releases/download/v0.0.4/bin.docker-machine-driver-brightbox_darwin-amd64 &&
sudo chmod 755 /usr/local/bin/docker-machine-driver-brightbox
```

### From Source

To build and install, first clone this repo onto a server running Docker,
then run:

```
$ make containerbuild && sudo make install
```

which will install the driver into `/usr/local/bin`

## Obtaining credentials

To use the driver you will need to create a 'Full Privileges' [API client
 within Brightbox
 Manager](http://manage.brightbox.com/account/api_clients). If you're
[collaborating](https://www.brightbox.com/docs/reference/collaboration/)
with other Brightbox Users, make sure you select the
 correct account first or, if you have insufficient privileges, obtain
 the API client details from the owner of the account.

## Using the driver

To use the driver first make sure you are running at least [version
0.5.1 of `docker-machine`](https://github.com/docker/machine/releases).

```
$ docker-machine -v
docker-machine version 0.5.1 (7e8e38e)
```

Check that `docker-machine` can see the Brightbox driver by asking for
the driver help.

```
$ docker-machine create -d brightbox | more
Usage: docker-machine create [OPTIONS] [arg...]

## Create a machine.

Specify a driver with --driver to include the create flags for that driver in the help text.

Options:

   --brightbox-account 								Brightbox Cloud Account to operate on [$BRIGHTBOX_ACCOUNT]
   --brightbox-api-url "https://api.gb1.brightbox.com/"				Brightbox Cloud Api URL for selected Region [$BRIGHTBOX_API_URL]
...
```

To create a machine you'll need the API Client details you obtained earlier. 

Then creating a Docker host is as simple as

```
$ docker-machine create -d brightbox --brightbox-client cli-xyzab \
--brightbox-client-secret Secret example
Running pre-create checks...
Creating machine...
Waiting for machine to be running, this may take a few minutes...
Machine is running, waiting for SSH to be available...
Detecting operating system of created instance...
Provisioning created instance...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
To see how to connect Docker to this machine, run: docker-machine env example
```

or if you don't want your client secret displayed anywhere:

```
$ (export BRIGHTBOX_CLIENT_SECRET;
read -s -p "Enter Secret: " BRIGHTBOX_CLIENT_SECRET && \
docker-machine create -d brightbox \
--brightbox-client cli-xyzab  example)
```

This creates a small server in the default
[server group](https://www.brightbox.com/docs/guides/cli/server-groups/) for the
account, and accesses the server over IPv6.

If you are running `docker-machine` on another server on the same Brightbox
Cloud account and are using the default firewall policy, then this will work
straight away. If you're running `docker-machine` from elsewhere then you'll
need to
[alter your firewall policy](https://www.youtube.com/watch?v=Q3eYMV_hbDk&hd=1)
first to allow access in to Docker on TCP port 2376. Make sure this rule is
tight as there is no authentication on the Docker port.

## Changing the settings

The driver has several options that you can use to get precisely the
Docker host you want. You can see them all in the help list by running
`docker-machine create -d brightbox | more`

Here are the most useful options:

*   `--brightbox-type`

    By default `docker-machine` creates a small 1gb SSD server as the
    Docker host. If you want a larger one, check the [server sizing
    page](https://www.brightbox.com/pricing/#full-pricing-table) for
    the available sizes, and then specify the memory size plus either
    `.ssd` or `.ssd-high-io` (for the larger disk version). So if you
    want a 4GB server just use `4gb.ssd` for this option.

    For more details on the available ids and handles [use the
    CLI](https://www.brightbox.com/docs/guides/cli/installation/)
    `brightbox types` command

*   `--brightbox-image`

    You can select the image you want to use for the Docker host by specifying
    the `img-xxxxx` id of the image you require. You can get the image id from
    the Image Library in [Brightbox Manager](https://manage.brightbox.com) or
    [via the CLI](https://www.brightbox.com/docs/guides/cli/image-library/). Remember
    that Docker requires a 64-bit operating system.

*   `--brightbox-group`

    You can add [server groups, and therefore firewall
    policies](https://www.brightbox.com/docs/guides/cli/firewall/)
    using the `--brightbox-group` option. Remember firewall policies
    are cumulative on Brightbox Cloud and specifying groups
    replaces the default option of putting the server in the default
    group.

*   `--brightbox-zone`

    Every
    [Region](https://www.brightbox.com/docs/reference/glossary/#region)
    on Brightbox Cloud has [multiple availability
    zones](https://www.brightbox.com/docs/reference/glossary/#zone)
    within it. Normally the default auto-allocation does the right thing
    but if you want specific placement then specify the zone id or handle
    with this option.

    For more details on the available ids and handles [use the
    CLI](https://www.brightbox.com/docs/guides/cli/installation/)
    `brightbox zones` command

*   `--brightbox-ipv4`

    This is a flag that makes `docker-machine` access the server over IPv4
    rather than IPv6. Brightbox servers run on a private IPv4 network by
    default, so this will stop external access to the server from unless you map
    a Cloud IP in [Brightbox Manager](https://manage.brightbox.com) or
    [via the CLI](https://www.brightbox.com/docs/guides/cli/cloud-ips/).

## Help

If you need help using this driver, drop an email to support at brightbox
dot com.

## License

This code is released under an MIT License.

Copyright (c) 2015 Brightbox Systems Ltd.
