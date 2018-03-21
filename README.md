# git-spliter
Splits a Git repo into separate repos maintaining history

https://help.github.com/articles/splitting-a-subfolder-out-into-a-new-repository/

```console
git-spliter init \
  --repo-dir="/home/tamal/go/src/k8s.io/kubernetes" \
  --pkg-dir="cluster/addons" \
  --github-username=kube-addons \
  --github-topics=kubernetes,addons,pharmer,appscode
```

To split multiple folders into one repo, use the command from [here](https://stackoverflow.com/a/17867910/244009):

```console
$ cd repo
$ git filter-branch --index-filter 'git rm --cached -qr --ignore-unmatch -- . && git reset -q $GIT_COMMIT -- admission registry runtime workload' --prune-empty -- --all
```
