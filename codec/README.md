## Related Interfaces

The following diagram shows the server-side protocol processing flow, which includes the related interfaces in the `codec` package.

```ascii
                              package                     req body                                                       req struct
+-------+        +-------+    []byte     +--------------+  []byte    +-----------------------+    +----------------------+
|       +------->+ Framer +------------->| Codec-Decode +----------->| Compressor-Decompress +--->| Serializer-Unmarshal +------------+
|       |        +-------+               +--------------+            +-----------------------+    +----------------------+            |
|       |                                                                                                                        +----v----+
|network|                                                                                                                        | Handler |
|       |                                                 rsp body                                                               +----+----+
|       |                                                  []byte                                                         rsp struct  |
|       |                                +---------------+           +---------------------+       +--------------------+             |
|       <--------------------------------+  Codec-Encode +<--------- + Compressor-Compress + <-----+ Serializer-Marshal +-------------+
+-------+                                +---------------+           +---------------------+       +--------------------+
```