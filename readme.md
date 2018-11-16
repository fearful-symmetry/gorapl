# goRAPL

A dead-simple, low-level API for accessing the Intel RAPL API.



This is a rather experimental library that's early on in it's life, so it isn't very useful yet! It has a lot of shortcomings, and a lot of things that needs to do, but can't. Most notably, it can only handle systems with one CPU socket.


# Big questions, TOD0s:

- msr_safe support
- How we want to handle multiple CPU packages? Should that be left to the user?
- exactly how low-level do we want this to be?

