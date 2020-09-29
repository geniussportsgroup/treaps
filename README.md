# Golang Oredered Set

This package implements ordered sets implemented through Treaps. It accepts any type of key. 

The user must specify a function `less(k1, k2)` which receives a pair of keys and determine whether `k1` is or not less than `k2`. In the `less()` function, the user can specify access to the key; this is the way as the package can handle arbitrary types of keys.

Given a set, the following operations are supported:
- Insertion, search, and deletion in O(log n) expected case. The package can handle repeated keys.
- To know the i-th key inside the order in O(log n) expected case.
- To know what is the inorder position of any key in O(log n) expected case.
- To split the set by a key's value or at any specific position in O(log n) expected case.
- To extract any rank [i, j] of the set in O(log n) where i and j are position respect to the order.
- Union and interception of sets in O(m log n).
- Join of rank disjoint sets in O(max(log n, log m)).

## Installation

    go get -u 
