## Introduction

This Golang package holds the data model. It consists of four parts:

* Data definition of the state data on the blockchain.
* Data definition of the transaction payload.
* Data definition of the client-side data, shared by the Major Tool and the Client.
* Data definition of events sent from the blockchain to the client side.

The remainder of this file motivates the chosen data storage and marshaling tools.

## State Data on the Blockchain

We had to choose between JSON and Google Protocol Buffers. JSON is easier to debug, but it turns out that marshaling Golang variables into JSON is not deterministic as pointed out here: https://stackoverflow.com/questions/44755089/does-serialized-content-strictly-follow-the-order-in-definition-use-encoding-jso. Google Protocol Buffers is not deterministic too in theory, but in practice it is as long as no maps are applied in the GPB data definitions. This is explained here: https://havoc.io/post/deterministic-protobuf/. The choice for Google Protocol Buffers is based on these sources.

Using UUIDs, we can arrange that every Sawtooth address holds exactly one value. Each item mentioned in requirement AX-5050 has its id calculated as follows. When the item is created, the client doing so should first create a UUID. Then this id is hashed and the hex representation of the hash is taken. The digits a-f should be small caps. Then the last 62 digits are taken. The final address becomes:

<6-digid transaction family> + <2-digid type code> + <62-digid remainder>

The exception to this scheme is the price list, which has a fixed address. This allows the blockchain to find its bootstrap information.

## Transaction Payload

We chose Google Protocol Buffers because we did for state data.

## Client-side data

On the client side, searching data is important. We hold the data in a SQLite 3 database. This way, no remote database is needed. The data resides in a local file and can be maintained with SQL statements.

## Events

We choose Google Protocol Buffers because we did so for server state data and transaction payloads.

Sawtooth events have the following fields:

* The event type, a string.
* Attributes, which are name/value pairs.
* A byte array, which is opaque to Hyperledger Sawtooth.

The event type and the attributes can be used to filter events.

## Miscelaneous

Timestamps are stored as integer values, the number of seconds since Epoch. Timestamps are stored in 64-bit signed integers.

Optional fields do not need to be handled with null values. The empty string is good enough, because an empty string is not a meaningful value itself.
