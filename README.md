Panic Room
===========

Panic room watches your sensitive files and panics when someone tries to mess with them.

Features
--------

* Watches named groups of files and sends alerts if they are modified
* Allows for fine-grained selection of watcher rules.
* Supports AWS SNS for sending notifications to multiple channels. (e.g. Email, SMS, Webhooks)

Configuration
-------------

````yaml
watchers:
  - name: My Wordpress Site Root
    paths:  ["/var/www/**/*"]
    excludes: ["/var/www/wp-uploads/**/*"]
    alerters: ["sysops"]
  - name: All wordpress code
    paths: ["/var/www/**/*.php"]
    alerters: ["weblog"]
  - name: /etc/passwd
    paths: ["/etc/passwd"]
    alerters: ["sysops"]
  
alerters:
  - name: sysops
    type: sns
    config:
      aws_region: eu-west-1
      topic_arn: arn:aws:sns:eu-west-1:12345678:panicroom-sysops
  - name: weblog
    type: log
    config:
      path: /data/log/wordpress_files.log
````


Usage
-----

    ./panicroom --config-path <CONFIG_PATH>
    

### Example: Filtering messages delivered to SNS

It is possible to filter SNS messages on a per-subscriber basis based on the operation name.

The following operation names exist:

* `REMOVE`
* `RENAME`
* `CHMOD`
* `CREATE`

These are sent in the `"operation"` message attribute. To filter a subscription to only created or removed files:

````json
{
    "operation": ["REMOVE","CREATE"]
}
````


Known Limitations
-----------------

1. Dot not support watching directories directly.
2. Does not properly handle newly created files.
3. SNS alerters will send an alert for each single event which can quickly flood notification channels.


Planned Features
----------------

[ ] Watch directories for new files and add matching files to watcher lists and handle their notifications
[ ] SNS Alerter: Group alerts together and push a single notification of all events after a short hysteresis (5s)