<a name="unreleased"></a>
## [Unreleased]

### Chore
- typo
- **core:** move admission UID retrieval ([#14](https://github.com/torbendury/kube-networkpolicy-denier/issues/14))

### Ci
- merge instead of rebase
- automatic resolve of conflicts
- integration/e2e test for helm chart Installing the Helm Chart and testing it against a Kubernetes cluster by applying a NetworkPolicy was a manual task until now. From now on, these steps will be done on every push to main automatically. Also, after a successful release, gh-pages branch will be automatically kept in-sync with main. Before, GitHub Actions denied triggering the pages CI to prevent recursion.
- rebase strategy
- explicit branch handling
- checkout separate main
- sync gh-pages with main
- explicit merge
- fast forward pull
- explicit pull of branches
- rename separate pipelines
- introduce changelog generation with git-chglog ([#13](https://github.com/torbendury/kube-networkpolicy-denier/issues/13))
- remove gh-pages sync
- rebase instead of merge
- allow unrelated main-gh-pages sync
- correct sync between main and gh-pages
- step naming
- skip existing helm chart version This bothered me for a while now, because not every CI run includes a new Helm release. However, main should be able to produce a release at every time and thus the goal is to keep CI green.
- correct Make handling
- correct diff exit code handling
- **docker:** change actual release stage

### Doc
- remove stability warning

### Feat
- shorten time interval for liveness and readiness probe ([#12](https://github.com/torbendury/kube-networkpolicy-denier/issues/12))
- **docker:** minify image size ([#15](https://github.com/torbendury/kube-networkpolicy-denier/issues/15))


<a name="kube-networkpolicy-denier-0.1.0"></a>
## [kube-networkpolicy-denier-0.1.0] - 2023-11-29
### Doc
- several status badges
- refine resource usage
- versioning advisory Added documentation for the projects' versioning approach. The Helm Chart uses SemVer and always contains a working and tested version of the controller. The controller image is still built with every pushed commit and is tested separately, so: a) not every image version will make it into a Helm release b) not every Helm release is expected to work with any version of the controller. Also added a quick hit on the supported and tested Kubernetes versions after some historical research on API releases by Kubernetes.

### Feat
- set default message on module load This is most due to the fact that unit tests - among other non-'main function'-invokers will not parse any actual flags to the 'respMsg'.
- **core:** allow custom response message The response message should be configurable for administrators. It provides a value that matches the past behavior and can be controlled via Helm values. If the message is empty, the response is handled by Kubernetes which tells the user that the controller denied the request without further explanation.


<a name="kube-networkpolicy-denier-0.0.5"></a>
## [kube-networkpolicy-denier-0.0.5] - 2023-11-25
### Chore
- bump helm chart version
- TODOs moved to GitHub Issues for better visibility of the projects state

### Ci
- introduce pre-commit hooks Writing code can be hard, so I want to take as much of the burden out of my brain as possible. pre-commit hooks allow me to implement ideas and still maintain a certain minimum level of standards that I might otherwise overlook. This commit introduces a set of pre-commit hooks I initially chose. They seem to work out for me and did not find anything worth a mention initially. However, the hooks run quite fast (<5s, faster than I can remember some of my Makefile steps) and don't stop me from committing fast.
- action for syncing gh-pages with main

### Feat
- **core:** timeout handling Golang net/http does not provide default timeouts, and so didn't the admission controller. Server timeouts for header/body reads and write as well as connection keepalive timeouts were defined and tested, as well as handler-specific timeouts were implemented. The handlers now listen on the handed request context which gets cancelled after a handler-specific time. The benchmark tests have been adjusted so we can handle big load situations. At about 15k req/s, the server starts to throw errors which I will have to investigate later.
- **helm:** faster readiness and liveness feedback loop

### Test
- better grown and sufficient load testing


<a name="kube-networkpolicy-denier-0.0.4"></a>
## [kube-networkpolicy-denier-0.0.4] - 2023-11-23
### Chore
- Helm Chart version

### Ci
- use negative matches on CI events
- ignore non-code and non-ci relevant repos for builds
- **make:** helm chart version

### Doc
- repo doc prettifying and helm chart doc adjustments

### Feat
- **core:** stop logging healthcheck requests In [#2](https://github.com/torbendury/kube-networkpolicy-denier/issues/2), I initially wanted to make health check logging optional. Now that I hat some time to think and investigate, the optimal solution would have been to implement a second optional handler which gets registered at start (depending on health-check logging being on/off). The overhead for the solution would have been to big, so I removed the logging functionatility frmo the health check completely.

### Testing
- add basic perf benchmarks


<a name="kube-networkpolicy-denier-0.0.3"></a>
## [kube-networkpolicy-denier-0.0.3] - 2023-11-23
### Chore
- conflict

### Ci
- **helm:** set resource namespace
- **make:** faster local testing with minikube and chained build

### Fix
- **core:** gracefully handle shutdowns When terminating Pods, Kubernetes sends a SIGTERM to the process. The net/http server running in the container did not handle such signals. Now, a goroutine in the background is listening for SIGTERM signals and tries to stop the server gracefully before the Pod is shutdown. This prevents flaky API calls from the Kubernetes apiserver and lets existing connections finish before the server is shutdown.


<a name="kube-networkpolicy-denier-0.0.2"></a>
## [kube-networkpolicy-denier-0.0.2] - 2023-11-22
### Chore
- update TODOs

### Ci
- refactor Makefile
- **helm:** adjust minor helm parts - Chart version for testing update - test Pod to use HTTPS and ignore the selfsigned cert - resource recommendation for cpu/mem

### Doc
- repo README adjustments
- Helm Chart installation commands

### Feat
- improve logging with deeper information

### Hack
- add simple networkpolicy for testing


<a name="kube-networkpolicy-denier-0.0.1"></a>
## kube-networkpolicy-denier-0.0.1 - 2023-11-22
### Chore
- TODOs
- update TODOs
- TODOs

### Ci
- build only release target in multi-stage docker build
- remove unneccesary git installation git seems to be installed in the ubuntu container anyway
- add helm chart linting to test section and push tags, docker containers and helm charts
- set github token for chart releaser as env var
- use correct chart-releaser config
- remove explicit chart repo url
- configure git
- rely on git short commit hash as release tag
- **gh-action:** configure newer version of chart releaser
- **gh-action:** use different action for version bumping and add release mechanism
- **gh-action:** newer version of chart releaser
- **github-actions:** initial throw
- **helm:** correct NOTES include
- **helm:** readme, notes, cleanup
- **helm:** first working helm chart

### Doc
- repo README
- code documentation for all function signatures

### Feat
- Initial commit This introduces my new project for experimenting with admission controllers in Kubernetes. It provides a very easy and understandable way of denying new NetworkPolicies by using the Kubernetes API. One does not need to touch the apiserver directly but can use the interface and effectively deny all new NetworkPolicies in the Kubernetes cluster. This can be adapted to basically every scenario.
- **Dockerfile:** refine build stage to correctly download third party modules
- **core:** logging and correct API implementation Implement some very basic logging for startup, errors and incoming requests. Since my first run of this on a minikube cluster and a deeper glance at the ValidationWebhookConfiguration API, I found out that it is not sufficient to just return a non-200 status code but one has to correctly implement the AdmissionReview API. I did this in the easiest + fastest possible way for now.


[Unreleased]: https://github.com/torbendury/kube-networkpolicy-denier/compare/kube-networkpolicy-denier-0.1.0...HEAD
[kube-networkpolicy-denier-0.1.0]: https://github.com/torbendury/kube-networkpolicy-denier/compare/kube-networkpolicy-denier-0.0.5...kube-networkpolicy-denier-0.1.0
[kube-networkpolicy-denier-0.0.5]: https://github.com/torbendury/kube-networkpolicy-denier/compare/kube-networkpolicy-denier-0.0.4...kube-networkpolicy-denier-0.0.5
[kube-networkpolicy-denier-0.0.4]: https://github.com/torbendury/kube-networkpolicy-denier/compare/kube-networkpolicy-denier-0.0.3...kube-networkpolicy-denier-0.0.4
[kube-networkpolicy-denier-0.0.3]: https://github.com/torbendury/kube-networkpolicy-denier/compare/kube-networkpolicy-denier-0.0.2...kube-networkpolicy-denier-0.0.3
[kube-networkpolicy-denier-0.0.2]: https://github.com/torbendury/kube-networkpolicy-denier/compare/kube-networkpolicy-denier-0.0.1...kube-networkpolicy-denier-0.0.2
