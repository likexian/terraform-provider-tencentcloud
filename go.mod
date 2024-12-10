module github.com/tencentcloudstack/terraform-provider-tencentcloud

go 1.17

require (
	cloud.google.com/go/iam v1.0.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230923063757-afb1ddc0824c
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/aws/aws-sdk-go v1.36.30
	github.com/beevik/etree v1.2.0
	github.com/bflad/tfproviderlint v0.14.0
	github.com/client9/misspell v0.3.4
	github.com/fatih/color v1.15.0
	github.com/golangci/golangci-lint v1.52.2
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/hcl/v2 v2.13.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.20.0
	github.com/katbyte/terrafmt v0.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mozillazg/go-httpheader v0.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.2
	github.com/tencentcloud/tencentcloud-sdk-go-intl-en v3.0.646+incompatible
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos v1.0.799
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/api v1.0.285
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway v1.0.763
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm v1.0.825
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as v1.0.1052
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi v1.0.824
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam v1.0.1051
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat v1.0.825
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs v1.0.591
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb v1.0.944
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn v1.0.1010
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch v1.0.843
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs v1.0.627
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw v1.0.1018
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs v1.0.600
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ciam v1.0.695
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka v1.0.748
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb v1.0.1034
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit v1.0.1033
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls v1.0.1046
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.1052
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm v1.0.1014
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp v1.0.762
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb v1.0.692
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu v1.0.335
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain v1.0.652
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc v1.0.633
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb v1.0.673
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc v1.0.797
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.831
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/domain v1.0.414
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts v1.0.628
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb v1.0.760
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr v1.0.1004
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es v1.0.777
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap v1.0.970
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms v1.0.563
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse v1.0.729
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live v1.0.777
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb v1.0.672
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb v1.0.949
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor v1.0.844
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps v1.0.853
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization v1.0.1038
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres v1.0.1010
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns v1.0.1038
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts v1.0.762
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis v1.0.1008
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum v1.0.744
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf v1.0.1034
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses v1.0.748
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.0.486
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver v1.0.816
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl v1.0.1015
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm v1.0.691
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts v1.0.524
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag v1.0.860
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat v1.0.634
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb v1.0.199
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm v1.0.547
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr v1.0.696
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg v1.0.533
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq v1.0.955
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem v1.0.578
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo v1.0.1013
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke v1.0.1038
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket v1.0.947
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse v1.0.857
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf v1.0.674
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod v1.0.860
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc v1.0.1045
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf v1.0.1037
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata v1.0.792
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss v1.0.199
	github.com/tencentyun/cos-go-sdk-v5 v0.7.58
	github.com/yangwenmai/ratelimit v0.0.0-20180104140304-44221c2292e1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/hashicorp/go-uuid v1.0.3
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc v1.0.970
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris v1.0.993
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg v1.0.772
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/controlcenter v1.0.993
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/csip v1.0.860
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb v1.0.970
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus v1.0.831
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region v1.0.920
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcss v1.0.1031
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/thpc v1.0.998
	github.com/wI2L/jsondiff v0.3.0
	k8s.io/apimachinery v0.22.4
)

require (
	4d63.com/gocheckcompilerdirectives v1.2.1 // indirect
	4d63.com/gochecknoglobals v0.2.1 // indirect
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/compute v1.18.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/storage v1.28.1 // indirect
	github.com/Abirdcfly/dupword v0.0.11 // indirect
	github.com/Antonboom/errname v0.1.9 // indirect
	github.com/Antonboom/nilnil v0.1.3 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/Djarvur/go-err113 v0.0.0-20210108212216-aea10b59be24 // indirect
	github.com/GaijinEntertainment/go-exhaustruct/v2 v2.3.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/OpenPeeDeeP/depguard v1.1.1 // indirect
	github.com/alexkohler/prealloc v1.0.0 // indirect
	github.com/alingse/asasalint v0.0.11 // indirect
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/apparentlymart/go-textseg v1.0.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/ashanbrown/forbidigo v1.5.1 // indirect
	github.com/ashanbrown/makezero v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bflad/gopaniccheck v0.1.0 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/bkielbasa/cyclop v1.2.0 // indirect
	github.com/blizzy78/varnamelen v0.8.0 // indirect
	github.com/bombsimon/wsl/v3 v3.4.0 // indirect
	github.com/breml/bidichk v0.2.4 // indirect
	github.com/breml/errchkjson v0.3.1 // indirect
	github.com/butuzov/ireturn v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/charithe/durationcheck v0.0.10 // indirect
	github.com/chavacava/garif v0.0.0-20230227094218-b8c73b2037b8 // indirect
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/curioswitch/go-reassign v0.2.0 // indirect
	github.com/daixiang0/gci v0.10.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denis-tingaikin/go-header v0.4.3 // indirect
	github.com/esimonov/ifshort v1.0.4 // indirect
	github.com/ettle/strcase v0.1.1 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/firefart/nonamedreturns v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/fzipp/gocyclo v0.6.0 // indirect
	github.com/go-critic/go-critic v0.7.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/go-toolsmith/astequal v1.1.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/go-xmlfmt/xmlfmt v1.1.2 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golangci/check v0.0.0-20180506172741-cfe4005ccda2 // indirect
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a // indirect
	github.com/golangci/go-misc v0.0.0-20220329215616-d24fe342adfe // indirect
	github.com/golangci/gofmt v0.0.0-20220901101216-f2edd75033f2 // indirect
	github.com/golangci/lint-1 v0.0.0-20191013205115-297bf364a8e0 // indirect
	github.com/golangci/maligned v0.0.0-20180506175553-b1d89398deca // indirect
	github.com/golangci/misspell v0.4.0 // indirect
	github.com/golangci/revgrep v0.0.0-20220804021717-745bb2f7c2e6 // indirect
	github.com/golangci/unconvert v0.0.0-20180507085042-28b1c447d1f4 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.7.1 // indirect
	github.com/gookit/color v1.5.2 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20230107090616-13ace0543b28 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.4.2 // indirect
	github.com/gostaticanalysis/forcetypeassert v0.1.0 // indirect
	github.com/gostaticanalysis/nilerr v0.1.1 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320 // indirect
	github.com/hashicorp/go-getter v1.4.0 // indirect
	github.com/hashicorp/go-hclog v1.2.1 // indirect
	github.com/hashicorp/go-plugin v1.4.4 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hc-install v0.4.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-config-inspect v0.0.0-20191115094559-17f92b0546e8 // indirect
	github.com/hashicorp/terraform-exec v0.17.2 // indirect
	github.com/hashicorp/terraform-json v0.14.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.12.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.7.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.7.0 // indirect
	github.com/hashicorp/terraform-plugin-test v1.2.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.0.0-20220623143253-7d51757b572c // indirect
	github.com/hashicorp/terraform-svchost v0.0.0-20200729002733-f050f53b9734 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jgautheron/goconst v1.5.1 // indirect
	github.com/jingyugao/rowserrcheck v1.1.1 // indirect
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julz/importas v0.1.0 // indirect
	github.com/junk1tm/musttag v0.5.0 // indirect
	github.com/kisielk/errcheck v1.6.3 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/kkHAIKE/contextcheck v1.1.4 // indirect
	github.com/kulti/thelper v0.6.3 // indirect
	github.com/kunwardeep/paralleltest v1.0.6 // indirect
	github.com/kyoh86/exportloopref v0.1.11 // indirect
	github.com/ldez/gomoddirectives v0.2.3 // indirect
	github.com/ldez/tagliatelle v0.4.0 // indirect
	github.com/leonklingele/grouper v1.1.1 // indirect
	github.com/lufeee/execinquery v1.2.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/maratori/testableexamples v1.0.0 // indirect
	github.com/maratori/testpackage v1.1.1 // indirect
	github.com/matoous/godox v0.0.0-20230222163458-006bad1f9d26 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mbilski/exhaustivestruct v1.2.0 // indirect
	github.com/mgechev/revive v1.3.1 // indirect
	github.com/mitchellh/cli v1.0.0 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/moricho/tparallel v0.3.1 // indirect
	github.com/nakabonne/nestif v0.3.1 // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354 // indirect
	github.com/nishanths/exhaustive v0.9.5 // indirect
	github.com/nishanths/predeclared v0.2.2 // indirect
	github.com/nunnatsa/ginkgolinter v0.9.0 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polyfloyd/go-errorlint v1.4.0 // indirect
	github.com/posener/complete v1.2.1 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/quasilyte/go-ruleguard v0.3.19 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20210819130434-b3f0c404a727 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/ryancurrah/gomodguard v1.3.0 // indirect
	github.com/ryanrolds/sqlclosecheck v0.4.0 // indirect
	github.com/sanposhiho/wastedassign/v2 v2.0.7 // indirect
	github.com/sashamelentyev/interfacebloat v1.1.0 // indirect
	github.com/sashamelentyev/usestdlibvars v1.23.0 // indirect
	github.com/securego/gosec/v2 v2.15.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shazow/go-diff v0.0.0-20160112020656-b6b7b6733b8c // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/sivchari/containedctx v1.0.2 // indirect
	github.com/sivchari/nosnakecase v1.7.0 // indirect
	github.com/sivchari/tenv v1.7.1 // indirect
	github.com/smartystreets/goconvey v1.8.0 // indirect
	github.com/sonatard/noctx v0.0.2 // indirect
	github.com/sourcegraph/go-diff v0.7.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/cobra v1.6.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/ssgreg/nlreturn/v2 v2.2.1 // indirect
	github.com/stbenjam/no-sprintf-host-port v0.1.1 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/t-yuki/gocover-cobertura v0.0.0-20180217150009-aaee18c8195c // indirect
	github.com/tdakkota/asciicheck v0.2.0 // indirect
	github.com/tetafro/godot v1.4.11 // indirect
	github.com/tidwall/gjson v1.17.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/timakin/bodyclose v0.0.0-20221125081123-e39cf3fc478e // indirect
	github.com/timonwong/loggercheck v0.9.4 // indirect
	github.com/tomarrell/wrapcheck/v2 v2.8.1 // indirect
	github.com/tommy-muehle/go-mnd/v2 v2.5.1 // indirect
	github.com/ulikunitz/xz v0.5.5 // indirect
	github.com/ultraware/funlen v0.0.3 // indirect
	github.com/ultraware/whitespace v0.0.5 // indirect
	github.com/uudashr/gocognit v1.0.6 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.1 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	github.com/yagipy/maintidx v1.0.0 // indirect
	github.com/yeya24/promlinter v0.2.0 // indirect
	github.com/zclconf/go-cty v1.10.0 // indirect
	github.com/zclconf/go-cty-yaml v1.0.1 // indirect
	gitlab.com/bosi/decorder v0.2.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	golang.org/x/exp/typeparams v0.0.0-20230224173230-c95f2b4c22f2 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.8.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.114.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230320184635-7606e756e683 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/tools v0.4.3 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	mvdan.cc/gofumpt v0.4.0 // indirect
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed // indirect
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20221223090309-7455f1af531d // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
