# CodeSync

I got interested in learning the internals of Git, so this is an attempt to recreate the internals of GIT ( mainly `git commit` ).

# Supported Commands

1. `cat-file` - Prints out the object given a SHA1 hash of a Git Object
2. `ls-tree` - Prints the contents of a Git Tree based on the SHA1 hash
3. `hash-object` - Creates a new Git Object ( Blob | Tree | Commit ) based on given contents
4. `show-ref` - Lists out Git References
5. `commit` - Commits the contents of the curret work directory [ Doesn't support `git add`. Commit considers every change in the work directory and creates Git Objects currently. ]
