# Kubeconfig-merger

Merges discrete kubeconfig files into `~/.kube/config`.
By default, this application will look for kubeconfig files in the `~/.kube/config.d` directory.

## Usage

```bash
./kubeconfig-merger
```

Specifying a source directory other than `~/.kube/config.d`

```bash
./kubeconfig-merger -source /path/to/config/dir
```
