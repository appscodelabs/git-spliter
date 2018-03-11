# git-spliter
Splits a Git repo into separate repos maintaining history

https://help.github.com/articles/splitting-a-subfolder-out-into-a-new-repository/

```
git-spliter init \
  --repo-dir="/home/tamal/go/src/k8s.io/kubernetes" \
  --pkg-dir="cluster/addons" \
  --github-username=kube-addons \
  --github-topics=kubernetes,addons,pharmer,appscode
```
