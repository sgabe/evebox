MikroTik
========

Limited RouterOS API Integration
--------------------------------

This fork of EveBox contains an experimental feature that allows to add or
remove the source IP address of an alert to or from an address list on a
MikroTik device. It does so by invoking the RouterOS API when an alert is
escalated or de-escalated. This feature is disabled by default and must be
explicitly enabled in the EveBox Server configuration file. Note that
currently only the built-in SQLite datastore is supported.

Requirements
------------

- Enable RouterOS API and set ``Allowed Address`` according to your needs.
- Create new user group with the following policies:
    - ``api`` to login via API.
    - ``read`` to get the ID of an IP address.
    - ``write`` to add the IP address to the address list.
- Add new user to the previously created user group and set ``Allowed Address``
  to the IP address of the EveBox Server.
- Update the EveBox Server configuration file with the API credentials.

RouterOS API Sentences
----------------------

Add IP address to address list:

::

  /ip/firewall/address-list/add
  =list=<LIST>
  =address=<ADDRESS>
  =comment=<COMMENT>
  =disabled=no

Get ID of IP address in address list:

::

  /ip/firewall/address-list/print
  ?list=<LIST>
  ?address=<ADDRESS>
  =.proplist=.id

Remove IP address from address list:

::

  /ip/firewall/address-list/remove
  "=.id=<ID>
