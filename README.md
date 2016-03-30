# Luzifer / mondash-opsgenie

This small helper allows you to get notified about bad things happening on your [MonDash](https://mondash.org/) dashboard using [OpsGenie](https://www.opsgenie.com/). To make use of it you need to have version 1.7.0 or higher of MonDash running or use the hosted version at [mondash.org](https://mondash.org/).

## Usage

You will need following ingrediences:

- An OpsGenie API key from an "API" integration
- The board ID of your MonDash dashboard
- The URL of the MonDash instance (optional)

To set up just create a simple cronjob:

```bash
$ mondash-opsgenie --mondash-board=<your-board-id> --opsgenie-key=<your-opsgenie-key>
```

This will fetch the current data of the dashboard, evaluate whether there are metrics needing an alert and then create / close the alert in your OpsGenie account.

For detailed settings please refer to the `--help`:

```bash
# mondash-opsgenie --help
Usage of mondash-opsgenie:
  -a, --alert=[Critical,Unknown]: List of status to trigger alerts for
      --mondash-board="": Board ID of the mondash board to monitor
      --mondash-url="https://mondash.org/": URL of the mondash installation
  -o, --ok=[OK]: List of status to resolve alerts for
      --opsgenie-key="": API-Key for OpsGenie API integration
```
