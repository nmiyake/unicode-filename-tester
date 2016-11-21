unicode-filename-tester
=======================
Tests whether the file system in which the program is run supports writing distinct files that are equivalent from a
Unicode perspective. Useful to determine whether or not a file system performs Unicode normalization for precomposed or
decomposed characters.

Usage
-----
Install:

```
go get github.com/nmiyake/unicode-filename-tester
```

Run:

```
unicode-filename-tester
```

Run with verbose output:

```
unicode-filename-tester -v
```

Implementation
--------------
Creates a temporary directory in the working directory and attempts to write out two files: `ö.txt` (`U+00F6`) with
content "composed" and `o¨.txt` (`U+006F` + `U+0308`) with content "decomposed". Reads the files after writing them and
exits with a non-zero exit code if the content is the same and exits with an exit code of 0 if the content differs.

Most file systems use the file name as supplied exactly, and thus files with the provided names should be distinct.
However, the HFS+ file system performs NFD normalization on Unicode in file names, and thus these two representations
are both normalized to the same name and this program "fails".

This was written as a test to see if the [Apple File System (APFS)](https://en.wikipedia.org/wiki/Apple_File_System)
performs Unicode normalization in the same manner as HFS+. Luckily, it looks like this is not the case :).
