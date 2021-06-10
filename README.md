# Stock-Watcher

Not a watcher of market stocks, but websites and inventory!

## Usage

No distributions as of now but you can clone it and build it for your OS.

### Config

The configuration files looks like so. It's an array of sites you want to monitor

```
[
  {
    "url": "https://www.amd.com/en/direct-buy/5458372800/ca",
    "key": "AMD Radeon™ RX 6800 XT Graphics",
    "vendor": "AMD",
    "description": "AMD Radeon™ RX 6800 XT Graphics",
    "lookFor": ".product-out-of-stock",
    "onSuccessMessage": "IN stock",
    "onFailureMessage": "OUT of stock",
    "mailTo": "",
    "smsTo": "",
    "checkType": "className",
    "isNegativeCheck": true,
    "headers":[{key:"", value:""}]
  }
]
```

Definitions:

"url": The URL of the product you're monitoring

"key": The key that will show in the cli

"vendor": The vendor that's selling your object

"description": Some description. It's not used now

"lookFor": A key to look for on the page

"onSuccessMessage": What you want printed when it's a successful test

"onFailureMessage": What you want printed when it's an unsuccessful test

"mailTo": TODO - who to mail notifications of success to

"smsTo": TODO - who to send an SMS notification to

"checkType": Either:

- className : find an object by a class name. Helpful for sites with stylized "out of stock" messages
- test: find a term on the page. Useful for pages that show or hide Out of Stock text on their site
  "isNegativeCheck": IF you're polling for "out of stock" and you want a notification of in stock, you're checking the negative of your result. Set your OnSuccess to be the "in stock" message and call this a negative check :)

"headers": A Key/Value pair identifying custom headers to apply to a GET request

### CLI Arguments

`-c` The config file. Defaults to config.json in current directory

` -t` The ticker duration in seconds

## Mailer

To use the SMTP mailer, create a .env file in the same directory as your executable
`touch .env`

Add the following contents:

```
SMTP_HOST="your smtp hosdt"
SMTP_PORT="default port"
SMTP_EMAIL="your email account"
SMTP_PASSWORD="email account password"
```
