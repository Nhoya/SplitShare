# SplitShare
Shamir's Secret Sharing Algorithm implementation in Golang combined with PGP and a mail delivery system

SplitShare uses [Shamir's Secret Sharing Algorithm](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) to split a secret in multiple parts and [PGP](https://it.wikipedia.org/wiki/Pretty_Good_Privacy) to encrypt each share with the public key of the owner and then send the encrypted share via mail.

## Usage
You can download and run [the pre-builted binary](https://github.com/Nhoya/SplitShare/releases/latest) or build it from source.

To work properly SplitShare needs:

- The public key of each user
- The secret stored in a file
- A working SMTP server

`./splitshare <SECRET FILE>`

Starting SplitShare will ask for:

- Maximum number of shares
- Minumum amount of shares needed to decrypt the secret

And for each user:
- The users mail
- The path of the public key

Alternatively you can load a public keyring file with users public key and SplitShare will automatically load keys from user's email.

`./splitshare --pub-keyring ./pubring.gpg <SECRET FILE>`

*NOTE*: you must configure the SMTP credentials inside the `config.toml` file.

## Building
To build SplitShare you need the following packages:

```
github.com/SSSaaS/sssa-golang
github.com/TheZ3ro/go-pgp/pgp
github.com/go-mail/mail
github.com/spf13/viper
```

You can then clone the repo and build it with:

```
git clone https://github.com/Nhoya/SplitShare
go build splitshare.go
```

## Decrypt
To decrypt the secret you need the tool located inside the `decrypt` directory

You can either build it or use [the pre-builted one](https://github.com/Nhoya/SplitShare/releases/latest).

### Building the decryption tool
To build the decryption tool you just need:

`github.com/SSSaaS/sssa-golang`

You can then clone the repo and build it with 

```
git clone https://github.com/Nhoya/SplitShare
cd SplitShare/decrypt
go build decrypt.go
```

