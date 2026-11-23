fix(sandboxes): make the worker's start check hold up under packet loss

The check now has to pass twice in a row and waits longer for each address, so a worker starts a few seconds later.
