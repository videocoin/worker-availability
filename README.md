Worker uptime monitoring
===

## Collector component

Periodically stores worker into the database. Currently mongodb.

## Reporting server

Reporting server deployed at https://uptimereport.videocoin.net/. Uses oauth2 for the authentication.

The easiest way to get the report would to open the page `https://uptimereport.videocoin.net/report?start=2020-09-25T10:00:00-05:00&end=2020-10-25T10:00:00-05:00` in the browser and wait for report to generate.

Timestamps are in [rfc3339 format](https://tools.ietf.org/html/rfc3339).

The expected output would be csv report in format below:

```
worker,client_id,worker_address,configuration_hash,duration_online
BluBlu #3,49f36e11-28fc-43a2-b1d1-5313b5e7d1d7,0x4861847518FC4F6Bb58eDef4b59D19D891298194,0x102a34171d2dddfeb5bf5ef766602ebf3f7ce2fc82b52b61fb697b01c6a93f0e,7s
Wolf Crypto,4d710f22-1658-4eec-925a-7f815e03d66b,0x54Ce711A465ab11fa50504e962636F51cFfE5f83,0x920d3235571f1cdb33770f47e41ed559b9eedaf609da740e52b1f558a745fdf9,7s
D-Redcore,4f4087e9-d705-4ccf-8412-b354817d16fc,0x6E3472a5a3de7A9f333Fef38f913285261eb55B8,0xd5f70cf8da7fe213e2d13dea07c4d76c4db69330e8d744da1b0b5641b02093df,7s
Video-Korea1,f442d6da-e8df-4f4a-94d8-1e72d669509d,0x752157106A572b779ce5A3e1Bd32D096a5E687cf,0xed65f5afd7225e1c48e706c5ba67360dd5f548a42d9eed0706d8107e961a615e,7s
SunRay Alpha,fc8e2ef8-4cdf-4fa6-8a40-f1f239ce102a,0x8edf5671E6bf370fb7C2A1Cb8bA360C56Bed1D04,0x540ccaa54931f6f851e898d1e8e3b6d7673bc1d44d22cdc0df88e9693e75c5b6,7s
mainstream01,3a940856-9fff-4fba-ba66-9544caa7ac65,0x313CCd0dDc016899296CCDB5Ad89961A2F29AEc1,0xe8edd94d99df284f5e676594a9d960e4a557a1ef1958455881c0b50d4232c3ef,7s
```

## Generation of Payment file

Command format:
```
job incentives <payment-file> <start-time> <end-time>
```
Timestamps are in [rfc3339 format](https://tools.ietf.org/html/rfc3339).

Example:
```
job incentives test1.csv 2020-12-25T10:00:00-08:00 2021-01-21T22:40:00-08:00
```

Expected oputput format:
```
0x17E69e6E218aFD89163a07A3594E7723Ef96C08d,0.0320312249228395
0xc79E58d92ef1baf5a9937bB159f82fa677378598,0.03203506790123457
0x893e542B597a5e471C624a16806a16A0d237Eb6F,28988.91599863394
0x6E3472a5a3de7A9f333Fef38f913285261eb55B8,0.032995172003600824
0xe1a42050829744B7BA8438815019dbF47C94B384,0.03205043981481481
0x5c6e052cad3E7281e1Cce60d2C5b8BDBE8c5DdDe,1.3503086419753085
0x1ee8A8AB4f536332dD62d581Be3C393693c2132B,1.3503086419753085
0x83674506133d4D8cd3C4c31E8A9937470604dd2d,1.3503086419753085
0xDbc1C1a9c492aEa9C463b5fD83Cb835836787De7,0.0320312249228395
```
