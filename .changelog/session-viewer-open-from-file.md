fix(sessions): keep the viewer's one-time link out of the browser opener's arguments

`ant beta:sessions connect --web` now opens a temporary file that only you can read, and that file sends the browser on to the viewer. If the browser shows an error instead, because it can't read the temp directory, paste the link printed in the terminal into it.
