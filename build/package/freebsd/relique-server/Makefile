PORTNAME=			relique-server
DISTVERSION=		0.0.19
DISTVERSIONPREFIX= 	v
CATEGORIES=			sysutils

MAINTAINER=		relique@mathieu.macarrie.fr
COMMENT=		Rsync based backup utility

BUILD_DEPENDS=	bash:shells/bash
RUN_DEPENDS=	rsync:net/rsync
USES=			go:modules sqlite:3

USE_GITHUB=		yes
GH_ACCOUNT=		macarrie
GH_PROJECT=		relique
# Generated with make gomod-vendor
GH_TUPLE=		macarrie:relique-module-generic:0.0.1:reliquemodulegeneric/configs/var/lib/relique/default_modules/relique-module-generic \
                Masterminds:squirrel:v1.5.2:masterminds_squirrel/vendor/github.com/Masterminds/squirrel \
                Microsoft:go-winio:v0.4.16:microsoft_go_winio/vendor/github.com/Microsoft/go-winio \
                ProtonMail:go-crypto:04723f9f07d7:protonmail_go_crypto/vendor/github.com/ProtonMail/go-crypto \
                acomagu:bufpipe:v1.0.3:acomagu_bufpipe/vendor/github.com/acomagu/bufpipe \
                emirpasic:gods:v1.12.0:emirpasic_gods/vendor/github.com/emirpasic/gods \
                fsnotify:fsnotify:v1.5.1:fsnotify_fsnotify/vendor/github.com/fsnotify/fsnotify \
                gin-contrib:sse:v0.1.0:gin_contrib_sse/vendor/github.com/gin-contrib/sse \
                gin-gonic:gin:v1.7.7:gin_gonic_gin/vendor/github.com/gin-gonic/gin \
                go-git:gcfg:v1.5.0:go_git_gcfg/vendor/github.com/go-git/gcfg \
                go-git:go-billy:v5.3.1:go_git_go_billy_v5/vendor/github.com/go-git/go-billy/v5 \
                go-git:go-git:v5.4.2:go_git_go_git_v5/vendor/github.com/go-git/go-git/v5 \
                go-ini:ini:v1.66.2:go_ini_ini/vendor/gopkg.in/ini.v1 \
                go-playground:locales:v0.13.0:go_playground_locales/vendor/github.com/go-playground/locales \
                go-playground:universal-translator:v0.17.0:go_playground_universal_translator/vendor/github.com/go-playground/universal-translator \
                go-playground:validator:v10.4.1:go_playground_validator_v10/vendor/github.com/go-playground/validator/v10 \
                go-warnings:warnings:v0.1.2:go_warnings_warnings/vendor/gopkg.in/warnings.v0 \
                go-yaml:yaml:v2.4.0:go_yaml_yaml/vendor/gopkg.in/yaml.v2 \
                golang:crypto:32db794688a5:golang_crypto/vendor/golang.org/x/crypto \
                golang:net:60bc85c4be6d:golang_net/vendor/golang.org/x/net \
                golang:protobuf:v1.5.2:golang_protobuf/vendor/github.com/golang/protobuf \
                golang:sys:1c1b9b1eba6a:golang_sys/vendor/golang.org/x/sys \
                golang:text:v0.3.7:golang_text/vendor/golang.org/x/text \
                google:uuid:v1.3.0:google_uuid/vendor/github.com/google/uuid \
                hashicorp:errwrap:v1.0.0:hashicorp_errwrap/vendor/github.com/hashicorp/errwrap \
                hashicorp:go-multierror:v1.1.1:hashicorp_go_multierror/vendor/github.com/hashicorp/go-multierror \
                hashicorp:hcl:v1.0.0:hashicorp_hcl/vendor/github.com/hashicorp/hcl \
                imdario:mergo:v0.3.12:imdario_mergo/vendor/github.com/imdario/mergo \
                inconshreveable:mousetrap:v1.0.0:inconshreveable_mousetrap/vendor/github.com/inconshreveable/mousetrap \
                jbenet:go-context:d14ea06fba99:jbenet_go_context/vendor/github.com/jbenet/go-context \
                json-iterator:go:v1.1.12:json_iterator_go/vendor/github.com/json-iterator/go \
                kennygrant:sanitize:v1.2.4:kennygrant_sanitize/vendor/github.com/kennygrant/sanitize \
                kevinburke:ssh_config:4977a11b4351:kevinburke_ssh_config/vendor/github.com/kevinburke/ssh_config \
                lann:builder:47ae307949d0:lann_builder/vendor/github.com/lann/builder \
                lann:ps:62de8c46ede0:lann_ps/vendor/github.com/lann/ps \
                leodido:go-urn:v1.2.0:leodido_go_urn/vendor/github.com/leodido/go-urn \
                magiconair:properties:v1.8.5:magiconair_properties/vendor/github.com/magiconair/properties \
                mattn:go-isatty:v0.0.14:mattn_go_isatty/vendor/github.com/mattn/go-isatty \
                mattn:go-sqlite3:v1.14.11:mattn_go_sqlite3/vendor/github.com/mattn/go-sqlite3 \
                mitchellh:go-homedir:v1.1.0:mitchellh_go_homedir/vendor/github.com/mitchellh/go-homedir \
                mitchellh:mapstructure:v1.4.3:mitchellh_mapstructure/vendor/github.com/mitchellh/mapstructure \
                modern-go:concurrent:bacd9c7ef1dd:modern_go_concurrent/vendor/github.com/modern-go/concurrent \
                modern-go:reflect2:v1.0.2:modern_go_reflect2/vendor/github.com/modern-go/reflect2 \
                pelletier:go-toml:v1.9.4:pelletier_go_toml/vendor/github.com/pelletier/go-toml \
                pkg:errors:v0.9.1:pkg_errors/vendor/github.com/pkg/errors \
                protocolbuffers:protobuf-go:v1.27.1:protocolbuffers_protobuf_go/vendor/google.golang.org/protobuf \
                sergi:go-diff:v1.1.0:sergi_go_diff/vendor/github.com/sergi/go-diff \
                sirupsen:logrus:v1.8.1:sirupsen_logrus/vendor/github.com/sirupsen/logrus \
                spf13:afero:v1.6.0:spf13_afero/vendor/github.com/spf13/afero \
                spf13:cast:v1.4.1:spf13_cast/vendor/github.com/spf13/cast \
                spf13:cobra:v1.3.0:spf13_cobra/vendor/github.com/spf13/cobra \
                spf13:jwalterweatherman:v1.1.0:spf13_jwalterweatherman/vendor/github.com/spf13/jwalterweatherman \
                spf13:pflag:v1.0.5:spf13_pflag/vendor/github.com/spf13/pflag \
                spf13:viper:v1.10.1:spf13_viper/vendor/github.com/spf13/viper \
                subosito:gotenv:v1.2.0:subosito_gotenv/vendor/github.com/subosito/gotenv \
                ugorji:go:v1.1.7:ugorji_go_codec/vendor/github.com/ugorji/go \
                xanzy:ssh-agent:v0.3.0:xanzy_ssh_agent/vendor/github.com/xanzy/ssh-agent

USERS=			relique
GROUPS=			relique

do-build:
	${SETENV} XDG_CACHE_HOME=${WRKDIR}/.cache GOMODCACHE=${WRKDIR}/gomodcache ${GO_ENV} make -C ${WRKSRC} clean build_server BUILD_OUTPUT_DIR=${WRKDIR}/package

do-install:
	make -C ${WRKSRC} install INSTALL_ROOT="${STAGEDIR}/" INSTALL_SRC=${WRKDIR}/package INSTALL_ARGS="--server --freebsd --skip-user-creation"

post-install:
	# Strip relique binaries
	${STRIP_CMD} ${STAGEDIR}${PREFIX}/bin/relique-server

.include <bsd.port.mk>
