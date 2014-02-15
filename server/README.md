User Roles
==========

Currently this is a intermediate state and just contains simple "roles" to ensure uploads can work.

Datasets
--------

**s_dataset.upload**
:   The user may upload new datasets.
    The owner field inside the given manifest will be replaced by the uploading user.
    Newly uploaded datasets will be added as **deactivated**.
    The dataset needs to be approved and activated.

**s_dataset.manage**
:   The user may upload new datasets or activate/deactivate already uploaded ones.
    Newly uploaded datasets will be added as **deactivated**.

**s_dataset.admin**
:   The user may upload datasets, delete old ones and can activate/deactivate other datasets.
    Newly uploaded datasets will be added as **activated**.
