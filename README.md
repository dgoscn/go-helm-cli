# Helm CLI

Helm CLI is a command-line created in Golang  that will take as input a list of Helm Charts GitHub repo url or a local folder. i.e:

path (Ex: https://github.com/epinio/helm-charts/tree/main/chart/epinio or c:/epinio/helm-
charts/chart) .

## Expected Outputs

The binary should behave in this manner.

**Command** - <binary> add <chart-location>

**Result** - Adds the helm chart information to the CLI’s internal list and a storage of your choice.

**Command** - <binary> index

**Result** – Generates Helm Repo Index file

**Command** - <binary> install chart <chart-name>

Result - Installs the given helm chart in the current Kubernetes cluster. Installation of the helm chart
must happen inside a Kubernetes pod.

**Command** - <binary> images

**Result** - Provides a list of all the container images used in all the charts added.

## Prerequisites

- Docker: Ensure that Docker is installed and running on your machine.

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/<username>/helm-cli.git
   ```

2. Build the Docker image:

   ```bash
   cd helm-cli
   docker build -t helm-cli .
   ```

3. Run the Helm CLI commands using Docker:

   ```bash
   docker run -it --rm -v ~/.kube/config:/root/.kube/config -v $(pwd)/charts:/app/charts helm-cli <command>
   ```

   Replace `<command>` with one of the supported commands: `add`, `index`, `install`, or `images`. See the Usage section below for more details on each command.

## Usage

- **Add a Helm Chart:**

  ```bash
  docker run -it --rm -v ~/.kube/config:/root/.kube/config -v $(pwd)/charts:/app/charts helm-cli add <chart-location>
  ```
  So, after filled:

  ```bash
  docker run -it --rm -v ~/.kube/config:/root/.kube/config -v $(pwd)/charts:/app/charts helm-cli add https://github.com/epinio/helm-charts.git
  ```

  This command adds the Helm chart information from the specified `<chart-location>` to the CLI's internal list and stores it in the `charts` directory.

- **Generate Repo Index:**

  ```bash
  docker run -it --rm -v ~/.kube/config:/root/.kube/config -v $(pwd)/charts:/app/charts helm-cli index
  ```

  This command generates the Helm repository index file based on the charts added using the `add` command.

- **Install a Helm Chart:**

  ```bash
  docker run -it --rm -v ~/.kube/config:/root/.kube/config -v $(pwd)/charts:/app/charts helm-cli install chart <chart-name>
  ```

  In that case, we have the example with epinio address right below

  ```bash
  docker run -it --rm -v ~/.kube/config:/root/.kube/config -v $(pwd)/charts:/app/charts helm-cli install chart epinio
  ```

- **List Container Images:**

  ```bash
  docker run -it --rm -v ~/.kube/config:/root/.kube/config -v $(pwd)/charts:/app/charts helm-cli images
  ```

  This command provides a list of all the container images used in the charts added.

## Run the tests

You just need to go to the root file of the project and run the following command:

  ```bash
   > go test
  ```
