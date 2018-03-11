# git-spliter
Splits a Git repo into separate repos maintaining history

https://help.github.com/articles/splitting-a-subfolder-out-into-a-new-repository/

```
git-spliter init \
  --github-topics=kubernetes,addons,pharmer,appscode \
  --github-username=kube-addons \
  --pkg-dir="cluster/addons" \
  --repo-dir="/home/tamal/go/src/k8s.io/kubernetes"
```
