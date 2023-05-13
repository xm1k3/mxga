![mxga](static/mxga_banner.png)

```
Usage:
  mxga [command]

Available Commands:
  help        Help about any command
  retrieve    retrieve back money from all wallets
  trx         Send multiple transactions
  wallet      Create new wallets
```

Important NOTE:

all the examples are referred to the testnet/devnet network, if you want yo test on devnet or mainnet you will need to change the mode on mxga.yaml configuration file

---

Mxga configuration file: `$HOME/.config/mxga/`

CONFIGURATION_PATH/mxga.yaml
```yaml
# available modes: mainnet, testnet, devnet, manual
mode: testnet

wallet:
  # Path of your pem files
  path: $HOME/.config/mxga/wallets/pem/
  # format (recommended .pem)
  ext: .pem
  # filename of the main address
  main: main_address
  other:
    # all other wallets generated
    # recommended to not change che file name and keep the address
    # this if for transactions optimizations
    - erd1k7sp2wv2wreq2hm52lqyws46rjpehj6fjvav4g5w5hrwvf8s4cxqwzjchc
    - erd1e96l2qkz8h9yuzupu05yeydjlrmy0fkdas98c5fp4x2gczh2lc3sh6qerc
    - erd1cy2eq9gsct6jlugrtemxj36q5vqdwmdzx83z4kcz3j3jfksryseq8fpzs6
    - erd19rhl5gt6u5k6xh5mkaf3kkduc2se023wpqhcn67s2kh5f8v6jqrqpl7hqc
    - erd1wdsu4tallng55s8ugnw7d9s7qvn2hegrnmu6mm5f68r59dzthmds45gp6g
    - erd1w2g8vlls0gup9jglmza7lks66ahsrzhp6k2dtyez0jtym5r4wqascpxxxw
    - erd1uhhrja9ymwpa3qrpkwaxk9ua3r0m00l4npuzdsa4t0vryp22rgtstnal39
    - erd1cjmvxmvzmet5q02tan35nnydpscen6swtmjz689mxu9dyjynkktq9l7ds8
    - erd1k72l5l4hxvl2yrar5qrz6yunluyakmc55qhlkdvy52k6lyjqdmhqvwvvqk
    - erd1qfzvxu439amtcf5ysl6hsck235yc0j3vwq3vrmgy69pdtr7e584qlz8jfp

```

## Create new wallets:

by default mxga will generate 1 wallet on your `CONFIGURATION_PATH/wallets`
- you can find all your .pem secrets on `/pem`
- you can find all your wallets in json format on `/json` (you can set a custom password)

```
Create new wallets

Usage:
  mxga wallet [flags]

Flags:
  -a, --amount int        Number of wallets to create (default 1)
  -h, --help              help for wallet
  -P, --password string   Default password for json wallet file (default "Password123")
```

examples:

generate 1 wallet
```
mxga wallet
```

generate 10 wallets
```
mxga wallet --amount 10
mxga wallet -a 10
```

generate 10 wallets with custom password:
```
mxga wallet --amount 10 --password YOUR_PASSWORD_HERE
mxga wallet -a 10 -p YOUR_PASSWORD_HERE
```

## Make multiples transactions

Mxga can do many transactions in seconds and manage them in the same time

```
Send multiple transactions

Usage:
  mxga trx [flags]

Flags:
  -d, --data string     data
  -h, --help            help for trx
  -v, --value float32   value (default 0.1)
```

Examples:

Send 0.5 EGLD to all your wallets present in "other" section on the mxga.yaml

-d or --data if you want to add custom data to the transaction

```
mxga trx
mxga trx -v 0.5 -d "PoC trx"
```

Output:
```
[ success ] Hash:  318bc7d0f49cef1f019bb731b0b0e0e990b774aa15c906cfbe316915f90dc963
[ success ] Hash:  869ec999854a21164eaf6af92972e87d36918c79fd8d47f430646281a1203bdf
[ success ] Hash:  1d02913525cce49a70b794a7a30623a7aea40efc151ad92a37053f8a7e12b950
[ success ] Hash:  02a00f1bea3d36a7cc8bb7e973909459922612bd6456fa53bd5c3b731f499aea
[ success ] Hash:  aabc1a58ab1b528f038ab105767f2882dd8560a269cb875e66cc1210ca93b48e
[ success ] Hash:  ed91164f34d682ea3c1ec3db2673ec22e4353534340e412488af2be394f446fa
[ success ] Hash:  47ca96239a65230277ecc31aa22dd5c8001a4bd507b75fe08cad25d48c59a73e
[ success ] Hash:  ae70e75232ac1155da67253be73a0142da81046cc6b5ed275f2ca30752b442a7
[ success ] Hash:  b4ddb3d858db04535281c826fabdaa01e793ea9cfce58e3328a573d254842321
[ success ] Hash:  5817777ce218463d964f81b8e263166726d5d695de3b6c8c61fd9fc9eac27624
```
You can check this hash transactions on the elrond explorer

## Retrieve back from your wallets

You can also retrieve back your funds from your wallets with mxga

```
retrieve back money from all wallets

Usage:
  mxga retrieve [flags]

Flags:
  -a, --all             retrieve all money from all wallets
  -d, --data string     data
  -h, --help            help for retrieve
  -v, --value float32   value (default 0.1)
```

if you want custom amount you can use -v 0.23 or if you want back all your money you can use -a or --all

IMPORTANT: this works only with egld, it's a direct transaction not a swap

```
mxga retrieve --all -d "retrieve money PoC"
```

Output:
```
[ success ] Hash:  5df9e736d8a6523a506eb2bc5894e446daecb89cec4e926fc33a0dce2bbdfeda
[ success ] Hash:  34beb346d0e339df2202c9816d53707a7c1c70455cbf8309eb5011b43b2b26f4
[ success ] Hash:  5b198d84ea1758ca06919304a553d37bc74c026e6c9c396705cd101d29e8ab82
[ success ] Hash:  b027444418a34bf78a97df3832e0b9bca37ee30f689a42740c250b26f53fa882
[ success ] Hash:  ed3397c7ea3a2608755bd87e0c5aa113e5b0bf9299597b88387edbe92d287b40
[ success ] Hash:  caba047f966fac4842d911de15d152d6be9b604e261badbabbc49907d62a20e7
[ success ] Hash:  5b49983eac11f6be30fccdf289f8a530ca80131f7c4eb8705863c78486d5f2e3
[ success ] Hash:  6646300fb347021807dc169a6ea601e856c7c619ab23ee3d492f970fd56efdb0
[ success ] Hash:  910b7372d36a184a7f4fea869085f7757d6ebcbd6f0582f1b7cb9dfffb5a44c4
[ success ] Hash:  b5b9f880f99dee434fe43b6004736aac18a70f3bbe7028ae33e1b76f35321202
```


# License
mxga is distributed under Apache-2.0 License
