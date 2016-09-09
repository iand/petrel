// Package petrel provides a distributed hash table
/*

The design of Petrel is biased toward low latency writes at the expense of
lower performance for reads. Although it supports dynamic membership it is
assumed that the topology and number of members is generally stable over the
medium term.

Nodes in a Petrel hash table are organised in a classic Chord-like ring. A
ring consists of k nodes each with a single predecessor and a single successor
such that the last node's successor is the first node in the system.

Data keys are consistently hashed to produce an address on the ring. There are
N possible addresses (implemented as a 64 bit value). Each node is assigned a
base address and is responsible for managing data assigned to any address
equal to or less than this address but greater than the preceding node's base
address.

Petrel can be organised in a static topology but also supports dynamic
addition and removal of nodes, automatically rebalancing data around the ring.
Rebalancing is simplified by dividing the address space up into a fixed number
of partitions (n). It is expected that n is much larger than k. The number of
partitions is fixed on initialization of the system and cannot be changed
during its lifetime without rehashing. Partitions are assigned an address plus
a number of preceding addresses (i.e. N/n addresses).

The base address of a node must correspond to the address of a partition and
the node is responsible for the management of all partitions that fall into
its section of the address space (i.e. n/k partitions).

The data for a single partition is held in a single file that can be copied to
another node when it takes responsibility for the partition during
rebalancing. This avoids costly rehashing of keys and simplifies the dynamic
aspects of Petrel.

When a node (B) joins the ring it retrieves the current membership list from a
pre-existing node and locates the node with the highest number of existing
partitions (A). It then joins the ring as A's predecessor by assigning itself a
base address before A's and adopts the partitions ending at this address. It
then requests that A send it copies of the partitions just adopted. Thereafter
it becomes responsible for managing those partitions. During this period reads
and writes will be sent to A which will continue to serve read requests as
usual. However no mutations will be made to the data and write requests will
respond with a redirect to B. Any read or write requests made to B while it is
adopting the relevant partition will receive a "retry" response with a
suggested time interval, indicating that the client should queue the request
and retry after the interval has elapsed.

// TODO: leaving ring - planned and unplanned

// TODO: replication

// TODO: EDRA routing table maintenance

The following client operations are supported:

set - assigns the value stored for a key
get - retrieves the value stored for a key
del - deletes the value stored for a key
cas - compares and swaps the value stored for a key iff it matches the supplied value
members - retrieves the membership list of the ring
// TODO: apply a function to the stored value, e.g. increment


Note that any get/set requests sent to the wrong node will result in a
redirect response indicating the owner of the relevant partition.


*/
package petrel
