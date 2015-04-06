# prime_finder
Demonstration web service for finding the prime number less than or equal to an upper bound.

Once deployed, requests should be in the form of:
```
http://localhost:8080/prime?limit=1000
```
The service will respond with a JSON structure:

```
[
   {
      "sieve":"sundaram",
      "prime":997,
      "limit":1000,
      "compute_time_sec":0.000244736
   },
   {
      "sieve":"eratosthenes",
      "prime":997,
      "limit":1000,
      "compute_time_sec":0.00018688
   }
]
```
