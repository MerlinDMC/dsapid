A public dsapi / imgapi server
==============================
to publish datasets and images for SmartOS

The primary goal here is it to use the server with `imgadm` on SmartOS. The API is functional for that purpose and every aspect of the official API that is not needed for `imgadm` has second priority.

Version **0.6.2** is tested with `imgadm` on platform 20140111T020931Z.

Differences from the official dsapi
-----------------------------------
- all images on the server are public and can be downloaded by everyone
- maybe not all API methods are implemented
- definitely the upload machanism is different

Why did we build this?
----------------------
*there is already the official repository for datasets so this one is useless*

**NO** - ever tried to publish own images onto the official server? The community builds datasets from time to time and it's a bummer that those are not usable by a larger group of the community

TODO
----
- Cleaner / easier upload mechanism
- ACL for images
- Stats

Version history
---------------
0.6.4: changes to the syncer to accept a provider tag from dsapi sources
0.6.2: make code repository public

Changes from earlier versions (0.1.x)
-------------------------------------
- rewritten in Go
- no database server anymore: all action saved in plain files
- added export possibility: download a tarball of an image to `imgadm install` somewhere
