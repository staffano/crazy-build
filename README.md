# crazy-build 
How complicated can you make building the simplistic _hello world_ program?

In this project we present a home brewed build system based on _go_ to build something we call _artifacts_.

It is conceived as a prototype for what could be done when combining go and docker in a build system. And as learning experience for a poor Ada-centric relic of a programmer.

## Description
The build system resembles systems like _bitbake_, where the artifact corresponds to recipes and the tasks in bitbake responds to what we call _commands_ in crazy_build.

In crazy-build we utilize the fact that the build structure of a module is fixed during most of its life cycle and could therefore be represented as a compiled program. During the initial phase of a module the build structure could however be more volatile, but since we're using _go_ for compiling the binary build structure, it doesn't matter so much due to the quick development cycle of go.

### Dependency handling

None.

### Testing

None.

### Usefulness

None.