## Limit Memory Usage For Each PHP-FPM Pool

To limit memory usage for each PHP-FPM pool, you need to manage resources at two levels: the PHP level and the PHP-FPM pool level. The `memory_limit` directive in the php.ini file or in your PHP scripts controls the maximum amount of memory that a script is allowed to allocate, which works at the individual script level. However, to control memory usage at the PHP-FPM pool level, you will have to use a combination of PHP-FPM configuration and potentially operating system-level controls.

### Step 1: Configure `memory_limit` in PHP
First, ensure that the `memory_limit` directive is appropriately set in your php.ini file or within your PHP scripts. This setting limits the amount of memory a single script can consume, which indirectly contributes to managing the pool's memory usage.

For example, in php.ini:
```ini
memory_limit = 128M
```

This setting restricts each PHP script to a maximum of 128 MB of memory.

### Step 2: Limit PHP-FPM Pool Processes
To manage the memory usage of an entire PHP-FPM pool, you can limit the number of child processes that the pool can spawn. This, combined with the `memory_limit` setting, can help you control the overall memory usage.

Edit your pool configuration file (usually found in /etc/php/7.x/fpm/pool.d/ where 7.x should be replaced with your PHP version, and each pool has its own configuration file, typically named after the pool). Adjust the following settings:

- `pm` - This directive defines the process manager to use. Options include static, dynamic, or ondemand.
- `pm.max_children` - This sets the maximum number of child processes to spawn. Limiting this number helps control memory usage.
- `pm.start_servers`, `pm.min_spare_servers`, `pm.max_spare_servers` - These settings are relevant for dynamic and ondemand process managers and help manage the number of idle and active processes.
For example, to limit a pool to use a maximum of 10 child processes:

```ini
pm = dynamic
pm.max_children = 10
pm.start_servers = 2
pm.min_spare_servers = 2
pm.max_spare_servers = 4
```

### Step 3: Operating System-Level Controls

For a more stringent control over memory usage, consider using operating system-level tools like `cgroups` on Linux. `cgroups` (Control Groups) allows you to allocate, prioritize, deny, manage, and monitor system resources, like CPU, memory, disk I/O, etc., for a group of processes.

To use `cgroups` to limit memory for a PHP-FPM pool, you would:

1. Create a new `cgroup` for your PHP-FPM pool.
2. Set memory limits on that `cgroup`.
3. Assign the PHP-FPM master process to that `cgroup`.

Example commands to set up a `cgroup` with a memory limit:

```bash
cgcreate -g memory:/phpfpm
cgset -r memory.limit_in_bytes=500M phpfpm
cgclassify -g memory:phpfpm $(pidof php-fpm)
```

These commands create a cgroup named phpfpm, set a 500 MB memory limit, and assign the PHP-FPM process to it. Note that pidof php-fpm should match the process name of your PHP-FPM master process, which might differ based on your PHP version and configuration.



## Limit a System User From Using More Than 10GB of Disk Storage

To limit a system user from using more than 10GB of disk storage, you can utilize the Linux disk quota system. This system allows you to allocate disk usage limits for users and groups, preventing them from exceeding a specified amount of disk space. Here's how to set it up:

### Step 1: Install Quota Tools
First, ensure that the quota tools are installed on your system. On a Debian-based system, you can install them using:

```bash
sudo apt-get update
sudo apt-get install quota
```

For Red Hat-based systems, use:

```bash
sudo yum install quota
```

### Step 2: Configure Filesystem for Quotas
You need to enable quotas on the filesystem where the user's files are stored. This is typically done by editing the `/etc/fstab` file.

1. Open `/etc/fstab` with your preferred text editor, e.g., sudo nano `/etc/fstab`.
2. Find the line corresponding to the filesystem (e.g., `/home` if the user's data is stored in their home directory).
3. Add `,usrquota` (or `usrjquota=aquota.user,jqfmt=vfsv0` for journaled quotas) to the mount options.

For example, if your `/home` partition is on `/dev/sda1`, the entry might look like this before:
```bash
/dev/sda1 /home ext4 defaults 0 2
```

And like this after:

```bash
/dev/sda1 /home ext4 defaults,usrquota 0 2
```

Remount the filesystem and check that quotas are enabled:
```bash
sudo mount -o remount /home
```

### Step 3: Initialize and Assign Quotas
Initialize the quota database files on the filesystem with the quotacheck command:

```bash
sudo quotacheck -cum /home
```

Turn on quotas with the quotaon command:
```bash
sudo quotaon -v /home
```

Set the quota for the specific user. Replace username with the actual system username:
```bash
sudo setquota -u username 10240 10240 0 0 /home
```

This command sets both the soft and hard limits to 10GB (specified in kilobytes), with no limits on inodes (file count). The soft limit allows temporary exceeding of the quota, while the hard limit is strictly enforced.

### Step 4: Verify Quota Assignment
Verify the quota is correctly assigned to the user:

```bash
sudo repquota /home
```

Or for a specific user:

```bash
quota -u username
```
