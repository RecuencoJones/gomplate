---
title: aws functions
menu:
  main:
    parent: functions
---

The functions in the `aws` namespace interface with various Amazon Web Services
APIs to make it possible for a template to render differently based on the AWS
environment and metadata.

### Configuring AWS

A number of environment variables can be used to control how gomplate communicates
with AWS APIs. A few are documented here for convenience. See [the `aws-sdk-go` documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)
for details.

| Environment Variable | Description |
| -------------------- | ----------- |
| `AWS_TIMEOUT` | _(Default `500`)_ Adjusts timeout for API requests, in milliseconds. Not part of the AWS SDK. |
| `AWS_PROFILE` | Profile name the SDK should use when loading shared config from the configuration files. If not provided `default` will be used as the profile name. |
| `AWS_REGION` | Specifies where to send requests. See [this list](https://docs.aws.amazon.com/general/latest/gr/rande.html). Note that the region must be set for AWS functions to work correctly, either through this variable, or a configuration profile. |

## `aws.EC2Meta`

**Alias:** `ec2meta`

Queries AWS [EC2 Instance Metadata](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) for information. This only retrieves data in the `meta-data` path -- for data in the `dynamic` path use `aws.EC2Dynamic`.

For times when running outside EC2, or when the metadata API can't be reached, a `default` value can be provided.

### Usage

```go
aws.EC2Meta key [default]
```

### Arguments

| name | description |
|------|-------------|
| `key` | _(required)_ the metadata key to query |
| `default` | _(optional)_ the default value |

### Examples

```console
$ echo '{{aws.EC2Meta "instance-id"}}' | gomplate
i-12345678
```

## `aws.EC2Dynamic`

**Alias:** `ec2dynamic`

Queries AWS [EC2 Instance Dynamic Metadata](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) for information. This only retrieves data in the `dynamic` path -- for data in the `meta-data` path use `aws.EC2Meta`.

For times when running outside EC2, or when the metadata API can't be reached, a `default` value can be provided.

### Usage

```go
aws.EC2Dynamic key [default]
```

### Arguments

| name | description |
|------|-------------|
| `key` | _(required)_ the dynamic metadata key to query |
| `default` | _(optional)_ the default value |

### Examples

```console
$ echo '{{ (aws.EC2Dynamic "instance-identity/document" | json).region }}' | gomplate
us-east-1
```

## `aws.EC2Region`

**Alias:** `ec2region`

Queries AWS to get the region. An optional default can be provided, or returns
`unknown` if it can't be determined for some reason.

### Usage

```go
aws.EC2Region [default]
```

### Arguments

| name | description |
|------|-------------|
| `default` | _(optional)_ the default value |

### Examples

_In EC2_
```console
$ echo '{{ aws.EC2Region }}' | ./gomplate
us-east-1
```
_Not in EC2_
```console
$ echo '{{ aws.EC2Region }}' | ./gomplate
unknown
$ echo '{{ aws.EC2Region "foo" }}' | ./gomplate
foo
```

## `aws.EC2Tag`

**Alias:** `ec2tag`

Queries the AWS EC2 API to find the value of the given [user-defined tag](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html). An optional default
can be provided.

### Usage

```go
aws.EC2Tag tag [default]
```

### Arguments

| name | description |
|------|-------------|
| `tag` | _(required)_ the tag to query |
| `default` | _(optional)_ the default value |

### Examples

```console
$ echo 'This server is in the {{ aws.EC2Tag "Account" }} account.' | ./gomplate
foo
```
```console
$ echo 'I am a {{ aws.EC2Tag "classification" "meat popsicle" }}.' | ./gomplate
I am a meat popsicle.
```

## `aws.KMSEncrypt`

Encrypt an input string with the AWS Key Management Service (KMS).

At most 4kb (4096 bytes) of data may be encrypted.

The resulting ciphertext will be base-64 encoded.

The `keyID` parameter is used to reference the Customer Master Key to use,
and can be:

- the key's ID (e.g. `1234abcd-12ab-34cd-56ef-1234567890ab`)
- the key's ARN (e.g. `arn:aws:kms:us-east-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab`)
- the alias name (aliases must be prefixed with `alias/`, e.g. `alias/ExampleAlias`)
- the alias ARN (e.g. `arn:aws:kms:us-east-2:111122223333:alias/ExampleAlias`)

For information on creating keys, see [_Creating Keys_](https://docs.aws.amazon.com/kms/latest/developerguide/create-keys.html)

See [the AWS documentation](https://docs.aws.amazon.com/kms/latest/developerguide/overview.html)
for more details.

See also [`aws.KMSDecrypt`](#aws-kmsdecrypt).

### Usage

```go
aws.KMSEncrypt keyID input
```
```go
input | aws.KMSEncrypt keyID
```

### Arguments

| name | description |
|------|-------------|
| `keyID` | _(required)_ the ID of the Customer Master Key (CMK) to use for encryption |
| `input` | _(required)_ the string to encrypt |

### Examples

```console
$ export CIPHER=$(gomplate -i '{{ aws.KMSEncrypt "alias/gomplate" "hello world" }}')
$ gomplate -i '{{ env.Getenv "CIPHER" | aws.KMSDecrypt }}'
```

## `aws.KMSDecrypt`

Decrypt ciphertext that was encrypted with the AWS Key Management Service
(KMS).

The ciphertext must be base-64 encoded.

See [the AWS documentation](https://docs.aws.amazon.com/kms/latest/developerguide/overview.html)
for more details.

See also [`aws.KMSEncrypt`](#aws-kmsencrypt).

### Usage

```go
aws.KMSDecrypt input
```
```go
input | aws.KMSDecrypt
```

### Arguments

| name | description |
|------|-------------|
| `input` | _(required)_ the base-64 encoded ciphertext to decrypt |

### Examples

```console
$ export CIPHER=$(gomplate -i '{{ aws.KMSEncrypt "alias/gomplate" "hello world" }}')
$ gomplate -i '{{ env.Getenv "CIPHER" | aws.KMSDecrypt }}'
```

## `aws.Account`

Returns the currently-authenticated AWS account ID number.

Wraps the [STS GetCallerIdentity API](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetCallerIdentity.html)

See also [`aws.UserID`](#aws-userid) and [`aws.ARN`](#aws-arn).

### Usage

```go
aws.Account
```


### Examples

```console
$ gomplate -i 'My account is {{ aws.Account }}'
My account is 123456789012
```

## `aws.ARN`

Returns the AWS ARN (Amazon Resource Name) associated with the current authentication credentials.

Wraps the [STS GetCallerIdentity API](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetCallerIdentity.html)

See also [`aws.UserID`](#aws-userid) and [`aws.Account`](#aws-account).

### Usage

```go
aws.ARN
```


### Examples

```console
$ gomplate -i 'Calling from {{ aws.ARN }}'
Calling from arn:aws:iam::123456789012:user/Alice
```

## `aws.UserID`

Returns the unique identifier of the calling entity. The exact value
depends on the type of entity making the call. The values returned are those
listed in the `aws:userid` column in the [Principal table](http://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_variables.html#principaltable)
found on the Policy Variables reference page in the IAM User Guide.

Wraps the [STS GetCallerIdentity API](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetCallerIdentity.html)

See also [`aws.ARN`](#aws-arn) and [`aws.Account`](#aws-account).

### Usage

```go
aws.UserID
```


### Examples

```console
$ gomplate -i 'I am {{ aws.UserID }}'
I am AIDACKCEVSQ6C2EXAMPLE
```
