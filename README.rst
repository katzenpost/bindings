
.. image:: https://travis-ci.org/katzenpost/bindings.svg?branch=master
  :target: https://travis-ci.org/katzenpost/bindings

.. image:: https://godoc.org/github.com/katzenpost/bindings?status.svg
  :target: https://godoc.org/github.com/katzenpost/bindings

Language binding libraries
==========================

This repo contains golang which can be used with
Java and Python bindings.


dependencies
------------

* golang 1.9 or later

* gopy

  go get github.com/go-python/gopy

usage
=====

Note that you have to export ``GODEBUG`` variable in the execution environment in order for the bindings to work properly::

  GODEBUG=cgocheck=0


license
=======

AGPL: see LICENSE file for details.


supported by
============

.. image:: https://katzenpost.mixnetworks.org/_static/images/eu-flag-tiny.jpg

This project has received funding from the European Union’s Horizon 2020
research and innovation programme under the Grant Agreement No 653497, Privacy
and Accountability in Networks via Optimized Randomized Mix-nets (Panoramix).
