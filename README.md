# Parqet to YNAB

A CLI tool that lets you sync your Parqet portfolio balance to a YNAB budget of your choice.

**You need to have your Parqet portfolio public for this method to work since there is currently no public API available.**

## Usage

### Setting up the config

```parqet-to-ynab setup```

The command will ask you for:

* Parqet portfolio ID
* YNAB API key
* YNAB budget ID
* YNAB account ID


### Running the sync

```parqet-to-ynab sync```

The command will check your portfolio balance and create a new adjusting transaction if the YNAB balance and Parqet balance differ.

## Installation

Just download the archive from the newest [release](https://github.com/lukkigi/parqet-to-ynab/releases), unpack it and run it from your terminal.